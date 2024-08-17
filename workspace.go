package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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
		[]byte("{{.Message}}{{- if .MessageStdin }}\n```\n{{.MessageStdin}}\n```{{- end }}{{- range .Files}}\n{{ .Name }}\n```\n{{ .Content }}\n```{{- end }}"),
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
	return path.Join(w.TemplateDir(role), filepath.Clean(fmt.Sprintf("%s.tmpl", name)))
}

func (w *WorkSpace) SessionsDir() string {
	return path.Join(w.CacheDir, "sessions")
}

func (w *WorkSpace) SidDir() string {
	return path.Join(w.CacheDir, "sid")
}

func (w *WorkSpace) SidPath(name string) string {
	return path.Join(w.SidDir(), filepath.Clean(fmt.Sprintf("%s.sid", name)))
}

func (w *WorkSpace) ConfigPath() string {
	return path.Join(w.ConfigDir, "config.json")
}

func (w *WorkSpace) SessionPath(name string) string {
	return path.Join(w.SessionsDir(), filepath.Clean(fmt.Sprintf("%s.json", name)))
}

func (w *WorkSpace) SetupSession(sessionPath, model string) error {
	history := NewHistory(model)
	jsonSession, err := json.Marshal(history)
	if err != nil {
		return err
	}
	return w.writeFileIfNotExist(sessionPath, jsonSession)
}

func (w *WorkSpace) SaveSession(sessionName, runsOn string, history *History) error {
	jsonSession, err := json.Marshal(history)
	if err != nil {
		return err
	}

	err = w.writeFile(w.SessionPath(sessionName), jsonSession)
	if err != nil {
		return err
	}

	return w.writeFile(w.SidPath(runsOn), []byte(sessionName))
}

func (w *WorkSpace) LoadHistory(path string) (*History, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var history History
	if err := json.Unmarshal(file, &history); err != nil {
		return nil, err
	}

	return &history, nil
}

func (w *WorkSpace) LoadConfig() (*Config, error) {
	path := w.ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
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
		return w.writeFile(path, content)
	}
	return nil
}

func (w *WorkSpace) writeFile(path string, content []byte) error {
	return os.WriteFile(path, content, w.FilePerm)
}
