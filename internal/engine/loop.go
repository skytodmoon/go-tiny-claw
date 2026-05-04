package engine

import (
	"context"
	"fmt"
	"log"

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
}

func NewAgentEngine(p provider.LLMProvider, r tools.Registry, workDir string, enableThinking bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       r,
		WorkDir:        workDir,
		EnableThinking: enableThinking,
		MaxTurns:       20,
		tokenBudget:    MaxContextTokens,
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

	log.Println("[Context] 触发自动压缩...")

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

	log.Printf("[Context] 压缩完成: %d 条消息 -> %d 条消息\n", len(messages), len(compacted))
	return compacted
}

func (e *AgentEngine) Run(ctx context.Context, userPrompt string) error {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════════╗")
	fmt.Printf("  🎯 任务: %s\n", truncate(userPrompt, 52))
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")

	log.Printf("[Engine] 启动 Agent Loop, 工作区: %s\n", e.WorkDir)
	log.Printf("[Engine] Token 预算: %d, 压缩阈值: %.0f%%\n", e.tokenBudget, CompactThreshold*100)

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

		fmt.Printf("\n┌──────────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│ 🔄 Turn %d                                                        │\n", turn)
		fmt.Println("└──────────────────────────────────────────────────────────────────┘")

		if e.shouldCompact(contextLayer.Conversation) {
			contextLayer.Conversation = e.compactContext(contextLayer.Conversation)
			contextLayer.Compacted = true
		}

		messages := []schema.Message{
			{Role: schema.RoleSystem, Content: contextLayer.SystemPrompt},
		}
		messages = append(messages, contextLayer.Conversation...)

		availableTools := e.registry.GetAvailableTools()

		fmt.Println("\n┌──────────────────────────────────────────────────────────────────┐")
		fmt.Println("│ 🧠 Phase 1: THINKING                                             │")
		fmt.Println("└──────────────────────────────────────────────────────────────────┘")
		state.Phase = "THINKING"
		log.Println("[Thinking] 分析当前状态...")

		fmt.Println("\n💭 思考中...")
		fmt.Println("   • 分析任务进展")
		fmt.Println("   • 规划下一步行动")
		fmt.Println("   • 选择合适的工具")

		fmt.Println("\n┌──────────────────────────────────────────────────────────────────┐")
		fmt.Println("│ 🚀 Phase 2: ACTING                                               │")
		fmt.Println("└──────────────────────────────────────────────────────────────────┘")
		state.Phase = "ACTING"
		log.Println("[Acting] 执行行动...")

		actionResp, err := e.provider.Generate(ctx, messages, availableTools)
		if err != nil {
			return fmt.Errorf("Acting 阶段失败: %w", err)
		}

		contextLayer.Conversation = append(contextLayer.Conversation, *actionResp)

		if actionResp.Content != "" {
			state.Thought = actionResp.Content
			if len(actionResp.ToolCalls) == 0 {
				fmt.Printf("\n💬 %s\n", actionResp.Content)
			}
		}

		if len(actionResp.ToolCalls) == 0 {
			fmt.Println("\n┌──────────────────────────────────────────────────────────────────┐")
			fmt.Println("│ ✅ 任务完成                                                      │")
			fmt.Println("└──────────────────────────────────────────────────────────────────┘")
			state.IsComplete = true
			states = append(states, state)
			break
		}

		fmt.Println("\n┌──────────────────────────────────────────────────────────────────┐")
		fmt.Println("│ 👁️ Phase 3: OBSERVATION                                          │")
		fmt.Println("└──────────────────────────────────────────────────────────────────┘")
		state.Phase = "OBSERVATION"
		log.Printf("[Observation] 执行 %d 个工具调用...\n", len(actionResp.ToolCalls))

		for i, toolCall := range actionResp.ToolCalls {
			record := ToolCallRecord{
				Name:      toolCall.Name,
				Arguments: string(toolCall.Arguments),
			}

			fmt.Printf("\n   🛠️  工具 #%d: %s\n", i+1, toolCall.Name)
			fmt.Printf("   📥 参数: %s\n", string(toolCall.Arguments))

			result := e.registry.Execute(ctx, toolCall)
			record.Result = result.Output
			record.IsError = result.IsError

			if result.IsError {
				fmt.Printf("   ❌ 失败: %s\n", result.Output)
				log.Printf("[Observation] 工具 %s 执行失败: %s\n", toolCall.Name, result.Output)
			} else {
				output := result.Output
				if len(output) > MaxObservationLen {
					output = output[:MaxObservationLen] + "\n...[已截断]"
				}
				fmt.Printf("   ✅ 成功: %s\n", truncate(output, 150))
				log.Printf("[Observation] 工具 %s 执行成功\n", toolCall.Name)
			}

			state.ToolCalls = append(state.ToolCalls, record)

			observationMsg := schema.Message{
				Role:       schema.RoleUser,
				Content:    result.Output,
				ToolCallID: toolCall.ID,
			}
			contextLayer.Conversation = append(contextLayer.Conversation, observationMsg)
		}

		fmt.Println("\n┌──────────────────────────────────────────────────────────────────┐")
		fmt.Println("│ 🔄 Phase 4: RE-THINKING                                          │")
		fmt.Println("└──────────────────────────────────────────────────────────────────┘")
		state.Phase = "RE-THINKING"
		log.Println("[Re-thinking] 根据观察结果调整策略...")

		fmt.Println("\n📊 执行结果分析:")
		for i, tc := range state.ToolCalls {
			status := "✅"
			if tc.IsError {
				status = "❌"
			}
			fmt.Printf("   %s 工具 %d (%s): %s\n", status, i+1, tc.Name, truncate(tc.Result, 50))
		}

		fmt.Println("\n🔄 准备下一轮循环...")
		states = append(states, state)
	}

	if turn >= e.MaxTurns {
		fmt.Println("\n⚠️ 达到最大回合数限制")
		log.Println("[Engine] 达到最大回合数限制")
	}

	e.printSummary(states)
	return nil
}

func (e *AgentEngine) printSummary(states []QueryState) {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("│ 📊 执行摘要                                                      │")
	fmt.Println("╠══════════════════════════════════════════════════════════════════╣")
	fmt.Printf("│ 总回合数: %d                                                     │\n", len(states))

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

	fmt.Printf("│ 工具调用: %d 次 (成功 %d 次)                                    │\n", totalTools, successTools)
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")

	fmt.Println("\n╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("│ 🎉 任务执行完毕                                                  │")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
