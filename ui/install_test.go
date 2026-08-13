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

	// 'x' must not emit goScreenMsg
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if cmd != nil {
		if msg := cmd(); msg != nil {
			if _, ok := msg.(goScreenMsg); ok {
				t.Fatal("unexpected goScreenMsg on non-Enter key in installDone")
			}
		}
	}

	// Enter: locally emits goScreenMsg, on CI (no croc) transitions to installNotDetected with no Cmd
	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		if _, ok := cmd().(goScreenMsg); !ok {
			t.Fatalf("expected goScreenMsg on Enter, got %#v", cmd())
		}
	} else {
		im, ok := nm.(installModel)
		if !ok {
			t.Fatalf("want installModel, got %T", nm)
		}
		if im.state != installNotDetected {
			t.Fatalf("Enter in installDone with no croc: want installNotDetected, got %v", im.state)
		}
	}
}
