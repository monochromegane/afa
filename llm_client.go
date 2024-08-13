package main

import "github.com/monochromegane/aiforall/internal/llm"

type LLMClient interface {
	ChatCompletion() error
}

func getLLMClient() LLMClient {
	// First, only OpenAI is supported
	return &llm.OpenAIClient{}
}
