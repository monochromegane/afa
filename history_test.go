package main

import "testing"

func TestIsNewSession(t *testing.T) {
	hist := NewHistory("", "", nil)
	if !hist.IsNewSession() {
		t.Errorf("IsNewSession should return false for a new History instance")
	}

	hist.AddMessage("", "")
	if hist.IsNewSession() {
		t.Errorf("IsNewSession should return true after a message is added")
	}
}
