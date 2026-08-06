package ui

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"crocgui/croc"
)

type screen int

const (
	screenWelcome screen = iota
	screenSend
	screenReceive
	screenInstall
	screenHistory
)

// goScreenMsg asks the root model to switch to screen s.
type goScreenMsg struct {
	s   screen
	bin string
}

// startSendMsg tells the root to start the croc send subprocess.
type startSendMsg struct{}

// dialogFailedMsg is sent when the file dialog returns an error.
type dialogFailedMsg struct{ err error }

// RootModel dispatches to the active screen model.
type RootModel struct {
	cur   screen
	bin   string
	prog  *tea.Program
	wel   welcomeModel
	send  sendModel
	recv  receiveModel
	inst  installModel
	hist  historyModel
	width int
}

// setSize stores the terminal size on the active screen so it can fill the
// window and center its content.
func (r *RootModel) setSize(w int) {
	r.width = w
	r.wel.width = w
	r.send.width = w
	r.recv.width = w
	r.inst.width = w
	r.hist.width = w
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

// recordSend appends a send entry to the in-session history.
func (r *RootModel) recordSend(err error) {
	r.hist.add(HistoryEntry{
		Dir:   histSend,
		Files: append([]string{}, r.send.files...),
		Code:  r.send.code,
		Time:  nowStamp(),
		OK:    err == nil,
		Err:   errStr(err),
	})
}

// recordRecv appends a receive entry to the in-session history.
func (r *RootModel) recordRecv(err error) {
	r.hist.add(HistoryEntry{
		Dir:   histRecv,
		Files: []string{r.recv.file},
		Code:  r.recv.code.Value(),
		Time:  nowStamp(),
		OK:    err == nil,
		Err:   errStr(err),
	})
}

// errStr returns "" for nil errors, else the error message.
func errStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// nowStamp returns a short HH:MM wall-clock stamp for history display.
func nowStamp() string {
	return time.Now().Format("15:04")
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
	case screenHistory:
		return r.hist.Init()
	}
	return nil
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.setSize(msg.Width)
		return r, nil
	case goScreenMsg:
		if msg.bin != "" {
			r.bin = msg.bin
		}
		switch msg.s {
		case screenReceive:
			r.recv = newReceiveModel(r.bin)
		case screenSend:
			r.send = newSendModel(r.bin)
		}
		r.cur = msg.s
		return r, r.Init()
	case startReceiveMsg:
		outDir, err := downloadsDir()
		if err != nil {
			return r, func() tea.Msg { return crocFinishedMsg{err: err} }
		}
		r.recv.outDir = outDir
		return r, runCroc(r.prog, r.bin, croc.ReceiveArgs(msg.code, outDir), outDir)
	case startSendMsg:
		dir, _ := os.Getwd()
		return r, runCroc(r.prog, r.bin, croc.SendArgs(r.send.files), dir)
	case dialogFailedMsg:
		r.send.err = msg.err
		r.send.state = sendError
		return r, nil
	case crocEventMsg, crocFinishedMsg:
		var cmd tea.Cmd
		if r.cur == screenReceive {
			var nm tea.Model
			nm, cmd = r.recv.Update(msg)
			if recv, ok := nm.(receiveModel); ok {
				r.recv = recv
			}
		} else if r.cur == screenSend {
			var nm tea.Model
			nm, cmd = r.send.Update(msg)
			if s, ok := nm.(sendModel); ok {
				r.send = s
			}
		}
		// Record history when a transfer finishes.
		if fin, ok := msg.(crocFinishedMsg); ok {
			if r.cur == screenSend {
				r.recordSend(fin.err)
			} else if r.cur == screenReceive {
				r.recordRecv(fin.err)
			}
		}
		return r, cmd
	}
	switch r.cur {
	case screenWelcome:
		var cmd tea.Cmd
		r.wel, cmd = updateWelcome(r.wel, msg)
		return r, cmd
	case screenReceive:
		var cmd tea.Cmd
		nm, cmd := r.recv.Update(msg)
		if recv, ok := nm.(receiveModel); ok {
			r.recv = recv
		}
		return r, cmd
	case screenSend:
		var cmd tea.Cmd
		nm, cmd := r.send.Update(msg)
		if s, ok := nm.(sendModel); ok {
			r.send = s
		}
		return r, cmd
	case screenHistory:
		var cmd tea.Cmd
		nm, cmd := r.hist.Update(msg)
		if h, ok := nm.(historyModel); ok {
			r.hist = h
		}
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
	case screenSend:
		return r.send.View()
	case screenInstall:
		return r.inst.View()
	case screenHistory:
		return r.hist.View()
	}
	return ""
}
