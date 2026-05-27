package testutils

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/skytodmoon/go-tiny-claw/internal/engine"
	"github.com/skytodmoon/go-tiny-claw/internal/schema"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

func SetupTestEnvironment(t *testing.T) (workDir string, cleanup func()) {
	t.Helper()
	workDir = t.TempDir()

	testFiles := map[string]string{
		"README.md": "# Test Project\n\nThis is a test project.",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(workDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	cleanup = func() {
	}

	return workDir, cleanup
}

func CreateTestEngine(workDir string, generateFunc func(ctx context.Context, messages []schema.Message, tools []schema.ToolDefinition) (*schema.Message, error)) *engine.AgentEngine {
	mockProvider := NewMockProvider(generateFunc)
	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))
	return engine.NewAgentEngine(mockProvider, registry, true)
}

func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("file %q does not exist", path)
	}
}

func AssertFileContent(t *testing.T, path, expected string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != expected {
		t.Fatalf("file %q content = %q, want %q", path, string(content), expected)
	}
}
