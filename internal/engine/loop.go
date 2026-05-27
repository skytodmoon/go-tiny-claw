package engine

import (
	"context"
	"fmt"
	"sync"

	ctxpkg "github.com/skytodmoon/go-tiny-claw/internal/context"
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
	EnableThinking bool
	logger         *logger.Logger
	compactor      *ctxpkg.Compactor // 【新增】压缩器实例
}

// 【注意】：移除了 Engine 层级的 WorkDir，因为 WorkDir 现在应该跟随 Session 走！
func NewAgentEngine(p provider.LLMProvider, r tools.Registry, enableThinking bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       r,
		EnableThinking: enableThinking,
		logger:         logger.WithModule("engine"),
		// 【初始化压缩器】：将水位线阈值设为 3000 字符用于极端测试，保护最近的 6 条消息
		compactor: ctxpkg.NewCompactor(3000, 6),
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

// 【核心改造】: 移除 userPrompt 参数，改为接收一个具体的 Session 实例
func (e *AgentEngine) Run(ctx context.Context, session *Session, reporter Reporter) error {
	e.logger.Debug("[Engine] ═══════════════════════════════════════════════════════════════════")
	e.logger.Debug("[Engine] 🎯 唤醒会话 [%s]，锁定工作区: %s", session.ID, session.WorkDir)
	e.logger.Debug("[Engine] ═══════════════════════════════════════════════════════════════════")

	e.logger.Info("[Engine] 启动 Agent Loop, 会话: %s, 工作区: %s", session.ID, session.WorkDir)

	// 根据当前 Session 的工作区，动态组装最新的 System Prompt
	composer := ctxpkg.NewPromptComposer(session.WorkDir)
	systemMsg := composer.Build()
	e.logger.Debug("[Engine] 动态组装 System Prompt 完成")

	states := []QueryState{}
	turn := 0

	for {
		turn++
		state := QueryState{Turn: turn}

		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🔄 Turn %d", turn)
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")

		availableTools := e.registry.GetAvailableTools()

		// 1. 从 Session 提取出近期的 Working Memory (最近 20 条，给压缩器留下充足的判断空间)
		workingMemory := session.GetWorkingMemory(20)

		var contextHistory []schema.Message
		contextHistory = append(contextHistory, systemMsg)
		contextHistory = append(contextHistory, workingMemory...)

		// 2. 【核心注入点】: 在向 Provider 发起推理前，过一遍内存压缩器！
		// 无论你带出了多少上下文，如果字符总数超标，早期日志将被掩码化，超大日志将被掐头去尾
		compactedContext := e.compactor.Compact(contextHistory)

		// 3. ================= Phase 1: Thinking =================
		state.Phase = "THINKING"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🧠 Phase 1: THINKING")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Thinking] 分析当前状态...")

		if e.EnableThinking {
			if reporter != nil {
				reporter.OnThinking(ctx)
			}

			e.logger.Debug("[Thinking] 💭 思考中...")
			e.logger.Debug("[Thinking]   • 分析任务进展")
			e.logger.Debug("[Thinking]   • 规划下一步行动")
			e.logger.Debug("[Thinking]   • 选择合适的工具")

			thinkResp, err := e.provider.Generate(ctx, compactedContext, nil)
			if err != nil {
				return fmt.Errorf("Thinking 阶段失败: %w", err)
			}
			if thinkResp.Content != "" {
				// 【驾驭精髓】：写入 Session 的永远是全量的真实响应，不受 Compact 影响
				session.Append(*thinkResp)
				compactedContext = append(compactedContext, *thinkResp)
				state.Thought = thinkResp.Content
			}
		}

		// 3. ================= Phase 2: Action =================
		state.Phase = "ACTING"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🚀 Phase 2: ACTING")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Acting] 执行行动...")

		actionResp, err := e.provider.Generate(ctx, compactedContext, availableTools)
		if err != nil {
			return fmt.Errorf("Action 阶段失败: %w", err)
		}

		// 【驾驭精髓】：写入 Session 的永远是全量的真实响应，不受 Compact 影响
		session.Append(*actionResp)
		compactedContext = append(compactedContext, *actionResp)

		if actionResp.Content != "" {
			state.Action = actionResp.Content
			e.logger.Debug("[Acting] 💬 %s", truncate(actionResp.Content, 100))
		}

		if actionResp.Content != "" && reporter != nil {
			reporter.OnMessage(ctx, actionResp.Content)
		}

		if len(actionResp.ToolCalls) == 0 {
			e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
			e.logger.Debug("[Engine] │ ✅ 任务完成")
			e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
			state.IsComplete = true
			states = append(states, state)
			break
		}

		// 4. ================= 并发执行底层工具 =================
		state.Phase = "OBSERVATION"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 👁️ Phase 3: OBSERVATION")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Observation] 执行 %d 个工具调用...", len(actionResp.ToolCalls))

		observationMsgs := make([]schema.Message, len(actionResp.ToolCalls))
		toolResults := make([]ToolCallRecord, len(actionResp.ToolCalls))
		var wg sync.WaitGroup

		e.logger.Debug("[Engine] 模型请求并发调用 %d 个工具...", len(actionResp.ToolCalls))

		for i, toolCall := range actionResp.ToolCalls {
			wg.Add(1)

			go func(idx int, call schema.ToolCall) {
				defer wg.Done()

				e.logger.Debug("[Observation]   -> [Go-%d] 🛠️ 触发并行执行: %s", idx+1, call.Name)

				if reporter != nil {
					reporter.OnToolCall(ctx, call.Name, string(call.Arguments))
				}

				result := e.registry.Execute(ctx, call)

				if reporter != nil {
					displayOutput := result.Output
					if len(displayOutput) > 200 {
						displayOutput = displayOutput[:200] + "... (已截断)"
					}
					reporter.OnToolResult(ctx, call.Name, displayOutput, result.IsError)
				}

				toolResults[idx] = ToolCallRecord{
					Name:      call.Name,
					Arguments: string(call.Arguments),
					Result:    result.Output,
					IsError:   result.IsError,
				}

				if result.IsError {
					e.logger.Debug("[Observation]   -> [Go-%d] ❌ 失败: %s", idx+1, truncate(result.Output, 50))
				} else {
					e.logger.Debug("[Observation]   -> [Go-%d] ✅ 成功", idx+1)
				}

				observationMsgs[idx] = schema.Message{
					Role:       schema.RoleUser,
					Content:    result.Output,
					ToolCallID: call.ID,
				}
			}(i, toolCall)
		}

		wg.Wait()
		e.logger.Debug("[Engine] 所有并发工具执行完毕")

		session.Append(observationMsgs...)
		state.ToolCalls = toolResults

		// 5. ================= Re-thinking =================
		state.Phase = "RE-THINKING"
		e.logger.Debug("[Engine] ┌──────────────────────────────────────────────────────────────────┐")
		e.logger.Debug("[Engine] │ 🔄 Phase 4: RE-THINKING")
		e.logger.Debug("[Engine] └──────────────────────────────────────────────────────────────────┘")
		e.logger.Info("[Re-thinking] 根据观察结果调整策略...")

		states = append(states, state)
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
