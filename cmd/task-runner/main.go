package main

import (
	"context"
	"fmt"
	"os"

	"github.com/skytodmoon/go-tiny-claw/internal/engine"
	"github.com/skytodmoon/go-tiny-claw/internal/logger"
	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/schema"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

type Task struct {
	ID          int
	Name        string
	Prompt      string
	Description string
	Expected    []string
}

func main() {
	if err := logger.Init(logger.DEBUG, "logs", true); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	log := logger.WithModule("task-runner")

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

	eng := engine.NewAgentEngine(llmProvider, registry, true)

	tasks := []Task{
		{
			ID:          1,
			Name:        "文件读写基础",
			Prompt:      "请读取 README.md 文件并根据内容创建一个简单的项目说明文档 SUMMARY.md",
			Description: "测试基本的文件读取和写入能力",
			Expected:    []string{"SUMMARY.md"},
		},
		{
			ID:          2,
			Name:        "文件编辑",
			Prompt:      "请读取 README.md 文件，将其中的 'AI Agent' 替换为 'AI智能助手'，并保存为 README_updated.md",
			Description: "测试文件内容编辑能力",
			Expected:    []string{"README_updated.md"},
		},
		{
			ID:          3,
			Name:        "多文件创建",
			Prompt:      "请创建三个文件：TODO.md（待办事项列表）、NOTES.md（笔记）、CHANGELOG.md（变更日志），每个文件至少包含3行内容",
			Description: "测试同时创建多个文件的能力",
			Expected:    []string{"TODO.md", "NOTES.md", "CHANGELOG.md"},
		},
		{
			ID:          4,
			Name:        "内容分析",
			Prompt:      "请读取 README.md 文件，分析项目的核心功能和架构特点，然后创建一个 ANALYSIS.md 文件记录分析结果",
			Description: "测试内容理解和分析能力",
			Expected:    []string{"ANALYSIS.md"},
		},
		{
			ID:          5,
			Name:        "代码生成",
			Prompt:      "请创建一个简单的 Go 语言 Hello World 程序，保存为 hello.go，然后编译运行验证",
			Description: "测试代码生成和命令执行能力",
			Expected:    []string{"hello.go"},
		},
		{
			ID:          6,
			Name:        "目录结构",
			Prompt:      "请查看当前目录结构，然后创建一个 STRUCTURE.md 文件记录项目目录结构和每个目录的用途",
			Description: "测试目录探索和结构化输出能力",
			Expected:    []string{"STRUCTURE.md"},
		},
		{
			ID:          7,
			Name:        "错误处理",
			Prompt:      "请尝试读取一个不存在的文件 MISSING_FILE.md，如果文件不存在则创建一个包含错误处理说明的文件",
			Description: "测试错误处理和异常情况处理能力",
			Expected:    []string{"MISSING_FILE.md"},
		},
		{
			ID:          8,
			Name:        "文档转换",
			Prompt:      "请读取 README.md 文件，将其内容转换为 Markdown 格式的表格形式，保存为 README_table.md",
			Description: "测试内容转换和格式处理能力",
			Expected:    []string{"README_table.md"},
		},
		{
			ID:          9,
			Name:        "配置文件生成",
			Prompt:      "请创建一个项目配置文件 config.yaml，包含数据库配置、日志配置和应用设置等部分",
			Description: "测试配置文件生成能力",
			Expected:    []string{"config.yaml"},
		},
		{
			ID:          10,
			Name:        "测试报告",
			Prompt:      "请运行项目的测试用例（执行 go test ./...），然后创建一个 TEST_REPORT.md 文件记录测试结果",
			Description: "测试命令执行和结果记录能力",
			Expected:    []string{"TEST_REPORT.md"},
		},
		{
			ID:          11,
			Name:        "文档摘要",
			Prompt:      "请读取 README.md 文件，提取关键信息创建一个简洁的项目摘要，保存为 ABSTRACT.md",
			Description: "测试信息提取和摘要能力",
			Expected:    []string{"ABSTRACT.md"},
		},
		{
			ID:          12,
			Name:        "代码审查",
			Prompt:      "请查看 internal/engine/loop.go 文件，分析代码结构和可能的优化点，创建 REVIEW.md 文件记录审查结果",
			Description: "测试代码理解和审查能力",
			Expected:    []string{"REVIEW.md"},
		},
		{
			ID:          13,
			Name:        "批量重命名",
			Prompt:      "请创建一个名为 batch_rename.sh 的 Shell 脚本，用于批量重命名项目中的测试文件",
			Description: "测试脚本生成能力",
			Expected:    []string{"batch_rename.sh"},
		},
		{
			ID:          14,
			Name:        "依赖分析",
			Prompt:      "请查看 go.mod 文件，分析项目依赖情况，创建 DEPENDENCIES.md 文件记录依赖列表和用途",
			Description: "测试依赖分析能力",
			Expected:    []string{"DEPENDENCIES.md"},
		},
		{
			ID:          15,
			Name:        "综合任务",
			Prompt:      "请完成以下任务：1) 读取 README.md；2) 创建项目概述文档；3) 生成 API 文档模板；4) 创建使用示例。所有输出保存到 docs/ 目录",
			Description: "测试多步骤复杂任务执行能力",
			Expected:    []string{"docs/"},
		},
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                TASK TEST SUITE - 15 种任务类型                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	passed := 0
	failed := 0

	for _, task := range tasks {
		fmt.Printf("┌──────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│ 📋 任务 %d: %s\n", task.ID, task.Name)
		fmt.Printf("│ %s\n", task.Description)
		fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")
		fmt.Printf("📝 Prompt: %s\n", task.Prompt)
		fmt.Println()

		session := engine.NewSession(fmt.Sprintf("task-%d", task.ID), workDir)
		session.Append(schema.Message{Role: schema.RoleUser, Content: task.Prompt})
		err := eng.Run(context.Background(), session, nil)
		if err != nil {
			fmt.Printf("❌ 任务 %d 失败: %v\n\n", task.ID, err)
			failed++
			continue
		}

		allExist := true
		for _, expected := range task.Expected {
			if _, err := os.Stat(expected); os.IsNotExist(err) {
				fmt.Printf("❌ 缺失预期文件: %s\n", expected)
				allExist = false
			}
		}

		if allExist {
			fmt.Printf("✅ 任务 %d 完成\n\n", task.ID)
			passed++
		} else {
			fmt.Printf("❌ 任务 %d 未完全完成\n\n", task.ID)
			failed++
		}
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     任务执行汇总                             ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  总任务数: %d                                                 ║\n", len(tasks))
	fmt.Printf("║  成功:     %d                                                 ║\n", passed)
	fmt.Printf("║  失败:     %d                                                 ║\n", failed)
	fmt.Printf("║  成功率:   %.1f%%                                               ║\n", float64(passed)/float64(len(tasks))*100)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	if failed > 0 {
		os.Exit(1)
	}
	fmt.Println("\n🎉 所有任务完成！")
}
