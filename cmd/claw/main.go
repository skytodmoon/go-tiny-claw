package main

import (
	"net/http"
	"os"

	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/skytodmoon/go-tiny-claw/internal/context"
	"github.com/skytodmoon/go-tiny-claw/internal/engine"
	"github.com/skytodmoon/go-tiny-claw/internal/feishu"
	"github.com/skytodmoon/go-tiny-claw/internal/logger"
	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

func main() {
	if err := logger.Init(logger.DEBUG, "logs", true); err != nil {
		logger.WithModule("main").Fatal("初始化日志失败: %v", err)
	}

	log := logger.WithModule("main")

	if os.Getenv("SILICONFLOW_API_KEY") == "" {
		log.Fatal("请先导出 SILICONFLOW_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()

	llmProvider := provider.NewSiliconFlowProvider("deepseek-ai/DeepSeek-V3")

	// 创建技能加载器
	skillLoader := context.NewSkillLoader(workDir)

	// 注册工具
	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))
	registry.Register(tools.NewReadSkillTool(skillLoader)) // 新增：技能懒加载工具

	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	// 1. 初始化飞书 Bot
	bot := feishu.NewFeishuBot(eng)

	// 2. 使用 httpserverext 创建事件处理函数
	handler := httpserverext.NewEventHandlerFunc(bot.GetEventDispatcher())

	// 3. 注册路由并启动 HTTP 服务
	http.HandleFunc("/webhook/event", handler)
	port := ":48080"
	log.Info("🚀 go-tiny-claw 飞书服务端已启动，正在监听 %s 端口", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("HTTP 服务器启动失败: %v", err)
	}
}
