# 代码审查报告：internal/engine/loop.go

## 代码结构

### 核心组件
1. **AgentEngine**: 主引擎，负责管理 Agent 的运行逻辑。
2. **QueryState**: 记录当前任务的执行状态（如阶段、思考、行动等）。
3. **ContextLayer**: 管理上下文层，包括系统提示、对话历史和压缩状态。

### 主要功能
- `Run`: 启动 Agent 的主循环。
- `buildSystemPrompt`: 根据当前阶段生成系统提示。
- `estimateTokens` 和 `shouldCompact`: 管理 token 使用。
- `compactContext`: 压缩历史对话。

## 优化建议

### 性能优化
1. **Token 压缩逻辑**: 
   - 改进为智能摘要生成（如 NLP 模型）。
   - 优先压缩高频或重复内容。
2. **Context Management**: 
   - 动态调整压缩触发条件。

### 可维护性
1. **模块化**: 
   - 拆分 `Run` 方法为小方法（如 `thinkingPhase`, `actingPhase`）。
   - 提取工具调用逻辑为独立方法。
2. **日志**: 
   - 结构化日志（如 JSON 格式）。

### 功能增强
1. **工具调用失败重试**: 提供重试逻辑或备用方案。
2. **动态调整**: 根据任务动态调整 `MaxTurns` 和 `tokenBudget`。
3. **测试覆盖率**: 补充测试用例（如 `compactContext`）。

### 可读性改进
1. **常量命名**: 更具体（如 `DefaultMaxContextTokens`）。
2. **注释补充**: 对复杂结构字段添加注释。