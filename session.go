package main

import (
	"bufio"
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
	Interactive              bool
	Client                   llm.LLMClient
}

func NewSession(config *Config, history *History, systemPromptTemplatePath, userPromptTemplatePath string, interactive bool) *Session {
	client := llm.GetLLMClient(history.Model)
	return &Session{
		Config:                   config,
		History:                  history,
		SystemPromptTemplatePath: systemPromptTemplatePath,
		UserPromptTemplatePath:   userPromptTemplatePath,
		Interactive:              interactive,
		Client:                   client,
	}
}

func (s *Session) Start(message, messageStdin string, ctx context.Context, r io.Reader, w io.Writer) error {
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

		lines := strings.Split(userPrompt, "\n")
		for _, line := range lines {
			fmt.Fprintln(w, fmt.Sprintf("> %s", line))
		}

		message, err := s.chatCompletion(ctx, userPrompt)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, message)
		fmt.Fprintln(w)
	}

	if !s.Interactive {
		return nil
	}

	fmt.Fprint(w, "> ")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		userPrompt := scanner.Text()
		if userPrompt == "" {
			fmt.Fprint(w, "> ")
			continue
		}
		if userPrompt == "exit" {
			break
		}

		message, err := s.chatCompletion(ctx, userPrompt)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, message)
		fmt.Fprint(w, "\n> ")
	}

	return nil
}

func (s *Session) chatCompletion(ctx context.Context, userPrompt string) (string, error) {
	ctx = context.WithValue(ctx, "openai-api-key", s.Config.OpenAIAPIKey)

	s.History.AddMessage("user", userPrompt)
	response, err := s.Client.ChatCompletion(s.History.Request, ctx)
	if err != nil {
		return "", err
	}
	s.History.AddMessage(response.Message.Role, response.Message.Content)

	return response.Message.Content, err
}
