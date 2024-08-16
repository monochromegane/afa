package main

import "github.com/monochromegane/aiforall/internal/payload"

type History struct {
	*payload.Request
}

type HistoryMessage struct {
	*payload.Message
}

func NewHistory(model string) *History {
	return &History{
		&payload.Request{
			Model:    model,
			Messages: []*payload.Message{},
		},
	}
}

func (h *History) IsNewSession() bool {
	return len(h.Messages) == 0
}

func (h *History) AddMessage(role, content string) {
	h.Messages = append(h.Messages, &payload.Message{Role: role, Content: content})
}
