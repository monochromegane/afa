package main

type Option struct {
	Script *ScriptOption `json:"script"`
	Init   *InitOption   `json:"init"`
	Chat   *ChatOption   `json:"chat"`
	Viewer *ViewerOption `json:"viewer"`
	List   *ListOption   `json:"list"`
}

type ScriptOption struct {
	Enabled bool `json:"enabled"`
}

type InitOption struct {
	NoInteraction bool `json:"no_interaction"`
}

type ChatOption struct {
	Model                string `json:"model"`
	SystemPromptTemplate string `json:"system_prompt_template"`
	UserPromptTemplate   string `json:"user_prompt_template"`
	Schema               string `json:"schema"`
	RunsOn               string `json:"runs_on"`
	Interactive          bool   `json:"interactive"`
	Stream               bool   `json:"stream"`
	WithHistory          bool   `json:"with_history"`
	DryRun               bool   `json:"dry_run"`
	MockRun              bool   `json:"mock_run"`
	Quote                bool   `json:"quote"`
	Save                 bool   `json:"save"`
}

type ListOption struct {
	Count         int  `json:"count"`
	OrderByModify bool `json:"order_by_modify"`
}

type ViewerOption struct {
	Enabled bool     `json:"enabled"`
	Command []string `json:"command"`
}

func NewOption() *Option {
	return &Option{
		Init: &InitOption{
			NoInteraction: false,
		},
		Script: &ScriptOption{
			Enabled: false,
		},
		Chat: &ChatOption{
			Model:                "gpt-4o-mini",
			SystemPromptTemplate: "default",
			UserPromptTemplate:   "default",
			Schema:               "",
			RunsOn:               "",
			Interactive:          false,
			Stream:               false,
			WithHistory:          false,
			DryRun:               false,
			MockRun:              false,
			Quote:                false,
			Save:                 true,
		},
		Viewer: &ViewerOption{
			Enabled: false,
			Command: []string{"afa-tui", "-a"},
		},
		List: &ListOption{
			Count:         10,
			OrderByModify: false,
		},
	}
}

func (o *Option) SetScriptOptions() {
	o.Chat.Interactive = false
	o.Chat.WithHistory = false
	o.Chat.Stream = false
	o.Viewer.Enabled = false
}
