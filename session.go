package main

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) Start() error {
	client := getLLMClient()
	return client.ChatCompletion()
}
