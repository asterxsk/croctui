package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type welcomeModel struct {
	selected int // 0 = Send, 1 = Receive, 2 = History
	width    int
}

func (w welcomeModel) Init() tea.Cmd { return nil }

func updateWelcome(w welcomeModel, msg tea.Msg) (welcomeModel, tea.Cmd) {
	switch msg := msg.(type) {
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

func (w welcomeModel) View() string {
	actions := []string{"Send", "Receive", "History"}
	var buttons []string
	for i, a := range actions {
		if i == w.selected {
			buttons = append(buttons, selectedButton.Render(a))
		} else {
			buttons = append(buttons, idleButton.Render(a))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)

	word := wordmark.Render("CROC")
	sub := dimStyle.Render("easily & securely transfer files  ·  peer-to-peer over the internet")
	foot := footerHint.Render("s send · r receive · t history · arrows move · enter confirm · q quit")

	content := lipgloss.JoinVertical(lipgloss.Center,
		word,
		sub,
		"",
		row,
		"",
		foot,
	)
	if w.width <= 0 {
		return content + "\n"
	}
	return lipgloss.Place(w.width, 12, lipgloss.Center, lipgloss.Center, content) + "\n"
}
