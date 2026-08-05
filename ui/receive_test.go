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