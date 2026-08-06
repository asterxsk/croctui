package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSendConfirmsFiles(t *testing.T) {
	m := newSendModel("croc")
	// simulate dialog returning files
	nm, cmd := m.Update(sendFilesPickedMsg{files: []string{"a.txt", "b.txt"}})
	sm := nm.(sendModel)
	if sm.state != sendConfirm {
		t.Fatalf("state = %v, want sendConfirm", sm.state)
	}
	if len(sm.files) != 2 {
		t.Fatalf("files = %v", sm.files)
	}
	if cmd != nil {
		t.Fatalf("expected nil cmd, got %v", cmd)
	}
	_ = sm
}

func TestSendViewListsFiles(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt", "b.txt"}
	m.state = sendConfirm
	if !strings.Contains(m.View(), "a.txt") || !strings.Contains(m.View(), "b.txt") {
		t.Fatalf("send view missing files:\n%s", m.View())
	}
}

func TestSendEnterEmitsStartSend(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	if _, ok := cmd().(startSendMsg); !ok {
		t.Fatalf("expected startSendMsg, got %#v", cmd())
	}
}