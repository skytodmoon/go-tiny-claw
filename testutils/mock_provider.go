package testutils

import (
	"context"
	"encoding/json"

	"github.com/skytodmoon/go-tiny-claw/internal/schema"
)

type MockProvider struct {
	GenerateFunc func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error)
}

func (m *MockProvider) Generate(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, messages, tools)
	}
	return &schema.Message{
		Role:    schema.RoleAssistant,
		Content: "Mock response",
	}, nil
}

func NewMockProvider(generateFunc func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error)) *MockProvider {
	return &MockProvider{
		GenerateFunc: generateFunc,
	}
}

func CreateToolCall(name string, args interface{}) schema.ToolCall {
	argsJSON, _ := json.Marshal(args)
	return schema.ToolCall{
		ID:        "test-call-1",
		Name:      name,
		Arguments: argsJSON,
	}
}

func CreateReadFileCall(path string) schema.ToolCall {
	return CreateToolCall("read_file", map[string]string{"path": path})
}

func CreateWriteFileCall(path, content string) schema.ToolCall {
	return CreateToolCall("write_file", map[string]string{"path": path, "content": content})
}

func CreateEditFileCall(path, oldStr, newStr string) schema.ToolCall {
	return CreateToolCall("edit_file", map[string]string{"path": path, "old_str": oldStr, "new_str": newStr})
}
