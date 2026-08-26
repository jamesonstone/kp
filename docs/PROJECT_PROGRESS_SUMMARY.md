# PROJECT PROGRESS SUMMARY

## FEATURE PROGRESS TABLE

| ID | FEATURE | PATH | PHASE | PAUSED | CREATED | SUMMARY |
| -- | ------- | ---- | ----- | ------ | ------- | ------- |
| 0001 | v0-init-utility | `docs/specs/0001-v0-init-utility` | implement | no | 2026-05-22 | Builds the initial Darwin-only `kp` Go CLI for an interactive root launcher, grouped `--help` output, prompt printing, exact clipboard verification, interactive `kp list`, whimsical rotating picker farewells, and local repo support scaffolding without a `prompt` subcommand. The feature embeds prompt assets, supports user prompt overrides, and stays limited to local build/install/test scope. |
| 0003 | merge-command | `docs/specs/0003-merge-command` | deliver | no | 2026-08-18 | Adds `kp merge` as a concise built-in prompt for evidence-backed Mermaid dependency graphs, exact-current merge readiness, topological waves, maximum safe independent concurrency, directional dependency proof, protected-workload gates, behavior-based recovery, wave revalidation, failure isolation, and separate merge, deployment, runtime, production, and rollback evidence. |
| 0004 | punchlist-command | `docs/specs/0004-punchlist-command` | deliver | no | 2026-08-25 | Adds `kp punchlist` as a built-in prompt for scanning a living punch list, clustering related observations, fixing shared causes, and keeping implemented, merged, deployed, and validated states distinct. |
| 0005 | handoff-prompts | `docs/specs/0005-handoff-prompts` | deliver | no | 2026-08-26 | Replaces ambiguous `kp handoff` with explicit chat-to-agent and agent-to-agent prompts that clarify at the origin, preserve zero-context task evidence and authority, reconcile at the destination, hydrate confirmed context, and request permission before implementation. |

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
- **APPROACH**: 1. Add one embedded `merge` prompt through the existing registry. 2. Synthesize LoopC's observe-act-remeasure discipline, Merge Controller's PR-forest and downstream-unlock model, and the repository's exact merge gate into six concise steps. 3. Keep routine, separately authorized remediation on the existing PR head between waves while invalidating old-head readiness and merge authority; reserve replacement PRs for material or unsafe changes. 4. Pin exact output and discovery behavior in prompt, registry, list, verbose-list, and help tests. 5. Update durable merge guidance, README, and testing guidance. 6. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: PRs #17, #19, and #23 are merged. Issue #26 tracks the current in-place remediation correction on `GH-26`. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
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
- **APPROACH**: 1. Replace the ambiguous embedded handoff asset with chat-handoff and agent-handoff. 2. Clarify implementation-changing questions at the origin before emission. 3. Require destination live-state reconciliation, destination clarification, context hydration, and explicit permission before implementation. 4. Pin exact prompt hashes and discovery behavior without adding command-specific runtime code. 5. Validate all local CLI, source-size, and hygiene gates before ready-PR delivery.
- **OPEN ITEMS**: Issue #32 tracks locally validated ready-PR delivery on GH-32. Hosted correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0005-handoff-prompts/SPEC.md`

## LAST UPDATED

2026-08-26 11:00:00 EDT
