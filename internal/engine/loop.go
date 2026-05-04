package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/schema"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

type AgentEngine struct {
	provider       provider.LLMProvider
	registry       tools.Registry
	WorkDir        string
	EnableThinking bool
	MaxTurns       int
}

func NewAgentEngine(p provider.LLMProvider, r tools.Registry, workDir string, enableThinking bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       r,
		WorkDir:        workDir,
		EnableThinking: enableThinking,
		MaxTurns:       20,
	}
}

type PlanStep struct {
	Step      int     `json:"step"`
	Action    string  `json:"action"`
	Tool      string  `json:"tool,omitempty"`
	Thought   string  `json:"thought"`
}

type StructuredPlan struct {
	Goal       string     `json:"goal"`
	Steps      []PlanStep `json:"steps"`
	Total      int        `json:"total_steps"`
	Confidence float64    `json:"confidence"`
}

func parsePlan(content string) (*StructuredPlan, error) {
	if !strings.Contains(content, "{") {
		return &StructuredPlan{
			Goal: content,
			Steps: []PlanStep{
				{Step: 1, Action: "直接回答", Thought: content},
			},
			Total: 1,
		}, nil
	}

	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}") + 1
	if startIdx == -1 || endIdx == -1 {
		return nil, fmt.Errorf("无法解析计划 JSON")
	}

	jsonStr := content[startIdx:endIdx]
	var plan StructuredPlan
	if err := json.Unmarshal([]byte(jsonStr), &plan); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	return &plan, nil
}

type JSONToolCall struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Content   string `json:"content"`
	OldStr    string `json:"old_str"`
	NewStr    string `json:"new_str"`
}

func extractJSONToolCalls(content string, availableTools []schema.ToolDefinition) []schema.ToolCall {
	var toolCalls []schema.ToolCall

	for _, tool := range availableTools {
		toolNameLower := strings.ToLower(tool.Name)

		switch toolNameLower {
		case "read_file":
			if strings.Contains(content, `"path"`) && strings.Contains(content, `"name"`) && strings.Contains(content, `"read_file"`) {
				start := strings.Index(content, "{")
				end := strings.LastIndex(content, "}") + 1
				if start != -1 && end > start {
					var tc JSONToolCall
					if json.Unmarshal([]byte(content[start:end]), &tc) == nil && tc.Path != "" {
						args, _ := json.Marshal(map[string]string{"path": tc.Path})
						toolCalls = append(toolCalls, schema.ToolCall{
							ID:        fmt.Sprintf("call_%d", len(toolCalls)+1),
							Name:      "read_file",
							Arguments: args,
						})
					}
				}
			}

		case "write_file":
			if strings.Contains(content, `"path"`) && strings.Contains(content, `"content"`) {
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					if strings.Contains(line, `"path"`) && strings.Contains(line, `"SUMMARY"`) {
						start := strings.Index(content, "{")
						end := strings.LastIndex(content, "}") + 1
						if start != -1 && end > start {
							var tc JSONToolCall
							if json.Unmarshal([]byte(content[start:end]), &tc) == nil && tc.Path != "" && tc.Content != "" {
								args, _ := json.Marshal(map[string]string{"path": tc.Path, "content": tc.Content})
								toolCalls = append(toolCalls, schema.ToolCall{
									ID:        fmt.Sprintf("call_%d", len(toolCalls)+1),
									Name:      "write_file",
									Arguments: args,
								})
							}
						}
						break
					}
				}
			}

		case "edit_file":
			if strings.Contains(content, `"old_str"`) && strings.Contains(content, `"new_str"`) {
				start := strings.Index(content, "{")
				end := strings.LastIndex(content, "}") + 1
				if start != -1 && end > start {
					var tc JSONToolCall
					if json.Unmarshal([]byte(content[start:end]), &tc) == nil && tc.OldStr != "" && tc.NewStr != "" {
						args, _ := json.Marshal(map[string]string{"old_str": tc.OldStr, "new_str": tc.NewStr, "path": tc.Path})
						toolCalls = append(toolCalls, schema.ToolCall{
							ID:        fmt.Sprintf("call_%d", len(toolCalls)+1),
							Name:      "edit_file",
							Arguments: args,
						})
					}
				}
			}
		}
	}

	return toolCalls
}

func (e *AgentEngine) Run(ctx context.Context, userPrompt string) error {
	log.Printf("[Engine] 引擎启动，锁定工作区: %s\n", e.WorkDir)
	log.Printf("[Engine] 慢思考模式: %v, 最大回合数: %d\n", e.EnableThinking, e.MaxTurns)
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ 🎯 用户任务: %-48s║\n", truncate(userPrompt, 48))
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	contextHistory := []schema.Message{
		{
			Role: schema.RoleSystem,
			Content: `你是 go-tiny-claw，一个专业的编码助手。

## 核心指令

### 1. 工具使用
可用工具（必须严格使用 JSON 格式调用）：
- read_file: {"path": "文件路径"}
- write_file: {"path": "文件路径", "content": "文件内容"}
- edit_file: {"path": "文件路径", "old_str": "原字符串", "new_str": "新字符串"}

### 2. 输出规则
- 如果需要调用工具，必须使用 JSON 格式的工具调用
- 不要在 content 中输出 JSON，只通过 tool_calls 调用工具
- 如果可以直接回答用户问题，直接回答即可

### 3. 总结回复
当任务完成时，用自然语言总结给用户。`,
		},
		{
			Role:    schema.RoleUser,
			Content: userPrompt,
		},
	}

	turnCount := 0

	for turnCount < e.MaxTurns {
		turnCount++
		fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("│ 🔄 Turn %d                                                    │\n", turnCount)
		fmt.Println("└──────────────────────────────────────────────────────────────┘")

		availableTools := e.registry.GetAvailableTools()

		// =====================================================
		// Phase 1: THINKING - 思考阶段
		// =====================================================
		fmt.Println("\n🧠 ═══ [Phase 1: THINKING] 思考阶段 ═══")
		fmt.Println("   分析当前状态，决定下一步行动...")
		log.Println("[Thinking] 开始思考当前任务状态...")

		// =====================================================
		// Phase 2: ACTING - 行动阶段
		// =====================================================
		fmt.Println("\n🚀 ═══ [Phase 2: ACTING] 行动阶段 ═══")
		fmt.Println("   选择并执行工具调用...")
		log.Println("[Acting] 生成行动响应...")

		actionResp, err := e.provider.Generate(ctx, contextHistory, availableTools)
		if err != nil {
			return fmt.Errorf("Action 阶段生成失败: %w", err)
		}

		contextHistory = append(contextHistory, *actionResp)

		if actionResp.Content != "" {
			if len(actionResp.ToolCalls) == 0 {
				extractedCalls := extractJSONToolCalls(actionResp.Content, availableTools)
				if len(extractedCalls) > 0 {
					log.Printf("[Acting] 🔍 从 Content 中检测到 %d 个隐式工具调用\n", len(extractedCalls))
					actionResp.ToolCalls = extractedCalls
				} else {
					fmt.Printf("\n💬 [模型回复]: %s\n", actionResp.Content)
				}
			}
		}

		if len(actionResp.ToolCalls) == 0 {
			fmt.Println("\n✅ [任务完成] 模型未请求更多工具调用")
			log.Println("[Engine] 模型未请求调用工具，任务宣告完成。")
			break
		}

		// =====================================================
		// Phase 3: OBSERVATION - 观察阶段
		// =====================================================
		fmt.Printf("\n👁️ ═══ [Phase 3: OBSERVATION] 观察阶段 ═══\n")
		fmt.Printf("   执行 %d 个工具调用并收集结果...\n", len(actionResp.ToolCalls))
		log.Printf("[Observation] 执行 %d 个工具调用...\n", len(actionResp.ToolCalls))

		for i, toolCall := range actionResp.ToolCalls {
			fmt.Printf("\n   🛠️  工具 #%d: %s\n", i+1, toolCall.Name)
			fmt.Printf("   📥 参数: %s\n", string(toolCall.Arguments))
			log.Printf("[Observation] 执行工具: %s, 参数: %s\n", toolCall.Name, string(toolCall.Arguments))

			result := e.registry.Execute(ctx, toolCall)

			if result.IsError {
				fmt.Printf("   ❌ 执行失败: %s\n", result.Output)
				log.Printf("[Observation] 工具执行报错: %s\n", result.Output)
			} else {
				outputPreview := result.Output
				if len(outputPreview) > 100 {
					outputPreview = outputPreview[:100] + "..."
				}
				fmt.Printf("   ✅ 执行成功: %s\n", outputPreview)
				log.Printf("[Observation] 工具执行成功 (返回 %d 字节)\n", len(result.Output))
			}

			observationMsg := schema.Message{
				Role:       schema.RoleUser,
				Content:    result.Output,
				ToolCallID: toolCall.ID,
			}
			contextHistory = append(contextHistory, observationMsg)
		}

		fmt.Println("\n   📝 观察结果已添加到上下文，准备下一轮循环...")
	}

	if turnCount >= e.MaxTurns {
		fmt.Println("\n⚠️ 达到最大回合数限制，任务终止")
		log.Println("[Engine] ⚠️ 达到最大回合数限制，任务终止")
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║ 🎉 任务执行完毕                                               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
