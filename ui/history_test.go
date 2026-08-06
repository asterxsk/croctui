package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHistoryEmptyState(t *testing.T) {
	h := newHistoryModel()
	v := h.View()
	if !strings.Contains(v, "No transfers yet") {
		t.Fatalf("history empty state missing: %q", v)
	}
}

func TestHistoryRecordsEntry(t *testing.T) {
	h := newHistoryModel()
	h.add(HistoryEntry{Dir: histSend, Files: []string{"a.txt"}, Code: "1234", Time: "12:00", OK: true})
	if len(h.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.entries))
	}
	v := h.View()
	if !strings.Contains(v, "a.txt") || !strings.Contains(v, "SEND") {
		t.Fatalf("history view missing entry details: %q", v)
	}
}

func TestHistoryShowsErrorStatus(t *testing.T) {
	h := newHistoryModel()
	h.add(HistoryEntry{Dir: histRecv, Files: []string{"b.txt"}, Code: "5678", Time: "12:05", OK: false, Err: "timeout"})
	v := h.View()
	if !strings.Contains(v, "RECV") {
		t.Fatalf("history view missing RECV direction: %q", v)
	}
}

func TestHistoryEscGoesBack(t *testing.T) {
	h := newHistoryModel()
	_, cmd := h.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command")
	}
	gs, ok := cmd().(goScreenMsg)
	if !ok || gs.s != screenWelcome {
		t.Fatalf("expected goScreenMsg(screenWelcome), got %#v", cmd())
	}
}

func TestHistoryRecordsMultipleAndNavigates(t *testing.T) {
	h := newHistoryModel()
	h.add(HistoryEntry{Dir: histSend, Files: []string{"a"}, Code: "1", Time: "12:00", OK: true})
	h.add(HistoryEntry{Dir: histRecv, Files: []string{"b"}, Code: "2", Time: "12:01", OK: true})
	if h.sel != 1 {
		t.Fatalf("expected sel=1 (newest), got %d", h.sel)
	}
	nm, _ := h.Update(tea.KeyMsg{Type: tea.KeyUp})
	if nm.(historyModel).sel != 0 {
		t.Fatalf("expected sel=0 after up, got %d", nm.(historyModel).sel)
	}
}
