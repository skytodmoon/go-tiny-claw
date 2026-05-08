package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/skytodmoon/go-tiny-claw/internal/logger"
	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/schema"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

const (
	MaxContextTokens  = 200000
	CompactThreshold  = 0.75
	MaxObservationLen = 8000
)

type AgentEngine struct {
	provider       provider.LLMProvider
	registry       tools.Registry
	WorkDir        string
	EnableThinking bool
	MaxTurns       int
	tokenBudget    int
	logger         *logger.Logger
}

func NewAgentEngine(p provider.LLMProvider, r tools.Registry, workDir string, enableThinking bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       r,
		WorkDir:        workDir,
		EnableThinking: enableThinking,
		MaxTurns:       20,
		tokenBudget:    MaxContextTokens,
		logger:         logger.WithModule("engine"),
	}
}

type QueryState struct {
	Turn        int
	Phase       string
	Thought     string
	Action      string
	Observation string
	ToolCalls   []ToolCallRecord
	IsComplete  bool
}

type ToolCallRecord struct {
	Name      string
	Arguments string
	Result    string
	IsError   bool
}

type ContextLayer struct {
	SystemPrompt string
	Memory       []schema.Message
	Conversation []schema.Message
	Compacted    bool
}

func (e *AgentEngine) buildSystemPrompt(phase string) string {
	basePrompt := `你是 go-tiny-claw，一个专业的 AI Agent 编码助手。

## 核心架构

你运行在一个 **Thinking → Acting → Observation → Re-thinking** 的循环中：

1. **Thinking (思考)**: 分析当前状态，规划下一步行动
2. **Acting (行动)**: 选择并调用合适的工具
3. **Observation (观察)**: 收集工具执行结果
4. **Re-thinking (再思考)**: 根据观察结果调整策略

## 工具系统

可用工具（通过 tool_calls 调用）：
- read_file: 读取文件内容 {"path": "文件路径"}
- write_file: 写入文件 {"path": "路径", "content": "内容"}
- edit_file: 编辑文件 {"path": "路径", "old_str": "原文本", "new_str": "新文本"}
- bash: 执行 bash 命令 {"command": "命令"}

## 执行原则

1. **渐进式探索**: 先理解，再行动，最后验证
2. **错误恢复**: 工具执行失败时，分析原因并重试
3. **任务分解**: 复杂任务分解为多个步骤
4. **上下文感知**: 根据之前的观察结果调整行为

## 输出规则

- 需要调用工具时，使用 tool_calls
- 任务完成时，用自然语言总结
- 遇到错误时，分析原因并提出解决方案`

	switch phase {
	case "thinking":
		return basePrompt + `

## 当前阶段: THINKING

请深入思考：
1. 当前任务进展如何？
2. 下一步应该做什么？
3. 需要使用哪些工具？
4. 可能遇到什么问题？`

	case "acting":
		return basePrompt + `

## 当前阶段: ACTING

请执行行动：
1. 根据思考结果选择工具
2. 准备正确的参数
3. 调用工具执行`

	default:
		return basePrompt
	}
}

func (e *AgentEngine) estimateTokens(messages []schema.Message) int {
	total := 0
	for _, msg := range messages {
		total += len(msg.Content) / 4
		for _, tc := range msg.ToolCalls {
			total += len(tc.Name) + len(tc.Arguments)
		}
	}
	return total
}

func (e *AgentEngine) shouldCompact(messages []schema.Message) bool {
	tokens := e.estimateTokens(messages)
	return float64(tokens)/float64(e.tokenBudget) > CompactThreshold
}

func (e *AgentEngine) compactContext(messages []schema.Message) []schema.Message {
	if len(messages) <= 4 {
		return messages
	}

	e.logger.Info("[Context] 触发自动压缩...")

	summary := "## 历史对话摘要\n\n"
	for _, msg := range messages[1 : len(messages)-2] {
		if msg.Role == schema.RoleUser {
			summary += fmt.Sprintf("- 用户: %s\n", truncate(msg.Content, 100))
		} else if msg.Role == schema.RoleAssistant {
			summary += fmt.Sprintf("- 助手: %s\n", truncate(msg.Content, 100))
		}
	}

	compacted := []schema.Message{
		messages[0],
		{Role: schema.RoleSystem, Content: summary},
	}
	compacted = append(compacted, messages[len(messages)-2:]...)

	e.logger.Info("[Context] 压缩完成: %d 条消息 -> %d 条消息\n", len(messages), len(compacted))
	return compacted
}

func (e *AgentEngine) Run(ctx context.Context, userPrompt string) error {
	e.logger.Debug("[Engine] ═══════════════════════════════════════════════════════════════════")
	e.logger.Debug("[Engine] 🎯 任务: %s", truncate(userPrompt, 52))
	e.logger.Debug("[Engine] ═══════════════════════════════════════════════════════════════════")

	e.logger.Info("[Engine] 启动 Agent Loop, 工作区: %s", e.WorkDir)
	e.logger.Info("[Engine] Token 预算: %d, 压缩阈值: %.0f%%", e.tokenBudget, CompactThreshold*100)

	contextLayer := &ContextLayer{
		SystemPrompt: e.buildSystemPrompt("thinking"),
		Conversation: []schema.Message{
			{Role: schema.RoleUser, Content: userPrompt},
		},
	}

	states := []QueryState{}
	turn := 0

	for turn < e.MaxTurns {
		turn++
		state := QueryState{Turn: turn}

		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🔄 Turn %d", turn)
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")

		if e.shouldCompact(contextLayer.Conversation) {
			contextLayer.Conversation = e.compactContext(contextLayer.Conversation)
			contextLayer.Compacted = true
		}

		messages := []schema.Message{
			{Role: schema.RoleSystem, Content: contextLayer.SystemPrompt},
		}
		messages = append(messages, contextLayer.Conversation...)

		availableTools := e.registry.GetAvailableTools()

		state.Phase = "THINKING"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🧠 Phase 1: THINKING")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Thinking] 分析当前状态...")

		e.logger.Debug("[Thinking] 💭 思考中...")
		e.logger.Debug("[Thinking]   • 分析任务进展")
		e.logger.Debug("[Thinking]   • 规划下一步行动")
		e.logger.Debug("[Thinking]   • 选择合适的工具")

		state.Phase = "ACTING"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🚀 Phase 2: ACTING")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Acting] 执行行动...")

		actionResp, err := e.provider.Generate(ctx, messages, availableTools)
		if err != nil {
			return fmt.Errorf("Acting 阶段失败: %w", err)
		}

		contextLayer.Conversation = append(contextLayer.Conversation, *actionResp)

		if actionResp.Content != "" {
			state.Thought = actionResp.Content
			if len(actionResp.ToolCalls) == 0 {
				e.logger.Debug("[Acting] 💬 %s", actionResp.Content)
			}
		}

		if len(actionResp.ToolCalls) == 0 {
			e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
			e.logger.Debug("[Engine] │ ✅ 任务完成")
			e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
			state.IsComplete = true
			states = append(states, state)
			break
		}

		state.Phase = "OBSERVATION"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 👁️ Phase 3: OBSERVATION")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Observation] 执行 %d 个工具调用...", len(actionResp.ToolCalls))

		// 【设计原则】同一回合内的工具调用被认为是无依赖的，可以并行执行
		// 有依赖关系的操作（如先读取再写入）会由模型拆分成不同回合完成
		// 这种设计确保了：
		//   - 无依赖操作：并行执行，提高性能
		//   - 有依赖操作：通过不同回合保证顺序执行

		// 预分配切片存放并发工具执行结果
		toolResults := make([]struct {
			record        ToolCallRecord
			observationMsg schema.Message
		}, len(actionResp.ToolCalls))

		// 使用 WaitGroup 等待所有协程完成
		var wg sync.WaitGroup

		e.logger.Debug("[Engine] 模型请求并发调用 %d 个工具...", len(actionResp.ToolCalls))

		// 遍历所有工具调用，为每个工具开启一个 Goroutine
		for i, toolCall := range actionResp.ToolCalls {
			wg.Add(1)

			// 开启协程，注意将索引和工具调用作为参数传入，避免闭包变量捕获陷阱
			go func(idx int, call schema.ToolCall) {
				defer wg.Done()

				e.logger.Debug("[Observation]   -> [Go-%d] 🛠️ 触发并行执行: %s", idx+1, call.Name)
				e.logger.Debug("[Observation]   -> [Go-%d] 📥 参数: %s", idx+1, string(call.Arguments))

				// 调用底层 Registry 执行工具
				result := e.registry.Execute(ctx, call)

				var output string
				if result.IsError {
					e.logger.Debug("[Observation]   -> [Go-%d] ❌ 失败: %s", idx+1, result.Output)
					e.logger.Info("[Observation] 工具 %s 执行失败: %s", call.Name, result.Output)
					output = result.Output
				} else {
					output = result.Output
					if len(output) > MaxObservationLen {
						output = output[:MaxObservationLen] + "\n...[已截断]"
					}
					e.logger.Debug("[Observation]   -> [Go-%d] ✅ 成功: %s", idx+1, truncate(output, 150))
					e.logger.Info("[Observation] 工具 %s 执行成功", call.Name)
				}

				// 封装结果，每个协程操作不同的索引，无需加锁
				toolResults[idx] = struct {
					record        ToolCallRecord
					observationMsg schema.Message
				}{
					record: ToolCallRecord{
						Name:      call.Name,
						Arguments: string(call.Arguments),
						Result:    result.Output,
						IsError:   result.IsError,
					},
					observationMsg: schema.Message{
						Role:       schema.RoleUser,
						Content:    result.Output,
						ToolCallID: call.ID,
					},
				}

			}(i, toolCall) // 闭包传参
		}

		// 阻塞等待所有并发协程执行完毕
		wg.Wait()
		e.logger.Debug("[Engine] 所有并发工具执行完毕，开始聚合观察结果...")

		// 聚合结果：按顺序追加到状态和上下文
		for _, result := range toolResults {
			state.ToolCalls = append(state.ToolCalls, result.record)
			contextLayer.Conversation = append(contextLayer.Conversation, result.observationMsg)
		}

		state.Phase = "RE-THINKING"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🔄 Phase 4: RE-THINKING")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Re-thinking] 根据观察结果调整策略...")

		e.logger.Debug("[Re-thinking] 📊 执行结果分析:")
		for i, tc := range state.ToolCalls {
			status := "✅"
			if tc.IsError {
				status = "❌"
			}
			e.logger.Debug("[Re-thinking]   %s 工具 %d (%s): %s", status, i+1, tc.Name, truncate(tc.Result, 50))
		}

		e.logger.Debug("[Re-thinking] 🔄 准备下一轮循环...")
		states = append(states, state)
	}

	if turn >= e.MaxTurns {
		e.logger.Debug("[Engine] ⚠️ 达到最大回合数限制")
		e.logger.Info("[Engine] 达到最大回合数限制")
	}

	e.printSummary(states)
	return nil
}

func (e *AgentEngine) printSummary(states []QueryState) {
	e.logger.Debug("[Engine] ═══════════════════════════════════════════════════════════════════")
	e.logger.Debug("[Engine] │ 📊 执行摘要")
	e.logger.Debug("[Engine] ╠══════════════════════════════════════════════════════════════════")
	e.logger.Debug("[Engine] │ 总回合数: %d", len(states))

	totalTools := 0
	successTools := 0
	for _, s := range states {
		for _, tc := range s.ToolCalls {
			totalTools++
			if !tc.IsError {
				successTools++
			}
		}
	}

	e.logger.Debug("[Engine] │ 工具调用: %d 次 (成功 %d 次)", totalTools, successTools)
	e.logger.Debug("[Engine] ╚══════════════════════════════════════════════════════════════════")
	e.logger.Debug("[Engine] 🎉 任务执行完毕")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
