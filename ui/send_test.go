package ui

import (
	"reflect"
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

// TestSendEnterAlwaysOpensDialog pins the new interact behavior: Enter in
// sendConfirm ALWAYS opens the add dialog, even when files are present — it
// never starts the send.
func TestSendEnterAlwaysOpensDialog(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	if reflect.ValueOf(cmd).Pointer() != reflect.ValueOf(pickFilesCmd).Pointer() {
		t.Fatalf("expected pickFilesCmd on Enter, got %#v", cmd)
	}
}

// TestSendSpaceSends pins that space in sendConfirm (with files) transitions
// to sendWaiting and emits startSendMsg — the keyboard equivalent of the Send
// button.
func TestSendSpaceSends(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	sm := nm.(sendModel)
	if sm.state != sendWaiting {
		t.Fatalf("state = %v, want sendWaiting after space", sm.state)
	}
	if cmd == nil {
		t.Fatal("expected a command")
	}
	if _, ok := cmd().(startSendMsg); !ok {
		t.Fatalf("expected startSendMsg after space, got %#v", cmd())
	}
}

// TestSendSpaceEmptyNoOp pins that space with no files is a no-op (the Send
// button is disabled until files exist).
func TestSendWaitingEscGoesBack(t *testing.T) {
	m := newSendModel("croc")
	m.state = sendWaiting

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command on Esc while waiting")
	}
	msg := cmd()
	gs, ok := msg.(goScreenMsg)
	if !ok || gs.s != screenWelcome {
		t.Fatalf("expected goScreenMsg(screenWelcome), got %#v", msg)
	}
}

func TestSendSpaceEmptyNoOp(t *testing.T) {
	m := newSendModel("croc")
	m.state = sendConfirm
	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if cmd != nil {
		t.Fatalf("space with no files should be a no-op, got cmd=%#v", cmd)
	}
	if nm.(sendModel).state != sendConfirm {
		t.Fatalf("state = %v, want sendConfirm", nm.(sendModel).state)
	}
}

// TestSendAOpensDialog pins that 'a' in sendConfirm returns pickFilesCmd (the
// command that opens the native file dialog) rather than a modal state.
func TestSendAOpensDialog(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	sm := nm.(sendModel)
	if cmd == nil {
		t.Fatal("expected a command after 'a'")
	}
	// It must be the same pick-files command (compared by func pointer) — we
	// never execute it in a test, or it would open the real native dialog.
	if reflect.ValueOf(cmd).Pointer() != reflect.ValueOf(pickFilesCmd).Pointer() {
		t.Fatalf("expected pickFilesCmd after 'a', got %#v", cmd)
	}
	if sm.state != sendConfirm {
		t.Fatalf("state = %v, want sendConfirm after 'a'", sm.state)
	}
}

// TestSendShiftAOpensFolderDialog pins that Shift+a ('A') in sendConfirm
// returns pickFolderCmd (the folder picker command).
func TestSendShiftAOpensFolderDialog(t *testing.T) {
	m := newSendModel("croc")
	m.state = sendConfirm
	nm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A")})
	if cmd == nil {
		t.Fatal("expected a command after Shift+a")
	}
	if reflect.ValueOf(cmd).Pointer() != reflect.ValueOf(pickFolderCmd).Pointer() {
		t.Fatalf("expected pickFolderCmd after Shift+a, got %#v", cmd)
	}
	if nm.(sendModel).state != sendConfirm {
		t.Fatalf("state = %v, want sendConfirm", nm.(sendModel).state)
	}
}

func TestSendAddMoreAppends(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	// dialog returns more files → append
	nm2, _ := m.Update(sendFilesPickedMsg{files: []string{"b.txt", "c.txt"}})
	sm2 := nm2.(sendModel)
	if len(sm2.files) != 3 {
		t.Fatalf("expected 3 files after append, got %v", sm2.files)
	}
	if sm2.state != sendConfirm {
		t.Fatalf("expected sendConfirm after append, got %v", sm2.state)
	}
}

func TestSendAddMoreDedupes(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	nm2, _ := m.Update(sendFilesPickedMsg{files: []string{"a.txt", "b.txt"}})
	sm2 := nm2.(sendModel)
	if len(sm2.files) != 2 {
		t.Fatalf("expected 2 files (a.txt deduped), got %v", sm2.files)
	}
}

func TestSendAddMoreCancelKeepsList(t *testing.T) {
	m := newSendModel("croc")
	m.files = []string{"a.txt"}
	m.state = sendConfirm
	// cancel the add dialog → stay on the list, files intact
	nm2, _ := m.Update(sendFilesPickedMsg{files: nil, cancelled: true})
	sm2 := nm2.(sendModel)
	if sm2.state != sendConfirm {
		t.Fatalf("expected sendConfirm after cancel, got %v", sm2.state)
	}
	if len(sm2.files) != 1 {
		t.Fatalf("expected original file kept, got %v", sm2.files)
	}
}

// TestSendEmptyEnterOpensDialog pins that Enter on an empty list opens the file
// dialog (pickFilesCmd) instead of sending.
func TestSendEmptyEnterOpensDialog(t *testing.T) {
	m := newSendModel("croc") // starts empty, state sendConfirm
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected a command on Enter with empty list")
	}
	// Must be the pick-files command (func-pointer compare, never execute).
	if reflect.ValueOf(cmd).Pointer() != reflect.ValueOf(pickFilesCmd).Pointer() {
		t.Fatalf("expected pickFilesCmd on empty Enter, got %#v", cmd)
	}
	if m.state != sendConfirm {
		t.Fatalf("state = %v, want sendConfirm", m.state)
	}
}

// TestSendInitDoesNotOpenDialog pins that entering the send screen does NOT
// auto-open the file dialog — it opens an empty list instead.
func TestSendInitDoesNotOpenDialog(t *testing.T) {
	m := newSendModel("croc")
	cmd := m.Init()
	if cmd != nil {
		t.Fatal("expected nil Init cmd (no auto file dialog)")
	}
	if m.state != sendConfirm {
		t.Fatalf("state = %v, want sendConfirm", m.state)
	}
	if len(m.files) != 0 {
		t.Fatalf("expected empty file list, got %v", m.files)
	}
}

// TestSendViewEmptyShowsAddHint pins the empty-list view: it shows the
// "Press Enter to add" hint rather than nothing.
func TestSendViewEmptyShowsAddHint(t *testing.T) {
	m := newSendModel("croc")
	v := m.View()
	if !strings.Contains(v, "Press Enter to add") {
		t.Fatalf("empty send view missing add hint:\n%s", v)
	}
}

// TestSendDirShowsFolderNotSize pins the "0 bytes" bug: a directory in the
// send list must show a folder marker and NOT a misleading size (dirs stat to
// ~0 bytes on Windows).
func TestSendDirShowsFolderNotSize(t *testing.T) {
	dir := t.TempDir()
	m := newSendModel("croc")
	m.state = sendConfirm
	m.files = []string{dir}
	v := stripANSI(m.View())
	if !strings.Contains(v, "📁") {
		t.Fatalf("send view missing folder marker for dir:\n%s", v)
	}
	if strings.Contains(v, "0 B") || strings.Contains(v, "(0") {
		t.Fatalf("send view shows misleading dir size:\n%s", v)
	}
}

// TestSendNavigateAndDelete pins arrow/w-s/j-k navigation and x deletion.
func TestSendNavigateAndDelete(t *testing.T) {
	m := newSendModel("croc")
	m.state = sendConfirm
	m.files = []string{"a.txt", "b.txt", "c.txt"}
	m.sel = 0

	// Move down with j (sel 0 -> 1).
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if nm.(sendModel).sel != 1 {
		t.Fatalf("sel after j = %d, want 1", nm.(sendModel).sel)
	}
	// Move down with 's' (1 -> 2).
	nm, _ = nm.(sendModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	if nm.(sendModel).sel != 2 {
		t.Fatalf("sel after s = %d, want 2", nm.(sendModel).sel)
	}
	// Move up with 'w' (2 -> 1).
	nm, _ = nm.(sendModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("w")})
	if nm.(sendModel).sel != 1 {
		t.Fatalf("sel after w = %d, want 1", nm.(sendModel).sel)
	}
	// Arrow down is also supported.
	nm, _ = nm.(sendModel).Update(tea.KeyMsg{Type: tea.KeyDown})
	if nm.(sendModel).sel != 2 {
		t.Fatalf("sel after down = %d, want 2", nm.(sendModel).sel)
	}
	// Arrow up.
	nm, _ = nm.(sendModel).Update(tea.KeyMsg{Type: tea.KeyUp})
	if nm.(sendModel).sel != 1 {
		t.Fatalf("sel after up = %d, want 1", nm.(sendModel).sel)
	}

	// Delete selected (index 1 = "b.txt") with x.
	nm, _ = nm.(sendModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	sm := nm.(sendModel)
	if len(sm.files) != 2 || sm.files[0] != "a.txt" || sm.files[1] != "c.txt" {
		t.Fatalf("files after delete = %v, want [a.txt c.txt]", sm.files)
	}
	// Cursor stays valid (sel now points at "c.txt").
	if sm.sel >= len(sm.files) {
		t.Fatalf("sel %d out of range after delete", sm.sel)
	}
}

// TestSendDeleteKeepsCursorInRange pins that deleting the last row clamps the
// cursor so it never points past the end.
func TestSendDeleteKeepsCursorInRange(t *testing.T) {
	m := newSendModel("croc")
	m.state = sendConfirm
	m.files = []string{"a.txt", "b.txt"}
	m.sel = 1
	nm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	sm := nm.(sendModel)
	if len(sm.files) != 1 {
		t.Fatalf("files = %v, want [a.txt]", sm.files)
	}
	if sm.sel != 0 {
		t.Fatalf("sel = %d, want 0 (clamped)", sm.sel)
	}
}
