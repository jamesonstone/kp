---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0011"
  slug: "picker-tui"
  dir: "0011-picker-tui"
references:
  - id: "github-issue-66"
    name: "Redesign the kp picker as a native TUI and add the kp v0 namespace"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/66"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "v0-init-utility"
    name: "v0 init utility"
    type: "feature_artifact"
    target: "docs/specs/0001-v0-init-utility/SPEC.md"
    relation: "constrains"
    read_policy: "must"
    used_for: "bare-command model, launcher, clipboard verification, exit codes, farewells"
    status: "active"
  - id: "init-command"
    name: "init command"
    type: "feature_artifact"
    target: "docs/specs/0010-init-command/SPEC.md"
    relation: "constrains"
    read_policy: "should"
    used_for: "Init launcher row and questionnaire entry"
    status: "active"
delivery_intent: "issue_branch_pr_ready"
---
# SPEC

## PURPOSE

Replace the `fzf`-driven picker behind `kp` and `kp list` with a small,
purpose-built two-pane TUI that keeps the same interaction, and move ten
legacy prompt commands under a `kp v0` compatibility namespace so their root
names can be reused.

## CONTEXT

- Before this change, `kp` built tab-separated rows with emoji, wrote every
  preview body to a temporary directory, and ran `fzf` with a `cat` preview.
  `kp list` ran a second, differently configured `fzf` invocation. Row
  alignment needed a hand-rolled emoji width table.
- Domain behavior (registry loading, `runPrompt`, clipboard copy plus verify,
  Init and Find port) already sat outside the `fzf` code. Presentation,
  selection parsing, preview transport, and cancellation classification were
  coupled to `fzf` exit codes and tab-delimited output.
- Tests injected `FZFRunner` and `LauncherRunner` hooks, so the picker was
  already a replaceable function from items to a selected ID.

## REQUIREMENTS

- `kp` and `kp list` open the same two-pane picker: list left, live preview
  right. Selection immediately updates the preview.
- Keys: `j`/`k`, arrows, Tab/Shift-Tab, Ctrl-N/Ctrl-P move and wrap at both
  ends; Enter selects; Esc clears an active filter, then cancels; Ctrl-C
  cancels; any other printable key (or `/`) starts a filter, as typing did in
  the `fzf` launcher; Home/End jump; Ctrl-D/Ctrl-U, PgDn/PgUp, and
  Shift-Up/Shift-Down scroll the preview.
- The TUI owns pane sizing, rendering, selection, scrolling, the preview
  viewport, keys, resize, and styles. Prompt discovery, loading, clipboard,
  and execution stay in `internal/cmd`.
- No emoji in the picker and no Nerd Font glyphs. Color comes only from the
  terminal's basic ANSI palette, so each theme shades it for its own
  background: magenta for the `kp` badge, pointer, selected row, and preview
  title; cyan for command names and key hints; blue for group headings and
  Markdown headings; yellow for tags, match counts, and scroll position;
  faint for secondary text. `NO_COLOR` leaves bold and faint only.
- Each list row shows its label and, in an aligned column, the command that
  runs it directly (`kp review`, `kp init`), with the subcommand name
  highlighted. Items sit under `prompts` and `commands` headings. The column
  drops only when titles would get too narrow. User prompts carry a
  `user prompt` tag in the preview.
- The preview tints Markdown structure (headings, list markers, inline and
  fenced code) without changing text or width. It is not a renderer.
- `kp init` on a TTY opens the template in the user's editor (`$KP_EDITOR`,
  `$EDITOR`, `nvim`, then `vi`) instead of the raw-mode questionnaire; see
  `0010-init-command`. Piped stdin keeps one line per section.
- The preview wraps to at most 88 columns with hanging indents for list items,
  expands tabs, strips control characters, and shows a scroll position only
  when the body overflows.
- Narrow terminals drop the preview instead of squeezing both panes. Every
  rendered line fits the terminal width and the frame fits its height.
- The picker draws on `/dev/tty`, so stdout carries only the selected output.
  With no terminal it exits `3` with an instruction. Cancel exits `130` with
  the existing farewell. `--no-fzf` keeps the numbered fallback.
- Direct commands must not query the terminal or slow down.
- `kp v0 <command>` serves `agent-handoff`, `chat-handoff`, `clarify`,
  `continue`, `goal`, `parentthread`, `plan`, `pr`, `punchlist`, and `status`
  with byte-identical bodies and the same `--print`/`--copy` behavior. `v0`
  is reserved. Root `kp <old-name>` reports the move unless a user prompt now
  owns that name. The launcher and `kp list` show only current root prompts.

## ACCEPTANCE CRITERIA

- AC1: Launcher items are the prompts, with `kp init` ("Start a task from
  scratch") among them in name order, then Find port and Help,
  with no emoji, each carrying its command, grouped under headings.
- AC7: `kp init` on a TTY preloads the template in the editor, keeps
  multi-line answers and spacing, reopens incomplete drafts with a notice,
  cancels on an empty file or editor abort, and prints only the prompt to
  stdout.
- AC2: Navigation, wrap, filter, Esc, Ctrl-C, Enter, and empty-result
  behavior match the key requirements (pure model tests through `Update`).
- AC3: Layout math, list windowing on resize, preview scroll clamping, and
  width/height bounds hold across terminal sizes.
- AC4: Selecting from `kp` or `kp list` prints and copies exactly the selected
  prompt; cancel exits `130` with no side effects; no terminal exits `3`.
- AC5: `kp v0 <name> --print` equals each moved file's body; copy and verify
  behave as before; root old names fail with a pointer; user prompts can
  reuse the names; `kp v0 <unknown>` fails; `kp new v0` is rejected.
- AC6: Direct-command startup is unchanged and sends no terminal queries.

## DECISIONS

- Adopt Bubble Tea v2 and Lip Gloss v2. Extending `fzf` could not remove the
  temp-file preview transport, emoji width tables, tab-delimited parsing, or
  exit-code guessing, and cannot style a preview header or control wrapping.
  A hand-written raw-mode loop would reimplement input decoding, resize, and
  diff rendering. Bubble Tea v1 was tried first and rejected: its package
  `init` sends an OSC 11 background query on every `kp` run, which stalls on
  terminals that do not answer. v2 has no such `init`; it raises the Go
  floor to 1.26.
- Do not use Bubbles. Its list and viewport components bring their own key
  maps, pagination, filtering UI, and status chrome; the picker needs about 60
  lines of list windowing and scroll clamping.
- Use substring filtering, not subsequence. Subsequence matched `rev` in
  "Pre-authorize task delivery" and, without ranking, reads as noise.
- A first pass used one accent and kept commands out of the list. Review
  feedback: it read as too plain, and people recognize actions by command
  name. The palette widened to basic ANSI roles and commands joined the list
  as an aligned column.
- Replace `kp init`'s raw-mode Shift+Enter input with the user's editor:
  spacing is easier to write and read, and the CSI key decoding goes away.
- Keep the whimsical cancel farewells. They print to stderr after the picker
  closes and belong to 0001's contract, not the picker's presentation.
- Keep the `--no-fzf` flag name for script compatibility; its help text now
  names the interactive picker.
- Serve `v0` prompts from `prompts/v0/` as Cobra subcommands so help,
  completion, and unknown-name errors come from Cobra. They are built-in only,
  not user-overridable; a user file named `clarify.md` now owns root
  `kp clarify`.
- Record topology as `single-lane, because tightly coupled and requires
  continuous design judgment: one picker package, its two call sites, the
  prompt split, and tests share state and visual decisions`.

## TASK CHECKLIST

- [x] Move ten prompts to `prompts/v0/`; add `V0BuiltIns`, `kp v0`, the
  reserved name, the root not-found pointer, and a help section.
- [x] Add `internal/picker` (state, keys, layout, view) on Bubble Tea v2.
- [x] Route `kp` and `kp list` through one `pick` helper; remove `fzf` code.
- [x] Update and add tests for AC1–AC6.
- [x] Update README, Constitution dependencies, testing reference, and the
  progress summary.
- [x] Independent review for unnecessary visuals and complexity.
- [x] Widen the palette, add the command column and group headings, and tint
  preview Markdown.
- [x] Move `kp init` on a TTY to an editor draft; remove the raw-mode input.
- [x] Run format, test, race, vet, build, and pty acceptance.

## VALIDATION

- AC1: `TestRootLauncherListsCurrentActions`, `TestLauncherItemsHaveNoEmoji`,
  `TestLauncherMarksUserPromptsAsSecondaryDetail`.
- AC2: `internal/picker/update_test.go` drives every key through `Update`.
- AC3: `TestComputeLayout`, `TestRenderFitsEveryTerminalSize` (3x2 through
  200x50, with and without a 120-character query), `TestResizeKeepsSelectionVisible`,
  `TestPreviewScrollClampsAndResetsOnMove`, wrap and sanitize tests.
- AC4: `TestListPickerSelection`, `TestRootLauncher*`,
  `TestPickerWithoutTerminalExitsConfig`, plus pty acceptance.
- AC5: `internal/cmd/v0_test.go`, moved prompt hash tests, registry and
  built-in tests.
- AC6: startup timing and pty terminal-query check below.
- AC7: `internal/cmd/init_editor_test.go` plus a real `nvim` pty session.

## DELIVERY DECISION

Issue #66, branch `GH-66`, canonical worktree
`~/worktrees/jamesonstone/kp/GH-66`, ready pull request to `main`.

## REVIEW

An independent review pass looked for unnecessary visuals, unnecessary
complexity, and bugs. Applied: header query truncation (it could overflow the
width), compact chrome below 8 rows with a final per-line width guard, a wrap
fix that no longer inserts a space after hyphen breaks, removal of an unused
layout field and a redundant `canceled` flag, one filtering predicate, a
cached longest title, the root command path in both call sites, a single
legacy help row instead of ten, and a README note that `v0` serves built-ins
only. Kept deliberately: the top padding row (dropped in compact mode), the
scroll percentage (the requested scroll-position cue, shown only on
overflow), the bold preview title (it shows the full label when the list
truncates), and `›`/`│`/`…` (ambiguous width only in CJK-width terminals).

## EVIDENCE

- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./...` — `PASS`. `go test -race ./...` — `PASS`.
- `go vet ./...` — `PASS`. `go build ./...` — `PASS`. `make build` — `PASS`.
- `git diff --check` — `PASS`. Source-file-size audit — `PASS` (no changed Go
  file above 300 lines).
- Isolated CLI acceptance (empty `--config`) — `PASS`: `kp review --print`
  and `kp v0 clarify --print` print their bodies; `kp list --plain` prints
  `merge review ship`; `kp clarify` exits `1` with
  `prompt not found: clarify (moved to "kp v0 clarify")`; `kp v0 review`
  exits `1` with `unknown command`.
- Pty acceptance (Python `pty` plus `pyte`, `xterm-256color`) — `PASS`:
  launcher at 110x30, 100x24, 64x16, and 44x6; j/k, arrows, Ctrl-D/Ctrl-U,
  typed and `/` filters, Esc clears then exits `130`, Ctrl-C exits `130`,
  resize reflows, `kp list` plus Enter wrote exactly the prompt body to
  redirected stdout and exited `0`, and the picker still works with stdin
  redirected. No controlling terminal exits `3` with an instruction.
- Terminal queries — `PASS`: `kp merge --print` in a pty sends no OSC 11 or
  cursor-position query. First picker frame in about 22 ms.
- Startup — `PASS`: `kp review --print` over two rounds of 100 runs:
  3.46 and 3.69 ms per run against 3.76 and 3.21 ms for the previous `fzf`
  build (noise-level difference).
- Color — `PASS`: SGR codes limited to bold, faint, reverse, and basic ANSI
  foregrounds `33`–`36`; `NO_COLOR=1` leaves only bold, faint, and reverse.
  Dark and light theme renders of the pty screen were reviewed visually.
- `kp init` editor — `PASS`: real `nvim --clean` in a pty opened the draft
  on the first answer line; saving one section reopened the file with a
  notice listing the four empty sections and kept the text; deleting
  everything cancelled with empty stdout. A scripted editor completing every
  section produced exactly the rendered prompt on redirected stdout, with
  blank lines and indented bullets intact.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has
  no hosted format, test, race, vet, or build workflow.
- Real-terminal visual check on light and dark themes — `PENDING` for the
  human reviewer; the pty run proves structure and SGR usage, not appearance.
- Production validation — `NOT_APPLICABLE`: local CLI.
