package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type WorkSpace struct {
	ConfigDir string
	CacheDir  string

	DirPerm  os.FileMode
	FilePerm os.FileMode
}

func NewWorkSpace(configDir, cacheDir string) *WorkSpace {
	return &WorkSpace{
		ConfigDir: configDir,
		CacheDir:  cacheDir,
		DirPerm:   os.FileMode(0700),
		FilePerm:  os.FileMode(0600),
	}
}

func (w *WorkSpace) Setup(config *Config) error {
	if err := w.setupDirs(); err != nil {
		return err
	}
	return w.setupFiles(config)
}

func (w *WorkSpace) setupDirs() error {
	for _, dir := range []string{
		w.ConfigDir,
		w.TemplateDir("system"),
		w.TemplateDir("user"),
		w.CacheDir,
		w.SessionsDir(),
		w.SidDir(),
	} {
		if err := w.mkDirAllIfNotExist(dir); err != nil {
			return err
		}
	}
	return nil
}

func (w *WorkSpace) setupFiles(config *Config) error {
	if err := w.writeFileIfNotExist(
		w.TemplatePath("system", "default"),
		[]byte("You are a helpful assistant."),
	); err != nil {
		return err
	}

	if err := w.writeFileIfNotExist(
		w.TemplatePath("user", "default"),
		[]byte("{{.Message}}{{- if .MessageStdin }}\n```\n{{.MessageStdin}}\n```{{- end }}"),
	); err != nil {
		return err
	}

	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err := w.writeFileIfNotExist(w.ConfigPath(), jsonConfig); err != nil {
		return err
	}

	return nil
}

func (w *WorkSpace) TemplateDir(role string) string {
	return path.Join(w.ConfigDir, "templates", role)
}

func (w *WorkSpace) TemplatePath(role, name string) string {
	return path.Join(w.TemplateDir(role), fmt.Sprintf("%s.tmpl", name))
}

func (w *WorkSpace) SessionsDir() string {
	return path.Join(w.CacheDir, "sessions")
}

func (w *WorkSpace) SidDir() string {
	return path.Join(w.CacheDir, "sid")
}

func (w *WorkSpace) ConfigPath() string {
	return path.Join(w.ConfigDir, "config.json")
}

func (w *WorkSpace) SessionPathFromTime(startedAt time.Time) string {
	return w.SessionPathFromName(startedAt.Format("2006-01-02_15-04-05"))
}

func (w *WorkSpace) SessionPathFromName(sessionName string) string {
	return path.Join(w.SessionsDir(), fmt.Sprintf("%s.json", sessionName))
}

func (w *WorkSpace) SetupSession(sessionPath, model string) error {
	history := History{Model: model, Messages: []*HistoryMessage{}}
	jsonSession, err := json.Marshal(history)
	if err != nil {
		return err
	}
	return w.writeFileIfNotExist(sessionPath, jsonSession)
}

func (w *WorkSpace) mkDirAllIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, w.DirPerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WorkSpace) writeFileIfNotExist(path string, content []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, content, w.FilePerm); err != nil {
			return err
		}
	}
	return nil
}
