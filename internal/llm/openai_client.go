package llm

type OpenAIClient struct {
	ApiKey string
}

func (c *OpenAIClient) ChatCompletion() error { return nil }
