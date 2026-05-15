# Go Tiny Claw

一个轻量级的 AI智能助手 框架，使用 Go 语言实现。

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

### 🟡 Level 2（已完成｜进阶）

在 Level 1 基础上实现：

- ✅ 让 Thinking 阶段更结构化（输出计划而不是一句话）
- ✅ 支持多步执行（一次任务多步完成）
- ✅ Agent 根据情况选择不同工具（动态选择）

👉 关键词：**让 Agent 开始"连续行动"**

### 🔴 Level 3（已完成｜对标业界 Claude Code 最佳实践）

结合 Claude Code 设计思路，实现：

- ✅ 四阶段循环：Thinking → Acting → Observation → Re-thinking
- ✅ Agent 根据执行结果调整下一步行为
- ✅ Prompt 结构优化（思考 vs 执行分离）
- ✅ Context Engineering（动态上下文组装）
- ✅ 自动压缩机制（Token 预算管理）
- ✅ 执行摘要统计

👉 关键词：**让 Agent 开始"像一个系统在运转"**

> Claude Code 源码泄漏存档： https://github.com/FlyAIBox/civil-engineering-cloud-claude-code-source-v2.1.88

## 当前已实现

### 核心功能
- ✅ 多 Provider 支持（OpenAI、Claude、MiniMax、DeepSeek、GLM47、SiliconFlow）
- ✅ 上下文管理（Token 监控、自动压缩）
- ✅ 工具注册表（ReadFile、WriteFile、EditFile、Bash）
- ✅ 四阶段循环（Thinking → Acting → Observation → Re-thinking）
- ✅ **工具调用并行化**：同一回合内无依赖的工具调用并发执行，提高性能

### 并发设计原则

框架采用智能并发策略：

| 场景 | 执行方式 | 说明 |
|------|----------|------|
| 同一回合内多个工具调用 | **并行执行** | 模型单次思考请求的多个工具调用被认为是无依赖的 |
| 有依赖关系的操作 | **串行执行** | 通过不同回合保证顺序（如先读取再写入） |

这种设计确保：
- **性能优化**：无依赖操作并行执行，减少总耗时
- **数据一致性**：有依赖操作通过回合机制保证顺序执行

### Level 1 功能
- ✅ 慢思考机制（Thinking Phase）
- ✅ 工具调用（Tool Registry）
- ✅ 单次任务执行

### Level 2 功能
- ✅ 结构化思考（输出计划）
- ✅ 多步执行
- ✅ 动态工具选择

### Level 3 功能
- ✅ Context Engineering（动态上下文组装）
- ✅ 自动压缩机制（Token 预算管理）
- ✅ Prompt 结构优化（思考 vs 执行分离）
- ✅ 执行结果反馈机制
- ✅ 执行摘要统计
- ✅ 工具调用并行化（并发执行）
- ✅ **技能懒加载机制（Progressive Disclosure）**

### 技能懒加载机制

框架采用"渐进式暴露"策略，大幅优化 Token 消耗：

| 阶段 | 加载内容 | Token 消耗 |
|------|----------|------------|
| **启动时** | 仅加载技能元数据（名称、触发描述、标签） | ~2,000 Token |
| **执行中** | 按需加载技能完整正文 | ~1,500 Token/技能 |

**核心设计**：
- **read_skill 工具**：模型判定需要某个技能时，主动调用此工具加载完整正文
- **智能注入**：技能正文以系统消息形式注入上下文，供后续对话使用
- **缓存机制**：已加载的技能会保留在上下文中，无需重复加载

**优势**：
- **Token 效率提升**：从启动时消耗 30,000+ Token（50 个技能）降至 ~2,000 Token
- **扩展性增强**：技能数量不再受 Token 限制
- **模型感知**：模型显式控制技能加载时机

### 日志系统
- ✅ 日志级别控制（DEBUG、INFO、WARN、ERROR、FATAL）
- ✅ 日志文件输出（按日期分割）
- ✅ 彩色日志输出（控制台）
- ✅ 模块化日志（每个模块独立日志）
- ✅ 日志格式化（时间戳、级别、模块名）

## 快速开始

```bash
# 设置 API Key
export SILICONFLOW_API_KEY="your-api-key"
export NVIDIA_API_KEY="your-nvidia-key"

# 配置 Provider 顺序（可选，默认 nvidia,siliconflow）
export PROVIDER_ORDER="nvidia,siliconflow"

# 配置模型（可选）
export NVIDIA_MODEL="minimaxai/minimax-m2.7"
export SILICONFLOW_MODEL="deepseek-ai/DeepSeek-V3"

# 运行
go run cmd/claw/main.go
```

### LLM 服务容错机制

框架支持多 Provider 故障转移，当主服务不可用时自动切换到备用服务：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `PROVIDER_ORDER` | Provider 优先级顺序，逗号分隔 | `nvidia,siliconflow` |
| `NVIDIA_API_KEY` | NVIDIA API Key | - |
| `NVIDIA_MODEL` | NVIDIA 使用的模型 | `minimaxai/minimax-m2.7` |
| `SILICONFLOW_API_KEY` | SiliconFlow API Key | - |
| `SILICONFLOW_MODEL` | SiliconFlow 使用的模型 | `deepseek-ai/DeepSeek-V3` |
| `LLM_TIMEOUT` | 单个 Provider 超时时间 | `30s` |
| `MAX_RETRIES` | 最大重试次数 | `2` |

**故障转移流程**：
```
用户请求 → NVIDIA（主）→ 失败/超时 → SiliconFlow（备）→ 失败 → 重试
```

**推荐模型**：

| Provider | 模型 | 参数 | 说明 |
|----------|------|------|------|
| NVIDIA | `minimaxai/minimax-m2.7` | 7B | 默认，平衡性能和效果 |
| NVIDIA | `meta/llama-3.3-70b-instruct` | 70B | 高质量 |
| SiliconFlow | `deepseek-ai/DeepSeek-V3` | 67B | 默认，高性能 |
| SiliconFlow | `Qwen/Qwen2-72B-Instruct` | 72B | 高质量 |

## 项目结构

```
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
│   │   ├── glm47.go            # GLM4.7/NVIDIA 适配器
│   │   ├── nvidia.go           # NVIDIA API 适配器
│   │   ├── siliconflow.go      # 硅基流动适配器
│   │   └── failover.go         # 故障转移 Provider
│   ├── context/                 # 上下文管理
│   ├── tools/                   # 工具注册表与内置工具
│   │   ├── registry.go         # 工具注册表
│   │   ├── read_file.go       # 文件读取工具
│   │   ├── write_file.go      # 文件写入工具
│   │   ├── edit_file.go       # 文件编辑工具
│   │   ├── bash.go            # Bash 命令工具
│   │   └── read_skill.go      # 技能懒加载工具
│   ├── memory/                  # 基于文件的记忆存储
│   ├── feishu/                  # 飞书机器人集成
│   └── logger/                  # 日志系统
├── integration/                 # 集成测试
├── testutils/                   # 测试工具
├── workspace/                   # 工作区目录
│   ├── .claw/skills/           # 技能文件目录
│   │   ├── weather/            # 天气查询技能
│   │   ├── news/               # 新闻查询技能
│   │   └── git-workflow/       # Git 提交流程技能
│   └── AGENTS.md               # 项目架构规范
└── README.md
```

## 测试任务套件

框架提供了 15 种不同类型的测试任务，用于验证稳定性：

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

运行测试套件：
```bash
go run cmd/task-runner/main.go
```

## License

MIT
