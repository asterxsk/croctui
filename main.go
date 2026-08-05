package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"crocgui/ui"
)

func main() {
	bin, err := detectStub()
	m := ui.New(bin, err)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

// TODO(task2): replace with croc.Detect().
func detectStub() (string, error) { return "croc", nil }
