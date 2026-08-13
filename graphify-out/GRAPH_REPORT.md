# Graph Report - D:\Apps\projects\crocgui  (2026-08-13)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 258 nodes · 510 edges · 16 communities (15 shown, 1 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 89 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `6ab766ed`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- github.com/charmbracelet/bubbletea.Cmd
- testing.T
- frame_test.go
- ParseLine
- package.json
- newHistoryModel
- dialog_test.go
- New
- newReceiveModel
- ui/install.go
- crocgui.js
- ui/runner.go
- root.go
- runner_integration_test.go
- crocgui

## God Nodes (most connected - your core abstractions)
1. `newSendModel()` - 26 edges
2. `RootModel` - 16 edges
3. `ParseLine()` - 12 edges
4. `receiveModel` - 11 edges
5. `sendModel` - 11 edges
6. `screenView()` - 10 edges
7. `installModel` - 9 edges
8. `EventWithMeta` - 9 edges
9. `files` - 8 edges
10. `newHistoryModel()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `crocEventMsg` --references--> `EventWithMeta`  [EXTRACTED]
  ui/runner.go → croc/parse.go
- `streamOutput()` --calls--> `ParseLine()`  [EXTRACTED]
  ui/runner.go → croc/parse.go
- `TestProgressScale()` --calls--> `newReceiveModel()`  [INFERRED]
  ui/progress_scale_test.go → ui/receive.go
- `TestAllScreensShowLogoAndFitWindow()` --calls--> `newSendModel()`  [INFERRED]
  ui/frame_test.go → ui/send.go
- `TestSendDirShowsFolderNotSize()` --calls--> `stripANSI()`  [INFERRED]
  ui/send_test.go → ui/frame_test.go

## Import Cycles
- None detected.

## Communities (16 total, 1 thin omitted)

### Community 0 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (23): Install(), context.CancelFunc, context.Context, github.com/charmbracelet/bubbles/progress.Model, github.com/charmbracelet/bubbles/spinner.Model, github.com/charmbracelet/bubbles/textinput.Model, github.com/charmbracelet/bubbletea.Cmd, github.com/charmbracelet/bubbletea.Model (+15 more)

### Community 1 - "testing.T"
Cohesion: 0.13
Nodes (31): testing.T, TestProgressScale(), TestSendCodeVisibleDuringTransfer(), TestSendConfirmEscGoesBack(), TestSendConfirmXRemovesLastFile(), TestSendDoneEnterStartsNewSend(), newSendModel(), TestSendAddMoreAppends() (+23 more)

### Community 2 - "frame_test.go"
Cohesion: 0.11
Nodes (21): stripANSI(), TestAllScreensShowLogoAndFitWindow(), TestFrameCentersVertically(), TestFrameFillsWidth(), TestFrameNoHiddenExtraHeight(), TestFrameZeroWidthPipesThrough(), TestScreenViewFitsWindow(), TestScreenViewRendersLogoHeader() (+13 more)

### Community 3 - "ParseLine"
Cohesion: 0.16
Nodes (21): Event, EventHistory, EventKind, EventMetadata, EventWithMeta, NewEventHistory(), ParseChunk(), ParseLine() (+13 more)

### Community 4 - "package.json"
Cohesion: 0.10
Nodes (19): bin, crocgui, description, engines, node, files, license, name (+11 more)

### Community 5 - "newHistoryModel"
Cohesion: 0.19
Nodes (11): histDir, itoa(), joinLines(), newHistoryModel(), TestHistoryEmptyState(), TestHistoryEscGoesBack(), TestHistoryRecordsEntry(), TestHistoryRecordsMultipleAndNavigates() (+3 more)

### Community 6 - "dialog_test.go"
Cohesion: 0.22
Nodes (12): dialogScript(), folderScript(), TestManualPick(), parsePaths(), PickFiles(), PickFolder(), TestFolderScriptOpensAtCwd(), TestFolderScriptUsesFolderBrowser() (+4 more)

### Community 7 - "New"
Cohesion: 0.18
Nodes (9): TestReceiveArgs(), TestSendArgs(), Detect(), ReceiveArgs(), SendArgs(), main(), New(), TestRootStartSendWithoutProgramShowsError() (+1 more)

### Community 8 - "newReceiveModel"
Cohesion: 0.23
Nodes (9): newReceiveModel(), contains(), findSubstring(), TestReceiveCodeTooShort(), TestReceiveEnterEmitsStartReceive(), TestReceiveErrorPrecedence(), recvState, startReceiveMsg (+1 more)

### Community 9 - "ui/install.go"
Cohesion: 0.22
Nodes (9): installProgressTick(), newInstallModel(), TestInstallDoneOnlyAdvancesOnEnter(), TestInstallEnterRunsInstall(), TestInstallViewPromptsForWinget(), installFinishedMsg, installStartedMsg, installState (+1 more)

### Community 10 - "crocgui.js"
Cohesion: 0.22
Nodes (6): BIN, fs, path, PKG_ROOT, res, { spawnSync }

### Community 11 - "ui/runner.go"
Cohesion: 0.22
Nodes (7): io.Reader, os/exec.Cmd, crocEventMsg, crocFinishedMsg, buildCrocCmd(), streamOutput(), TestBuildCrocCmd()

### Community 12 - "root.go"
Cohesion: 0.31
Nodes (6): dialogFailedMsg, goScreenMsg, errStr(), nowStamp(), screen, startSendMsg

### Community 13 - "runner_integration_test.go"
Cohesion: 0.73
Nodes (5): buildFakeCroc(), exeSuffix(), runCrocInProgram(), TestRunCrocNonBlocking(), TestRunCrocStreamsCodeAndFile()

## Knowledge Gaps
- **32 isolated node(s):** `sendFilesPickedMsg`, `welcomePulseMsg`, `BIN`, `fs`, `path` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `RootModel` connect `github.com/charmbracelet/bubbletea.Cmd` to `root.go`, `newHistoryModel`?**
  _High betweenness centrality (0.084) - this node is a cross-community bridge._
- **Why does `newSendModel()` connect `testing.T` to `github.com/charmbracelet/bubbletea.Cmd`, `frame_test.go`?**
  _High betweenness centrality (0.079) - this node is a cross-community bridge._
- **Why does `ParseLine()` connect `ParseLine` to `ui/runner.go`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Are the 24 inferred relationships involving `newSendModel()` (e.g. with `TestAllScreensShowLogoAndFitWindow()` and `TestProgressScale()`) actually correct?**
  _`newSendModel()` has 24 INFERRED edges - model-reasoned connections that need verification._
- **What connects `sendFilesPickedMsg`, `welcomePulseMsg`, `BIN` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `github.com/charmbracelet/bubbletea.Cmd` be split into smaller, more focused modules?**
  _Cohesion score 0.08888888888888889 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.13445378151260504 - nodes in this community are weakly interconnected._