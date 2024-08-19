package main

import (
	"encoding/json"

	"github.com/monochromegane/afa/internal/payload"
)

type History struct {
	*payload.Request
}

type HistoryMessage struct {
	*payload.Message
}

func NewHistory(model, schema string, rawSchema *json.RawMessage) *History {
	request := &payload.Request{
		Model:    model,
		Messages: []*payload.Message{},
	}
	if schema != "" {
		request.JsonSchema = &payload.JsonSchema{
			Name:   schema,
			Schema: rawSchema,
		}
	}
	return &History{request}
}

func (h *History) IsNewSession() bool {
	return len(h.Messages) == 0
}

func (h *History) AddMessage(role, content string) {
	h.Messages = append(h.Messages, &payload.Message{Role: role, Content: content})
}

func (h *History) FirstPrompt() string {
	for _, message := range h.Messages {
		if message.Role == "user" {
			return message.Content
		}
	}
	return ""
}
