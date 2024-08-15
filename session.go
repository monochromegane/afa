package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/monochromegane/aiforall/internal/llm"
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
	client := llm.GetLLMClient(s.History.Model)

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

	ctx = context.WithValue(ctx, "openai-api-key", s.Config.OpenAIAPIKey)
	_, err := client.ChatCompletion(s.History.Request, ctx)
	if err != nil {
		return err
	}
	return nil
}
