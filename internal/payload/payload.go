package payload

import "encoding/json"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model      string      `json:"model"`
	Messages   []*Message  `json:"messages"`
	JsonSchema *JsonSchema `json:"json_schema,omitempty"`
}

type JsonSchema struct {
	Name   string           `json:"name"`
	Schema *json.RawMessage `json:"schema"`
}

type Response struct {
	Message *Message `json:"message"`
}
