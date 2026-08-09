package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRootStartSendWithoutProgramShowsError(t *testing.T) {
	r := New("croc", nil).(*RootModel)
	r.send.files = []string{"file.txt"}
	r.send.state = sendWaiting
	nm, cmd := r.Update(startSendMsg{})
	if cmd != nil {
		t.Fatal("expected no command without a program")
	}
	if nm.(*RootModel).send.state != sendError {
		t.Fatal("expected visible send error without a program")
	}
}

// TestScreenSwitchKeepsSize pins the top-left bug: switching to a screen
// (e.g. send) must carry the last-known window size into the fresh model, so
// the screen renders framed/centered instead of top-left until the next
// WindowSizeMsg.
func TestScreenSwitchKeepsSize(t *testing.T) {
	r := New("croc", nil).(*RootModel)
	r.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	r.Update(goScreenMsg{s: screenSend})
	if r.send.width != 80 || r.send.height != 24 {
		t.Fatalf("send model size = %dx%d, want 80x24", r.send.width, r.send.height)
	}
	v := stripANSI(r.View())
	if !strings.HasPrefix(strings.TrimSpace(v), "┌") {
		t.Fatalf("send view not framed/centered after switch:\n%s", v)
	}

	// Same for receive.
	r.Update(goScreenMsg{s: screenReceive})
	if r.recv.width != 80 || r.recv.height != 24 {
		t.Fatalf("recv model size = %dx%d, want 80x24", r.recv.width, r.recv.height)
	}
	v = stripANSI(r.View())
	if !strings.HasPrefix(strings.TrimSpace(v), "┌") {
		t.Fatalf("receive view not framed/centered after switch:\n%s", v)
	}
}
