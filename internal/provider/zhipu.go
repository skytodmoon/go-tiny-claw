package provider

type ZhipuProvider struct {
	apiKey string
}

func NewZhipuProvider(apiKey string) *ZhipuProvider {
	return &ZhipuProvider{apiKey: apiKey}
}

func (z *ZhipuProvider) Chat(prompt string) (string, error) {
	return "", nil
}

func (z *ZhipuProvider) Name() string {
	return "zhipu"
}