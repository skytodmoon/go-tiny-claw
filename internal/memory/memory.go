package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Memory struct {
	Plan  []string `json:"plan"`
	Todo  []string `json:"todo"`
	Tasks []Task   `json:"tasks"`
}

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Content string `json:"content"`
}

func NewMemory() *Memory {
	return &Memory{
		Plan:  []string{},
		Todo:  []string{},
		Tasks: []Task{},
	}
}

func (m *Memory) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (m *Memory) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, m)
}