package main

import "fmt"

type AIForAll struct {
	Project string
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
