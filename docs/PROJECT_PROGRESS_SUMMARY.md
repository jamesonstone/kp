# PROJECT PROGRESS SUMMARY

## FEATURE PROGRESS TABLE

| ID | FEATURE | PATH | PHASE | PAUSED | CREATED | SUMMARY |
| -- | ------- | ---- | ----- | ------ | ------- | ------- |
| 0001 | v0-init-utility | `docs/specs/0001-v0-init-utility` | implement | no | 2026-05-22 | Builds the initial Darwin-only `kp` Go CLI for an interactive root launcher, grouped `--help` output, prompt printing, exact clipboard verification, interactive `kp list`, whimsical rotating picker farewells, and local repo support scaffolding without a `prompt` subcommand. The feature embeds prompt assets, supports user prompt overrides, and stays limited to local build/install/test scope. |
| 0003 | merge-command | `docs/specs/0003-merge-command` | deliver | no | 2026-08-18 | Adds `kp merge` as a concise built-in prompt for evidence-backed Mermaid dependency graphs, exact-current merge readiness, topological waves, maximum safe independent concurrency, directional dependency proof, protected-workload gates, behavior-based recovery, wave revalidation, failure isolation, and separate merge, deployment, runtime, production, and rollback evidence. |
| 0004 | punchlist-command | `docs/specs/0004-punchlist-command` | deliver | no | 2026-08-25 | Adds `kp punchlist` as a built-in prompt for scanning a living punch list, clustering related observations, fixing shared causes, and keeping implemented, merged, deployed, and validated states distinct. |
| 0005 | handoff-prompts | `docs/specs/0005-handoff-prompts` | deliver | no | 2026-08-26 | Replaces ambiguous `kp handoff` with explicit chat-to-agent and agent-to-agent prompts that clarify at the origin, preserve zero-context task evidence and authority, reconcile at the destination, hydrate confirmed context, and request permission before implementation. |
| 0006 | plan-command | `docs/specs/0006-plan-command` | deliver | no | 2026-08-28 | Adds `kp plan` as a medium-length built-in prompt for driving an implementation plan to an evidence-backed local maximum in planning mode, without expanding `clarify`. |
| 0007 | goal-command | `docs/specs/0007-goal-command` | deliver | no | 2026-09-01 | Adds `kp goal` as a built-in prompt for Evidence-Backed Goal Convergence: research before asking, accumulate one goal model, challenge material ambiguity, and return a user-confirmed executable `/goal` before `kp plan`. |
| 0008 | ship-command | `docs/specs/0008-ship-command` | deliver | no | 2026-09-01 | Adds `kp ship` as a built-in prompt that pre-authorizes a coding agent to complete the current task thread through branch, pull-request, review, CI, in-scope merge, and established-workflow deployment without expanding `continue`, `goal`, or `merge`. |
| 0009 | review-command | `docs/specs/0009-review-command` | deliver | no | 2026-09-29 | Adds `kp review` as a built-in prompt for an independent, read-only correctness review of the current branch or pull request against the target branch, without expanding `pr`, `punchlist`, `ship`, or `merge`. |
| 0010 | init-command | `docs/specs/0010-init-command` | deliver | no | 2026-09-29 | Adds `kp init` as a dedicated interactive questionnaire that constructs a coding-agent prompt from objective, context, invariants, constraints, and definition of done, then prints and copies the generated body. TTY Shift+Enter inserts newlines; generated section headers use emoji prefixes. The focused `kp` launcher includes Init. |
| 0011 | picker-tui | `docs/specs/0011-picker-tui` | deliver | no | 2026-09-30 | Replaces the `fzf` picker behind `kp` and `kp list` with a restrained two-pane Bubble Tea v2 TUI that keeps the same keys, preview, stdout, clipboard, and exit semantics, and moves ten legacy prompts under a `kp v0` compatibility namespace. |

## PROJECT INTENT

Kit is a document-first workflow harness for disciplined thought work. It keeps durable project context in canonical markdown artifacts so humans and coding agents can move from research to specification, planning, tasks, implementation, reflection, and completion with explicit traceability.

## GLOBAL CONSTRAINTS

See `docs/CONSTITUTION.md` for project-wide constraints and principles.

## FEATURE SUMMARIES

### v0-init-utility

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Prompt insertion now favors a traditional CLI shape with an interactive default: `kp` opens an emoji-enhanced launcher for prompts and common safe command entries, `kp --help` is the low-friction static help entrypoint, `kp <name>` prints and copies prompt text with exact clipboard verification, `kp list` owns the prompt-only selector, and `kp scaffold` creates reusable repo support files without creating direct Kit project state. Built-ins now include the issue/branch/pull-request handoff prompt at `kp pr`. The project remains a local macOS utility with no `prompt` namespace and no release packaging in this feature.
- **APPROACH**: 1. Keep command parsing in `cmd/kp` and `internal/cmd`; keep path resolution in `internal/config`; keep prompt parsing, validation, and registry behavior in `internal/prompt`; keep clipboard copy/read/verify behavior in `internal/clipboard`; keep repo support file generation in `internal/scaffold`. 2. Use Cobra for predictable command/flag behavior, help output, command aliases where needed, and future top-level command registration. 3. Model bare prompt commands as root behavior: known command names route to command handlers, while any non-reserved single positional argument routes to prompt lookup. 4. Render root help with a Kit-style grouped layout behind `kp --help`, separating direct prompt commands such as `kp clarify`, `kp agent-handoff`, `kp chat-handoff`, and `kp pr` from prompt-library commands such as `kp list`, `kp new`, `kp edit`, and `kp rm`, plus utility commands such as `kp scaffold`. 5. Make bare `kp` an `fzf` launcher with Tab/Shift-Tab and arrow navigation, including user prompts plus safe command entries; prompt selections execute copy/print behavior, while side-effecting commands show help instead of writing files. 6. Render root launcher and picker cancellation as a random-start, in-process rotation of whimsical farewells while preserving typed exit `130` semantics and operational diagnostics. 7. Render root launcher rows as fixed-width table-like rows so prompt, command, and action columns stay visually aligned. 8. Reserve command names before registry lookup so user prompts cannot shadow `list`, `new`, `edit`, `rm`, `scaffold`, `prompt`, `help`, or `version`. 9. Strip YAML frontmatter once during prompt parsing and store metadata separately from the body that copy, print, and preview consume. 10. Implement exact clipboard verification as string equality after reading `pbpaste`; keep checksums only as optional diagnostics in verbose logs. 11. Use a compact emoji-enhanced `fzf` selector for `kp list`, with `--no-fzf` as the explicit numbered fallback and `--plain`/`--verbose` for non-interactive listing. 12. Make `kp scaffold` skip existing files by default, append missing `.gitignore` patterns, support `--dry-run`, and support `--force` for scaffold files while excluding `.kit.yaml`, `.kit/`, specs, notes, progress summary, and Constitution docs. 13. Keep Darwin-only clipboard behavior behind build-tagged files and provide a non-Darwin unsupported implementation for compile/test feedback outside macOS. 14. Treat README, Makefile, and local install flow as part of the implementation surface because release packaging is out of scope.
- **OPEN ITEMS**: T001-T011 and T013 are complete. T012 remains blocked only for optional performance/RSS evidence; the previous Cmd+V paste validation gap no longer applies after the no-paste default change. Issue #24 tracks the locally validated whimsical picker farewells on `GH-24`; merge remains unauthorized.
- **POINTERS**: `docs/specs/0001-v0-init-utility/BRAINSTORM.md`, `docs/specs/0001-v0-init-utility/SPEC.md`, `docs/specs/0001-v0-init-utility/PLAN.md`, `docs/specs/0001-v0-init-utility/TASKS.md`

### merge-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that turns an explicitly authorized PR set into an evidence-backed graphical dependency plan and safe merge waves without overstating merge, deployment, runtime, or production state.
- **APPROACH**: 1. Add one embedded `merge` prompt through the existing registry. 2. Synthesize LoopC's observe-act-remeasure discipline, Merge Controller's PR-forest and downstream-unlock model, and the repository's exact merge gate into six concise steps. 3. Keep routine, separately authorized remediation on the existing PR head between waves while invalidating old-head readiness and merge authority; reserve replacement PRs for material or unsafe changes. 4. Default to squash and merge with merge-commit fallback in the merge rule; both methods are authorized by default. 5. Apply a deadline validation budget of operational correctness with `NOT_RUN_BY_INSTRUCTION` for excluded suites. 6. Pin exact output and discovery behavior in prompt, registry, list, verbose-list, and help tests. 7. Update durable merge guidance, README, and testing guidance. 8. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Issue #40 tracks the deadline validation budget on `GH-40`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0003-merge-command/SPEC.md`

### punchlist-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that turns a living punch list into clustered, evidence-backed engineering work without treating items as an independent ticket queue or overstating deployed or validated state.
- **APPROACH**: 1. Add one embedded `punchlist` prompt through the existing registry. 2. Encode environment discovery, whole-list clustering, a 95% clarification gate, worklane reuse, engineering-note conventions, and implemented/merged/deployed/validated separation. 3. Pin exact output, approved-body hash, required contract phrases, and discovery behavior in prompt, registry, list, verbose-list, and help tests. 4. Update README and the project progress summary. 5. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Issue #28 tracks the command on `GH-28`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0004-punchlist-command/SPEC.md`

### handoff-prompts

- **STATUS**: implement
- **PAUSED**: no
- **INTENT**: Provide explicit, provider-neutral zero-context handoffs for chat-to-agent and agent-to-agent transfers without losing decisions, evidence, authority, validation, or the next safe action.
- **APPROACH**: 1. Replace the ambiguous embedded handoff asset with chat-handoff and agent-handoff. 2. Clarify only origin questions that change destination authority, landing lane, or remaining-work contract and cannot be recorded as native evidence states; record unreachable third-party artifacts as `UNKNOWN`/`BLOCKED` and emit. 3. Require destination live-state reconciliation, destination clarification, context hydration, and explicit permission before implementation. 4. Pin exact prompt hashes and discovery behavior without adding command-specific runtime code. 5. Validate all local CLI, source-size, and hygiene gates before ready-PR delivery.
- **OPEN ITEMS**: Issue #32 is delivered. Issue #34 tracks the unreachable-evidence origin-clarification correction on `GH-34`. Hosted correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0005-handoff-prompts/SPEC.md`

### plan-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that drives a current implementation plan to a practical, evidence-backed local maximum without leaving planning mode or expanding `clarify`.
- **APPROACH**: 1. Add one embedded `plan` prompt through the existing registry. 2. Encode progressive research, silent adversarial iteration, a 95% evidence-backed stop rule, and one replacement-plan output. 3. Pin exact output, approved-body hash, required contract phrases, and discovery behavior in prompt, registry, list, verbose-list, and help tests. 4. Update README and the project progress summary. 5. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Issue #44 tracks the command on `GH-44`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0006-plan-command/SPEC.md`

### goal-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that converges incomplete engineering intent into a user-confirmed executable `/goal` before implementation planning, without expanding `clarify` or replacing `plan`.
- **APPROACH**: 1. Add one embedded `goal` prompt through the existing registry. 2. Encode Evidence-Backed Goal Convergence: one accumulated model, research before questions, material questions only, synthesis, adversarial challenge, and a confirmed `/goal` contract. 3. Teach `kp plan` to consume an accepted `/goal` as its intent contract. 4. Pin exact output, approved-body hashes, required contract phrases, and discovery behavior in prompt, registry, list, verbose-list, and help tests. 5. Update README and the project progress summary. 6. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Issue #46 tracks the command on `GH-46`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0007-goal-command/SPEC.md`

### ship-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that pre-authorizes shipping the current task thread through the full delivery lifecycle, scoped only to that task, without expanding `continue`, `goal`, or `merge`.
- **APPROACH**: 1. Add one embedded `ship` prompt through the existing registry. 2. Preserve the supplied authorization contract, including the leading `/goal`. 3. Pin exact output, approved-body hash, required contract phrases, and discovery behavior in prompt, registry, list, verbose-list, and help tests. 4. Update README and the project progress summary. 5. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Issue #48 tracks the command on `GH-48`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0008-ship-command/SPEC.md`

### review-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that independently reviews the current branch or pull request for correctness against the target branch, without making changes or expanding `pr`, `punchlist`, `ship`, or `merge`.
- **APPROACH**: 1. Add one embedded `review` prompt through the existing registry. 2. Preserve the supplied untrusted-implementation review contract. 3. Pin exact output, approved-body hash, required contract phrases, and discovery behavior in prompt, registry, list, verbose-list, and help tests. 4. Update README and the project progress summary. 5. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Issue #58 tracks the command on `GH-58`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0009-review-command/SPEC.md`

### init-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction interactive path from five pieces of human knowledge to a complete coding-agent prompt, plus a pipeable blank template via `--output-only`.
- **APPROACH**: 1. Add a dedicated `init` Cobra command rather than an embedded prompt, because the body is constructed at runtime. 2. Reuse existing stdin line reading for pipes, clipboard copy/verify, exit codes, and reserved-name enforcement. 3. Keep questionnaire chrome on stderr so stdout remains the product. 4. On a TTY, enter raw mode so Shift+Enter inserts a newline and Enter continues; prefix generated section headers with emojis and give the questionnaire the same aligned `›` chrome as the launcher. 5. Expose Init on the focused `kp` launcher. 6. Pin rendering, clipboard isolation, cancellation, key decoding, launcher, and reserved-name tests. 7. Document the command in help, README, and the progress summary.
- **OPEN ITEMS**: Issue #64 tracks launcher discovery and questionnaire layout on `GH-64`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0010-init-command/SPEC.md`

### picker-tui

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Make the `kp` picker feel native to `kp` rather than configured `fzf`, without changing the interaction model, and free ten root prompt names by serving them from `kp v0`.
- **APPROACH**: 1. Add `internal/picker` on Bubble Tea v2 and Lip Gloss v2 (v1 rejected for its init-time terminal query); it owns layout, keys, scrolling, resize, and styles, and returns a selected ID. 2. Keep registry, clipboard, and execution in `internal/cmd` behind one `pick` helper. 3. Draw on `/dev/tty` so stdout stays clean. 4. Serve moved prompts from `prompts/v0/` as Cobra subcommands with unchanged bodies and flags. 5. Use the basic ANSI palette by role, show each item's command in an aligned list column under group headings, and tint preview Markdown. 6. Move `kp init` on a TTY to an editor draft (`$KP_EDITOR`, `$EDITOR`, `nvim`, `vi`); piped stdin keeps line mode.
- **OPEN ITEMS**: Issue #66 tracks delivery on `GH-66`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0011-picker-tui/SPEC.md`

## LAST UPDATED

2026-09-30 12:00:00 EDT
