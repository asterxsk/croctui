package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"crocgui/croc"
)

// installTickMsg advances the pulsing progress bar while the install runs.
type installTickMsg time.Time

type installState int

const (
	installPrompt installState = iota
	installRunning
	installDone
	installFailed
	installNotDetected
)

type installModel struct {
	state    installState
	spinner  spinner.Model
	bar      progress.Model // bubbles progress bar shown while installing
	progress float64        // 0..1 pulsing during install
	progUp   bool           // direction of the pulse
	err      error
	width    int
	height   int
}

func newInstallModel() installModel {
	return installModel{
		state:   installPrompt,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		bar:     progress.New(progress.WithDefaultGradient()),
	}
}

type installStartedMsg struct{}
type installFinishedMsg struct{ err error }

// installProgressTick is a Cmd that schedules the next pulse step of the
// install progress bar (roughly 20 fps while winget runs).
func installProgressTick() tea.Msg {
	time.Sleep(50 * time.Millisecond)
	return installTickMsg{}
}

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
			return m, tea.Batch(m.spinner.Tick, installProgressTick, func() tea.Msg {
				return installFinishedMsg{err: croc.Install()}
			})
		case installTickMsg:
			// Pulse the progress bar toward 100% (it never quite arrives; the
			// installFinishedMsg lands when winget actually finishes).
			const step = 0.06
			if m.progUp {
				m.progress += step
				if m.progress >= 0.95 {
					m.progress = 0.95
					m.progUp = false
				}
			} else {
				m.progress -= step
				if m.progress <= 0.15 {
					m.progress = 0.15
					m.progUp = true
				}
			}
			return m, installProgressTick
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
		box := panel.Render(
			"croc is not installed.\n\n" +
				"Install it via winget?\n\n" +
				dimStyle.Render("This will run: winget install schollz.croc"),
		)
		content = box
	case installRunning:
		bar := ""
		if m.progress >= 0 {
			bar = m.bar.ViewAs(m.progress)
		}
		content = m.spinner.View() + " Installing croc via winget…\n\n" +
			bar + "\n\n" +
			dimStyle.Render("Downloading the croc binary…")
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
	return screenView(m.width, m.height, "set up croc", content)
}
