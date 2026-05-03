package provider

type Provider interface {
	Chat(prompt string) (string, error)
	Name() string
}

type Message struct {
	Role    string
	Content string
}