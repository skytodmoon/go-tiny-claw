package provider

type ClaudeProvider struct {
	apiKey string
}

func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{apiKey: apiKey}
}

func (c *ClaudeProvider) Chat(prompt string) (string, error) {
	return "", nil
}

func (c *ClaudeProvider) Name() string {
	return "claude"
}