package feishu

type FeishuBot struct {
	appID     string
	appSecret string
}

func NewFeishuBot(appID, appSecret string) *FeishuBot {
	return &FeishuBot{
		appID:     appID,
		appSecret: appSecret,
	}
}

func (f *FeishuBot) HandleCallback(payload []byte) (string, error) {
	return "", nil
}

func (f *FeishuBot) SendMessage(chatID, message string) error {
	return nil
}