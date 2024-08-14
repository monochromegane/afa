package main

type Session struct {
	SessionPath              string
	SystemPromptTemplatePath string
	UserPromptTemplatePath   string
}

func NewSession(sessionPath, systemPromptTemplatePath, userPromptTemplatePath string) *Session {
	return &Session{
		SessionPath:              sessionPath,
		SystemPromptTemplatePath: systemPromptTemplatePath,
		UserPromptTemplatePath:   userPromptTemplatePath,
	}
}

func (s *Session) Start(message, messageStdin string) error {
	history, err := LoadHistory(s.SessionPath)
	if err != nil {
		return err
	}

	client := getLLMClient(history.Model)

	if history.IsNewSession() {
		systemPrompt, err := NewPrompt(s.SystemPromptTemplatePath, "", message, messageStdin)
		if err != nil {
			return err
		}
		history.AddMessage("system", systemPrompt)
	}

	if message != "" || messageStdin != "" {
		userPrompt, err := NewPrompt(s.UserPromptTemplatePath, "", message, messageStdin)
		if err != nil {
			return err
		}
		history.AddMessage("user", userPrompt)
	}

	return client.ChatCompletion()
}
