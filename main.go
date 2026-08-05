package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"crocgui/croc"
	"crocgui/ui"
)

func main() {
	bin, err := croc.Detect()
	m := ui.New(bin, err)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}
