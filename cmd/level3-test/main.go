package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skytodmoon/go-tiny-claw/internal/engine"
	"github.com/skytodmoon/go-tiny-claw/internal/logger"
	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

type TestCase struct {
	Name        string
	Prompt      string
	CheckFiles  []string
	Description string
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    LEVEL 3 VALIDATION TESTS                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	log := logger.WithModule("level3-test")

	if os.Getenv("SILICONFLOW_API_KEY") == "" {
		log.Fatal("SILICONFLOW_API_KEY not set")
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("无法获取工作目录: %v", err)
	}

	testCases := []TestCase{
		{
			Name:        "Read and Write",
			Prompt:      "Read the README.md file and create a SUMMARY.md with a brief overview",
			CheckFiles:  []string{"SUMMARY.md"},
			Description: "Tests file reading and writing with thinking phases",
		},
	}

	llmProvider := provider.NewSiliconFlowProvider("deepseek-ai/DeepSeek-V3")
	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))

	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	passed := 0
	failed := 0

	for i, tc := range testCases {
		fmt.Printf("\n╔══════════════════════════════════════════════════════════════╗")
		fmt.Printf("\n║ Test %d: %s", i+1, tc.Name)
		fmt.Printf("\n║ %s", tc.Description)
		fmt.Printf("\n╚══════════════════════════════════════════════════════════════╝\n")

		for _, f := range tc.CheckFiles {
			os.Remove(f)
		}

		fmt.Printf("\n📝 Running with prompt: %s\n", truncate(tc.Prompt, 60))

		ctx := context.Background()
		err := eng.Run(ctx, tc.Prompt, nil)

		if err != nil {
			fmt.Printf("\n❌ Test failed: %v\n", err)
			failed++
			continue
		}

		allExist := true
		for _, f := range tc.CheckFiles {
			if _, err := os.Stat(f); os.IsNotExist(err) {
				fmt.Printf("❌ Missing file: %s\n", f)
				allExist = false
			} else {
				fmt.Printf("✅ Found file: %s\n", f)
			}
		}

		if allExist {
			fmt.Printf("\n✅ Test passed!\n")
			passed++
		} else {
			fmt.Printf("\n❌ Test failed!\n")
			failed++
		}
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         TEST SUMMARY                          ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  Total:  %d                                                     ║\n", len(testCases))
	fmt.Printf("║  Passed: %d                                                     ║\n", passed)
	fmt.Printf("║  Failed: %d                                                     ║\n", failed)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	if failed > 0 {
		os.Exit(1)
	}
	fmt.Println("\n🎉 All tests passed!")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
