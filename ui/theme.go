package ui

import "github.com/charmbracelet/lipgloss"

// theme.go defines the single visual vocabulary shared by every screen.
// One blue accent carried by structural regions (wordmark, borders, selected
// states, progress); neutrals for the surface; semantic colors reserved for
// real states (error, done, dim). Every screen builds on these styles so the
// app reads as one coherent surface across send / receive / history.

var (
	// accent is the primary blue used for identity and emphasis.
	accent = lipgloss.Color("39")

	// Semantic / state colors.
	colErr   = lipgloss.Color("196") // red
	colDone  = lipgloss.Color("42")  // green
	colDim   = lipgloss.Color("8")   // gray
	colText  = lipgloss.Color("15")  // white
	colTitle = lipgloss.Color("6")   // cyan, for titles

	// Box border style — the one border vocabulary.
	boxBorder = lipgloss.RoundedBorder()

	// panel is a bordered, padded region used for lists, code, progress.
	panel = lipgloss.NewStyle().
		Border(boxBorder).
		BorderForeground(accent).
		Padding(1, 2).
		Margin(1, 2)

	// subtlePanel is a dimmer bordered region (secondary info).
	subtlePanel = lipgloss.NewStyle().
			Border(boxBorder).
			BorderForeground(lipgloss.Color("8")).
			Padding(1, 2).
			Margin(1, 2)

	// codeBox is the emphasized region that holds the share code.
	codeBox = lipgloss.NewStyle().
		Border(boxBorder).
		BorderForeground(accent).
		Padding(1, 3).
		Margin(1, 0).
		Bold(true)

	// wordmark is the CROC title on the welcome screen.
	wordmark = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			Margin(0, 0, 1, 0)

	// headerLine / footerLine frame every screen.
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colTitle).
			Margin(0, 0, 1, 0)

	// selectedButton is the active action button.
	selectedButton = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(accent).
			Padding(0, 2).
			Margin(0, 1)

	// idleButton is a non-selected action button.
	idleButton = lipgloss.NewStyle().
			Border(boxBorder).
			BorderForeground(lipgloss.Color("8")).
			Foreground(colText).
			Padding(0, 2).
			Margin(0, 1)

	// footerHint is the dim keybinding hint bar.
	footerHint = lipgloss.NewStyle().
			Foreground(colDim).
			Margin(1, 0, 0, 0)

	// statusText is normal body text on a panel.
	statusText = lipgloss.NewStyle().Foreground(colText)
)
