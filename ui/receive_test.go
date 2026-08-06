package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestReceiveCodeTooShort(t *testing.T) {
	m := newReceiveModel("croc")
	m.code.SetValue("abc")
	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("expected nil cmd for short code")
	}
	recv, ok := got.(receiveModel)
	if !ok {
		t.Fatal("expected receiveModel from Update")
	}
	if recv.status == "" {
		t.Fatal("expected a status error message")
	}
}

func TestReceiveErrorPrecedence(t *testing.T) {
	m := newReceiveModel("croc")
	// Simulate an error state with both err and status set
	m.state = recvError
	m.err = &testError{"croc subprocess failed"}
	m.status = "connection timeout"

	// View should show both with status first
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
	// Should contain both messages
	if !contains(view, "connection timeout") {
		t.Errorf("expected status 'connection timeout' in view, got %q", view)
	}
	if !contains(view, "croc subprocess failed") {
		t.Errorf("expected err 'croc subprocess failed' in view, got %q", view)
	}
}

func TestReceiveEnterEmitsStartReceive(t *testing.T) {
	m := newReceiveModel("croc")
	m.code.SetValue("1234-5678-9012")
	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	// The command is a Batch of spinner.Tick + startReceiveMsg
	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected BatchMsg, got %#v", cmd())
	}
	found := false
	for _, c := range batch {
		if _, ok := c().(startReceiveMsg); ok {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected startReceiveMsg in batch, got %#v", cmd())
	}
	recv, ok := got.(receiveModel)
	if !ok {
		t.Fatal("expected receiveModel from Update")
	}
	if recv.state != recvWaiting {
		t.Fatalf("expected recvWaiting state, got %v", recv.state)
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}