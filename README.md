# CrocGUI 🐊

A lightweight terminal GUI for the [croc](https://github.com/schollz/croc) file-transfer CLI — built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Send files to any computer peer-to-peer (over the internet, not just LAN), securely end-to-end encrypted, from a friendly terminal interface:

- **Send** — pick files (native Windows dialog, multi-select) → get a share code → copy it → watch live progress
- **Receive** — enter a code → files land in your Downloads → watch live progress
- **Auto-install** — if `croc` isn't installed, the app offers to install it via `winget`

## Run with npx

The npm wrapper builds (once) and runs the Go binary. Requires [Go](https://go.dev/dl/) (1.26+) installed.

```bash
# from this repo
npx .
```

Or install it and run the binary anywhere:

```bash
npm install
npx crocgui
```

The first run compiles the Go source; subsequent runs reuse the built binary (it rebuilds automatically when source files change).

## Run directly with Go

```bash
go run .
# or build once
go build -o crocgui.exe .
./crocgui.exe
```

## Requirements

- **Go 1.26+** (for building)
- **croc** CLI — auto-installed via winget on first launch if missing (or `winget install schollz.croc`)
- **Windows** — the native file dialog and winget install are Windows-only (the rest of the app is portable)
- A real terminal (the TUI needs a TTY; Windows Terminal recommended)

## Tests

```bash
go test ./...
```

## How it works

CrocGUI shells out to the real `croc` CLI as a subprocess — no embedded transfer logic, always behavior-compatible with croc. The `--pass` flag is the relay password (default `pass123`), not the share code; croc generates the code itself and auto-copies it to your clipboard.
