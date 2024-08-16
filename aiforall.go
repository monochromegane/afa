package main

import (
	"context"
	"fmt"
	"io"
	"os"
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
	Interactive          bool
}

func NewAIForAll(configDir, cacheDir string) *AIForAll {
	return &AIForAll{
		WorkSpace: NewWorkSpace(configDir, cacheDir),
		Input:     os.Stdin,
		Output:    os.Stdout,
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
	sessionPath := ai.WorkSpace.SessionPathFromTime(time.Now())
	if err := ai.WorkSpace.SetupSession(sessionPath, ai.Model); err != nil {
		return err
	}
	return ai.startSession(sessionPath)
}

func (ai *AIForAll) Source() error {
	fmt.Println("Run as source mode.")
	sessionPath := ai.WorkSpace.SessionPathFromName(ai.SessionName)
	return ai.startSession(sessionPath)
}

func (ai *AIForAll) Resume() error {
	fmt.Println("Run as resume mode.")
	sessionPath := ai.WorkSpace.SessionPathFromName(ai.SessionName)
	return ai.startSession(sessionPath)
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
	)
	err = session.Start(ai.Message, ai.MessageStdin, context.Background(), ai.Input, ai.Output)
	if err != nil {
		return err
	}

	return ai.WorkSpace.SaveSession(sessionPath, session.History)
}
