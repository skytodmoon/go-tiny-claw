# Go Tiny Claw 项目摘要

## 项目概述

**Go Tiny Claw** 是一个轻量级的 AI 智能助手框架，使用 Go 语言实现。其核心理念是让 Agent 具备"慢思考、会使用工具、能推进任务"的能力。

## 核心功能

- **四阶段循环**: Thinking → Acting → Observation → Re-thinking，确保任务执行的系统性与动态调整能力。
- **多工具支持**: 内置文件读写、编辑、Bash 命令等工具，并通过统一注册机制动态调用。
- **上下文管理**: 支持动态上下文组装与 Token 预算管理，优化资源使用。
- **日志系统**: 提供模块化、多级别、文件分割与彩色输出的日志功能。

## 阶段实现

- **Level 1**: 支持慢思考、工具调用与单次任务执行。
- **Level 2**: 实现结构化思考、多步执行与动态工具选择。
- **Level 3**: 对标业界 Claude Code 最佳实践，优化 Prompt 结构、执行摘要统计与反馈机制。

## 项目结构

- **Provider 适配器**: 支持 OpenAI、Claude、MiniMax、DeepSeek、GLM47、SiliconFlow 等多种 LLM 提供商。
- **核心模块**: 包括 MainLoop 实现、上下文管理、工具注册表与内置工具。
- **测试套件**: 提供 15 种任务类型，覆盖文件读写、代码生成、文档转换等场景。

## 快速开始

```bash
# 设置 API Key
# 运行程序
go run cmd/claw/main.go
```

## 测试任务

运行测试套件：
```bash
go run cmd/task-runner/main.go
```

## 许可证

MIT License