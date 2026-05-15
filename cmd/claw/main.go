package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	// 获取配置的 Provider 顺序（逗号分隔）
	// 示例: PROVIDER_ORDER="nvidia,siliconflow" 或 "siliconflow,nvidia"
	providerOrder := os.Getenv("PROVIDER_ORDER")
	if providerOrder == "" {
		providerOrder = "nvidia,siliconflow" // 默认顺序：nvidia 为主，siliconflow 为备
	}

	workDir, _ := os.Getwd()

	// 构建 Provider 列表
	providers, names := buildProviders(providerOrder, log)

	if len(providers) == 0 {
		log.Fatal("没有可用的 LLM Provider，请配置 NVIDIA_API_KEY 或 SILICONFLOW_API_KEY")
	}

	var llmProvider provider.LLMProvider

	if len(providers) == 1 {
		// 只有一个 Provider，直接使用
		llmProvider = providers[0]
		log.Info("使用单一 Provider: %s", names[0])
	} else {
		// 使用故障转移 Provider
		timeoutStr := os.Getenv("LLM_TIMEOUT")
		timeout := 30 * time.Second
		if timeoutStr != "" {
			if t, err := time.ParseDuration(timeoutStr); err == nil {
				timeout = t
			}
		}

		maxRetries := 2
		if retriesStr := os.Getenv("MAX_RETRIES"); retriesStr != "" {
			if r, err := strconv.Atoi(retriesStr); err == nil && r > 0 {
				maxRetries = r
			}
		}

		llmProvider = provider.NewFailoverProvider(provider.FailoverConfig{
			Providers:  providers,
			Names:      names,
			Timeout:    timeout,
			MaxRetries: maxRetries,
		})
		log.Info("使用故障转移模式，Provider 顺序: %s", strings.Join(names, " → "))
		log.Info("超时时间: %v, 最大重试次数: %d", timeout, maxRetries)
	}

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

// buildProviders 根据配置顺序构建 Provider 列表
func buildProviders(order string, log *logger.Logger) ([]provider.LLMProvider, []string) {
	var providers []provider.LLMProvider
	var names []string

	// 读取模型配置（支持环境变量）
	nvidiaModel := os.Getenv("NVIDIA_MODEL")
	if nvidiaModel == "" {
		nvidiaModel = "minimaxai/minimax-m2.7" // 默认模型
	}

	siliconFlowModel := os.Getenv("SILICONFLOW_MODEL")
	if siliconFlowModel == "" {
		siliconFlowModel = "deepseek-ai/DeepSeek-V3" // 默认模型
	}

	orderParts := strings.Split(strings.ToLower(order), ",")

	for _, part := range orderParts {
		part = strings.TrimSpace(part)
		switch part {
		case "nvidia":
			if key := os.Getenv("NVIDIA_API_KEY"); key != "" {
				providers = append(providers, provider.NewNvidiaProvider(nvidiaModel))
				names = append(names, "NVIDIA")
				log.Info("已配置 NVIDIA Provider，模型: %s", nvidiaModel)
			} else {
				log.Warn("NVIDIA_API_KEY 未设置，跳过 NVIDIA Provider")
			}
		case "siliconflow":
			if key := os.Getenv("SILICONFLOW_API_KEY"); key != "" {
				providers = append(providers, provider.NewSiliconFlowProvider(siliconFlowModel))
				names = append(names, "SiliconFlow")
				log.Info("已配置 SiliconFlow Provider，模型: %s", siliconFlowModel)
			} else {
				log.Warn("SILICONFLOW_API_KEY 未设置，跳过 SiliconFlow Provider")
			}
		default:
			log.Warn("未知的 Provider 类型: %s", part)
		}
	}

	return providers, names
}
