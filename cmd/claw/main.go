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

func main() {
	if err := logger.Init(logger.DEBUG, "logs", true); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	log := logger.WithModule("main")

	if os.Getenv("SILICONFLOW_API_KEY") == "" {
		log.Fatal("请先导出 SILICONFLOW_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()

	llmProvider := provider.NewSiliconFlowProvider("deepseek-ai/DeepSeek-V3")

	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))

	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	//prompt := "请读取 README.md 文件并根据内容创建一个简单的项目说明文档 SUMMARY.md"
	prompt := "请执行一个 Level 3 Agent 演示任务：总结当前项目的 Agent 架构。"

	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatal("引擎运行崩溃: %v", err)
	}
}
