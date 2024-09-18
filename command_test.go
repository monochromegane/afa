package main

import (
	"strings"
	"testing"
)

func TestGetXdgHomeDirNoEnv(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	_, err := getXdgHomeDir("XDG_CONFIG_HOME")
	if err == nil {
		t.Errorf("getXdgHomeDir should return error when env is not set")
	}
}

func TestGetXdgHomeDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "~/.config")
	xdgHome, err := getXdgHomeDir("XDG_CONFIG_HOME")
	if err != nil {
		t.Errorf("getXdgHomeDir should not return error when env is set")
	}

	if strings.HasPrefix(xdgHome, "~") {
		t.Errorf("getXdgHomeDir should replace tilde to user's home dir")
	}
}
