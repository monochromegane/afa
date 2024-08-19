package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

type AIForAll struct {
	WorkSpace *WorkSpace
	Input     io.Reader
	Output    io.Writer
	Option    *Option

	SessionName  string
	Message      string
	MessageStdin string
	Files        []string
}

func NewAIForAll(configDir, cacheDir string) (*AIForAll, error) {
	workSpace := NewWorkSpace(configDir, cacheDir)
	option, err := workSpace.LoadOption()
	if err != nil {
		return nil, err
	}
	if option.Chat.RunsOn == "" {
		option.Chat.RunsOn = strconv.Itoa(os.Getppid())
	}
	return &AIForAll{
		WorkSpace: workSpace,
		Input:     os.Stdin,
		Output:    os.Stdout,
		Option:    option,
	}, nil
}

func (ai *AIForAll) Init() error {
	fmt.Print("Enter your OpenAI API key: ")
	apiKey, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("Failed to read OpenAI API key: %v", err)
	}
	return ai.WorkSpace.Setup(NewOption(), NewSecret(string(apiKey)))
}

func (ai *AIForAll) New() error {
	ai.SessionName = ai.sessionNameFromTime(time.Now())
	sessionPath := ai.WorkSpace.SessionPath(ai.SessionName)
	if err := ai.WorkSpace.SetupSession(sessionPath, ai.Option.Chat.Model, ai.Option.Chat.Schema); err != nil {
		return err
	}
	return ai.startSession(sessionPath)
}

func (ai *AIForAll) Source() error {
	sessionPath := ai.WorkSpace.SessionPath(ai.SessionName)
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return fmt.Errorf("%s: no such session log", sessionPath)
	}
	return ai.startSession(sessionPath)
}

func (ai *AIForAll) Resume() error {
	sidPath := ai.WorkSpace.SidPath(ai.Option.Chat.RunsOn)
	if _, err := os.Stat(sidPath); os.IsNotExist(err) {
		return fmt.Errorf("%s: no such sid", sidPath)
	}

	data, err := os.ReadFile(sidPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	ai.SessionName = lines[0]
	return ai.Source()
}

func (ai *AIForAll) List() error {
	names, histories, err := ai.WorkSpace.ListSessions(ai.Option.List.Count, ai.Option.List.OrderByModify)
	if err != nil {
		return err
	}
	for i, name := range names {
		fmt.Fprintf(ai.Output, "%s\t%s\n", name, strings.Split(histories[i].FirstPrompt(), "\n")[0])
	}
	return nil
}

func (ai *AIForAll) startSession(sessionPath string) error {
	history, err := ai.WorkSpace.LoadHistory(sessionPath)
	if err != nil {
		return err
	}
	secret, err := ai.WorkSpace.LoadSecret()
	if err != nil {
		return err
	}
	session := NewSession(
		secret,
		history,
		ai.WorkSpace.TemplatePath("system", ai.Option.Chat.SystemPromptTemplate),
		ai.WorkSpace.TemplatePath("user", ai.Option.Chat.UserPromptTemplate),
		ai.Option.Chat.Interactive,
		ai.Option.Chat.Stream,
	)
	err = session.Start(ai.Message, ai.MessageStdin, ai.Files, context.Background(), ai.Input, ai.Output)
	if err != nil {
		return err
	}

	return ai.WorkSpace.SaveSession(ai.SessionName, ai.Option.Chat.RunsOn, session.History)
}

func (ai *AIForAll) sessionNameFromTime(startedAt time.Time) string {
	return startedAt.Format("2006-01-02_15-04-05")
}
