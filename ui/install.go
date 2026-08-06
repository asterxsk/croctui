package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"crocgui/croc"
)

type installState int

const (
	installPrompt installState = iota
	installRunning
	installDone
	installFailed
	installNotDetected
)

type installModel struct {
	state   installState
	spinner spinner.Model
	err     error
	width   int
}

func newInstallModel() installModel {
	return installModel{
		state:   installPrompt,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}

type installStartedMsg struct{}
type installFinishedMsg struct{ err error }

func (m installModel) Init() tea.Cmd { return nil }

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case installPrompt:
		if k, ok := msg.(tea.KeyMsg); ok {
			switch k.String() {
			case "enter", "y":
				return m, func() tea.Msg { return installStartedMsg{} }
			case "q", "ctrl+c", "n":
				return m, tea.Quit
			}
		}
		return m, nil
	default:
		switch msg := msg.(type) {
		case installStartedMsg:
			m.state = installRunning
			return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
				return installFinishedMsg{err: croc.Install()}
			})
		case installFinishedMsg:
			if msg.err != nil {
				m.state = installFailed
				m.err = msg.err
			} else {
				m.state = installDone
			}
			return m, nil
		case tea.KeyMsg:
			if m.state == installDone {
				// re-detect and move to welcome ONLY on Enter
				if k := msg.String(); k == "enter" {
					if bin, err := croc.Detect(); err == nil {
						return m, func() tea.Msg { return goScreenMsg{s: screenWelcome, bin: bin} }
					}
					// croc still not detected — show retry state
					m.state = installNotDetected
					return m, nil
				}
			}
			if m.state == installNotDetected {
				if k := msg.String(); k == "enter" {
					// Retry install
					return m, func() tea.Msg { return installStartedMsg{} }
				}
			}
			if k := msg.String(); k == "q" || k == "ctrl+c" || k == "esc" {
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m installModel) View() string {
	var content string
	switch m.state {
	case installPrompt:
		content = panel.Render(
			"croc is not installed.\n\n" +
				"Install it via winget?\n\n" +
				"[Y]es   [N]o (quit)\n\n" +
				dimStyle.Render("This will run: winget install schollz.croc"),
		)
	case installRunning:
		content = m.spinner.View() + " Installing croc via winget…"
	case installDone:
		content = lipgloss.NewStyle().Foreground(colDone).Render("✓ croc installed successfully!") +
			"\n\n" + dimStyle.Render("Press Enter to continue.")
	case installNotDetected:
		content = "croc still not detected (PATH may need refresh).\n\n" +
			dimStyle.Render("Enter retry install · Esc quit")
	case installFailed:
		content = lipgloss.NewStyle().Foreground(colErr).Render("✗ Failed to install croc:\n"+m.err.Error()) +
			"\n\n" + dimStyle.Render("Press Esc to quit.")
	}
	if m.width <= 0 {
		return content + "\n"
	}
	return lipgloss.Place(m.width, 10, lipgloss.Center, lipgloss.Center, content) + "\n"
}
