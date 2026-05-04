# Go Tiny Claw - 项目说明文档

## 概述

Go Tiny Claw 是一个轻量级的 AI 智能助手框架，使用 Go 语言实现。其核心理念是**"能慢思考、会用工具、能推进任务"**的 Agent。

## 核心功能

1. **慢思考机制**：显式区分 Thinking / Acting 阶段，先思考再执行。
2. **工具调用**：支持文件读取、写入、编辑以及 Bash 命令执行等工具。
3. **多步任务执行**：支持动态工具选择和多步任务完成。
4. **上下文管理**：包括 Token 监控、自动压缩和动态上下文组装。
5. **日志系统**：模块化日志、文件输出和彩色控制台输出。

## 快速开始

1. 设置 API Key（选择以下一种）：
   ```bash
   export SILICONFLOW_API_KEY="your-api-key"
   export ZHIPU_API_KEY="your-api-key"
   export NVIDIA_API_KEY="your-api-key"
   ```

2. 运行程序：
   ```bash
   go run cmd/claw/main.go
   ```

## 项目结构

```
go-tiny-claw/
├── cmd/               # 程序入口
├── internal/          # 核心实现（引擎、Provider、工具、日志等）
├── integration/       # 集成测试
├── testutils/         # 测试工具
├── README.md          # 详细文档
```

## 测试任务套件

框架提供 15 种测试任务，涵盖文件读写、代码生成、错误处理等多种类型。运行测试套件：
```bash
   go run cmd/task-runner/main.go
```

## License

MIT