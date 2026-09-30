---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0010"
  slug: "init-command"
  dir: "0010-init-command"
references:
  - id: "github-issue"
    name: "Add kp init prompt questionnaire"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/60"
    relation: "supports"
    read_policy: "must"
    used_for: "original issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "github-issue-62"
    name: "kp init multiline Shift+Enter and emoji section headers"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/62"
    relation: "supports"
    read_policy: "must"
    used_for: "follow-up issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "github-issue-64"
    name: "Put kp init on the launcher and fix its questionnaire layout"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/64"
    relation: "supports"
    read_policy: "must"
    used_for: "follow-up issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "v0-init-utility"
    name: "v0 init utility"
    type: "feature_artifact"
    target: "docs/specs/0001-v0-init-utility/SPEC.md"
    relation: "constrains"
    read_policy: "must"
    used_for: "bare-command model, clipboard verification, stderr interaction, reserved names, help, exit codes"
    status: "active"
delivery_intent: "issue_branch_pr_ready"
---
# SPEC

## PURPOSE

Add `kp init` as a small interactive questionnaire that turns five pieces of
human knowledge into a complete coding-agent prompt, prints that prompt to
stdout, and copies the identical body to the macOS clipboard.

## CONTEXT

- `kp` already has dedicated Cobra commands for interactive utilities
  (`list`, `new`, `find-port`, `scaffold`) and a separate registry of stored
  prompts invoked as bare names (`clarify`, `goal`, `review`).
- `kp init` is not a stored prompt. It constructs a prompt at runtime from
  operator input.
- Clipboard copy already uses `pbcopy` plus exact `pbpaste` verification.
  Interactive line input already lives on stderr through `app.inputReader`.
- Cancellation already maps EOF and empty picker input to exit `130` via
  `errPickerCanceled`.
- Feature `0001-v0-init-utility` is the prior `init` name in this repo. That
  feature built the CLI itself. This feature adds the `kp init` command.

## REQUIREMENTS

- `kp init` starts a sequential questionnaire covering exactly:
  1. 🎯 Objective
  2. 🧭 Known business/domain context
  3. 🔒 Invariants
  4. 🚧 Constraints
  5. ✅ Definition of done
- Each step asks the user to enter one or more sentences. Answers have no
  arbitrary practical length limit beyond ordinary stdin reading.
- On a TTY, Shift+Enter inserts a newline in the current section and Enter
  submits that section. Piped or non-TTY stdin stays one line per section.
- Empty or whitespace-only answers are rejected until a non-empty sentence is
  provided. EOF or other cancellation uses the existing exit `130` path and
  must not copy to the clipboard.
- After answers are collected, render this body and no additional commentary
  beyond the emoji prefixes on section headers:

```text
🎯  Objective

<objective>

🧭  Known business/domain context

<context>

🔒  Invariants

<invariants>

🚧  Constraints

<constraints>

✅  Definition of done

<definition of done>

Independently investigate. Do not assume my suspected implementation or root cause is correct.
```

- On success, write that exact body to stdout and copy it through the existing
  clipboard `Copy` plus `Verify` path. Stdout and clipboard bytes must match.
- `kp init --output-only` prints the blank template to stdout only. It must
  not start the questionnaire and must not construct or call a clipboard.
- Preserve the bare-command CLI model. Do not add a `prompt` namespace.
- Reserve `init` so user prompts cannot shadow the command.
- Document the command in grouped `--help` and README.
- Add `init` to the focused `kp` launcher as a command row. Selecting it runs
  the same questionnaire as `kp init`.
- Questionnaire chrome uses aligned emoji titles, an indented question, the
  launcher `›` prompt, and a single TTY session hint. Do not repeat the key
  hint on every field.
- Do not add persistence, history, templates, AI calls, network behavior, a
  general form framework, or unrelated refactors.

## ASSUMPTIONS

- Piped and non-TTY stdin keep one line per section so scripts and tests stay
  stable.
- On a TTY, Shift+Enter is distinguishable from Enter only in raw mode with
  xterm `modifyOtherKeys` (or equivalent CSI) sequences.
- Questionnaire chrome belongs on stderr so stdout stays pipeable.
- Reusing `errPickerCanceled` keeps cancellation messaging consistent.

## ACCEPTANCE CRITERIA

- `kp init` with five non-empty answers prints the rendered prompt on stdout
  and copies the same bytes after clipboard verification.
- Rendered output preserves user text including internal newlines, keeps the
  specified emoji-prefixed section order, and always ends with the
  independent-investigation sentence.
- On a TTY, Shift+Enter inserts a newline in the current section and Enter
  continues. Piped stdin does not show the TTY key hint.
- Bare `kp` includes an Init launcher row. Selecting it starts the
  questionnaire. Other secondary commands stay hidden.
- `kp init --output-only` prints only the blank template and does not touch
  the clipboard.
- EOF during the questionnaire exits `130` with no clipboard side effects.
- `kp new init` is rejected as a reserved name.
- Existing prompt commands, list output, and help sections besides the new
  `init` rows remain unchanged.

## ACCEPTED PLAN

1. Add a dedicated `init` Cobra command that reuses `app.inputReader` and the
   existing clipboard factory.
2. Keep prompt rendering in a pure function so tests can pin the exact body
   without driving the full CLI.
3. Reserve `init`, add help/README discovery, and pin focused command tests.
4. Run the repository validation suite and open a ready PR for the active
   issue (`#60` originally; `#62` for TTY multiline and emoji headers).

## TASK CHECKLIST

- [x] Write this spec in `docs/specs/0010-init-command/SPEC.md`.
- [x] Implement `kp init` and `--output-only`.
- [x] Reserve `init` and update help, README, and progress summary.
- [x] Add focused tests for rendering, clipboard, errors, and cancellation.
- [x] Add TTY Shift+Enter multiline input and emoji-prefixed section headers.
- [x] Add Init to the focused launcher and restyle questionnaire chrome.
- [x] Run format, test, race, vet, and build validation.

## VALIDATION

- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./internal/cmd ./internal/prompt` — `PASS`.
- `go test ./...` — `PASS` across every package.
- `go test -race ./...` — `PASS` across every package.
- `go vet ./...` — `PASS`.
- `go build ./...` — `PASS`.
- `make build` — `PASS`; produced `bin/kp`.
- Isolated CLI acceptance with an empty temporary config directory — `PASS`:
  `kp init --output-only` printed the emoji-prefixed blank template, `list --plain` did not
  include `init`, `--help` showed `Construct a coding-agent prompt`, `kp init --help`
  mentioned Shift+Enter, and `kp new init` exited `1` with `reserved prompt name: init`.
- `git diff --check` — `PASS`.
- Source-file-size audit of version-eligible Go files — `PASS`: changed Go
  files remain at or under 300 physical lines (`init.go`, `init_tty.go`,
  `init_keys.go`, `init_test.go`, `init_keys_test.go`).
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.

Deterministic rendering, section order, user-text preservation, and the
required investigation sentence are proven by pure render tests plus an
interactive CLI test. `--output-only` rendering and clipboard isolation are
proven by a CLI test whose clipboard factory must not be called. Interactive
stdout/clipboard equality is proven by comparing stdout to the fake clipboard
body. Cancellation and empty-then-retry behavior are proven by stdin fixtures.
Existing-command regressions are proven by unchanged list/help/prompt tests.
Launcher discovery and `command:init` execution are proven by focused launcher
tests. Questionnaire chrome uses aligned emoji headings, an indented question,
and the `›` prompt.

## DECISIONS

- Use a Cobra command instead of an embedded prompt because the product is
  constructed at runtime from operator input.
- Do not reuse root `--print`. The requested flag is `--output-only`, and root
  `--print` applies only to stored prompt execution.
- Expose `init` on the focused launcher. Discovery is launcher, help, and
  README.
- Record topology as `single-lane, because tightly coupled and high-overlap:
  one command, shared clipboard/stdin/help/reserved-name paths, and docs in a
  single delivery lane`.
- Enable TTY Shift+Enter by entering raw mode and decoding CSI sequences rather
  than adding a form framework. Keep non-TTY input as one line per section.

## DISCOVERIES

- `runPrompt` copies, verifies, then prints. `init` should do the same after
  answers are collected so a clipboard failure does not emit a successful body.
- `promptPort` and the numbered picker already treat EOF as exit `130`.
- `init` must be reserved; `find-port` currently is not, but that is out of
  scope.

## DOCUMENTATION UPDATES

- README command guide, quick start, and reserved-name list.
- Grouped root help Prompt Library rows.
- `docs/PROJECT_PROGRESS_SUMMARY.md`.

## DELIVERY DECISION

Issue #64, branch `GH-64`, canonical worktree
`~/worktrees/jamesonstone/kp/GH-64`, ready pull request to `main`. Earlier
delivery: issue #60 / `GH-60` and issue #62 / `GH-62`.

## OUTCOME

- `kp init` is a dedicated Cobra command that collects five answers, renders
  the specified prompt with emoji-prefixed section headers, prints it to
  stdout, and copies the identical body through existing clipboard
  verification.
- On a TTY, Shift+Enter inserts a newline and Enter continues to the next
  section. Piped stdin remains one line per section.
- `kp init --output-only` prints only the blank template and does not touch
  the clipboard.
- `init` is reserved. Stored prompt commands and list output remain unchanged.
  The focused launcher now includes Init. Issue #64 tracks launcher discovery
  and questionnaire layout on `GH-64`.

## REPOSITORY MEMORY

- Create this living specification because `kp init` is a dedicated
  constructor, not a stored prompt, and must stay distinct from
  `0001-v0-init-utility`.
- Constitution Current Project State now names `kp init` as part of the
  shipped CLI surface. No new invariant, non-goal, or workflow rule was
  required.
