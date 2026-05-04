package integration

import (
	"context"
	"testing"

	"github.com/skytodmoon/go-tiny-claw/internal/schema"
	"github.com/skytodmoon/go-tiny-claw/testutils"
)

func TestLevel3ReadWriteCycle(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "I'll read the README.md file first.",
				ToolCalls: []schema.ToolCall{testutils.CreateReadFileCall("README.md")},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Now I'll create a SUMMARY.md file.",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("SUMMARY.md", "Project Summary\n\nThis is a test.")},
			}, nil
		case 3:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Task completed successfully!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(), "Create a summary of the project")
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 3 {
		t.Errorf("callCount = %d, want 3", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/SUMMARY.md")
}

func TestLevel3EditCycle(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Reading file first...",
				ToolCalls: []schema.ToolCall{testutils.CreateReadFileCall("README.md")},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Editing file...",
				ToolCalls: []schema.ToolCall{testutils.CreateEditFileCall("README.md", "Test", "Updated Test")},
			}, nil
		case 3:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Done!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(), "Update the README")
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 3 {
		t.Errorf("callCount = %d, want 3", callCount)
	}
}

func TestLevel3MultipleTools(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "First step: create file1",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("file1.txt", "Content 1")},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Second step: create file2",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("file2.txt", "Content 2")},
			}, nil
		case 3:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "All files created!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(), "Create two test files")
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 3 {
		t.Errorf("callCount = %d, want 3", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/file1.txt")
	testutils.AssertFileExists(t, workDir+"/file2.txt")
}
