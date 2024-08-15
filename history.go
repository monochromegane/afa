package main

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

func (h *History) LastMessage() *HistoryMessage {
	return h.Messages[len(h.Messages)-1]
}

func (m *HistoryMessage) IsAsUser() bool {
	return m.Role == "user"
}
