package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	welcomeStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("32"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	selectedStyle = lipgloss.NewStyle().Reverse(true)
)

type welcomeModel struct {
	selected int // 0 = Send, 1 = Receive
}

func (w welcomeModel) Init() tea.Cmd { return nil }

func updateWelcome(w welcomeModel, msg tea.Msg) (welcomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "left", "k", "h", "tab":
			w.selected = 1 - w.selected
		case "down", "right", "j", "l", "shift+tab":
			w.selected = 1 - w.selected
		case "enter", "s", "r":
			if msg.String() == "r" {
				w.selected = 1
			} else if msg.String() == "s" {
				w.selected = 0
			}
			if w.selected == 0 {
				return w, func() tea.Msg { return goScreenMsg{s: screenSend} }
			}
			return w, func() tea.Msg { return goScreenMsg{s: screenReceive} }
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
	send := "  Send  "
	recv := "  Receive  "
	if w.selected == 0 {
		send = selectedStyle.Render(send)
	} else {
		recv = selectedStyle.Render(recv)
	}
	line := lipgloss.JoinHorizontal(lipgloss.Center, "  ", send, recv, "  ")
	return lipgloss.JoinVertical(lipgloss.Center,
		welcomeStyle.Render("CROC"),
		dimStyle.Render("easily & securely transfer files\npeer-to-peer over the internet"),
		"",
		line,
		"",
		dimStyle.Render("s / arrows select · enter confirm · q / ctrl-c quit"),
	) + "\n"
}
