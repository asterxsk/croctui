package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWelcomeViewShowsCroc(t *testing.T) {
	m := welcomeModel{}
	if !strings.Contains(m.View(), "CROC") {
		t.Fatalf("welcome view missing CROC: %q", m.View())
	}
}

func TestWelcomeViewShowsActions(t *testing.T) {
	m := welcomeModel{}
	v := m.View()
	for _, want := range []string{"Send", "Receive"} {
		if !strings.Contains(v, want) {
			t.Fatalf("welcome view missing %q: %q", want, v)
		}
	}
}

func TestWelcomeEnterSendsGoSend(t *testing.T) {
	m := welcomeModel{selected: 0} // 0 = Send
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command, got nil")
	}
	msg := cmd()
	gs, ok := msg.(goScreenMsg)
	if !ok || gs.s != screenSend {
		t.Fatalf("expected goScreenMsg(screenSend), got %#v", msg)
	}
}

func TestWelcomeEnterSendsGoReceive(t *testing.T) {
	m := welcomeModel{selected: 1} // 1 = Receive
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	msg := cmd()
	gs, ok := msg.(goScreenMsg)
	if !ok || gs.s != screenReceive {
		t.Fatalf("expected goScreenMsg(screenReceive), got %#v", msg)
	}
}
