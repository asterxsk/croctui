# Clickable TUI Design

**Date:** 2026-02-13

## Goal

Add clickable controls to the CrocGUI Bubble Tea TUI. The first mouse event automatically switches the presentation and interaction layer to mouse UI; the first keyboard event switches it back to keyboard UI. All existing screens and actionable states remain usable from the keyboard.

## Interaction architecture

Bubble Tea mouse reporting is enabled with `tea.WithMouseAllMotion()`. `RootModel` owns an `inputMode` value with keyboard and mouse states. On every `tea.MouseMsg`, the root enters mouse mode and forwards the event to the active screen. On every `tea.KeyMsg`, the root enters keyboard mode before forwarding the key to the active screen.

A small local component layer provides reusable clickable buttons and list rows. Components own their rendered bounds, hover state, disabled state, and conversion of left-button releases into semantic actions. Screen models continue to own transfer state machines and keyboard handlers; mouse actions are translated into those same existing transitions and commands instead of duplicating subprocess logic.

If the terminal does not emit mouse messages, the application remains keyboard-driven. No startup capability probe or keybind is required.

## Screen behavior

- Welcome has Send, Receive, and History buttons.
- Send confirmation has clickable file rows, Add files, Add folder, Delete selected, Send, and Back. Empty-list Send is disabled; wheel events move file selection. Waiting/transferring retain a clickable Back/Cancel action. Done exposes Send another and Back; error exposes Back.
- Receive input keeps the existing keyboard-editable text input and adds clickable Start and Back. Waiting/transferring expose Back/Cancel. Done and error expose Back.
- History has clickable rows, wheel selection, and Back. Empty history has Back only.
- Install exposes clickable install/retry and quit/back actions according to state.

Mouse mode replaces verbose keyboard footer hints with concise mouse affordances. Controls use the existing square-border/accent design language, with padded hit targets and dim disabled states. Keyboard handlers and semantics remain intact when switching back.

## Data flow

1. Bubble Tea emits `tea.MouseMsg` or `tea.KeyMsg`.
2. `RootModel.Update` updates the input mode and dispatches the message.
3. The active screen/component updates hover or selection state.
4. A mouse release inside an enabled component emits a semantic action or invokes the screen's existing transition path.
5. Root-owned screen switching, subprocess launching, cancellation, and history recording remain unchanged.

## Testing

Unit tests cover component bounds, hover, left-release activation, disabled controls, wheel selection, root mode switching and forwarding, representative mouse actions on every screen/state, and mouse-mode versus keyboard-mode rendering. The full suite runs with `go test ./...`, plus `go vet ./...` and `go build ./...`.

A Windows Terminal smoke test verifies mouse reporting, every screen's clickable actions, switching back to keyboard by typing, and preservation of native dialogs and croc subprocess transfers.

## Constraints

- Keep Go 1.26.5, Bubble Tea v1.3.10, lipgloss v1.1.0, and the existing dependency set.
- Keep the application keyboard-capable; mouse mode is automatic and does not remove key handling.
- Do not embed croc transfer logic or alter the existing subprocess contract.
- Keep history in-session only and preserve the keyboard-only/no-mouse-capture behavior outside this feature's explicit mouse reporting option.
- Maintain Windows-first behavior while keeping component and root tests portable.
