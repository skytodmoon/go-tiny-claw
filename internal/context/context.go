package context

import "strings"

type ContextManager struct {
	messages []string
	maxTokens int
}

func NewContextManager(maxTokens int) *ContextManager {
	return &ContextManager{
		messages:  []string{},
		maxTokens: maxTokens,
	}
}

func (c *ContextManager) AddMessage(msg string) {
	c.messages = append(c.messages, msg)
}

func (c *ContextManager) ComposePrompt() string {
	return strings.Join(c.messages, "\n")
}

func (c *ContextManager) TokenCount() int {
	return len(c.ComposePrompt())
}

func (c *ContextManager) Compact() {
	for c.TokenCount() > c.maxTokens && len(c.messages) > 0 {
		c.messages = c.messages[1:]
	}
}