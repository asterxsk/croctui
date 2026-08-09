# AGENTS.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project overview

CrocGUI is a Windows-first terminal GUI for the [`croc`](https://github.com/schollz/croc) file-transfer CLI. Go 1.26.5 + Bubble Tea v1.3.10 + lipgloss v1.1.0. Module name is `crocgui`. It shells out to the real `croc` binary as a subprocess — there is no embedded transfer logic, and this is a deliberate product constraint (`--pass` is the relay password, not the share code; croc generates the code itself and copies it to the clipboard).

## Commands

Run from the repo root. PowerShell is the primary shell.

```powershell
go build ./...          # build all packages
go build -o bin/crocgui.exe .   # build the app binary
go run .                # run the TUI (needs a real TTY; Windows Terminal recommended)
go vet ./...
go test ./...           # full test suite
go test ./ui -run TestRunCrocNonBlocking -v   # single test
go test ./ui -run 'TestRunCroc|TestSend'      # tests matching a pattern
```

npm wrapper (builds on first run, rebuilds when sources are newer than the binary):

```powershell
npm run build   # == go build -o bin/crocgui.exe .
npx .           # build-if-stale, then run
```

Debug tooling:

- No mouse capture: the app runs `tea.NewProgram` with only `tea.WithAltScreen()` — it is keyboard-only by design.

The suite includes integration tests that compile `ui/fake_croc.go` (a standalone helper with `// +build ignore`, built into a temp dir at test time — never commit a binary). Full suite takes ~5s; `ui` package is the slow one.

## Architecture

Three packages plus `main`:

- `croc` — everything about the croc CLI: `Detect()` (LookPath), `SendArgs`/`ReceiveArgs` argv builders, `Install()` (winget), and `parse.go` which regex-parses croc's stdout into typed `Event`s (`EventCode`, `EventFile`, `EventProgress`, `EventError`, `EventStatus`).
- `dialog` — native Windows file/folder pickers. Shells out to `powershell.exe -NoProfile -STA` running WinForms scripts; output is UTF-8, one absolute path per line; empty output = cancelled.
- `ui` — all Bubble Tea models and rendering.
- `main.go` — detects croc, creates the root model, runs `tea.NewProgram` with `tea.WithAltScreen()`. Keyboard-only; no mouse capture.

### Message flow (read `ui/root.go` first)

`RootModel` owns the screen enum (welcome / send / receive / install / history) and dispatches every message to the active screen model. Screen models are **not** persistent: switching to send/receive constructs a fresh model (`newSendModel`/`newReceiveModel`). The root re-applies the last-known window size after a screen switch because fresh models start at 0×0.

Transfers are subprocess-driven. The root handles `startSendMsg`/`startReceiveMsg` by calling `runCroc` (`ui/runner.go`), which spawns croc with stdout+stderr wired to an `io.Pipe`, streams lines into `crocEventMsg` via `p.Send` from a goroutine, and returns `crocFinishedMsg` when the process exits. Screen models consume `crocEventMsg`/`crocFinishedMsg` to drive their state machines; the root records history entries on finish. Keep the non-blocking design: `runCroc` must never block the Bubble Tea event loop (pinned by `TestRunCrocNonBlocking`).

### Selection animation (`ui/welcome.go`)

The app is keyboard-only — there is no mouse handling anywhere. The one animated element is the welcome screen's selection: a `welcomePulseMsg` tick (300 ms) flips `welcomeModel.pulse` between 0 and 1, and `welcomeButton` renders the selected button with a breathing chevron (`▸` ↔ `▶`, border/text pulsing accent ↔ color 33). List selections (send/history) use a static `›` accent marker; `selectedButton` is an accent outline (border + bold accent text), not a solid fill.

### Rendering and design language (`ui/theme.go`)

One shared visual vocabulary: `accent` = color 39 blue; square `NormalBorder`; `selectedButton` is an accent outline (bold accent text in an accent-bordered box); `panel`/`codeBox`/`idleButton` for structure. Every screen composes via `screenView(w, h, caption, content)` = CROC ASCII-art header + caption + centered content, wrapped by `frame()` (full-window square border, computed inside-out: border 2 cols + `Padding(0,1)` gutter). `w <= 0` (unit tests) renders unframed so tests assert on content, not chrome.

### Session model

History is in-session only (`ui/history.go`, entries appended by `RootModel.recordSend`/`recordRecv`) — nothing persists to disk. Receive always targets the Downloads folder (`downloadsDir()` in the ui package). If `croc.Detect()` fails at startup, the app opens on the install screen which offers `winget install schollz.croc`.

## Windows-specific constraints

- `dialog` and `install.go` are Windows-only (PowerShell/WinForms, winget). The rest is portable.
- The TUI requires a real TTY; it will not behave correctly under a non-interactive runner.
- Binary output goes to `bin/crocgui.exe` (gitignored via `*.exe`).
