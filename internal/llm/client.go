package llm

import (
	"context"

	"github.com/monochromegane/afa/internal/llm/openai"
	"github.com/monochromegane/afa/internal/payload"
)

type LLMClient interface {
	ChatCompletion(*payload.Request, context.Context) (*payload.Response, error)
	ChatCompletionStream(*payload.Request, context.Context, func(*payload.Response) error) error
}

func GetLLMClient(model string) LLMClient {
	// First, only OpenAI is supported
	return openai.NewClient()
}
