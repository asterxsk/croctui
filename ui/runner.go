package ui

import (
	"bufio"
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
// shared pipe so output can be streamed back as messages. Kept pure for tests.
func buildCrocCmd(bin string, args []string, dir string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	_, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	// TODO: read from r in a goroutine inside runCroc
	return cmd
}

// runCroc starts croc and returns a tea.Cmd that completes when croc exits.
// While running, each output line is parsed and dispatched as crocEventMsg.
func runCroc(p *tea.Program, bin string, args []string, dir string) tea.Cmd {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w

	// Stream output: split on \r and \n so progress frames arrive live.
	go func() {
		defer r.Close()
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
	}()

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		w.Close()
		return crocFinishedMsg{err: err}
	})
}