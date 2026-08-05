# CrocGUI — Design Spec

Date: 2026-08-05
Status: Approved (design direction)

## Overview

A lightweight terminal GUI for [croc](https://github.com/schollz/croc) — a peer-to-peer file transfer tool with end-to-end encryption. CrocGUI wraps the existing `croc` CLI in a friendly terminal interface: a Croc-style welcome screen, a **Send** flow (pick files via a native Windows dialog, get a share code, watch live progress), and a **Receive** flow (paste a code, watch the transfer, files land in Downloads).

The app is intentionally tiny: it shells out to the installed `croc` binary and adds no file-transfer logic of its own.

## Goals / Non-goals

**Goals**
- Send files: pick multiple files → get a code → share it → live progress.
- Receive files: enter a code → live progress → files saved.
- As light as possible: one small Go static binary, minimal dependencies.
- Friendly, Croc-styled terminal UI.

**Non-goals (YAGNI)**
- No own relay/PAKE/encryption logic — that is croc's job.
- No embedded croc library — we shell out (see Architecture).
- No cross-platform file dialog (Windows-only for the native dialog; see Trade-offs).
- No transfer history / resume management beyond what croc does by itself.
- No QR-code display, no transfer settings screen (croc defaults are used).
- No background/queued transfers — one transfer at a time, foreground.

## Tech stack

| Concern | Choice | Why |
|---------|--------|-----|
| Language | Go 1.26 (installed) | Installed, single static binary, tiny footprint |
| TUI | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Go's Ratatui; mature, minimal, async via `Cmd`/`Msg` |
| Components | [bubbles](https://github.com/charmbracelet/bubbles) (progress, textinput, list) | Ready-made progress bar + text input |
| Styling | [lipgloss](https://github.com/charmbracelet/lipgloss) | Composable styling for the Croc look |
| Engine | Shell out to `croc` CLI subprocess | Always matches real croc; zero transfer logic |
| File picker | Native Windows dialog via PowerShell bridge | User-selected; familiar Explorer feel, multi-select |
| Clipboard | `clip` (Windows built-in) | TUI can't touch clipboard directly |
| Auto-install | `winget install schollz.croc` | Matches how croc is already installed; prompts user |

## Architecture

### Components

- **`main.go`** — entrypoint, starts the Bubble Tea program.
- **`ui/`** — TUI screens (Bubble Tea models):
  - `welcome.go` — splash + Send/Receive buttons.
  - `send.go` — file list confirmation → launch → code display → progress.
  - `receive.go` — code input → progress → completion.
  - `progress.go` — shared live-progress view (progress bar + rate/status).
- **`croc/`** — the engine adapter (pure Go, no TUI dependencies):
  - `runner.go` — builds and runs `croc` subprocess commands.
  - `parse.go` — streams and parses croc stdout/stderr into progress events.
  - `install.go` — detects missing binary, prompts/triggers winget install.
- **`dialog/`** — native Windows file dialog bridge (PowerShell).
- **`clipboard/`** — copy code to clipboard via `clip`.

### Data flow

```
User picks "Send" ──► dialog: native multi-select ──► file list confirmed
     └─► croc send <files> --yes  (subprocess; croc generates the code)
          └─► croc prints code + progress lines
               └─► parser ──► progress Msgs ──► TUI live progress bar
                    └─► done/error ──► completion screen
```

### Key detail: driving croc non-interactively

croc's normal TUI output is interactive and terminal-graphics-heavy (progress bars, prompts, colors). To parse it reliably:

- **Send:** `croc send <file1> <file2> … --yes` — `--yes` skips prompts; croc generates and prints the code itself (see Code generation).
- **Receive:** `croc <code> --yes` — `--yes` auto-accepts; downloads into the current working directory (we set it to `Downloads`).
- Progress lines are parsed from stdout/stderr with a regex for the percent and rate; anything unparsed is passed through as status text.

We rely on croc's stable CLI flags (`--code`, `--yes`, `send` subcommand) — not on internal library APIs.

### Code generation

`croc send` generates its own random code and prints it (e.g. `Code is: 1234-5678-9012`). We capture that line and offer a "copy code" key that pipes it to `clip`. No code generation on our side.

### Missing-binary flow

At startup, `exec.LookPath("croc")`. If absent, show a screen: *"croc is not installed. Install it via winget?"* → runs `winget install schollz.croc` (streaming output) → re-check → continue or retry.

### Error handling

- Every screen keeps a status/error line; errors are never silent.
- Subprocess failures (non-zero exit, missing binary, timeout) surface as a clear message with a back/retry option.
- A global "force quit" key (Ctrl+C twice / `q`) restores the terminal even if a subprocess is stuck — we always `defer` terminal cleanup.

## UI specification

### Screen 1 — Welcome

```
┌────────────────────────────────────────────┐
│                                            │
│    🐊  CROC                                │
│    easily & securely transfer files        │
│    peer-to-peer over the internet          │
│                                            │
│       [ Send ]    [ Receive ]              │
│                                            │
│    q / ctrl-c to quit                      │
└────────────────────────────────────────────┘
```

- `S` / tab / arrows select, Enter confirms.
- `q` / Ctrl+C quits.

### Screen 2 — Send

1. Press Send → native Windows dialog (multi-select) opens.
2. Cancel → back to welcome.
3. Files chosen → confirmation list in the TUI (name + size), with `[Send]` and `[Back]`.
4. Send → croc subprocess starts → **code shown as soon as croc prints it** (during the relay handshake) with hint `c = copy to clipboard`.
5. Live progress: animated bar (`bubbles/progress`) + bytes/s + status.
6. Done → "Transfer complete" + `enter = send another`, `esc = back`.

### Screen 3 — Receive

1. Text input for the code (paste-friendly, e.g. `Ctrl+V`).
2. Confirm → croc subprocess starts in `Downloads`.
3. Live progress + incoming file name.
4. Done → "Saved to Downloads" + `esc = back`.

### Copy affordance

A `c` key anywhere a code is displayed pipes the code to `clip` and flashes "Copied!". Windows-only, by design.

## Trade-offs & decisions

- **Shell out instead of embedding croc library.** croc's Go API is not a stable public contract, and embedding drags in the whole codebase + websocket/PAKE deps. Shelling out keeps this app a thin wrapper and behavior-identical to `croc` on the CLI. Cost: relies on the `croc` binary being installed (auto-install covers this) and on parsing CLI output.
- **Native dialog over in-terminal browser.** User chose the native Windows dialog. Cost: this feature is Windows-only (PowerShell + `System.Windows.Forms`); the rest of the app is portable. A future in-terminal `filepicker` component could make it cross-platform.
- **One transfer at a time.** Keeps the app simple; croc itself doesn't parallelize well from a single process anyway.

## Testing

- **Parser unit tests** (`croc/parse_test.go`): feed representative croc output lines (code line, progress frames, done/error) → assert parsed events. Pure functions, no TUI.
- **Runner tests**: fake croc binary (a stub script) exercised against the runner to verify argument construction and exit handling.
- **Manual QA checklist** (documented in the plan): real send/receive loop, cancel flows, missing-binary screen, clipboard copy.

## File layout

```
crocgui/
├── main.go
├── go.mod
├── ui/
│   ├── welcome.go
│   ├── send.go
│   ├── receive.go
│   └── progress.go
├── croc/
│   ├── runner.go
│   ├── parse.go
│   ├── parse_test.go
│   └── install.go
├── dialog/
│   └── dialog.go        (Windows PowerShell bridge)
├── clipboard/
│   └── clipboard.go     (clip bridge)
└── docs/superpowers/specs/2026-08-05-crocgui-design.md
```

## Milestones (build order)

1. Skeleton: `go mod init`, Bubble Tea app boots, welcome screen renders, clean quit.
2. `croc` adapter: runner + parser + unit tests (no UI).
3. Receive flow end-to-end (code input → subprocess → progress → done).
4. Send flow end-to-end (dialog → list → subprocess → code + copy → progress → done).
5. Missing-binary install screen + polish (colors, key hints, error states).
