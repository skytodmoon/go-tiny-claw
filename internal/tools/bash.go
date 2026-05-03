package tools

import "os/exec"

type BashTool struct{}

func NewBashTool() *BashTool {
	return &BashTool{}
}

func (b *BashTool) Name() string {
	return "bash"
}

func (b *BashTool) Description() string {
	return "Execute a bash command"
}

func (b *BashTool) Execute(args map[string]interface{}) (string, error) {
	cmd, ok := args["command"].(string)
	if !ok {
		return "", nil
	}
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}