package main

type Option struct {
	Chat *ChatOption `json:"chat"`
	List *ListOption `json:"list"`
}

type ChatOption struct {
	Model                string `json:"model"`
	SystemPromptTemplate string `json:"system_prompt_template"`
	UserPromptTemplate   string `json:"user_prompt_template"`
	Schema               string `json:"schema"`
	RunsOn               string `json:"runs_on"`
	Interactive          bool   `json:"interactive"`
	Stream               bool   `json:"stream"`
}

type ListOption struct {
	Count         int  `json:"count"`
	OrderByModify bool `json:"order_by_modify"`
}

func NewOption(runsOn string) *Option {
	return &Option{
		Chat: &ChatOption{
			Model:                "gpt-4o-mini",
			SystemPromptTemplate: "default",
			UserPromptTemplate:   "default",
			Schema:               "",
			RunsOn:               runsOn,
			Interactive:          false,
			Stream:               false,
		},
		List: &ListOption{
			Count:         10,
			OrderByModify: false,
		},
	}
}
