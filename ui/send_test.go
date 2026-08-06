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

// TestSendEnterTransitionsToWaiting pins the fix: Enter in sendConfirm must
// transition to sendWaiting (so the code box is reachable) rather than stay
// in sendConfirm — the root launches croc on startSendMsg, and the model must
// be in sendWaiting to render the code when EventCode arrives.
func TestSendEnterTransitionsToWaiting(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	sm := nm.(sendModel)
	if sm.state != sendWaiting {
		t.Fatalf("state = %v, want sendWaiting", sm.state)
	}
}

func TestSendAMarksAdding(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	sm := nm.(sendModel)
	if !sm.adding {
		t.Fatalf("expected adding=true after 'a', got %+v", sm)
	}
}

func TestSendAddMoreAppends(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	// press 'a' to enter adding mode
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	sm := nm.(sendModel)
	if !sm.adding {
		t.Fatal("expected adding=true")
	}
	// dialog returns more files → append
	nm2, _ := sm.Update(sendFilesPickedMsg{files: []string{"b.txt", "c.txt"}})
	sm2 := nm2.(sendModel)
	if len(sm2.files) != 3 {
		t.Fatalf("expected 3 files after append, got %v", sm2.files)
	}
	if sm2.adding {
		t.Fatal("expected adding=false after append")
	}
	if sm2.state != sendConfirm {
		t.Fatalf("expected sendConfirm after append, got %v", sm2.state)
	}
}

func TestSendAddMoreDedupes(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	sm := nm.(sendModel)
	nm2, _ := sm.Update(sendFilesPickedMsg{files: []string{"a.txt", "b.txt"}})
	sm2 := nm2.(sendModel)
	if len(sm2.files) != 2 {
		t.Fatalf("expected 2 files (a.txt deduped), got %v", sm2.files)
	}
}

func TestSendAddMoreCancelKeepsList(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	sm := nm.(sendModel)
	// cancel the add-more dialog → back to confirm, list intact
	nm2, _ := sm.Update(sendFilesPickedMsg{files: nil, cancelled: true})
	sm2 := nm2.(sendModel)
	if sm2.state != sendConfirm {
		t.Fatalf("expected sendConfirm after cancel, got %v", sm2.state)
	}
	if len(sm2.files) != 1 {
		t.Fatalf("expected original file kept, got %v", sm2.files)
	}
}
