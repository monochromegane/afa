package main

type Option struct {
	Model                string `json:"model"`
	SystemPromptTemplate string `json:"system_prompt_template"`
	UserPromptTemplate   string `json:"user_prompt_template"`
	Schema               string `json:"schema"`
	RunsOn               string `json:"runs_on"`
	Interactive          bool   `json:"interactive"`
	Stream               bool   `json:"stream"`
}

func NewOption(runsOn string) *Option {
	return &Option{
		Model:                "gpt-4o-mini",
		SystemPromptTemplate: "default",
		UserPromptTemplate:   "default",
		Schema:               "",
		RunsOn:               runsOn,
		Interactive:          false,
		Stream:               false,
	}
}
