package ui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// joinLines joins rendered lines with newlines (used to compose list rows).
func joinLines(rows []string) string {
	return strings.Join(rows, "\n")
}

// itoa formats an int as a decimal string (small helper for counts).
func itoa(n int) string {
	return strconv.Itoa(n)
}

// dirName returns the display label for a history direction.
func (d histDir) String() string {
	if d == histRecv {
		return "RECV"
	}
	return "SEND"
}

// history.go — the in-session record of completed transfers.
//
// A HistoryEntry is one completed (or failed) transfer. Direction is send or
// receive. Time is a short wall-clock stamp for display. The history screen
// lists entries oldest-first with the most recent highlighted.

type histDir int

const (
	histSend histDir = iota
	histRecv
)

// HistoryEntry is a single recorded transfer.
type HistoryEntry struct {
	Dir   histDir
	Files []string
	Code  string
	Time  string
	OK    bool
	Err   string
}

type historyModel struct {
	entries []HistoryEntry
	width   int
	height  int
	sel     int // selected row index (for scrolling emphasis)
}

func newHistoryModel() historyModel {
	return historyModel{}
}

// add records a completed transfer. The newest entry is appended to the end.
func (h *historyModel) add(e HistoryEntry) {
	h.entries = append(h.entries, e)
	h.sel = len(h.entries) - 1
}

func (h historyModel) Init() tea.Cmd { return nil }

func (h historyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return h, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
		case "up", "k":
			if h.sel > 0 {
				h.sel--
			}
		case "down", "j":
			if h.sel < len(h.entries)-1 {
				h.sel++
			}
		}
	}
	return h, nil
}

func (h historyModel) View() string {
	footer := footerPad("up/down or j/k move · esc or q back")
	if len(h.entries) == 0 {
		empty := panel.Render("No transfers yet.\n\nTransfers you make will appear here.")
		content := lipgloss.JoinVertical(lipgloss.Center, empty, separatorLine(), footer)
		return screenView(h.width, h.height, "your transfer history", content)
	}
	var rows []string
	for i, e := range h.entries {
		mark := "  "
		if i == h.sel {
			mark = "› "
		}
		status := lipgloss.NewStyle().Foreground(colDone).Render("✓")
		if !e.OK {
			status = lipgloss.NewStyle().Foreground(colErr).Render("✗")
		}
		files := ""
		if len(e.Files) > 0 {
			files = e.Files[0]
			if len(e.Files) > 1 {
				files += " +" + itoa(len(e.Files)-1)
			}
		}
		rows = append(rows, mark+status+" "+e.Dir.String()+"  "+files+"  code="+e.Code+"  "+e.Time)
	}
	list := panel.Render(joinLines(rows))
	content := lipgloss.JoinVertical(lipgloss.Center, list, separatorLine(), footer)
	return screenView(h.width, h.height, "your transfer history", content)
}
