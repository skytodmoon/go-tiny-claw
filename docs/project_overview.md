# Go Tiny Claw 项目概述

## 核心理念

"能慢思考、会用工具、能推进任务" 的 AI Agent 框架。

## 功能实现

### 核心功能
- 多 Provider 支持（OpenAI、Claude、MiniMax、DeepSeek、GLM47、SiliconFlow）
- 上下文管理（Token 监控、自动压缩）
- 工具注册表（ReadFile、WriteFile、EditFile）
- 四阶段循环（Thinking → Acting → Observation → Re-thinking）

### Level 1 功能
- 慢思考机制（Thinking Phase）
- 工具调用（Tool Registry）
- 单次任务执行

### Level 2 功能
- 结构化思考（输出计划）
- 多步执行
- 动态工具选择

### Level 3 功能
- Context Engineering（动态上下文组装）
- 自动压缩机制（Token 预算管理）
- Prompt 结构优化（思考 vs 执行分离）
- 执行结果反馈机制
- 执行摘要统计

### 日志系统
- 日志级别控制（DEBUG、INFO、WARN、ERROR、FATAL）
- 日志文件输出（按日期分割）
- 彩色日志输出（控制台）
- 模块化日志（每个模块独立日志）
- 日志格式化（时间戳、级别、模块名）
