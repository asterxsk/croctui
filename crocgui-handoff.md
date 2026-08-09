# CrocGUI — Mouse clicks + UI spacing handoff

**Date:** 2026-08-07
**Repo:** `D:\Apps\projects\crocgui` (module `crocgui`, Go + Bubble Tea v1.3.10 + lipgloss v1.1.0, Windows-first)
**Branch:** `master` (only branch; **no worktrees exist right now**)

## What the next session must do

The user's standing request (from the previous session): **mouse clicks still don't work in a real terminal, and there are "quite a few UI spacing and size issues" — fix both, ideally in parallel via subagents.** No progress was made on this last round; read this doc, then start fresh.

## Current state of the working tree

`go build ./...`, `go vet ./...`, and `go test ./...` **all pass right now**. There is a large body of **uncommitted** work (see `git status`):

- **Modified:** `dialog/dialog.go`, `dialog/dialog_test.go`, `main.go`, `ui/history.go`, `ui/install.go`, `ui/receive.go`, `ui/receive_test.go`, `ui/root.go`, `ui/send.go`, `ui/send_fix_test.go`, `ui/send_test.go`, `ui/theme.go`, `ui/welcome.go`, `ui/welcome_test.go`
- **Untracked:** `ui/frame_test.go`, `ui/mouse.go`, `ui/mouse_test.go`, `ui/root_test.go`

Nothing is committed. If a fresh agent is unsure what's baseline vs. new, `git diff` + the untracked files tell the story. The last commit (`eae54b8`) is "feat: polished full-window TUI with accent colors, borders, and history".

## What has already been tried (do NOT re-litigate, but DO verify)

All of this is in the current working tree and is covered by passing tests:

1. **`Click()` no longer requires `Action == Press`.** On Windows Terminal/ConPTY a left-button press can arrive as `Press`, `Release`, or `Motion` (Bubble Tea's own comment: "Motion can be reported as a release event in some terminals"). `ui/mouse.go:Click` now matches **any** event carrying `MouseButtonLeft`; screens dedupe via `clickStart(&held, msg)` / `releaseClick(&held, msg)` so one physical click fires once.
2. **Fixed-size footers so layout never shifts between mouse/keyboard modes** (a shifted layout moves hit boxes out from under the cursor → clicks miss). `footerPad`/`footerChips` are fixed 70×3; send uses `footerLine`/`chipsLine` fixed 70×1. Pinned by `TestFooterModesSameSize`, `TestSendFooterModesSameSize`, `TestFooterChipsCentered`.
3. **Screens capture the view BEFORE flipping `mouse=true`** (chips footer shifts the centered layout), per `ui/send.go:168` `view := m.View()`.
4. **`scanHitBox` is ANSI-aware** and returns 0-based display-cell coords to match Bubble Tea's `MouseMsg.X/Y`.
5. **Chip centering** was fixed so chips render centered, not left-aligned.

**Despite all this, the user reports clicks still don't fire in the real terminal.** The most likely remaining causes (hypotheses to test):
- The actual event shape on the user's terminal differs from every assumption (e.g. `Button == MouseButtonNone` on Motion, or coordinates offset by the window border / padding).
- Bubble Tea's `WithMouseCellMotion()` on Windows may need the mode to be re-armed after a screen switch, or the coordinate space differs between SGR and ConPTY records.
- The hit-box scan might be off-by-one relative to what the terminal actually renders (border padding vs. cursor position).

## The critical next step (do this FIRST)

**Get ground truth on the real event stream.** Set `CROCGUI_MOUSE_DEBUG=1`, run `bin\crocgui.exe` (or `go run .`), click around, and read `crocgui_mouse.log` (written by `ui/mouse.go:mouseDebugLog` via `ui/root.go:214`, format: `[screen] X=.. Y=.. Action=.. Button=..`). That log reveals the exact Action/Button/X/Y shape the user's terminal delivers. **The previous session's fix was made without this log** — that is the gap. Do not guess again; instrument first.

Also compare the mouse coordinates in the log against what the UI actually renders: for a clicked button, cross-check `scanHitBox` output (write a tiny test or add a debug print) against the logged X/Y. Suspect a systematic offset (padding/border consumes cells that the hit-box math may double-count).

## UI spacing / sizing issues

The user says there are "quite a few UI spacing and size issues" but didn't enumerate them. Suspects worth auditing (from the code as read):
- `frame()` inside-out width math (border 2 cols + `Padding(0,1)` gutter 2 cols) — see `ui/theme.go:179`. Content is placed at `Width(width-2)` then padded, which may produce a **ragged right edge** or an off-by-two centered box. Verified-lipgloss-behavior notes live in memory `[[full-window-frame]]`.
- The 70-col fixed footers vs. a narrower/wide window: if the terminal is < 74 cols, `footerPad`/`footerChips` overflow.
- Send screen (content-heavy): the 24-row window may clip the list + buttons + footer at small heights (`TestAllScreensShowLogoAndFitWindow` only asserts `<= 24`, not that content is visible).
- Vertical centering interplay between `frame()`'s `Place` and `screenView` — content can end up top-cramped or bottom-cramped depending on height.
- Panel margins (`Margin(1,2)` on `panel`/`codeBox`) vs. the frame gutter — double-gutter / inconsistent alignment between screens.

Run the app (`go run .`) at several terminal sizes (e.g. 80×24, 100×30, 60×20, and a tall one) and screenshot/describe each screen; fix obvious misalignments while preserving the existing design language and all tests.

## Prior design language (do not break)

Full detail is in `C:\Users\asterxsk\.claude\projects\D--Apps-projects-crocgui\memory\crocgui-windowed-tui.md` (design language + the 9 "CRITICAL mouse gotchas") and `...\memory\full-window-frame.md` (frame width math). Summary:
- `accent` = color 39 blue for wordmark/borders/selection; square `NormalBorder`; `selectedButton` = bold black-on-blue block (border cells included); `idleButton` = bordered gray.
- Every screen renders the **borderless** CROC ASCII-art header (`crocHeader`/`logoBox`) + caption via `screenView(w,h,caption,content)`; `frame()` wraps in the full square window border.
- Screen captions: welcome "CROC — easily & securely transfer files…", send "send files to any device", receive "receive files from any device", history "your transfer history", install "set up croc".
- Send screen: Enter = interact (always opens add dialog), space/Send button = send. `sendConfirm` list navigable by up/down/w/s/j/k, x/backspace deletes, a adds files, A (Shift+a) adds folder.

## Plan / reference artifacts (read, don't duplicate)

- Implementation plan (from earlier plan mode): `C:\Users\asterxsk\.claude\plans\dynamic-meandering-quiche.md` — send-page interaction + app-wide mouse support architecture.
- `ui/mouse.go` — HitBox/Click/Wheel/clickStart/releaseClick + debug log.
- `ui/root.go` — screen enum, dispatch, `tea.MouseMsg` forwarding (line ~211).
- `ui/theme.go` — frame/footer/chip/scanHitBox helpers.
- `ui/mouse_test.go`, `ui/frame_test.go`, `ui/send.go` (sendList/footer), `main.go` (`tea.WithMouseCellMotion()`).
- Fake croc binary for tests: `ui/fake_croc.go`; runner tests in `ui/runner_test.go`.

## Suggested skills for the next session

- `systematic-debugging` — hard bug (clicks not firing): hypothesis → instrument (`CROCGUI_MOUSE_DEBUG=1`) → verify → fix. Do this before touching code.
- `impeccable` — UI spacing/sizing audit and polish pass across all screens (this is exactly its remit).
- `dispatching-parallel-agents` — user explicitly wants parallel subagents; mouse + spacing are independent and can run concurrently.
- `using-git-worktrees` — keep the two workstreams isolated if running agents in parallel.
- `executing-plans` — the existing plan (`dynamic-meandering-quiche.md`) describes the intended mouse architecture.
- `tdd` / `test-driven-development` — every mouse fix so far shipped with a passing regression test; keep that discipline.

## Constraints / environment notes

- Windows 11, PowerShell primary. Repo must build with `go build ./...` and pass `go test ./...` before/after changes.
- Binary output at `bin\crocgui.exe` (rebuilt via `go build -o bin\crocgui.exe .`).
- Redacted: none — no secrets in this project.
