package ui

import tea "github.com/charmbracelet/bubbletea"

type screen int

const (
	screenWelcome screen = iota
	screenSend
	screenReceive
	screenInstall
)

// goScreenMsg asks the root model to switch to screen s.
type goScreenMsg struct{ s screen }

// rootModel dispatches to the active screen model.
type rootModel struct {
	cur  screen
	bin  string
	wel  welcomeModel
	send sendModel
	recv receiveModel
	inst installModel
}

// New returns the root model. If installErr != nil the app starts on the
// install screen (croc is missing).
func New(bin string, installErr error) tea.Model {
	r := &rootModel{cur: screenWelcome, bin: bin, wel: welcomeModel{selected: 0}}
	if installErr != nil {
		r.cur = screenInstall
		r.inst = newInstallModel()
	}
	return r
}

func (r *rootModel) Init() tea.Cmd {
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

func (r *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case goScreenMsg:
		r.cur = msg.s
		r.bin = "" // cleared; re-set on route entry (see Tasks 4-6)
		return r, r.Init()
	}
	switch r.cur {
	case screenWelcome:
		var cmd tea.Cmd
		r.wel, cmd = updateWelcome(r.wel, msg)
		return r, cmd
	}
	return r, nil
}

func (r *rootModel) View() string {
	switch r.cur {
	case screenWelcome:
		return r.wel.View()
	}
	return ""
}

// Placeholder screen models, implemented in later tasks.

type sendModel struct{}    // implemented in Task 5
type receiveModel struct{} // implemented in Task 4
type installModel struct{} // implemented in Task 6

func (sendModel) Init() tea.Cmd     { return nil }
func (receiveModel) Init() tea.Cmd  { return nil }
func (installModel) Init() tea.Cmd  { return nil }
func newInstallModel() installModel { return installModel{} }
