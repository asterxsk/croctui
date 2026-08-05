package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"crocgui/croc"
)

type screen int

const (
	screenWelcome screen = iota
	screenSend
	screenReceive
	screenInstall
)

// goScreenMsg asks the root model to switch to screen s.
type goScreenMsg struct{ s screen }

// RootModel dispatches to the active screen model.
type RootModel struct {
	cur  screen
	bin  string
	prog *tea.Program
	wel  welcomeModel
	send sendModel
	recv receiveModel
	inst installModel
}

// New returns the root model. If installErr != nil the app starts on the
// install screen (croc is missing).
func New(bin string, installErr error) tea.Model {
	r := &RootModel{cur: screenWelcome, bin: bin, wel: welcomeModel{selected: 0}}
	if installErr != nil {
		r.cur = screenInstall
		r.inst = newInstallModel()
	}
	return r
}

// SetProgram sets the program pointer so the root can spawn croc commands.
func (r *RootModel) SetProgram(p *tea.Program) {
	r.prog = p
}

func (r *RootModel) Init() tea.Cmd {
	switch r.cur {
	case screenWelcome:
		return r.wel.Init()
	case screenSend:
		return r.send.Init()
	case screenReceive:
		return r.recv.Init()
	case screenInstall:
		return r.inst.Init()
	}
	return nil
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case goScreenMsg:
		switch msg.s {
		case screenReceive:
			r.recv = newReceiveModel(r.bin)
		case screenSend:
			// Task 5
		}
		r.cur = msg.s
		return r, r.Init()
	case startReceiveMsg:
		outDir, err := downloadsDir()
		if err != nil {
			return r, nil // surface via a message in Task 6 polish
		}
		return r, runCroc(r.prog, r.bin, croc.ReceiveArgs(msg.code, outDir), outDir)
	case crocEventMsg, crocFinishedMsg:
		var cmd tea.Cmd
		if r.cur == screenReceive {
			var nm tea.Model
			nm, cmd = r.recv.Update(msg)
			r.recv = nm.(receiveModel)
		}
		return r, cmd
	}
	switch r.cur {
	case screenWelcome:
		var cmd tea.Cmd
		r.wel, cmd = updateWelcome(r.wel, msg)
		return r, cmd
	}
	return r, nil
}

func (r *RootModel) View() string {
	switch r.cur {
	case screenWelcome:
		return r.wel.View()
	case screenReceive:
		return r.recv.View()
	}
	return ""
}

// Placeholder screen models, implemented in later tasks.

type sendModel struct{}    // implemented in Task 5
type installModel struct{} // implemented in Task 6

func (sendModel) Init() tea.Cmd     { return nil }
func (installModel) Init() tea.Cmd  { return nil }
func newInstallModel() installModel { return installModel{} }