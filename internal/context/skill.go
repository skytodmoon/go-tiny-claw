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
func (s *SkillLoader) LoadSkillBody(skillName string) (string, error) {
	meta, ok := s.skillMetadata[skillName]
	if !ok {
		return "", fmt.Errorf("技能不存在: %s", skillName)
	}

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
