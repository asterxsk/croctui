package ui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"crocgui/croc"
	"crocgui/dialog"
)

type sendState int

const (
	sendConfirm sendState = iota
	sendWaiting
	sendTransferring
	sendDone
	sendError
)

type sendModel struct {
	bin      string
	state    sendState
	files    []string
	sel      int // cursor into the file list (sendConfirm)
	code     string
	file     string
	progress progress.Model
	spinner  spinner.Model
	pct      float64
	rate     string
	status   string
	err      error
	width    int
	height   int
}

func newSendModel(bin string) sendModel {
	return sendModel{
		bin:      bin,
		state:    sendConfirm,
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

// pickFilesCmd opens the native file dialog and reports the selection.
var pickFilesCmd = func() tea.Msg {
	files, err := dialog.PickFiles()
	if err != nil {
		return dialogFailedMsg{err: err}
	}
	return sendFilesPickedMsg{files: files, cancelled: len(files) == 0}
}

// pickFolderCmd opens the native folder dialog and reports the selection.
var pickFolderCmd = func() tea.Msg {
	files, err := dialog.PickFolder()
	if err != nil {
		return dialogFailedMsg{err: err}
	}
	return sendFilesPickedMsg{files: files, cancelled: len(files) == 0}
}

func (m sendModel) Init() tea.Cmd {
	// No file dialog on entry — the send screen opens an empty navigable list
	// and the user presses Enter to add files.
	return nil
}

// sendFiles starts the transfer: transition to sendWaiting NOW so the model
// shows the code box while the root launches croc (the root handles
// startSendMsg to spawn the subprocess; the model must already be in
// sendWaiting to render the code when EventCode arrives). No-op when there are
// no files.
func (m sendModel) sendFiles() (sendModel, tea.Cmd) {
	if len(m.files) == 0 {
		return m, nil
	}
	m.state = sendWaiting
	m.status = "connecting to relay…"
	return m, func() tea.Msg { return startSendMsg{} }
}

func (m sendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case sendFilesPickedMsg:
		if msg.cancelled || len(msg.files) == 0 {
			// Cancel just closes the dialog — stay on the list.
			return m, nil
		}
		// Append new files, dedupe exact paths.
		for _, f := range msg.files {
			dup := false
			for _, existing := range m.files {
				if f == existing {
					dup = true
					break
				}
			}
			if !dup {
				m.files = append(m.files, f)
			}
		}
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
			switch k {
			case "enter":
				// Enter is the "interact" key — it always opens the add dialog
				// (files when empty, or more files). Sending happens via the
				// spacebar.
				return m, pickFilesCmd
			case "a":
				// Add more files: re-open the dialog and append (dedupe).
				return m, pickFilesCmd
			case "A":
				// Shift+a: pick a folder instead of files.
				return m, pickFolderCmd
			case " ":
				// Space sends.
				return m.sendFiles()
			case "x", "backspace":
				// Delete the selected file (or the last one if none selected).
				i := m.sel
				if i >= len(m.files) {
					i = len(m.files) - 1
				}
				if i >= 0 {
					m.files = append(m.files[:i], m.files[i+1:]...)
					if m.sel >= len(m.files) && m.sel > 0 {
						m.sel--
					}
				}
			case "up", "k", "w":
				if m.sel > 0 {
					m.sel--
				}
			case "down", "j", "s":
				if m.sel < len(m.files)-1 {
					m.sel++
				}
			case "esc":
				return m, func() tea.Msg { return goScreenMsg{s: screenWelcome} }
			}
		case sendWaiting:
			if k := msg.String(); k == "esc" {
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
	if m.state == sendWaiting || m.state == sendTransferring {
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, cmd
}

var (
	codeBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		Margin(1, 0)
)

// sendList renders the navigable file-list screen: an empty list shows a
// centered "Press Enter to add" prompt; a populated list shows each file with
// the selected row highlighted. Keyboard-driven only — no on-screen buttons;
// the footer describes the controls (move, delete, add via Enter, send via
// space).
func (m sendModel) sendList() string {
	var b strings.Builder
	if len(m.files) == 0 {
		b.WriteString("Files to send:\n\n" +
			dimStyle.Render("Press Enter to add\n"))
	} else {
		b.WriteString("Files to send:\n\n")
		for i, f := range m.files {
			mark := "  "
			icon := "📄"
			row := "  " + f
			if info, err := os.Stat(f); err == nil {
				if info.IsDir() {
					// A directory's Size() is meaningless (often 0) — show a
					// folder marker instead of a misleading "0 B".
					icon = "📁"
				} else {
					row += " (" + formatSize(info.Size()) + ")"
				}
			}
			row = icon + " " + row
			if i == m.sel {
				mark = "› "
				row = selectedRow.Render(row)
			}
			b.WriteString(mark + row + "\n")
		}
	}
	list := panel.Render(b.String())

	// Keyboard-only — no on-screen buttons; the footer is the control hint.
	footer := footerLine("↑/↓ or w/s or j/k move · x delete · enter add · shift+a folder · space send · esc back")
	return lipgloss.JoinVertical(lipgloss.Center, list, "", footer)
}

func (m sendModel) View() string {
	var content string
	switch m.state {
	case sendConfirm:
		content = m.sendList()

	case sendWaiting:
		content = m.status + "\n\n" +
			codeBoxStyle.Render("Code: "+m.code+"\n"+dimStyle.Render("(code copied to clipboard by croc)")) +
			"\n\n" + dimStyle.Render("Esc to go back")

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
			b.WriteString("\n\n" + codeBoxStyle.Render("Code: "+m.code+"\n"+dimStyle.Render("(code copied to clipboard by croc)")))
		}
		content = b.String()

	case sendDone:
		content = lipgloss.NewStyle().Foreground(colDone).Render("✓ Transfer complete") +
			"\n\n" + dimStyle.Render("enter send another · esc back")

	case sendError:
		errMsg := "Error"
		if m.err != nil {
			errMsg = m.err.Error()
		}
		if m.status != "" && m.status != errMsg {
			errMsg = m.status + ": " + errMsg
		}
		content = lipgloss.NewStyle().Foreground(colErr).Render("✗ "+errMsg) +
			"\n\n" + dimStyle.Render("Esc to go back")
	}
	return screenView(m.width, m.height, "send files to any device", content)
}
