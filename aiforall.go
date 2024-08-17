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

	SystemPromptTemplate string
	UserPromptTemplate   string
	Model                string
	SessionName          string
	Message              string
	MessageStdin         string
	RunsOn               string
	Interactive          bool
	Stream               bool
}

func NewAIForAll(configDir, cacheDir string) *AIForAll {
	return &AIForAll{
		WorkSpace: NewWorkSpace(configDir, cacheDir),
		Input:     os.Stdin,
		Output:    os.Stdout,
		RunsOn:    strconv.Itoa(os.Getppid()),
	}
}

func (ai *AIForAll) Init() error {
	fmt.Print("Enter your OpenAI API key: ")
	apiKey, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("Failed to read OpenAI API key: %v", err)
	}
	config := &Config{
		OpenAIAPIKey: string(apiKey),
	}
	return ai.WorkSpace.Setup(config)
}

func (ai *AIForAll) New() error {
	ai.SessionName = ai.sessionNameFromTime(time.Now())
	sessionPath := ai.WorkSpace.SessionPath(ai.SessionName)
	if err := ai.WorkSpace.SetupSession(sessionPath, ai.Model); err != nil {
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
	sidPath := ai.WorkSpace.SidPath(ai.RunsOn)
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

func (ai *AIForAll) startSession(sessionPath string) error {
	history, err := ai.WorkSpace.LoadHistory(sessionPath)
	if err != nil {
		return err
	}
	config, err := ai.WorkSpace.LoadConfig()
	if err != nil {
		return err
	}
	session := NewSession(
		config,
		history,
		ai.WorkSpace.TemplatePath("system", ai.SystemPromptTemplate),
		ai.WorkSpace.TemplatePath("user", ai.UserPromptTemplate),
		ai.Interactive,
		ai.Stream,
	)
	err = session.Start(ai.Message, ai.MessageStdin, context.Background(), ai.Input, ai.Output)
	if err != nil {
		return err
	}

	return ai.WorkSpace.SaveSession(ai.SessionName, ai.RunsOn, session.History)
}

func (ai *AIForAll) sessionNameFromTime(startedAt time.Time) string {
	return startedAt.Format("2006-01-02_15-04-05")
}
