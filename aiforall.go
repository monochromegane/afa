package main

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

type AIForAll struct {
	ConfigDir string
	CacheDir  string
	WorkSpace *WorkSpace
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
	return ai.startSession()
}

func (ai *AIForAll) Source() error {
	fmt.Println("Run as source mode.")
	return ai.startSession()
}

func (ai *AIForAll) Resume() error {
	fmt.Println("Run as resume mode.")
	return ai.startSession()
}

func (ai *AIForAll) startSession() error {
	session := NewSession()
	return session.Start()
}
