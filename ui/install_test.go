package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInstallViewPromptsForWinget(t *testing.T) {
	m := newInstallModel()
	v := m.View()
	if !strings.Contains(v, "croc") || !strings.Contains(v, "winget") {
		t.Fatalf("install view missing croc/winget prompt:\n%s", v)
	}
}

func TestInstallEnterRunsInstall(t *testing.T) {
	m := newInstallModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	if _, ok := cmd().(installStartedMsg); !ok {
		t.Fatalf("expected installStartedMsg, got %#v", cmd())
	}
}

func TestInstallDoneOnlyAdvancesOnEnter(t *testing.T) {
	m := newInstallModel()
	m.state = installDone // simulate installDone state

	// Press 'x' (not Enter) - should NOT re-detect or emit goScreenMsg
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(goScreenMsg); ok {
			t.Fatal("unexpected goScreenMsg on non-Enter key in installDone")
		}
	}

	// Press Enter - should re-detect and emit goScreenMsg
	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command on Enter")
	}
	if _, ok := cmd().(goScreenMsg); !ok {
		t.Fatalf("expected goScreenMsg on Enter, got %#v", cmd())
	}
}
