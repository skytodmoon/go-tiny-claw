# Go Tiny Claw

一个轻量级的 AI Agent 框架，使用 Go 语言实现。

## 核心理念

**"能慢思考、会用工具、能推进任务"** 的 Agent

## 版本规划

### 🟢 Level 1（已完成｜基于前 6 讲）

让你的 Agent 具备以下能力：

1. **慢思考（Thinking 阶段）**
   - 在 Main Loop 中，显式区分 Thinking / Acting
   - 先思考"下一步做什么"，再执行

2. **工具调用（Tool Registry）**
   - 至少接入 1 个工具（如：文件读取 / 简单编辑）
   - 通过统一注册与分发机制调用（而不是写死）

3. **YOLO 执行一次完整任务**
   - 给定一个目标（如：读取一个文件并总结）
   - Agent 自己决定调用工具并完成任务

👉 关键词：**慢思考 → 动手做 → 得到结果**

### 🟡 Level 2（进阶｜可提前探索）

在 Level 1 基础上尝试：

- 让 Thinking 阶段更结构化（比如输出"计划"而不是一句话）
- 支持多步执行（一次任务不一定一步完成）
- 尝试让 Agent 根据情况选择不同工具（而不是固定一个）

👉 关键词：**让 Agent 开始"连续行动"**

### 🔴 Level 3（挑战｜对标业界 Claude Code 最佳实践）

如果你有兴趣，可以结合 Claude Code 设计思路，尝试：

- 做一个更清晰的循环：Thinking → Acting → Observation → 再思考
- 让 Agent 根据执行结果调整下一步行为
- 尝试优化 Prompt 结构（思考 vs 执行分离）

👉 关键词：**让 Agent 开始"像一个系统在运转"**

> Claude Code 源码泄漏存档： https://github.com/FlyAIBox/civil-engineering-cloud-claude-code-source-v2.1.88

## 当前已实现

- ✅ 多 Provider 支持（OpenAI、Claude、MiniMax、DeepSeek、GLM47、SiliconFlow）
- ✅ 上下文管理（Token 监控）
- ✅ 工具注册表（ReadFile、WriteFile）
- ✅ 慢思考机制（Thinking Phase）
- ✅ 基于文件的记忆存储

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

```
go-tiny-claw/
├── cmd/
│   └── claw/
│       └── main.go              # 程序入口
├── internal/
│   ├── engine/
│   │   └── loop.go              # MainLoop 核心实现
│   ├── provider/                 # LLM Provider 适配器
│   │   ├── interface.go         # Provider 接口定义
│   │   ├── openai.go           # OpenAI/Zhipu 适配器
│   │   ├── claude.go           # Anthropic Claude 适配器
│   │   ├── minimax.go          # MiniMax 适配器
│   │   ├── deepseek.go         # DeepSeek 适配器
│   │   ├── glm47.go            # GLM4.7 适配器
│   │   └── siliconflow.go      # 硅基流动适配器
│   ├── context/                 # 上下文管理
│   ├── tools/                   # 工具注册表与内置工具
│   │   ├── registry.go         # 工具注册表
│   │   ├── read_file.go       # 文件读取工具
│   │   └── write_file.go      # 文件写入工具
│   ├── memory/                  # 基于文件的记忆存储
│   └── feishu/                  # 飞书机器人集成
└── README.md
```

## License

MIT
