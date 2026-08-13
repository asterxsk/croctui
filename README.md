# croctui 🐊

[![npm version](https://img.shields.io/npm/v/@asterxsk/croctui?color=cb3837&label=npm)](https://www.npmjs.com/package/@asterxsk/croctui)
[![npm downloads](https://img.shields.io/npm/dm/@asterxsk/croctui)](https://www.npmjs.com/package/@asterxsk/croctui)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)](https://go.dev)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v1.3-ff6b6b)](https://github.com/charmbracelet/bubbletea)
[![platform](https://img.shields.io/badge/platform-windows%20%E2%80%A2%20linux%20%E2%80%A2%20macOS-lightgrey)](https://github.com/asterxsk/croctui/releases)
[![license](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

**Keyboard-driven terminal GUI for [`croc`](https://github.com/schollz/croc) — send anything, anywhere, peer-to-peer.**

```sh
npx @asterxsk/croctui
```

No Go required. `2.8 kB` npm package — prebuilt binary fetched on install, then runs offline.

```
 ██████╗██████╗  ██████╗  ██████╗
██╔════╝██╔══██╗██╔═══██╗██╔════╝
██║     ██████╔╝██║   ██║██║
██║     ██╔══██╗██║   ██║██║
╚██████╗██║  ██║╚██████╔╝╚██████╗
 ╚═════╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝
      securely transfer files peer-to-peer
```

---

## Features

- **Send** — native file picker (multi-select) → share code (auto-copied by `croc`) → live progress
- **Receive** — paste code → files land in `Downloads` → live progress
- **History** — in-session record of every transfer
- **Auto-install** — offers `winget install schollz.croc` if `croc` is missing
- **Keyboard-only** — no mouse capture, no click targets. Arrows / `hjkl` / `tab` / `enter` / `esc`.
- **Full-window TUI** — single accent (`#39` blue), square borders, framed layout that fills the terminal
- **Zero bloat** — shells out to the real `croc` CLI. No embedded crypto, no relay, always behavior-compatible.

## Quick start

### npx (recommended — no Go needed)

```bash
npx @asterxsk/croctui          # latest
npx @asterxsk/croctui@0.1.0     # pinned

# also aliased as crocgui
npx @asterxsk/croctui --help
```

The npm package is `2.8 kB` (unpacked `8.2 kB`). On first run it downloads a prebuilt binary (`croctui-<platform>-<arch>[.exe]`, ~5 MB, stripped) from [GitHub Releases](https://github.com/asterxsk/croctui/releases) and caches it in `bin/`. Subsequent runs are offline. If no release is available and Go is installed, it builds from source.

### Global install

```bash
npm i -g @asterxsk/croctui
croctui
# or
crocgui
```

### From source (development)

```bash
git clone https://github.com/asterxsk/croctui.git
cd croctui
go run .
# or
go build -o croctui.exe .
./croctui.exe
go test ./...
```

## Usage

```
Welcome:  s Send  ·  r Receive  ·  t History  ·  arrows move  ·  enter confirm  ·  q quit

Send:     enter Add files  ·  shift+a Add folder  ·  ↑/↓ move  ·  x delete  ·  space Send  ·  esc Back
          → code boxed on screen: "Code: 1234-5678-9012  (code copied to clipboard by croc)"
Receive:  type/paste code  ·  enter Start  ·  esc Back
History:  ↑/↓ move  ·  esc Back
Install:  enter Install  ·  q Quit
```

All transfers stream live — percent, rate, file name — parsed from `croc` stdout.

## How it works

```
Pick files → croc send <files> --yes → code printed → parser → TUI progress
Enter code → croc --yes --rename <code> --out <Downloads> → parser → TUI progress
```

`croctui` **shells out** to the real `croc` binary. There is no embedded transfer logic — every flag and behavior is `croc`'s own. `--pass` is the relay password (default `pass123`), not the share code; `croc` generates and auto-copies the code.

Parser: `croc/parse.go` regex-parses `Code is:`, `Sending`/`Receiving`, `%` and `B/s` lines into typed `Event`s. Runner: `ui/runner.go` pipes `stdout+stderr` via `io.Pipe` without blocking the Bubble Tea loop.

## Requirements

| | Users (`npx`) | Developers |
|---|---|---|
| Node | 18+ | 18+ |
| Go | not needed | 1.26+ |
| `croc` | auto-installed via `winget` | `winget install schollz.croc` |
| Terminal | Real TTY (Windows Terminal recommended) | same |
| OS | Windows 11 primary · Linux / macOS work (dialog is Windows-only) | same |

Prebuilt binaries: `windows-amd64`, `windows-arm64`, `linux-amd64`, `linux-arm64`, `darwin-amd64`, `darwin-arm64` — built with `-ldflags="-s -w"` on every `v*` tag.

## Tech stack

- **Go 1.26** · **Bubble Tea v1.3** · **lipgloss v1.1** · **bubbles**
- Single static binary, minimal deps, `go vet` / `go test` clean
- Native Windows dialog via PowerShell / WinForms (`dialog/`), `winget` installer (`croc/install.go`)

## Project structure

```
croctui/
├── main.go            # Detect croc, run tea.NewProgram(WithAltScreen)
├── bin/
│   ├── croctui.js     # npx entry — fetch-or-build, spawn binary
│   └── install.js     # postinstall fetcher
├── croc/              # CLI adapter: Detect, SendArgs/ReceiveArgs, ParseLine, Install
├── dialog/            # Windows file/folder picker
├── ui/                # Bubble Tea screens: welcome / send / receive / history / install
├── docs/
└── .github/workflows/release.yml  # v* → 6-platform build + Release
```

## Credits

The transfer engine is [**`croc`**](https://github.com/schollz/croc) by [**Zack Schollz (@schollz)**](https://github.com/schollz) — a peer-to-peer, end-to-end encrypted file transfer tool and the de facto standard for CLI file sharing. `croctui` is a thin, keyboard-driven interface over it; all crypto, relay, PAKE and transfer logic is `croc`'s.

If `croctui` helps you, please star [`schollz/croc`](https://github.com/schollz/croc).

Other thanks: [Charm](https://charm.sh) for Bubble Tea, lipgloss and bubbles.

## License

MIT — see [LICENSE](LICENSE) (same as `croc` when applicable for vendored behavior).

---

<p align="center">
  <sub>Built by <a href="https://github.com/asterxsk">@asterxsk</a> · Not affiliated with <a href="https://github.com/schollz">schollz</a> · PRs welcome</sub>
</p>
