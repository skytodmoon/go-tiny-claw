# Go Tiny Claw

| 部分 | 内容 |
|------|------|
| 标题 | Go Tiny Claw |
| 描述 | 一个轻量级的 AI智能助手 框架，使用 Go 语言实现。 |
| 核心理念 | **"能慢思考、会用工具、能推进任务"** 的 Agent |

## 版本规划

| 版本 | 状态 | 内容 |
|------|------|------|
| Level 1 | 🟢 已完成 | - 慢思考（Thinking 阶段）<br>- 工具调用（Tool Registry）<br>- YOLO 执行一次完整任务 |
| Level 2 | 🟡 已完成 | - 结构化思考<br>- 多步执行<br>- 动态工具选择 |
| Level 3 | 🔴 已完成 | - 四阶段循环<br>- 动态上下文组装<br>- Token 预算管理<br>- 执行摘要统计 |

## 当前已实现

| 功能类型 | 实现内容 |
|----------|----------|
| 核心功能 | - 多 Provider 支持<br>- 上下文管理<br>- 工具注册表<br>- 四阶段循环 |
| Level 1 功能 | - 慢思考机制<br>- 工具调用<br>- 单次任务执行 |
| Level 2 功能 | - 结构化思考<br>- 多步执行<br>- 动态工具选择 |
| Level 3 功能 | - Context Engineering<br>- 自动压缩机制<br>- Prompt 结构优化 |
| 日志系统 | - 日志级别控制<br>- 日志文件输出<br>- 彩色日志输出 |

## 快速开始

```bash
# 设置 API Key
export SILICONFLOW_API_KEY="your-api-key"

# 运行
go run cmd/claw/main.go
```

## 项目结构

```
/go-tiny-claw/
├── cmd/
├── internal/
├── integration/
└── testutils/
```

## 测试任务套件

| 任务 | 类型 | 描述 |
|------|------|------|
| 1 | 文件读写基础 | 读取 README.md 并创建项目说明文档 |
| 2 | 文件编辑 | 读取文件并替换内容 |
| 3 | 多文件创建 | 同时创建多个文件 |
| ... | ... | ... |

## License

MIT