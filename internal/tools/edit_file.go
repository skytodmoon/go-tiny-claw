package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/skytodmoon/go-tiny-claw/internal/schema"
)

type EditFileTool struct {
	workDir string
}

func NewEditFileTool(workDir string) *EditFileTool {
	return &EditFileTool{workDir: workDir}
}

func (t *EditFileTool) Name() string {
	return "edit_file"
}

func (t *EditFileTool) Definition() schema.ToolDefinition {
	return schema.ToolDefinition{
		Name:        t.Name(),
		Description: "编辑文件内容，支持在指定位置插入、替换或删除文本。请提供相对于工作区的路径。",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path": map[string]interface{}{
					"type":        "string",
					"description": "要编辑的文件路径，如 cmd/claw/main.go",
				},
				"old_str": map[string]interface{}{
					"type":        "string",
					"description": "需要被替换的原字符串（精确匹配）",
				},
				"new_str": map[string]interface{}{
					"type":        "string",
					"description": "替换后的新字符串",
				},
			},
			"required": []string{"path", "old_str", "new_str"},
		},
	}
}

type editFileArgs struct {
	Path   string `json:"path"`
	OldStr string `json:"old_str"`
	NewStr string `json:"new_str"`
}

func (t *EditFileTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var input editFileArgs
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	fullPath := filepath.Join(t.workDir, input.Path)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	if !strings.Contains(string(content), input.OldStr) {
		return "", fmt.Errorf("未找到指定的 old_str，编辑失败")
	}

	newContent := strings.Replace(string(content), input.OldStr, input.NewStr, 1)

	if err := os.WriteFile(fullPath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return fmt.Sprintf("成功编辑文件 %s", input.Path), nil
}
