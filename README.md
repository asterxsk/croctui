# croctui 🐊 — CrocGUI

Keyboard-driven terminal GUI for [croc](https://github.com/schollz/croc) — peer-to-peer, end-to-end encrypted file transfer (Go + Bubble Tea).

- **Send** — native Windows file picker → share code (auto-copied) → live progress
- **Receive** — paste code → files to Downloads → live progress
- **Auto-install** — offers `winget install schollz.croc` if missing

No mouse — keyboard only (`WithAltScreen`, no mouse capture).

## npx (no Go needed)

```bash
npx @asterxsk/croctui
# or pinned
npx @asterxsk/croctui@0.1.0
```

Zero runtime deps. The postinstall fetches a prebuilt binary from GitHub Releases (`croctui-<platform>-<arch>[.exe]`); the 8.5 kB npm tarball contains only `bin/croctui.js` + `bin/install.js`. First run downloads once, then runs offline. If no release is present and Go exists, it builds from source.

Also aliased: `npx @asterxsk/croctui` == `croctui` == `crocgui` binaries.

## From source (dev)

```bash
go run .
go build -o croctui.exe .
go test ./...
```

Requires Go 1.26+ only for development.

## Requirements

- Node 18+ for `npx` (no Go needed for users)
- `croc` CLI — auto-installed via winget on first launch if missing
- Windows-first (file dialog + winget are Windows-only; rest is portable)
- Real TTY (Windows Terminal recommended)

## How it works

Shells out to the real `croc` CLI — no embedded transfer logic. `--pass` is the relay password, not the code; croc generates and copies the code.
