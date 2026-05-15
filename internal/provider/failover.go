// internal/provider/failover.go
package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/skytodmoon/go-tiny-claw/internal/logger"
	"github.com/skytodmoon/go-tiny-claw/internal/schema"
)

var log = logger.WithModule("provider")

// FailoverProvider 故障转移 Provider，支持主备切换
type FailoverProvider struct {
	providers  []LLMProvider // Provider 列表，按优先级排序
	names      []string      // Provider 名称列表，用于日志
	timeout    time.Duration // 单个 Provider 超时时间
	maxRetries int           // 最大重试次数
}

// FailoverConfig 故障转移配置
type FailoverConfig struct {
	Providers  []LLMProvider // Provider 列表，顺序决定优先级（第一个为主服务）
	Names      []string      // Provider 名称，与 Providers 一一对应
	Timeout    time.Duration // 单个 Provider 超时时间，默认 30 秒
	MaxRetries int           // 最大重试次数，默认 2 次
}

// NewFailoverProvider 创建故障转移 Provider
func NewFailoverProvider(config FailoverConfig) *FailoverProvider {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 2
	}
	if len(config.Names) == 0 {
		// 默认命名
		config.Names = make([]string, len(config.Providers))
		for i := range config.Providers {
			config.Names[i] = fmt.Sprintf("Provider-%d", i+1)
		}
	}

	return &FailoverProvider{
		providers:  config.Providers,
		names:      config.Names,
		timeout:    config.Timeout,
		maxRetries: config.MaxRetries,
	}
}

// Generate 实现 LLMProvider 接口，支持故障转移
func (f *FailoverProvider) Generate(ctx context.Context, messages []schema.Message, availableTools []schema.ToolDefinition) (*schema.Message, error) {
	var lastErr error

	// 尝试所有 Provider
	for attempt := 0; attempt < f.maxRetries; attempt++ {
		for i, provider := range f.providers {
			name := f.names[i]

			log.Info("[Failover] 尝试 Provider: %s (尝试 %d/%d)", name, attempt+1, f.maxRetries)

			// 创建带超时的上下文
			timeoutCtx, cancel := context.WithTimeout(ctx, f.timeout)

			// 执行请求
			result, err := provider.Generate(timeoutCtx, messages, availableTools)

			cancel() // 取消上下文

			if err == nil && result != nil {
				log.Info("[Failover] Provider %s 成功", name)
				return result, nil
			}

			// 记录错误
			if err != nil {
				log.Warn("[Failover] Provider %s 失败: %v", name, err)
				lastErr = err
			}

			// 如果不是最后一个 Provider，继续尝试下一个
			if i < len(f.providers)-1 {
				log.Info("[Failover] 切换到下一个 Provider")
			}
		}
	}

	// 所有 Provider 都失败了
	return nil, fmt.Errorf("所有 Provider 均失败，最后错误: %w", lastErr)
}

// GetProviderNames 获取所有 Provider 名称
func (f *FailoverProvider) GetProviderNames() []string {
	return f.names
}

// GetProviderCount 获取 Provider 数量
func (f *FailoverProvider) GetProviderCount() int {
	return len(f.providers)
}
