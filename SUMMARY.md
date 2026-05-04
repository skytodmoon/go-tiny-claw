# Go Tiny Claw 项目说明

## 项目概述

Go Tiny Claw 是一个轻量级的 AI Agent 框架，使用 Go 语言实现。其核心理念是构建一个 **能慢思考、会用工具、能推进任务** 的 Agent。

## 版本规划

### Level 1（已完成）
- 实现慢思考机制（Thinking 阶段）
- 工具调用（Tool Registry）
- 单次任务执行能力

### Level 2（已完成）
- 结构化思考（输出计划）
- 多步执行能力
- 动态工具选择

### Level 3（已完成）
- 四阶段循环（Thinking → Acting → Observation → Re-thinking）
- 动态上下文组装
- Token 预算管理与自动压缩
- 执行结果反馈机制

## 快速开始

```bash
# 设置 API Key（选择其中一个）
export SILICONFLOW_API_KEY="your-api-key"
export ZHIPU_API_KEY="your-api-key"
export NVIDIA_API_KEY="your-api-key"

# 运行
go run cmd/claw/main.go
```

## 项目结构

- **cmd/claw/main.go**: 程序入口
- **internal/engine/loop.go**: MainLoop 核心实现
- **internal/provider/**: LLM Provider 适配器
- **internal/context/**: 上下文管理
- **internal/tools/**: 工具注册表与内置工具
- **internal/memory/**: 基于文件的记忆存储
- **internal/feishu/**: 飞书机器人集成

## 当前功能

- 多 Provider 支持（OpenAI、Claude、MiniMax、DeepSeek 等）
- 上下文管理（Token 监控、自动压缩）
- 工具注册表（ReadFile、WriteFile、EditFile）
- 四阶段循环（Thinking → Acting → Observation → Re-thinking）

## License

MIT