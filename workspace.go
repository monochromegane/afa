package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
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

func (w *WorkSpace) IsNotExist() bool {
	_, err := os.Stat(w.SecretPath())
	return os.IsNotExist(err)
}

func (w *WorkSpace) Setup(option *Option, secret *Secret) error {
	if err := w.setupDirs(); err != nil {
		return err
	}
	return w.setupFiles(option, secret)
}

func (w *WorkSpace) setupDirs() error {
	for _, dir := range []string{
		w.ConfigDir,
		w.TemplateDir("system"),
		w.TemplateDir("user"),
		w.SchemaDir(),
		w.CacheDir,
		w.SessionsDir(),
		w.SidDir(),
		w.SocketDir(),
	} {
		if err := w.mkDirAllIfNotExist(dir); err != nil {
			return err
		}
	}
	return nil
}

func (w *WorkSpace) setupFiles(option *Option, secret *Secret) error {
	if err := w.writeFileIfNotExist(
		w.TemplatePath("system", "default"),
		[]byte("You are a helpful assistant."),
	); err != nil {
		return err
	}

	if err := w.writeFileIfNotExist(
		w.TemplatePath("user", "default"),
		[]byte("{{ .Message }}\n{{ if .MessageStdin }}\n```\n{{ .MessageStdin }}```\n{{- end }}\n{{ range .Files }}\n- File: {{ .Name }}\\n```\n{{ .Content }}```\n{{ end -}}"),
	); err != nil {
		return err
	}

	if err := w.writeFileIfNotExist(
		w.SchemaPath("command_suggestion"),
		[]byte("{\n  \"type\": \"object\",\n  \"properties\": {\n    \"suggested_command\": {\n      \"type\": \"string\"\n    }\n  },\n  \"additionalProperties\": false,\n  \"required\": [\n    \"suggested_command\"\n  ]\n}"),
	); err != nil {
		return err
	}

	jsonOption, err := json.MarshalIndent(option, "", "  ")
	if err != nil {
		return err
	}
	if err := w.writeFileIfNotExist(w.OptionPath(), jsonOption); err != nil {
		return err
	}

	jsonSecret, err := json.MarshalIndent(secret, "", "  ")
	if err != nil {
		return err
	}
	if err := w.writeFileIfNotExist(w.SecretPath(), jsonSecret); err != nil {
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

func (w *WorkSpace) SchemaDir() string {
	return path.Join(w.ConfigDir, "schemas")
}

func (w *WorkSpace) SchemaPath(name string) string {
	return path.Join(w.SchemaDir(), filepath.Clean(fmt.Sprintf("%s.json", name)))
}

func (w *WorkSpace) SessionsDir() string {
	return path.Join(w.CacheDir, "sessions")
}

func (w *WorkSpace) SessionPath(name string) string {
	return path.Join(w.SessionsDir(), filepath.Clean(fmt.Sprintf("%s.json", name)))
}

func (w *WorkSpace) SidDir() string {
	return path.Join(w.CacheDir, "sid")
}

func (w *WorkSpace) SidPath(name string) string {
	return path.Join(w.SidDir(), filepath.Clean(fmt.Sprintf("%s.sid", name)))
}

func (w *WorkSpace) SocketDir() string {
	return path.Join(w.CacheDir, "sockets")
}

func (w *WorkSpace) SocketPath(name string) string {
	return path.Join(w.SocketDir(), filepath.Clean(fmt.Sprintf("%s.sock", name)))
}

func (w *WorkSpace) OptionPath() string {
	return path.Join(w.ConfigDir, "option.json")
}

func (w *WorkSpace) SecretPath() string {
	return path.Join(w.ConfigDir, "secret.json")
}

func (w *WorkSpace) SetupSession(sessionPath, model, schema string) error {
	rawSchema, err := w.LoadSchema(schema)
	if schema != "" && err != nil {
		return err
	}

	history := NewHistory(model, schema, rawSchema)
	jsonSession, err := json.Marshal(history)
	if err != nil {
		return err
	}
	return w.writeFileIfNotExist(sessionPath, jsonSession)
}

func (w *WorkSpace) LoadSchema(schema string) (*json.RawMessage, error) {
	path := w.SchemaPath(schema)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil, err
	}

	return &raw, nil
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

func (w *WorkSpace) RemoveSession(sessionName string) error {
	return os.Remove(w.SessionPath(sessionName))
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

func (w *WorkSpace) LoadOption() (*Option, error) {
	option := NewOption()

	path := w.OptionPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return option, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, option); err != nil {
		return nil, err
	}

	return option, nil
}

func (w *WorkSpace) LoadSecret() (*Secret, error) {
	path := w.SecretPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var secret Secret
	if err := json.Unmarshal(file, &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}

func (w *WorkSpace) ListSessions(count int, orderByModify bool) ([]string, []*History, error) {
	names := []string{}
	histories := []*History{}

	dirEntories, err := os.ReadDir(w.SessionsDir())
	if err != nil {
		return nil, nil, err
	}
	files := []fs.FileInfo{}
	for _, dirEntry := range dirEntories {
		if dirEntry.IsDir() {
			continue
		}
		info, err := dirEntry.Info()
		if err != nil {
			return nil, nil, err
		}
		files = append(files, info)
	}

	if orderByModify {
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})
	} else {
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() > files[j].Name()
		})
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		sessionName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		history, err := w.LoadHistory(w.SessionPath(sessionName))
		if err != nil {
			return nil, nil, err
		}

		names = append(names, sessionName)
		histories = append(histories, history)

		if len(names) >= count {
			break
		}
	}

	return names, histories, nil
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
