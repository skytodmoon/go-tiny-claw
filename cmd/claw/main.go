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
		log.Fatal("请先导出 SILICONFLOW_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()

	llmProvider := provider.NewSiliconFlowProvider("deepseek-ai/DeepSeek-V3")

	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))

	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	prompt := "请读取 README.md 文件并根据内容创建一个简单的项目说明文档 SUMMARY.md"

	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
