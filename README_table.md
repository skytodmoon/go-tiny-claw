| 标题 | 内容 |
|------|------|
| **项目名称** | Go Tiny Claw |
| **简介** | 一个轻量级的 AI智能助手 框架，使用 Go 语言实现。 |
| **核心理念** | **"能慢思考、会用工具、能推进任务"** 的 Agent |
| **版本规划** |
| 🟢 Level 1（已完成｜基于前 6 讲） |
| 功能 |
| 1. **慢思考（Thinking 阶段）** | 在 Main Loop 中，显式区分 Thinking / Acting；先思考"下一步做什么"，再执行 |
| 2. **工具调用（Tool Registry）** | 至少接入 1 个工具（如：文件读取 / 简单编辑）；通过统一注册与分发机制调用（而不是写死） |
| 3. **YOLO 执行一次完整任务** | 给定一个目标（如：读取一个文件并总结）；Agent 自己决定调用工具并完成任务 |
| 关键词 | **慢思考 → 动手做 → 得到结果** |
| 🟡 Level 2（已完成｜进阶） |
| 功能 |
| - 结构化思考 | 让 Thinking 阶段更结构化（输出计划而不是一句话） |
| - 多步执行 | 支持多步执行（一次任务多步完成） |
| - 动态工具选择 | Agent 根据情况选择不同工具（动态选择） |
| 关键词 | **让 Agent 开始"连续行动"** |
| 🔴 Level 3（已完成｜对标业界 Claude Code 最佳实践） |
| 功能 |
| - 四阶段循环 | Thinking → Acting → Observation → Re-thinking |
| - 执行结果反馈 | Agent 根据执行结果调整下一步行为 |
| - Prompt 结构优化 | 思考 vs 执行分离 |
| - Context Engineering | 动态上下文组装 |
| - 自动压缩机制 | Token 预算管理 |
| - 执行摘要统计 | |
| 关键词 | **让 Agent 开始"像一个系统在运转"** |
| **当前已实现** |
| **核心功能** |
| - 多 Provider 支持 | OpenAI、Claude、MiniMax、DeepSeek、GLM47、SiliconFlow |
| - 上下文管理 | Token 监控、自动压缩 |
| - 工具注册表 | ReadFile、WriteFile、EditFile |
| - 四阶段循环 | Thinking → Acting → Observation → Re-thinking |
| **Level 1 功能** |
| - 慢思考机制 | Thinking Phase |
| - 工具调用 | Tool Registry |
| - 单次任务执行 | |
| **Level 2 功能** |
| - 结构化思考 | 输出计划 |
| - 多步执行 |  |
| - 动态工具选择 |  |
| **Level 3 功能** |
| - Context Engineering | 动态上下文组装 |
| - 自动压缩机制 | Token 预算管理 |
| - Prompt 结构优化 | 思考 vs 执行分离 |
| - 执行结果反馈机制 | |
| - 执行摘要统计 | |
| **日志系统** |
| - 日志级别控制 | DEBUG、INFO、WARN、ERROR、FATAL |
| - 日志文件输出 | 按日期分割 |
| - 彩色日志输出 | 控制台 |
| - 模块化日志 | 每个模块独立日志 |
| - 日志格式化 | 时间戳、级别、模块名 |
| **快速开始** |
| 命令 | ```bash
export SILICONFLOW_API_KEY="your-api-key"
export ZHIPU_API_KEY="your-api-key"
export NVIDIA_API_KEY="your-api-key"

go run cmd/claw/main.go
``` |
| **项目结构** |
| 结构 | ```
go-tiny-claw/
├── cmd/
│   ├── claw/
│   │   └── main.go              # 程序入口
│   ├── level3-test/
│   │   └── main.go              # Level 3 验证测试
│   └── task-runner/
│       └── main.go              # 多任务测试套件
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
│   │   ├── write_file.go      # 文件写入工具
│   │   ├── edit_file.go       # 文件编辑工具
│   │   └── bash.go            # Bash 命令工具
│   ├── memory/                  # 基于文件的记忆存储
│   ├── feishu/                  # 飞书机器人集成
│   └── logger/                  # 日志系统
├── integration/                 # 集成测试
├── testutils/                   # 测试工具
└── README.md
``` |
| **测试任务套件** |
| 任务 | 类型 | 描述 |
|------|------|------|
| 1 | 文件读写基础 | 读取 README.md 并创建项目说明文档 |
| 2 | 文件编辑 | 读取文件并替换内容 |
| 3 | 多文件创建 | 同时创建多个文件 |
| 4 | 内容分析 | 分析项目功能和架构 |
| 5 | 代码生成 | 创建 Go 程序并验证 |
| 6 | 目录结构 | 探索目录并记录结构 |
| 7 | 错误处理 | 处理不存在文件的情况 |
| 8 | 文档转换 | 内容格式转换 |
| 9 | 配置文件生成 | 创建 YAML 配置文件 |
| 10 | 测试报告 | 运行测试并记录结果 |
| 11 | 文档摘要 | 提取关键信息创建摘要 |
| 12 | 代码审查 | 分析代码结构 |
| 13 | 脚本生成 | 创建 Shell 脚本 |
| 14 | 依赖分析 | 分析项目依赖 |
| 15 | 综合任务 | 多步骤复杂任务 |
| **运行测试套件** | `go run cmd/task-runner/main.go` |
| **License** | MIT |