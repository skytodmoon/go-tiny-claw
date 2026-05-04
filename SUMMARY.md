# Go Tiny Claw 项目说明

## 项目概述

Go Tiny Claw 是一个轻量级的 AI Agent 框架，采用 Go 语言实现。核心理念是设计一个"能慢思考、会用工具、能推进任务"的 Agent。

## 版本规划

### 🟢 Level 1（已完成）
- **慢思考**：显式区分 Thinking 和 Acting 阶段。
- **工具调用**：至少接入 1 个工具。
- **执行任务**：Agent 能够完成简单任务，如读取文件并总结。

### 🟡 Level 2（进阶）
- 结构化 Thinking 阶段。
- 支持多步任务执行。
- 动态选择工具。

### 🔴 Level 3（挑战）
- 动态调整行为。
- 优化 Prompt 结构。
- 完善循环机制。

## 已实现功能

- 多 Provider 支持（OpenAI、Claude 等）。
- 上下文管理。
- 工具注册表（如文件读写）。
- 慢思考机制。
- 基于文件的记忆存储。

## 快速开始

```bash
# 设置 API Key
export SILICONFLOW_API_KEY="your-api-key"

go run cmd/claw/main.go
```

## 项目结构
```
go-tiny-claw/
├── cmd/           # 程序入口
├── internal/      # 核心实现（引擎、Provider、工具等）
└── README.md
```

## License

MIT