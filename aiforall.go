package main

import "fmt"

type AIForAll struct {
	Project string
}

func (ai *AIForAll) New() error {
	fmt.Println("Run as new mode.")
	return ai.start()
}

func (ai *AIForAll) Source() error {
	fmt.Println("Run as source mode.")
	return ai.start()
}

func (ai *AIForAll) Resume() error {
	fmt.Println("Run as resume mode.")
	return ai.start()
}

func (ai *AIForAll) start() error {
	return nil
}
