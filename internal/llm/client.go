package llm

import (
	"context"

	"github.com/monochromegane/aiforall/internal/llm/openai"
	"github.com/monochromegane/aiforall/internal/payload"
)

type LLMClient interface {
	ChatCompletion(*payload.Request, context.Context) (*payload.Response, error)
	ChatCompletionStream(*payload.Request, context.Context, func(*payload.Response) error) error
}

func GetLLMClient(model string) LLMClient {
	// First, only OpenAI is supported
	return openai.NewClient()
}
