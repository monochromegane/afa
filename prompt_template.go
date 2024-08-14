package main

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"
)

func NewPrompt(promptTemplatePath, ctxString, message, messageStdin string) (string, error) {
	if _, err := os.Stat(promptTemplatePath); os.IsNotExist(err) {
		return "", err
	}
	promptTemplate, err := os.ReadFile(promptTemplatePath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("prompt").Parse(string(promptTemplate))
	if err != nil {
		return "", err
	}

	promptContext, err := newPromptContext(ctxString, message, messageStdin)
	if err != nil {
		return "", err
	}

	var prompt bytes.Buffer
	err = tmpl.Execute(&prompt, promptContext)
	if err != nil {
		return "", err
	}
	return prompt.String(), nil
}

func newPromptContext(ctxString, message, messageStdin string) (map[string]string, error) {
	if ctxString == "" {
		ctxString = "{}"
	}
	var ctx map[string]string
	if err := json.Unmarshal([]byte(ctxString), &ctx); err != nil {
		return nil, err
	}
	ctx["Message"] = message
	ctx["MessageStdin"] = messageStdin
	return ctx, nil
}
