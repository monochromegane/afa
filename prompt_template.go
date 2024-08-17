package main

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"
)

type PromptContext struct {
	Message      string
	MessageStdin string
	Files        []*PromptFile
	Context      map[string]string
}

type PromptFile struct {
	Name    string
	Content string
}

func NewPrompt(promptTemplatePath, ctxString, message, messageStdin string, files []string) (string, error) {
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

	promptContext, err := newPromptContext(ctxString, message, messageStdin, files)
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

func newPromptContext(ctxString, message, messageStdin string, files []string) (*PromptContext, error) {
	if ctxString == "" {
		ctxString = "{}"
	}
	var ctx map[string]string
	if err := json.Unmarshal([]byte(ctxString), &ctx); err != nil {
		return nil, err
	}

	var fileData []*PromptFile
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		fileData = append(fileData, &PromptFile{
			Name:    file,
			Content: string(content),
		})
	}

	return &PromptContext{
		Message:      message,
		MessageStdin: messageStdin,
		Files:        fileData,
		Context:      ctx,
	}, nil
}
