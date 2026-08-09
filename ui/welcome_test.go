package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWelcomeViewShowsTransferCaption(t *testing.T) {
	m := welcomeModel{}
	if !strings.Contains(m.View(), "securely transfer files peer-to-peer over the internet") {
		t.Fatalf("welcome view missing transfer caption: %q", m.View())
	}
}

func TestWelcomeViewShowsActions(t *testing.T) {
	m := welcomeModel{}
	v := m.View()
	for _, want := range []string{"Send", "Receive", "History"} {
		if !strings.Contains(v, want) {
			t.Fatalf("welcome view missing %q: %q", want, v)
		}
	}
}

func TestWelcomePulseTicks(t *testing.T) {
	m := welcomeModel{}
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() must return a pulse tick command")
	}
	if _, ok := cmd().(welcomePulseMsg); !ok {
		t.Fatalf("Init() returned %T, want welcomePulseMsg", cmd())
	}
	nm, cmd2 := m.Update(welcomePulseMsg{})
	if nm.(welcomeModel).pulse != 1 {
		t.Fatalf("pulse = %d, want 1 after first tick", nm.(welcomeModel).pulse)
	}
	if cmd2 == nil {
		t.Fatal("expected a follow-up pulse tick command")
	}
}

func TestWelcomeSelectedButtonAnimated(t *testing.T) {
	m := welcomeModel{selected: 0}
	v := m.View()
	if !strings.Contains(v, "▸ Send") && !strings.Contains(v, "▶ Send") {
		t.Fatalf("selected button missing animated chevron marker:\n%s", v)
	}
}

func TestWelcomeEnterSendsGoHistory(t *testing.T) {
	m := welcomeModel{selected: 2} // 2 = History
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command, got nil")
	}
	msg := cmd()
	gs, ok := msg.(goScreenMsg)
	if !ok || gs.s != screenHistory {
		t.Fatalf("expected goScreenMsg(screenHistory), got %#v", msg)
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
