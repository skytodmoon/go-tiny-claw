#!/bin/bash

set -e

echo "🚀 Starting e2e tests..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

mkdir -p logs

echo "📁 Project root: $PROJECT_ROOT"

check_api_key() {
    if [ -z "$SILICONFLOW_API_KEY" ]; then
        echo "⚠️  SILICONFLOW_API_KEY not set, skipping e2e tests"
        exit 0
    fi
}

test_quick_start() {
    echo ""
    echo "🧪 Test 1: Quick start test"
    echo "─────────────────────────────────"

    local log_file="logs/test-quick-start.log"
    rm -f "$log_file"

    go run cmd/claw/main.go 2>&1 | tee "$log_file"

    if [ -f "SUMMARY.md" ]; then
        echo "✅ Test passed: SUMMARY.md created"
    else
        echo "❌ Test failed: SUMMARY.md not created"
        return 1
    fi

    echo "✅ Test 1 passed"
}

test_with_custom_prompt() {
    echo ""
    echo "🧪 Test 2: Custom prompt test"
    echo "─────────────────────────────────"

    local log_file="logs/test-custom-prompt.log"

    cat > cmd/claw/test_prompt.go << 'EOF'
package main

import (
	"context"
	"log"
	"os"

	"github.com/skytodmoon/go-tiny-claw/internal/engine"
	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

func main() {
	if os.Getenv("SILICONFLOW_API_KEY") == "" {
		log.Fatal("SILICONFLOW_API_KEY not set")
	}

	workDir, _ := os.Getwd()

	llmProvider := provider.NewSiliconFlowProvider("deepseek-ai/DeepSeek-V3")
	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))

	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)
	prompt := "Create a test file called 'test.txt' with content 'Hello, Level 3!'"

	if err := eng.Run(context.Background(), prompt); err != nil {
		log.Fatal(err)
	}
}
EOF

    go run cmd/claw/test_prompt.go 2>&1 | tee "$log_file"

    if [ -f "test.txt" ]; then
        local content=$(cat test.txt)
        if [[ "$content" == *"Hello"* ]]; then
            echo "✅ Test passed: test.txt created with correct content"
        else
            echo "❌ Test failed: test.txt content incorrect: $content"
            return 1
        fi
    else
        echo "❌ Test failed: test.txt not created"
        return 1
    fi

    echo "✅ Test 2 passed"
}

cleanup() {
    echo ""
    echo "🧹 Cleaning up..."
    rm -f test.txt
    rm -f cmd/claw/test_prompt.go
    echo "✅ Cleanup complete"
}

main() {
    check_api_key

    trap cleanup EXIT

    test_quick_start
    test_with_custom_prompt

    echo ""
    echo "🎉 All e2e tests passed!"
}

main "$@"
