# Go Tiny Claw 项目说明

## 项目概述

Go Tiny Claw 是一个轻量级的 AI 智能助手框架，使用 Go 语言实现。其核心理念是让 Agent 具备"慢思考、工具调用和任务推进"的能力。

## 核心功能

- **慢思考机制（Thinking Phase）**: 显式区分 Thinking / Acting，先思考"下一步做什么"，再执行。
- **工具调用（Tool Registry）**: 支持多种工具（如文件读取、编辑、Bash 命令等），通过统一注册与分发机制调用。
- **四阶段循环**: Thinking → Acting → Observation → Re-thinking，确保任务逐步推进。

## 版本规划

### Level 1（已完成）
- 慢思考机制
- 工具调用
- 单次任务执行

### Level 2（已完成）
- 结构化思考
- 多步执行
- 动态工具选择

### Level 3（已完成）
- Context Engineering（动态上下文组装）
- 自动压缩机制（Token 预算管理）
- Prompt 结构优化

## 项目结构

```
go-tiny-claw/
├── cmd/          # 程序入口
├── internal/     # 核心实现（引擎、Provider适配器、工具等）
├── integration/  # 集成测试
└── testutils/    # 测试工具
```

## 快速开始

1. 设置 API Key：
   ```bash
   export SILICONFLOW_API_KEY="your-api-key"
   ```
2. 运行：
   ```bash
   go run cmd/claw/main.go
   ```

## 测试任务

框架提供了 15 种测试任务，包括文件读写、代码生成、错误处理等。运行测试：
```bash
go run cmd/task-runner/main.go
```

## 许可证
MIT