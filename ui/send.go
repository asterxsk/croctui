package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"crocgui/croc"
	"crocgui/dialog"
)

type sendState int

const (
	sendIdle sendState = iota
	sendConfirm
	sendWaiting
	sendTransferring
	sendDone
	sendError
)

type sendModel struct {
	bin      string
	state    sendState
	files    []string
	code     string
	file     string
	progress progress.Model
	spinner  spinner.Model
	pct      float64
	rate     string
	status   string
	err      error
}

func newSendModel(bin string) sendModel {
	return sendModel{
		bin:      bin,
		state:    sendIdle,
		progress: progress.New(progress.WithDefaultGradient()),
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot)),
		pct:      -1,
	}
}

// sendFilesPickedMsg carries the dialog result into the model.
type sendFilesPickedMsg struct {
	files     []string
	cancelled bool
}

func (m sendModel) Init() tea.Cmd {
	return func() tea.Msg {
		files, err := dialog.PickFiles()
		if err != nil {
			return dialogFailedMsg{err: err}
		}
		return sendFilesPickedMsg{files: files, cancelled: len(files) == 0}
	}
}

func (m sendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sendFilesPickedMsg:
		if msg.cancelled || len(msg.files) == 0 {
			m.state = sendIdle
			m.status = ""
			return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
		}
		m.files = msg.files
		m.state = sendConfirm
		return m, nil

	// NOTE: startSendMsg is handled by the root (it launches the croc
	// subprocess). The send model transitions to sendWaiting itself in the
	// sendConfirm Enter handler, so no case is needed here.

	case crocEventMsg:
		ev := msg.ev
		switch ev.Kind {
		case croc.EventCode:
			m.code = ev.Code
			m.status = "waiting for receiver…"
			return m, nil
		case croc.EventFile:
			m.file = ev.File
			return m, nil
		case croc.EventProgress:
			m.pct = ev.Percent
			m.rate = ev.Rate
			m.state = sendTransferring
			return m, nil
		case croc.EventError:
			m.status = ev.Text
			m.state = sendError
			return m, nil
		default:
			if ev.Text != "" {
				m.status = ev.Text
			}
			return m, nil
		}

	case crocFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = sendError
		} else {
			m.state = sendDone
		}
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case sendConfirm:
			k := msg.String()
			if k == "enter" {
				// Transition to sendWaiting NOW so the model shows the code box
				// while the root launches croc. (The root handles startSendMsg
				// to spawn the subprocess; the model must already be in
				// sendWaiting to render the code when EventCode arrives.)
				m.state = sendWaiting
				m.status = "connecting to relay…"
				return m, func() tea.Msg { return startSendMsg{} }
			}
			if k == "esc" {
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
			if k == "x" || k == "backspace" {
				if len(m.files) > 0 {
					m.files = m.files[:len(m.files)-1]
				}
				if len(m.files) == 0 {
					return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
				}
			}
		case sendIdle:
			if k := msg.String(); k == "esc" || k == "q" || k == "ctrl+c" {
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
		case sendDone, sendError:
			if k := msg.String(); k == "enter" && m.state == sendDone {
				return m, func() tea.Msg { return goScreenMsg{s: screenSend} }
			}
			if k := msg.String(); k == "esc" || k == "q" || k == "ctrl+c" {
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
		}
	}

	var cmd tea.Cmd
	if m.state == sendIdle || m.state == sendWaiting || m.state == sendTransferring {
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, cmd
}

var (
	codeBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(1, 0)
)

func (m sendModel) View() string {
	switch m.state {
	case sendIdle:
		return m.spinner.View() + " Opening file dialog…"

	case sendConfirm:
		var b strings.Builder
		b.WriteString("Files to send:\n\n")
		for _, f := range m.files {
			info, err := os.Stat(f)
			size := ""
			if err == nil {
				size = " (" + formatSize(info.Size()) + ")"
			}
			b.WriteString("  " + f + size + "\n")
		}
		b.WriteString("\nEnter to send, Esc to go back, x/backspace to remove last file")
		return b.String()

	case sendWaiting:
		return m.spinner.View() + " " + m.status + "\n\n" +
			codeBoxStyle.Render("Code: "+m.code+"\n(code copied to clipboard by croc)") +
			"\n\nEsc to go back"

	case sendTransferring:
		bar := ""
		if m.pct >= 0 {
			bar = m.progress.ViewAs(m.pct / 100)
		}
		rate := ""
		if m.rate != "" {
			rate = " @ " + m.rate
		}
		file := m.file
		if file == "" && len(m.files) > 0 {
			file = m.files[0]
		}
		var b strings.Builder
		b.WriteString(m.spinner.View() + " Sending " + file + rate + "\n" + bar)
		if m.code != "" {
			b.WriteString("\n\n" + codeBoxStyle.Render("Code: "+m.code+"\n(code copied to clipboard by croc)"))
		}
		return b.String()

	case sendDone:
		return "✓ Transfer complete\n\nEnter to send another, Esc to go back"

	case sendError:
		errMsg := "Error"
		if m.err != nil {
			errMsg = m.err.Error()
		}
		if m.status != "" && m.status != errMsg {
			errMsg = m.status + ": " + errMsg
		}
		return "✗ " + errMsg + "\n\nEsc to go back"
	}
	return ""
}