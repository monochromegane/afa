package main

import (
	"fmt"
	"syscall"
	"time"

	"golang.org/x/term"
)

type AIForAll struct {
	ConfigDir string
	CacheDir  string
	WorkSpace *WorkSpace

	SystemPromptTemplate string
	UserPromptTemplate   string
	Model                string
	SessionName          string
	Message              string
	MessageStdin         string
}

func NewAIForAll(configDir, cacheDir string) *AIForAll {
	return &AIForAll{
		WorkSpace: NewWorkSpace(configDir, cacheDir),
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
	fmt.Println("Run as new mode.")
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
	session := NewSession(
		sessionPath,
		ai.WorkSpace.TemplatePath("system", ai.SystemPromptTemplate),
		ai.WorkSpace.TemplatePath("user", ai.UserPromptTemplate),
	)
	return session.Start(ai.Message, ai.MessageStdin)
}
