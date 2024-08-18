package main

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/monochromegane/afa/internal/llm"
	"github.com/monochromegane/afa/internal/payload"
)

type Session struct {
	Config                   *Config
	History                  *History
	SystemPromptTemplatePath string
	UserPromptTemplatePath   string
	Interactive              bool
	Stream                   bool
	Client                   llm.LLMClient
}

func NewSession(config *Config, history *History, systemPromptTemplatePath, userPromptTemplatePath string, interactive, stream bool) *Session {
	client := llm.GetLLMClient(history.Model)
	return &Session{
		Config:                   config,
		History:                  history,
		SystemPromptTemplatePath: systemPromptTemplatePath,
		UserPromptTemplatePath:   userPromptTemplatePath,
		Interactive:              interactive,
		Stream:                   stream,
		Client:                   client,
	}
}

func (s *Session) Start(message, messageStdin string, files []string, ctx context.Context, r io.Reader, w io.Writer) error {
	if s.History.IsNewSession() {
		systemPrompt, err := NewPrompt(s.SystemPromptTemplatePath, "", message, messageStdin, files)
		if err != nil {
			return err
		}
		s.History.AddMessage("system", systemPrompt)
	}

	runWithInput := false
	if message != "" || messageStdin != "" || len(files) > 0 {
		userPrompt, err := NewPrompt(s.UserPromptTemplatePath, "", message, messageStdin, files)
		if err != nil {
			return err
		}

		err = s.chatCompletionAndPrint(ctx, userPrompt, w)
		if err != nil {
			return err
		}
		runWithInput = true
	}

	if runWithInput && !s.Interactive {
		return nil
	}
	if runWithInput {
		fmt.Fprintln(w)
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

		err := s.chatCompletionAndPrint(ctx, userPrompt, w)
		if err != nil {
			return err
		}

		if !s.Interactive {
			break
		}
		fmt.Fprint(w, "\n> ")
	}

	return nil
}

func (s *Session) chatCompletionAndPrint(ctx context.Context, userPrompt string, w io.Writer) error {
	ctx = context.WithValue(ctx, "openai-api-key", s.Config.OpenAIAPIKey)

	s.History.AddMessage("user", userPrompt)

	role := ""
	message := ""
	if s.Stream {
		err := s.Client.ChatCompletionStream(s.History.Request, ctx, func(response *payload.Response) error {
			if r := response.Message.Role; r != "" {
				role = r
			}
			chunk := response.Message.Content
			message += chunk
			fmt.Fprint(w, chunk)
			return nil
		})
		if err != nil {
			return err
		}
		fmt.Fprintln(w)
	} else {
		response, err := s.Client.ChatCompletion(s.History.Request, ctx)
		if err != nil {
			return err
		}
		role = response.Message.Role
		message = response.Message.Content
		fmt.Fprintln(w, message)
	}
	s.History.AddMessage(role, message)

	return nil
}
