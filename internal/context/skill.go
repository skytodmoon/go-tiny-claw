// internal/context/skill.go
package context

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// SkillMetadata 技能元数据（用于懒加载的轻量级信息）
type SkillMetadata struct {
	Name        string   // 技能名称
	Description string   // 触发描述
	FilePath    string   // 技能文件完整路径
	Tags        []string // 标签（可选）
}

// Skill 完整技能结构
type Skill struct {
	Metadata SkillMetadata
	Body     string // Markdown 正文指令（懒加载）
}

// SkillLoader 负责从本地文件系统中加载并解析符合规范的技能模板
type SkillLoader struct {
	workDir       string
	skillMetadata map[string]SkillMetadata // 技能名称 -> 元数据
}

func NewSkillLoader(workDir string) *SkillLoader {
	loader := &SkillLoader{
		workDir:       workDir,
		skillMetadata: make(map[string]SkillMetadata),
	}
	loader.scanSkillDirectory()
	return loader
}

// scanSkillDirectory 扫描技能目录，只解析 YAML 元数据（懒加载第一步）
func (s *SkillLoader) scanSkillDirectory() {
	skillBaseDir := filepath.Join(s.workDir, ".claw", "skills")

	// 如果目录不存在，说明当前工作区没有配置技能
	if _, err := os.Stat(skillBaseDir); os.IsNotExist(err) {
		return
	}

	// 遍历查找 SKILL.md
	filepath.WalkDir(skillBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 仅处理名为 SKILL.md 的文件
		if !d.IsDir() && d.Name() == "SKILL.md" {
			content, err := os.ReadFile(path)
			if err == nil {
				meta := s.parseSkillMetadata(string(content), path)
				s.skillMetadata[meta.Name] = meta
			}
		}
		return nil
	})
}

// parseSkillMetadata 仅解析 YAML Frontmatter 中的元数据（不加载正文）
func (s *SkillLoader) parseSkillMetadata(content, filePath string) SkillMetadata {
	meta := SkillMetadata{
		Name:        "Unknown Skill",
		Description: "No description provided.",
		FilePath:    filePath,
	}

	// 简单解析 YAML Frontmatter (以 --- 包裹)
	if strings.HasPrefix(content, "---\n") || strings.HasPrefix(content, "---\r\n") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) == 3 {
			frontmatter := parts[1]

			// 逐行提取 metadata
			lines := strings.Split(frontmatter, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "name:") {
					meta.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
				} else if strings.HasPrefix(line, "description:") {
					meta.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
				} else if strings.HasPrefix(line, "tags:") {
					tagsStr := strings.TrimSpace(strings.TrimPrefix(line, "tags:"))
					tagsStr = strings.Trim(tagsStr, "[] ")
					if tagsStr != "" {
						tags := strings.Split(tagsStr, ",")
						for i, tag := range tags {
							tags[i] = strings.TrimSpace(tag)
						}
						meta.Tags = tags
					}
				}
			}
		}
	}

	return meta
}

// GetSkillsSkeleton 获取技能索引骨架（只含元数据，用于 System Prompt）
func (s *SkillLoader) GetSkillsSkeleton() string {
	if len(s.skillMetadata) == 0 {
		return ""
	}

	var skillsBuilder strings.Builder
	skillsBuilder.WriteString("\n## 可用技能（按需加载）\n")
	skillsBuilder.WriteString("当任务需要特定技能时，使用 `read_skill` 工具获取完整技能正文。\n\n")

	for name, meta := range s.skillMetadata {
		skillsBuilder.WriteString(fmt.Sprintf("### %s\n", name))
		skillsBuilder.WriteString(fmt.Sprintf("**触发条件**: %s\n", meta.Description))
		if len(meta.Tags) > 0 {
			skillsBuilder.WriteString(fmt.Sprintf("**标签**: %s\n", strings.Join(meta.Tags, ", ")))
		}
		skillsBuilder.WriteString("\n")
	}

	return skillsBuilder.String()
}

// LoadSkillBody 按需加载技能正文（懒加载核心）
// 支持精确匹配和模糊匹配
func (s *SkillLoader) LoadSkillBody(skillName string) (string, error) {
	// 1. 先尝试精确匹配
	meta, ok := s.skillMetadata[skillName]
	if ok {
		return s.loadSkillFile(meta)
	}

	// 2. 尝试模糊匹配（基于关键词）
	foundMeta, matchedName := s.findSkillByKeyword(skillName)
	if foundMeta != nil {
		// 记录匹配信息，便于调试
		fmt.Printf("[SkillLoader] 模糊匹配成功: '%s' -> '%s'\n", skillName, matchedName)
		return s.loadSkillFile(*foundMeta)
	}

	return "", fmt.Errorf("技能不存在: %s", skillName)
}

// findSkillByKeyword 通过关键词模糊查找技能
func (s *SkillLoader) findSkillByKeyword(keyword string) (*SkillMetadata, string) {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil, ""
	}

	// 优先级：
	// 1. 名称包含关键词（完全匹配的词）
	// 2. 描述包含关键词
	// 3. 标签包含关键词

	// 收集所有可能的匹配
	type match struct {
		name  string
		meta  SkillMetadata
		score int // 匹配分数
	}
	var matches []match

	for name, meta := range s.skillMetadata {
		score := 0
		lowerName := strings.ToLower(name)
		lowerDesc := strings.ToLower(meta.Description)

		// 名称匹配（高优先级）
		if strings.Contains(lowerName, keyword) {
			score += 10
			// 完全词匹配额外加分
			nameParts := strings.Fields(lowerName)
			for _, part := range nameParts {
				if part == keyword {
					score += 5
					break
				}
			}
		}

		// 描述匹配
		if strings.Contains(lowerDesc, keyword) {
			score += 5
		}

		// 标签匹配
		if len(meta.Tags) > 0 {
			for _, tag := range meta.Tags {
				if strings.Contains(strings.ToLower(tag), keyword) {
					score += 3
					break
				}
			}
		}

		if score > 0 {
			matches = append(matches, match{name, meta, score})
		}
	}

	// 找到最高分的匹配
	if len(matches) > 0 {
		bestMatch := matches[0]
		for _, m := range matches[1:] {
			if m.score > bestMatch.score {
				bestMatch = m
			}
		}
		return &bestMatch.meta, bestMatch.name
	}

	return nil, ""
}

// loadSkillFile 读取技能文件内容
func (s *SkillLoader) loadSkillFile(meta SkillMetadata) (string, error) {
	content, err := os.ReadFile(meta.FilePath)
	if err != nil {
		return "", fmt.Errorf("读取技能文件失败: %w", err)
	}

	// 提取正文部分（去掉 YAML Frontmatter）
	contentStr := string(content)
	if strings.HasPrefix(contentStr, "---\n") || strings.HasPrefix(contentStr, "---\r\n") {
		parts := strings.SplitN(contentStr, "---", 3)
		if len(parts) == 3 {
			return strings.TrimSpace(parts[2]), nil
		}
	}

	return contentStr, nil
}

// GetSkillNames 获取所有技能名称
func (s *SkillLoader) GetSkillNames() []string {
	names := make([]string, 0, len(s.skillMetadata))
	for name := range s.skillMetadata {
		names = append(names, name)
	}
	return names
}

// GetSkillMetadata 获取指定技能的元数据
func (s *SkillLoader) GetSkillMetadata(skillName string) (SkillMetadata, bool) {
	meta, ok := s.skillMetadata[skillName]
	return meta, ok
}
