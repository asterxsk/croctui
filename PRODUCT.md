# Product

<!-- impeccable:product-schema 1 -->

## Platform

terminal (TUI — Bubble Tea, Go)

## Users

Developers and technical users who need to send files to another computer peer-to-peer, over the internet (not just LAN), with end-to-end encryption, and prefer a terminal workflow over a heavy desktop app. The user's real scene: a terminal session (Windows Terminal recommended), where a lightweight, keyboard-driven tool that "just works" is the point.

## Product Purpose

CrocGUI is a lightweight terminal GUI for the `croc` file-transfer CLI. It lets anyone send files to any other computer — securely, end-to-end encrypted, over the internet — with a friendly, discoverable interface. Success: a user picks files, gets a share code, shares it, and the transfer completes with live progress — without learning croc's CLI flags.

## Positioning

The transfer engine is `croc` itself (the de facto standard for CLI peer-to-peer file transfer); CrocGUI's mechanism is the *interface*: a polished, keyboard-driven terminal GUI that hides croc's complexity behind Send/Receive flows. It is the same engine, made approachable.

## Operating Context

- Runs in a real terminal (TUI; needs a TTY). Windows Terminal recommended.
- Windows-first: native file dialog (PowerShell/WinForms) and `winget` auto-install of croc are Windows-only; the rest is portable.
- `croc` binary is required — the app offers to install it via winget if missing.
- Keyboard-driven: send/receive/quit flows, no mouse needed.
- One transfer at a time, foreground, with live progress.
- History is in-session (not persisted to disk).

## Capabilities and Constraints

- Send flow: native multi-select dialog → confirm file list → croc subprocess → share code (boxed) → live progress → done.
- Receive flow: enter a code → croc subprocess → live progress → files land in Downloads.
- Install screen: if croc is missing, offer `winget install schollz.croc`.
- History screen: in-session record of completed transfers.
- Shells out to the real `croc` CLI — no embedded transfer logic, always behavior-compatible.
- `--pass` is the relay password (default `pass123`), not the share code; croc generates and auto-copies the code.
- Terminal UI must fill the full terminal window (user requirement).
- Go + Bubble Tea + lipgloss + bubbles. Single static binary, minimal deps.

## Brand Commitments

- App name: **CrocGUI**; the crocodile/croc branding of the underlying `croc` tool is the natural identity anchor.
- The user asked for: accent colors, borders, a history feature, and a layout that fills the full shell window. (User-stated, confirmed.)

## Evidence on Hand

- Working TUI already exists (welcome/send/receive/install screens, 29 passing tests). The visual layer is being enriched, not rebuilt.
- Real croc 10.6.0 CLI output captured during planning; the parser is verified against it.
- Manual QA items (real 2-device transfer, real winget install, quit-during-transfer) documented but not yet run headlessly.

## Product Principles

- **The tool disappears into the task**: send and receive should feel immediate, not require navigating a maze.
- **Engine truth**: never reimplement croc's logic — always shell out to the real CLI.
- **Lightness**: single static binary, minimal dependencies, keyboard-first.
- **Honest progress**: show real transfer state (code, percent, rate), never fake it.
- **One transfer, foreground**: keep the model simple; no queueing.

## Accessibility & Inclusion

- Keyboard-only operation (no mouse required).
- Terminal-native: colors should remain legible in common terminal themes; avoid relying on color alone for meaning (icons/text accompany color cues where it matters).
