package payload

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
}

type Response struct {
	Message *Message `json:"message"`
}
