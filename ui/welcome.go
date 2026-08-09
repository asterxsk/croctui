package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type welcomeModel struct {
	selected int // 0 = Send, 1 = Receive, 2 = History
	pulse    int // animation frame: 0 or 1
	width    int
	height   int
}

// welcomePulseMsg advances the selection pulse animation (marker flips ▸/▶).
type welcomePulseMsg time.Time

// welcomePulseTick fires every 300ms to flip the selection chevron between its
// two animation frames.
func welcomePulseTick() tea.Msg {
	time.Sleep(300 * time.Millisecond)
	return welcomePulseMsg{}
}

func (w welcomeModel) Init() tea.Cmd { return welcomePulseTick }

func updateWelcome(w welcomeModel, msg tea.Msg) (welcomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case welcomePulseMsg:
		w.pulse = (w.pulse + 1) % 2
		return w, welcomePulseTick
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "left", "k", "h", "tab":
			w.selected = (w.selected + 2) % 3
		case "down", "right", "j", "l", "shift+tab":
			w.selected = (w.selected + 1) % 3
		case "enter", "s", "r", "t":
			if msg.String() == "s" {
				w.selected = 0
			} else if msg.String() == "r" {
				w.selected = 1
			} else if msg.String() == "t" {
				w.selected = 2
			}
			switch w.selected {
			case 0:
				return w, func() tea.Msg { return goScreenMsg{s: screenSend} }
			case 1:
				return w, func() tea.Msg { return goScreenMsg{s: screenReceive} }
			case 2:
				return w, func() tea.Msg { return goScreenMsg{s: screenHistory} }
			}
		case "q", "ctrl+c":
			return w, tea.Quit
		}
	}
	return w, nil
}

func (w welcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := updateWelcome(w, msg)
	return m, cmd
}

// welcomeButton renders the action label at a fixed single-line height.
// Inactive items are dim text; the active item shows a pulsing ▸/▶ marker in
// accent. Both states are the same height so the row never grows vertically.
func welcomeButton(label string, active bool, pulse int) string {
	if !active {
		return lipgloss.NewStyle().
			Foreground(colDim).
			Width(14).
			Height(1).
			MaxHeight(1).
			Align(lipgloss.Center).
			Render(label)
	}
	marker := "▸"
	if pulse == 1 {
		marker = "▶"
	}
	return lipgloss.NewStyle().
		Foreground(accent).
		Bold(true).
		Width(14).
		Height(1).
		MaxHeight(1).
		Align(lipgloss.Center).
		Render(marker + " " + label)
}

func (w welcomeModel) View() string {
	actions := []string{"Send", "Receive", "History"}
	var buttons []string
	for i, a := range actions {
		buttons = append(buttons, welcomeButton(a, i == w.selected, w.pulse))
	}
	// The row is left-anchored and expands to fill the content width, so the
	// buttons grow from the left like expanding bars.
	row := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)

	// The footer is always the same one-line keyboard hint at a fixed size, so
	// the centered layout never shifts between renders.
	foot := footerPad("s send · r receive · t history · arrows move · enter confirm · q quit")

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		"",
		row,
		separatorLine(),
		"",
		foot,
	)
	return screenView(w.width, w.height, "securely transfer files peer-to-peer over the internet", content)
}
