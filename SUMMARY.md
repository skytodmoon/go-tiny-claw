# Go Tiny Claw 项目说明

## 项目概述

Go Tiny Claw 是一个轻量级的 AI Agent 框架，使用 Go 语言实现。核心目标是构建一个具备**慢思考、工具调用和任务推进能力**的 Agent 系统。

## 核心能力

- **慢思考（Thinking Phase）**：在 Main Loop 中显式区分思考与行动阶段。
- **工具调用（Tool Registry）**：支持注册与分发工具，例如文件读取和编辑功能。  
- **任务执行**：Agent 能够根据目标（如读取文件并总结）自主调用工具完成任务。

## 版本规划

### 🟢 Level 1（已完成）
- 实现基础功能：慢思考、工具调用、单次任务执行。
- 关键词：**慢思考 → 动手做 → 得到结果**

### 🟡 Level 2（进阶）
- 结构化思考、多步执行、动态工具选择。
- 关键词：**让 Agent 开始"连续行动"**

### 🔴 Level 3（挑战）
- 优化 Agent 循环：Thinking → Acting → Observation → 再思考。
- 根据执行结果调整行为，优化 Prompt 结构。
- 关键词：**让 Agent 开始"像一个系统在运转"**

## 已实现功能

- 多 Provider 支持（OpenAI、Claude、MiniMax 等）。
- 上下文管理与 Token 监控。
- 工具注册表（ReadFile、WriteFile）。
- 慢思考机制与基于文件的记忆存储。

## 快速开始

```bash
# 设置 API Key
export SILICONFLOW_API_KEY="your-api-key"

go run cmd/claw/main.go
```

## 项目结构

```
go-tiny-claw/
├── cmd/
│   └── claw/
│       └── main.go
├── internal/
│   ├── engine/
│   │   └── loop.go
│   ├── provider/
│   │   ├── interface.go
│   │   └── ...
│   ├── tools/
│   │   ├── registry.go
│   │   ├── read_file.go
│   │   └── write_file.go
│   └── memory/
└── README.md
```

## License

MIT