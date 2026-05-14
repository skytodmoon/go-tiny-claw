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
	err := eng.Run(context.Background(),  "Create a summary of the project", nil)
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
	err := eng.Run(context.Background(),  "Update the README", nil)
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
	err := eng.Run(context.Background(),  "Create two test files", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 3 {
		t.Errorf("callCount = %d, want 3", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/file1.txt")
	testutils.AssertFileExists(t, workDir+"/file2.txt")
}

func TestLevel3ComplexWorkflow(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Step 1: Create initial file",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("report.md", "# Project Report\n\n## Introduction\n")},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Step 2: Read existing README for reference",
				ToolCalls: []schema.ToolCall{testutils.CreateReadFileCall("README.md")},
			}, nil
		case 3:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Step 3: Append content to report",
				ToolCalls: []schema.ToolCall{testutils.CreateEditFileCall("report.md", "## Introduction", "## Introduction\n\nBased on README analysis...")},
			}, nil
		case 4:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Step 4: Verify the final report",
				ToolCalls: []schema.ToolCall{testutils.CreateReadFileCall("report.md")},
			}, nil
		case 5:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Complex workflow completed!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(),  "Create a comprehensive project report", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 5 {
		t.Errorf("callCount = %d, want 5", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/report.md")
}

func TestLevel3ErrorHandling(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Trying to read non-existent file",
				ToolCalls: []schema.ToolCall{testutils.CreateReadFileCall("nonexistent.txt")},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "File not found, creating new file instead",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("fallback.txt", "Created as fallback")},
			}, nil
		case 3:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Error handled successfully!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(),  "Read a file that may not exist, handle error gracefully", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 3 {
		t.Errorf("callCount = %d, want 3", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/fallback.txt")
}

func TestLevel3EmptyContentHandling(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Creating file with empty content",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("empty.txt", "")},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Reading empty file",
				ToolCalls: []schema.ToolCall{testutils.CreateReadFileCall("empty.txt")},
			}, nil
		case 3:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Empty content handling test completed!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(),  "Test handling empty file content", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 3 {
		t.Errorf("callCount = %d, want 3", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/empty.txt")
}

func TestLevel3MaxTurnsLimit(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0
	maxTurns := 5

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		if callCount <= maxTurns {
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Continue working...",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("step"+string(rune(callCount+'0'))+".txt", "Step content")},
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Should not reach here",
			ToolCalls: nil,
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	eng.MaxTurns = maxTurns

	err := eng.Run(context.Background(),  "Test maximum turns limit", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != maxTurns {
		t.Errorf("callCount = %d, want %d", callCount, maxTurns)
	}
}

func TestLevel3ContextCompression(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		if callCount <= 6 {
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Creating file " + string(rune(callCount+'0')) + " to test context compression",
				ToolCalls: []schema.ToolCall{testutils.CreateWriteFileCall("file"+string(rune(callCount+'0'))+".txt", "Content for file " + string(rune(callCount+'0')))},
			}, nil
		} else {
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Context compression test completed!",
				ToolCalls: nil,
			}, nil
		}
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(),  "Create multiple files to test context compression mechanism", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 7 {
		t.Errorf("callCount = %d, want 7", callCount)
	}
}

func TestLevel3NoToolCalls(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "No tools needed, task is simple!",
			ToolCalls: nil,
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(),  "Answer a simple question without tools", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}
}

func TestLevel3MultipleToolCallsPerTurn(t *testing.T) {
	workDir, cleanup := testutils.SetupTestEnvironment(t)
	defer cleanup()

	callCount := 0

	generateFunc := func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error) {
		callCount++

		switch callCount {
		case 1:
			return &schema.Message{
				Role:    schema.RoleAssistant,
				Content: "Creating multiple files in parallel",
				ToolCalls: []schema.ToolCall{
					testutils.CreateWriteFileCall("parallel1.txt", "Parallel file 1"),
					testutils.CreateWriteFileCall("parallel2.txt", "Parallel file 2"),
					testutils.CreateWriteFileCall("parallel3.txt", "Parallel file 3"),
				},
			}, nil
		case 2:
			return &schema.Message{
				Role:      schema.RoleAssistant,
				Content:   "Multiple files created successfully!",
				ToolCalls: nil,
			}, nil
		}

		return &schema.Message{
			Role:    schema.RoleAssistant,
			Content: "Done",
		}, nil
	}

	eng := testutils.CreateTestEngine(workDir, generateFunc)
	err := eng.Run(context.Background(),  "Create three files simultaneously", nil)
	if err != nil {
		t.Fatalf("Engine.Run() error = %v", err)
	}

	if callCount != 2 {
		t.Errorf("callCount = %d, want 2", callCount)
	}

	testutils.AssertFileExists(t, workDir+"/parallel1.txt")
	testutils.AssertFileExists(t, workDir+"/parallel2.txt")
	testutils.AssertFileExists(t, workDir+"/parallel3.txt")
}
