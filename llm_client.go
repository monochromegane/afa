package main

import "github.com/monochromegane/aiforall/internal/llm"

type LLMClient interface {
	ChatCompletion() error
}

func getLLMClient(model string) LLMClient {
	// First, only OpenAI is supported
	return &llm.OpenAIClient{}
}
