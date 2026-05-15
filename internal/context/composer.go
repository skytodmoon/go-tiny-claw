// internal/context/composer.go
package context

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/skytodmoon/go-tiny-claw/internal/schema"
)

// PromptComposer 负责根据工作区环境动态生成 System Prompt
type PromptComposer struct {
	workDir     string
	skillLoader *SkillLoader
}

func NewPromptComposer(workDir string) *PromptComposer {
	return &PromptComposer{
		workDir:     workDir,
		skillLoader: NewSkillLoader(workDir),
	}
}

// Build 组装并返回一条完整的 RoleSystem 消息
func (c *PromptComposer) Build() schema.Message {
	var promptBuilder strings.Builder

	promptBuilder.WriteString(`# 核心身份
你名叫 go-tiny-claw，一个由驾驭工程驱动的骨灰级研发助手。
你具备极简主义哲学，拒绝废话。你能通过系统提供的内置工具，创建、读取、修改和执行工作区中的代码。

# 核心纪律 (CRITICAL)
1. 如需检查文件是否存在，请使用 bash 的 ls 或 test -f，而不是对目录使用 read_file。
2. 创建新文件时，务必使用 write_file，并同时提供 path 和 content 参数。
3. 编辑文件前务必先读取现有文件，以理解上下文。
4. 无论何时你需要写代码或创建文件，都要直接使用 write_file 工具。
5. 遇到工具执行报错时，仔细阅读 stderr，尝试自己修正命令并重试。
6. 始终用中文回复，以便传达你的进展和想法。

# 技能懒加载机制
系统采用"渐进式暴露"策略：启动时只加载技能元数据（名称和触发描述），
当你判定当前任务需要某个技能时，请使用 read_skill 工具加载完整技能正文。

# 技能使用规则
1. 技能加载后，使用 bash 工具执行技能中描述的命令，不要尝试调用不存在的工具。
2. 技能中会提供具体的 bash 命令示例，直接复制使用即可。
3. 如果工具执行失败，请检查命令格式并重新尝试。
`)

	agentsMDPath := filepath.Join(c.workDir, "AGENTS.md")
	content, err := os.ReadFile(agentsMDPath)
	if err == nil {
		promptBuilder.WriteString("\n# 项目专属指南 (来自 AGENTS.md)\n")
		promptBuilder.WriteString("以下是当前工作区特有的架构规范与注意事项，你的行为必须绝对符合以下要求：\n")
		promptBuilder.WriteString("```markdown\n")
		promptBuilder.WriteString(string(content))
		promptBuilder.WriteString("\n```\n")
	}

	skillsContent := c.skillLoader.GetSkillsSkeleton()
	if skillsContent != "" {
		promptBuilder.WriteString(skillsContent)
	}

	return schema.Message{
		Role:    schema.RoleSystem,
		Content: promptBuilder.String(),
	}
}

func (c *PromptComposer) GetSkillLoader() *SkillLoader {
	return c.skillLoader
}