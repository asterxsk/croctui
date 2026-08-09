package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// theme.go defines the single visual vocabulary shared by every screen.
// One blue accent carried by structural regions (wordmark, borders, selected
// states, progress); neutrals for the surface; semantic colors reserved for
// real states (error, done, dim). Every screen builds on these styles so the
// app reads as one coherent surface across send / receive / history.
//
// The outer window frame is drawn by frame() below. It is the one element
// every screen shares, so it lives here next to the box styles it wraps.

var (
	// accent is the primary blue used for identity and emphasis.
	accent = lipgloss.Color("39")

	// Semantic / state colors.
	colErr  = lipgloss.Color("196") // red
	colDone = lipgloss.Color("42")  // green
	colDim  = lipgloss.Color("8")   // gray
	colText = lipgloss.Color("15")  // white

	// Box border style — the one border vocabulary. Square (┌─┐││└─┘).
	boxBorder = lipgloss.NormalBorder()

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

	// logoBox frames the CROC ASCII-art header shown on every screen.
	logoBox = lipgloss.NewStyle().
		Padding(0, 2)

	// wordmark is the CROC ASCII-art title on the welcome screen.
	wordmark = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)

	// selectedButton is the active action button: an accent OUTLINE around the
	// button's shape (bold accent text inside an accent border), not a solid
	// fill. Selection reads as a bright accent outline rather than a solid
	// glowing block.
	selectedButton = lipgloss.NewStyle().
			Bold(true).
			Border(boxBorder).
			BorderForeground(accent).
			Foreground(accent).
			Padding(0, 3).
			Margin(0, 1)

	// idleButton is a non-selected action button.
	idleButton = lipgloss.NewStyle().
			Border(boxBorder).
			BorderForeground(lipgloss.Color("8")).
			Foreground(colText).
			Padding(0, 2).
			Margin(0, 1)

	// footerHint is the dim keybinding hint text (no vertical margin; screens
	// control spacing via footerPad).
	footerHint = lipgloss.NewStyle().
			Foreground(colDim)

	// statusText is normal body text on a panel.
	statusText = lipgloss.NewStyle().Foreground(colText)

	// selectedRow highlights the current row in a navigable list (accent fg,
	// bold) — used by the send file list.
	selectedRow = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)
)

// separator is a thin horizontal rule that divides sections of a screen
// (e.g. the buttons from the keybinding hint).
var (
	separator = lipgloss.NewStyle().
		Foreground(colDim).
		Width(28).
		MaxWidth(28)
)

// separatorLine renders a dim rule of ─ characters (capped at 28 cols) that
// reads as a lightweight section divider between screen regions.
func separatorLine() string {
	return separator.Render(strings.Repeat("─", 28))
}

// crocWordmark is the ASCII-art CROC splash shown at the top of every screen.
// It is kept in sync with the brand: block letters, accent blue, no emoji.
const crocWordmark = ` ██████╗██████╗  ██████╗  ██████╗
██╔════╝██╔══██╗██╔═══██╗██╔════╝
██║     ██████╔╝██║   ██║██║
██║     ██╔══██╗██║   ██║██║
╚██████╗██║  ██║╚██████╔╝╚██████╗
 ╚═════╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝`

// crocHeader is the bordered CROC logo: the ASCII art wrapped in a square
// accent frame, with a dim caption beneath it. Every screen renders this at
// the top, and its content sits centered beneath it.
func crocHeader(caption string) string {
	logo := logoBox.Render(wordmark.Render(crocWordmark))
	lines := []string{logo}
	if caption != "" {
		lines = append(lines, dimStyle.Render(caption))
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

// screenView composes a full-screen view for a given terminal size: the shared
// bordered CROC header on top, the screen's content centered beneath it, and
// the whole thing centered (and framed) in the terminal window.
//
//   - w <= 0 (unit tests / pre-resize): returns the header + content stack
//     uncentered and unframed, so tests assert on content, not chrome.
//   - w > 0, h > 0: renders a full square window frame that fills the terminal
//     and vertically centers the header+content block.
//   - w > 0, h <= 0 (only width known): renders the frame top-aligned in a
//     stub so resizes don't jitter.
func screenView(w, h int, caption, content string) string {
	header := crocHeader(caption)
	if w <= 0 {
		return lipgloss.JoinVertical(lipgloss.Center, header, content)
	}
	return frame(w, h, lipgloss.JoinVertical(lipgloss.Center, header, content))
}

// frame wraps a screen's content in the main window border, filling the full
// terminal width and height. When there is spare vertical room (content is
// shorter than the window) content is centered inside the bordered window;
// when content is taller than the window it renders at its natural height so
// the border always wraps the content and the top border never scrolls off.
// When width is not known at all (w <= 0, e.g. unit tests) it renders content
// alone so tests assert on content, not chrome.
//
// The frame is computed inside-out: the border consumes 2 columns, so content
// is placed at Width(w-2) and the border style wraps it. Padding(0,1) leaves a
// 1-col breathing gutter inside the border, which narrows the content box by a
// further 2. The outer frame is never given a Margin — a margin would shrink
// Width() and the frame would no longer fill the window.
func frame(width, height int, content string) string {
	if width <= 0 {
		return content
	}
	padded := lipgloss.NewStyle().Padding(0, 1).Render(content)

	// The border adds 2 rows to the final height. Only expand vertically (via
	// Place) when the natural bordered height leaves room to spare; otherwise
	// keep the natural height so the border always wraps the content exactly
	// and never overflows the terminal (which would scroll the top border off).
	boxH := lipgloss.Height(padded)
	if height > 2 && boxH+2 < height {
		boxH = height - 2
	}
	inner := lipgloss.Place(width-2, boxH, lipgloss.Center, lipgloss.Center, padded)
	return lipgloss.NewStyle().
		Border(boxBorder).
		BorderForeground(accent).
		Render(inner)
}

// footerPad renders the text key-hint footer at a FIXED width (70 cols) and
// FIXED height (3 rows), so the centered layout never shifts between screens.
// The hint text is centered within the 70-col block so it reads centered (not
// shoved to the left edge) on every screen.
func footerPad(text string) string {
	return footerHint.
		Width(70).
		MaxWidth(70).
		Align(lipgloss.Center).
		Render("\n\n" + text)
}

// footerLine renders a plain hint line at a FIXED width (70 cols) and height
// (1 row), centered. Used by screens with tight vertical space (the send list)
// where the full 3-row footer doesn't fit.
func footerLine(text string) string {
	return footerHint.
		Width(70).
		MaxWidth(70).
		Align(lipgloss.Center).
		Render(text)
}
