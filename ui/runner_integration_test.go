package ui

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// runnerTestModel captures all messages sent to it by a running tea program.
// Init returns the runCroc Cmd so the tea runtime executes it, mimicking how
// the root model launches croc. The program pointer is injected after
// construction (like the real app's RootModel.SetProgram).
type runnerTestModel struct {
	prog     *tea.Program
	events   []tea.Msg
	done     chan struct{}
	fakePath string
}

func (m *runnerTestModel) Init() tea.Cmd {
	return runCroc(m.prog, m.fakePath, []string{}, ".")
}

func (m *runnerTestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.events = append(m.events, msg)
	if _, ok := msg.(crocFinishedMsg); ok {
		close(m.done)
	}
	return m, nil
}

func (m *runnerTestModel) View() string { return "" }

// exeSuffix returns the platform executable suffix.
func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// buildFakeCroc compiles ui/fake_croc.go (a standalone helper, excluded from
// the package build via the +build ignore tag) into a temp executable. This
// avoids committing a platform binary to the repo.
func buildFakeCroc(t *testing.T, outPath string) {
	t.Helper()
	// fake_croc.go lives next to this test file (in the ui package dir).
	src := filepath.Join(".", "fake_croc.go")
	cmd := exec.Command("go", "build", "-o", outPath, src)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build fake_croc: %v\n%s", err, out)
	}
}

// runCrocInProgram runs the runCroc Cmd through a real tea program and returns
// the messages the model received, ordered.
func runCrocInProgram(t *testing.T, fakePath string) []tea.Msg {
	t.Helper()
	m := &runnerTestModel{done: make(chan struct{}), fakePath: fakePath}
	// WithoutRenderer: no TTY needed; WithInput(nil): no keyboard input.
	p := tea.NewProgram(m, tea.WithoutRenderer(), tea.WithInput(nil))
	m.prog = p // inject the program pointer before Run (like SetProgram)
	go func() { _, _ = p.Run() }()

	select {
	case <-m.done:
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for crocFinishedMsg")
	}
	p.Quit()
	return m.events
}

// TestRunCrocNonBlocking pins the fix for Critical 1: runCroc must NOT block
// the Bubble Tea event loop. It runs a fake croc binary that emits output
// lines (Code, then progress) before exiting, and asserts that crocEventMsg
// messages arrive BEFORE crocFinishedMsg — proving the streaming goroutine's
// p.Send calls are delivered while the subprocess is still running, rather
// than being batched until exit (which is what tea.ExecProcess did).
func TestRunCrocNonBlocking(t *testing.T) {
	fakePath := filepath.Join(t.TempDir(), "fake_croc"+exeSuffix())
	buildFakeCroc(t, fakePath)

	events := runCrocInProgram(t, fakePath)

	var gotEvent, gotFinished bool
	for _, msg := range events {
		switch msg.(type) {
		case crocEventMsg:
			gotEvent = true
		case crocFinishedMsg:
			gotFinished = true
			if !gotEvent {
				t.Fatal("crocFinishedMsg arrived before any crocEventMsg — event loop was blocked")
			}
		}
	}
	if !gotEvent {
		t.Fatal("no crocEventMsg received")
	}
	if !gotFinished {
		t.Fatal("no crocFinishedMsg received")
	}
}

// TestRunCrocStreamsCodeAndFile pins that the fake croc's output is parsed
// into the expected events (Code and File), verifying the streaming path.
func TestRunCrocStreamsCodeAndFile(t *testing.T) {
	fakePath := filepath.Join(t.TempDir(), "fake_croc"+exeSuffix())
	buildFakeCroc(t, fakePath)

	events := runCrocInProgram(t, fakePath)

	var sawCode, sawFile bool
	for _, msg := range events {
		if ev, ok := msg.(crocEventMsg); ok {
			if strings.Contains(ev.ev.Code, "fakecode") {
				sawCode = true
			}
			if strings.Contains(ev.ev.File, "file.txt") {
				sawFile = true
			}
		}
	}
	if !sawCode {
		t.Fatal("did not receive the code event from fake croc")
	}
	if !sawFile {
		t.Fatal("did not receive the file event from fake croc")
	}
}
