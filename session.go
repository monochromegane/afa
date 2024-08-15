package main

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type Session struct {
	Config                   *Config
	History                  *History
	SystemPromptTemplatePath string
	UserPromptTemplatePath   string
}

func NewSession(config *Config, history *History, systemPromptTemplatePath, userPromptTemplatePath string) *Session {
	return &Session{
		Config:                   config,
		History:                  history,
		SystemPromptTemplatePath: systemPromptTemplatePath,
		UserPromptTemplatePath:   userPromptTemplatePath,
	}
}

func (s *Session) Start(message, messageStdin string, ctx context.Context, r io.Reader, w io.Writer) error {
	client := getLLMClient(s.Config, s.History.Model)

	if s.History.IsNewSession() {
		systemPrompt, err := NewPrompt(s.SystemPromptTemplatePath, "", message, messageStdin)
		if err != nil {
			return err
		}
		s.History.AddMessage("system", systemPrompt)
	}

	if message != "" || messageStdin != "" {
		userPrompt, err := NewPrompt(s.UserPromptTemplatePath, "", message, messageStdin)
		if err != nil {
			return err
		}
		s.History.AddMessage("user", userPrompt)
	}

	if lastMessage := s.History.LastMessage(); lastMessage.IsAsUser() {
		lines := strings.Split(lastMessage.Content, "\n")
		for _, line := range lines {
			fmt.Fprintln(w, fmt.Sprintf("> %s", line))
		}
	}

	return client.ChatCompletion()
}
