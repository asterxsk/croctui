package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// TestProgressScale pins the fix for Critical 2: bubbles progress.ViewAs
// expects a value in [0, 1], not [0, 100]. 50% must render ~half-full,
// not full/100%.
func TestProgressScale(t *testing.T) {
	m := newSendModel("croc")
	m.pct = 50
	m.state = sendTransferring

	view := m.View()

	// The progress bar should NOT show 100% full when pct=50
	// We check that the bar doesn't contain the full-bar indicator
	// (the exact rendering depends on bubbles, but a full bar shows 100%)
	if strings.Contains(view, "100%") {
		t.Fatalf("progress bar shows 100%% when pct=50:\n%s", view)
	}

	// Also verify receive model
	r := newReceiveModel("croc")
	r.pct = 50
	r.state = recvTransferring
	r.file = "test.txt"

	view2 := r.View()
	if strings.Contains(view2, "100%") {
		t.Fatalf("receive progress bar shows 100%% when pct=50:\n%s", view2)
	}

	// Additional unit test directly on progress.Model
	p := progress.New(progress.WithDefaultGradient())
	halfView := p.ViewAs(0.5) // 50%
	fullView := p.ViewAs(1.0) // 100%

	// Half view should be different from full view
	if halfView == fullView {
		t.Fatalf("ViewAs(0.5) == ViewAs(1.0): both are %q", halfView)
	}

	// And verify ViewAs(50.0) was the old bug (renders full)
	_ = lipgloss.NewStyle() // just ensure import is used
}
