# 项目概述

## 核心理念

一个轻量级的 AI Agent 框架，使用 Go 语言实现，核心理念为 **"能慢思考、会用工具、能推进任务"**。

## 功能列表

### 核心功能
- 多 Provider 支持（OpenAI、Claude、MiniMax、DeepSeek、GLM47、SiliconFlow）
- 上下文管理（Token 监控、自动压缩）
- 工具注册表（ReadFile、WriteFile、EditFile、Bash）
- 四阶段循环（Thinking → Acting → Observation → Re-thinking）

### 高级功能
- 动态上下文组装
- 自动压缩机制（Token 预算管理）
- 执行结果反馈机制
- 执行摘要统计

## 快速开始

```bash
# 设置 API Key（选择其中一个）
export SILICONFLOW_API_KEY="your-api-key"
export ZHIPU_API_KEY="your-api-key"
export NVIDIA_API_KEY="your-api-key"

# 运行
go run cmd/claw/main.go
```
