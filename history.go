package main

import (
	"encoding/json"
	"os"
)

type History struct {
	Model    string            `json:"model"`
	Messages []*HistoryMessage `json:"messages"`
}

type HistoryMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (h *History) IsNewSession() bool {
	return len(h.Messages) == 0
}

func (h *History) AddMessage(role, content string) {
	h.Messages = append(h.Messages, &HistoryMessage{Role: role, Content: content})
}

func LoadHistory(path string) (*History, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var history History
	if err := json.Unmarshal(file, &history); err != nil {
		return nil, err
	}

	return &history, nil
}
