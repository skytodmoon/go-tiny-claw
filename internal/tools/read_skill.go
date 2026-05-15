// internal/tools/read_skill.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	ctxpkg "github.com/skytodmoon/go-tiny-claw/internal/context"
	"github.com/skytodmoon/go-tiny-claw/internal/schema"
)

// ReadSkillTool 用于按需加载技能正文的内置工具
type ReadSkillTool struct {
	skillLoader *ctxpkg.SkillLoader
}

// NewReadSkillTool 创建 read_skill 工具实例
func NewReadSkillTool(skillLoader *ctxpkg.SkillLoader) *ReadSkillTool {
	return &ReadSkillTool{
		skillLoader: skillLoader,
	}
}

// Name 返回工具名称
func (t *ReadSkillTool) Name() string {
	return "read_skill"
}

// Definition 返回工具定义
func (t *ReadSkillTool) Definition() schema.ToolDefinition {
	return schema.ToolDefinition{
		Name:        t.Name(),
		Description: "按需加载指定技能的完整正文内容。当任务需要使用某个技能时，调用此工具获取技能的完整说明和执行指南。",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"skill_name": map[string]interface{}{
					"type":        "string",
					"description": "要加载的技能名称",
				},
			},
			"required": []string{"skill_name"},
		},
	}
}

// ReadSkillArgs 定义工具参数
type ReadSkillArgs struct {
	SkillName string `json:"skill_name"`
}

// Execute 执行工具逻辑
func (t *ReadSkillTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params ReadSkillArgs
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	if params.SkillName == "" {
		return "", fmt.Errorf("skill_name 参数不能为空")
	}

	// 调用 SkillLoader 加载技能正文
	body, err := t.skillLoader.LoadSkillBody(params.SkillName)
	if err != nil {
		return "", err
	}

	// 返回格式化的技能正文
	result := fmt.Sprintf("## 技能: %s\n\n%s", params.SkillName, body)
	return result, nil
}
