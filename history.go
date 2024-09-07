package main

import (
	"bytes"
	"encoding/json"
	"fmt"

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

func (h *History) RemoveLastMessage() {
	h.Messages = h.Messages[:len(h.Messages)-1]
}

func (h *History) FirstUserPrompt() string {
	for _, message := range h.Messages {
		if message.Role == "user" {
			return message.Content
		}
	}
	return ""
}

func (h *History) View(detail bool) string {
	var buf bytes.Buffer
	for _, message := range h.Messages {
		switch message.Role {
		case "system":
			if detail {
				buf.WriteString(fmt.Sprintf("# System\n\n%s\n\n", message.Content))
			}
		case "assistant":
			buf.WriteString(fmt.Sprintf("# Assistant\n\n%s\n\n", message.Content))
		case "user":
			buf.WriteString(fmt.Sprintf("# You\n\n%s\n\n", message.Content))
		}
	}
	return buf.String()
}
