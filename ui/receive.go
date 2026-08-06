package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"crocgui/croc"
)

type recvState int

const (
	recvInput recvState = iota
	recvWaiting
	recvTransferring
	recvDone
	recvError
)

// startReceiveMsg is sent from receiveModel to root to begin the receive.
type startReceiveMsg struct{ code string }

type receiveModel struct {
	bin      string
	state    recvState
	code     textinput.Model
	progress progress.Model
	spinner  spinner.Model
	pct      float64
	file     string
	rate     string
	status   string
	err      error
	outDir   string // cached downloads directory
	width    int
}

func newReceiveModel(bin string) receiveModel {
	ti := textinput.New()
	ti.Placeholder = "e.g. 1234-5678-9012"
	ti.Focus()
	ti.CharLimit = 64
	return receiveModel{
		bin:      bin,
		state:    recvInput,
		code:     ti,
		progress: progress.New(progress.WithDefaultGradient()),
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot)),
		pct:      -1,
	}
}

func (m receiveModel) Init() tea.Cmd {
	if m.state == recvInput {
		return textinput.Blink
	}
	return m.spinner.Tick
}

func (m receiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case recvInput:
		if k, ok := msg.(tea.KeyMsg); ok {
			if k.String() == "enter" {
				code := m.code.Value()
				if len(code) < 6 {
					m.status = "code must be at least 6 characters"
					return m, nil
				}
				m.state = recvWaiting
				return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
					return startReceiveMsg{code: code}
				})
			}
			if k.String() == "esc" {
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
		}
		var cmd tea.Cmd
		m.code, cmd = m.code.Update(msg)
		return m, cmd

	case recvWaiting, recvTransferring:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if cmd != nil {
			return m, cmd
		}
		switch msg := msg.(type) {
		case crocEventMsg:
			ev := msg.ev
			switch ev.Kind {
			case croc.EventFile:
				m.file = ev.File
			case croc.EventProgress:
				m.pct = ev.Percent
				m.rate = ev.Rate
				m.state = recvTransferring
			case croc.EventError:
				m.status = ev.Text
				m.state = recvError
			default:
				if ev.Text != "" {
					m.status = ev.Text
				}
			}
			return m, nil
		case crocFinishedMsg:
			if msg.err != nil {
				m.err = msg.err
				m.state = recvError
			} else {
				m.state = recvDone
			}
			return m, nil
		case tea.KeyMsg:
			if k := msg.String(); k == "q" || k == "ctrl+c" {
				return m, tea.Quit
			}
		}
		return m, nil

	case recvDone:
		if k, ok := msg.(tea.KeyMsg); ok {
			if k.String() == "esc" || k.String() == "q" || k.String() == "enter" {
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
		}
		return m, nil

	case recvError:
		if k, ok := msg.(tea.KeyMsg); ok {
			if k.String() == "esc" || k.String() == "q" {
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
		}
		return m, nil
	}
	return m, nil
}

func (m receiveModel) View() string {
	var content string
	switch m.state {
	case recvInput:
		content = panel.Render(
			"Enter the code to receive:\n\n" +
				m.code.View() + "\n\n" +
				dimStyle.Render("enter start · esc back"),
		)

	case recvWaiting:
		content = m.spinner.View() + " " + m.status + "\n\n" + dimStyle.Render("Waiting for sender…")

	case recvTransferring:
		bar := ""
		if m.pct >= 0 {
			bar = m.progress.ViewAs(m.pct / 100)
		}
		rate := ""
		if m.rate != "" {
			rate = " @ " + m.rate
		}
		content = m.spinner.View() + " Receiving " + m.file + rate + "\n" + bar

	case recvDone:
		dir := m.outDir
		if dir == "" {
			var err error
			dir, err = downloadsDir()
			if err != nil {
				dir = "Downloads"
			}
		}
		content = lipgloss.NewStyle().Foreground(colDone).Render("✓ Done! Saved to "+dir) +
			"\n\n" + dimStyle.Render("Esc to go back")

	case recvError:
		errMsg := "Error"
		// Prefer specific error from croc, fall back to status
		if m.err != nil {
			errMsg = m.err.Error()
		}
		if m.status != "" && m.status != errMsg {
			errMsg = m.status + ": " + errMsg
		}
		content = lipgloss.NewStyle().Foreground(colErr).Render("✗ "+errMsg) +
			"\n\n" + dimStyle.Render("Esc to go back")
	}
	if m.width <= 0 {
		return content + "\n"
	}
	return lipgloss.Place(m.width, 10, lipgloss.Center, lipgloss.Center, content) + "\n"
}
