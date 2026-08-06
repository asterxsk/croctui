package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"crocgui/croc"
)

// TestSendCodeVisibleDuringTransfer pins the fix for Important 3:
// when EventProgress arrives first (transitioning to sendTransferring)
// and then EventCode arrives, the code box must still be visible in
// the sendTransferring view.
func TestSendCodeVisibleDuringTransfer(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"test.txt"}
	m.state = sendTransferring

	// First: progress arrives (simulates hashing phase)
	nm, _ := m.Update(crocEventMsg{ev: croc.Event{Kind: croc.EventProgress, Percent: 50, Rate: "1 MB/s"}})
	m = nm.(sendModel)
	if m.state != sendTransferring {
		t.Fatalf("state should be sendTransferring, got %v", m.state)
	}

	// Then: code arrives (real croc order)
	nm, _ = m.Update(crocEventMsg{ev: croc.Event{Kind: croc.EventCode, Code: "fakecode123"}})
	m = nm.(sendModel)

	// Code should be in the model
	if m.code != "fakecode123" {
		t.Fatalf("code not set: %q", m.code)
	}

	// View should contain the code box
	view := m.View()
	if !strings.Contains(view, "fakecode123") {
		t.Fatalf("code box not visible in sendTransferring view:\n%s", view)
	}
	if !strings.Contains(view, "Code:") {
		t.Fatalf("code label not visible in sendTransferring view:\n%s", view)
	}
}

// TestSendDoneEnterStartsNewSend pins the fix for Important 7:
// Enter in sendDone should start a fresh send flow (go to screenSend).
func TestSendDoneEnterStartsNewSend(t *testing.T) {
	m := newSendModel("croc")
	m.state = sendDone

	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command on Enter in sendDone")
	}
	msg := cmd()
	if _, ok := msg.(goScreenMsg); !ok {
		t.Fatalf("expected goScreenMsg, got %T: %v", msg, msg)
	}
	gsm := msg.(goScreenMsg)
	if gsm.s != screenSend {
		t.Fatalf("goScreenMsg screen = %v, want screenSend", gsm.s)
	}
	_ = nm
}

// TestSendConfirmEscGoesBack pins the fix for Important 7:
// Esc in sendConfirm should go back to welcome (not remove last file).
func TestSendConfirmEscGoesBack(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt", "b.txt"}
	m.state = sendConfirm

	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command on Esc in sendConfirm")
	}
	msg := cmd()
	if _, ok := msg.(goScreenMsg); !ok {
		t.Fatalf("expected goScreenMsg, got %T: %v", msg, msg)
	}
	gsm := msg.(goScreenMsg)
	if gsm.s != screenWelcome {
		t.Fatalf("goScreenMsg screen = %v, want screenWelcome", gsm.s)
	}
	_ = nm
}

// TestSendConfirmXRemovesLastFile pins the fix for Important 7:
// x/backspace in sendConfirm should remove the last file.
func TestSendConfirmXRemovesLastFile(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt", "b.txt"}
	m.state = sendConfirm

	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	sm := nm.(sendModel)
	if len(sm.files) != 1 || sm.files[0] != "a.txt" {
		t.Fatalf("files = %v, want [a.txt]", sm.files)
	}
}