package main

type Secret struct {
	OpenAI *OpenAISecret `json:"openai"`
}

type OpenAISecret struct {
	ApiKey string `json:"api_key"`
}

func NewSecret(openai_api_key string) *Secret {
	return &Secret{
		OpenAI: &OpenAISecret{
			ApiKey: openai_api_key,
		},
	}
}
