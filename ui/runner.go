package ui

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"crocgui/croc"
)

// crocEventMsg carries one parsed line of croc output to the model.
type crocEventMsg struct{ ev croc.Event }

// crocFinishedMsg is delivered when the croc subprocess exits.
type crocFinishedMsg struct{ err error }

// buildCrocCmd constructs the subprocess command with stdout+stderr wired to a
// shared buffer so output is captured without leaking a pipe reader.
// Kept pure for tests; production output streaming is handled by runCroc.
func buildCrocCmd(bin string, args []string, dir string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	return cmd
}

// streamOutput reads croc's output from r, splitting on \r and \n so progress
// frames arrive live, and sends one crocEventMsg per parsed line. It returns
// when r hits EOF (the pipe writer is closed on process exit).
func streamOutput(p *tea.Program, r io.Reader) {
	defer func() { _ = r.(io.Closer).Close() }()
	br := bufio.NewReader(r)
	var frag []byte
	flush := func() {
		if len(strings.TrimSpace(string(frag))) == 0 {
			frag = frag[:0]
			return
		}
		p.Send(crocEventMsg{ev: croc.ParseLine(string(frag))})
		frag = frag[:0]
	}
	for {
		b, err := br.ReadByte()
		if err != nil {
			break
		}
		if b == '\r' || b == '\n' {
			flush()
			continue
		}
		frag = append(frag, b)
	}
	flush()
}

// runCroc starts croc and returns a tea.Cmd that completes when croc exits.
// While running, each output line is parsed and dispatched as crocEventMsg.
//
// The streaming goroutine is started inside the returned tea.Cmd — i.e. only
// when Bubble Tea actually executes it. If the command is never run (e.g. the
// app quits), no goroutine is spawned and nothing leaks.
func runCroc(p *tea.Program, bin string, args []string, dir string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(bin, args...)
		cmd.Dir = dir
		r, w := io.Pipe()
		cmd.Stdout = w
		cmd.Stderr = w

		// Start streaming before the process runs so its first output is read.
		go streamOutput(p, r)

		// Block until the process exits; the goroutine is draining the pipe.
		// w.Close() signals EOF, unblocking streamOutput's final flush.
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			w.Close()
			return crocFinishedMsg{err: err}
		})()
	}
}