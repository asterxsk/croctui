package ui

import (
	"strings"
	"testing"
)

// width returns the display width of a rendered line (plain ASCII in tests).
func width(s string) int { return len([]rune(s)) }

// stripANSI removes ANSI escape sequences so assertions can check plain text.
func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			for i < len(s) && s[i] != 'm' {
				i++
			}
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// TestFrameFillsWidth pins the full-window requirement: frame() at a terminal
// width must render exactly that many columns, with a border on the outer edge
// (never a margin-shortened box). The border consumes 2 columns, the padding
// gutter another 2.
func TestFrameFillsWidth(t *testing.T) {
	out := frame(80, 0, "hello")
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("frame rendered no lines")
	}
	// Top border must span the full width with the correct corners.
	top := lines[0]
	if !strings.HasPrefix(top, "┌") || !strings.HasSuffix(top, "┐") {
		t.Fatalf("top border malformed: %q", top)
	}
	if width(top) != 80 {
		t.Fatalf("frame width = %d, want 80", width(top))
	}
	// A content line must be bracketed by the side borders.
	if !strings.HasPrefix(lines[1], "│") || !strings.HasSuffix(lines[1], "│") {
		t.Fatalf("content line missing side borders: %q", lines[1])
	}
	// The border must be a rectangle: every line the same width.
	for _, ln := range lines {
		if width(ln) != 80 {
			t.Fatalf("ragged frame line: width %d: %q", width(ln), ln)
		}
	}
}

// TestFrameZeroWidthPipesThrough makes sure the no-chrome fallback (used by
// unit tests and pre-resize renders) keeps content unadorned.
func TestFrameZeroWidthPipesThrough(t *testing.T) {
	out := frame(0, 0, "plain")
	if out != "plain" {
		t.Fatalf("frame(0,0,..) = %q, want content passthrough", out)
	}
}

// TestFrameCentersVertically checks that when height is known the content sits
// vertically centered inside the window (not crammed at the top).
func TestFrameCentersVertically(t *testing.T) {
	out := frame(40, 8, "X")
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != 8 {
		t.Fatalf("frame height = %d, want 8", len(lines))
	}
	// Border box of 8 lines: row[0] top, row[7] bottom, 6 interior rows.
	// Content "X" should be on one of the two middle interior rows (3 or 4).
	found := false
	for _, ln := range lines[1:7] {
		if strings.Contains(ln, "X") {
			found = true
		}
	}
	if !found {
		t.Fatalf("content X not centered in frame:\n%s", out)
	}
}

// TestScreenViewRendersLogoHeader pins that screenView always composes the
// borderless CROC logo at the top, centered, above the screen's content.
func TestScreenViewRendersLogoHeader(t *testing.T) {
	out := screenView(80, 24, "caption here", "body")
	joined := stripANSI(out)
	// The CROC art must appear WITHOUT its own box around it.
	if !strings.Contains(joined, "██████╗") {
		t.Fatalf("screenView missing CROC wordmark:\n%s", joined)
	}
	if !strings.Contains(joined, "caption here") {
		t.Fatalf("screenView missing caption:\n%s", joined)
	}
	// Body must be present below the logo.
	if !strings.Contains(joined, "body") {
		t.Fatalf("screenView missing body content:\n%s", joined)
	}
	// The logo line must be inside the window frame (centered), not inside a
	// nested logo box.
	for _, ln := range strings.Split(joined, "\n") {
		if strings.Contains(ln, "██████╗") {
			if !strings.HasPrefix(ln, "│") || !strings.HasSuffix(ln, "│") {
				t.Fatalf("logo line not inside window frame: %q", ln)
			}
			break
		}
	}
}

// TestScreenViewFitsWindow pins that the whole screen (logo + content + frame)
// never exceeds the terminal height, so the top border stays visible.
func TestScreenViewFitsWindow(t *testing.T) {
	out := screenView(80, 24, "caption", "body\n\nmore")
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) > 24 {
		t.Fatalf("screenView rendered %d lines (> 24):\n%s", len(lines), out)
	}
	if !strings.HasPrefix(stripANSI(lines[0]), "┌") {
		t.Fatalf("screenView top border missing: %q", lines[0])
	}
}

// TestFrameNoHiddenExtraHeight pins the "top border doesn't show" fix: the
// frame must never render taller than the terminal window (that's what scrolls
// the top border off the alt screen). When the window is taller than the
// content the frame fills it; when it's shorter the frame stays at the natural
// content height and the top border must still be present.
func TestFrameNoHiddenExtraHeight(t *testing.T) {
	content := strings.Repeat("line\n", 15) + "line" // 16 content rows
	natural := 18                                    // 16 + 2 border rows
	for _, h := range []int{12, 14, 16, 20, 24} {
		out := frame(80, h, content)
		lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
		want := natural
		if h > natural {
			want = h
		}
		if len(lines) != want {
			t.Fatalf("frame at height %d rendered %d lines, want %d:\n%s", h, len(lines), want, out)
		}
		// Top border must be present.
		if !strings.HasPrefix(stripANSI(lines[0]), "┌") {
			t.Fatalf("frame at height %d missing top border, first line: %q", h, lines[0])
		}
	}
}

// TestWelcomeContainsCrocWordmark verifies the ASCII-art header renders.
func TestWelcomeContainsCrocWordmark(t *testing.T) {
	w := welcomeModel{width: 80, height: 24}
	v := w.View()
	if !strings.Contains(v, "██████╗") {
		t.Fatalf("welcome view missing CROC block wordmark:\n%s", v)
	}
	if !strings.Contains(v, "securely transfer files peer-to-peer over the internet") {
		t.Fatalf("welcome view missing transfer caption:\n%s", v)
	}
}

// TestAllScreensShowLogoAndFitWindow pins the shared-logo layout: every screen
// renders the square-bordered CROC logo header, and the whole screen (logo +
// content + window frame) fits within the terminal height.
func TestAllScreensShowLogoAndFitWindow(t *testing.T) {
	build := func() []struct {
		name string
		out  string
	} {
		w := welcomeModel{selected: 0, width: 80, height: 24}
		r := newReceiveModel("croc")
		r.width, r.height = 80, 24
		s := newSendModel("croc")
		s.width, s.height = 80, 24
		s.state = sendConfirm
		s.files = []string{"a.txt"}
		h := newHistoryModel()
		h.width, h.height = 80, 24
		h.add(HistoryEntry{Dir: histSend, Files: []string{"a.txt"}, Code: "1", Time: "12:00", OK: true})
		im := newInstallModel()
		im.width, im.height = 80, 24
		im.state = installRunning
		im.progress = 0.5
		return []struct {
			name string
			out  string
		}{
			{"welcome", w.View()},
			{"receive", r.View()},
			{"send", s.View()},
			{"history", h.View()},
			{"install", im.View()},
		}
	}

	for _, sc := range build() {
		clean := stripANSI(sc.out)
		// Logo box + CROC art must be present.
		if !strings.Contains(clean, "┌") || !strings.Contains(clean, "██████╗") {
			t.Fatalf("%s screen missing logo header:\n%s", sc.name, clean)
		}
		// The whole screen must fit the 24-row window (top border visible).
		if n := len(strings.Split(strings.TrimSuffix(clean, "\n"), "\n")); n > 24 {
			t.Fatalf("%s screen rendered %d lines (> 24 window):\n%s", sc.name, n, clean)
		}
		if !strings.HasPrefix(clean, "┌") {
			t.Fatalf("%s screen top border missing:\n%s", sc.name, clean)
		}
	}
}
