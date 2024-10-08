package main

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/monochromegane/afa/internal/llm"
	"github.com/monochromegane/afa/internal/payload"
)

type MessageReader interface {
	io.Reader
}

type MessageWriter interface {
	io.Writer
	Disconnect() error
	Prompt() error
	Error() error
}

type DefaultMessageWriter struct {
	io.Writer
}

func (w *DefaultMessageWriter) Disconnect() error {
	return nil
}

func (w *DefaultMessageWriter) Prompt() error {
	fmt.Fprint(w, "> ")
	return nil
}

func (w *DefaultMessageWriter) Error() error {
	return nil
}

type Session struct {
	Secret                   *Secret
	History                  *History
	SystemPromptTemplatePath string
	UserPromptTemplatePath   string
	Interactive              bool
	Stream                   bool
	WithHistory              bool
	DryRun                   bool
	MockRun                  bool
	Verb                     string
	Client                   llm.LLMClient
}

func NewSession(secret *Secret, history *History, systemPromptTemplatePath, userPromptTemplatePath string, interactive, stream, withHistory, dryRun, mockRun, quote bool) *Session {
	client := llm.GetLLMClient(history.Model)
	verb := "%s"
	if quote {
		verb = "%q"
	}
	return &Session{
		Secret:                   secret,
		History:                  history,
		SystemPromptTemplatePath: systemPromptTemplatePath,
		UserPromptTemplatePath:   userPromptTemplatePath,
		Interactive:              interactive,
		Stream:                   stream,
		WithHistory:              withHistory,
		DryRun:                   dryRun,
		MockRun:                  mockRun,
		Verb:                     verb,
		Client:                   client,
	}
}

func (s *Session) Start(message, messageStdin string, files []string, ctx context.Context, r MessageReader, w MessageWriter) error {
	if s.History.IsNewSession() {
		systemPrompt, err := NewPrompt(s.SystemPromptTemplatePath, "", message, messageStdin, []string{})
		if err != nil {
			return err
		}
		s.History.AddMessage("system", systemPrompt)
	} else if s.WithHistory {
		fmt.Fprint(w, s.History.View(false))
	}

	runWithInput := false
	if message != "" || messageStdin != "" || len(files) > 0 {
		userPrompt, err := NewPrompt(s.UserPromptTemplatePath, "", message, messageStdin, files)
		if err != nil {
			return err
		}
		if s.DryRun {
			s.History.AddMessage("user", userPrompt)
			fmt.Fprint(w, s.History.View(true))
			// To prevent saving logs in dry-run mode, the last message is removed.
			s.History.RemoveLastMessage()
			return nil
		}

		err = s.chatCompletionAndPrint(ctx, userPrompt, w)
		if err != nil {
			return err
		}
		runWithInput = true
	}

	if runWithInput || s.DryRun {
		return nil
	}

	w.Prompt()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		userPrompt := scanner.Text()
		if userPrompt == "" {
			w.Prompt()
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
		w.Prompt()
	}

	return nil
}

func (s *Session) chatCompletionAndPrint(ctx context.Context, userPrompt string, w io.Writer) error {
	if s.MockRun {
		fmt.Fprintf(w, s.Verb, s.History.LastAssistantMessage())
		fmt.Fprintln(w)
		return nil
	}

	ctx = context.WithValue(ctx, "openai-api-key", s.Secret.OpenAI.ApiKey)

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
			fmt.Fprintf(w, s.Verb, chunk)
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
		fmt.Fprintf(w, s.Verb, message)
		fmt.Fprintln(w)
	}
	s.History.AddMessage(role, message)

	return nil
}
