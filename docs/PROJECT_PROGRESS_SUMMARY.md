# PROJECT PROGRESS SUMMARY

## FEATURE PROGRESS TABLE

| ID | FEATURE | PATH | PHASE | PAUSED | CREATED | SUMMARY |
| -- | ------- | ---- | ----- | ------ | ------- | ------- |
| 0001 | v0-init-utility | `docs/specs/0001-v0-init-utility` | implement | no | 2026-05-22 | Builds the initial Darwin-only `kp` Go CLI for an interactive root launcher, grouped `--help` output, prompt printing, exact clipboard verification, interactive `kp list`, and local repo support scaffolding without a `prompt` subcommand. The feature embeds prompt assets including `clarify`, `handoff`, `parentthread`, and `pr`, supports user prompt overrides, and stays limited to local build/install/test scope. |
| 0003 | merge-command | `docs/specs/0003-merge-command` | deliver | no | 2026-08-18 | Adds `kp merge` as a concise built-in prompt for evidence-backed Mermaid dependency graphs, exact-current merge readiness, topological waves, maximum safe independent concurrency, downstream-unlock priority, wave revalidation, failure isolation, and separate merge, deployment, runtime, production, and rollback evidence. |

## PROJECT INTENT

Kit is a document-first workflow harness for disciplined thought work. It keeps durable project context in canonical markdown artifacts so humans and coding agents can move from research to specification, planning, tasks, implementation, reflection, and completion with explicit traceability.

## GLOBAL CONSTRAINTS

See `docs/CONSTITUTION.md` for project-wide constraints and principles.

## FEATURE SUMMARIES

### v0-init-utility

- **STATUS**: implement
- **PAUSED**: no
- **INTENT**: Prompt insertion now favors a traditional CLI shape with an interactive default: `kp` opens an emoji-enhanced launcher for prompts and common safe command entries, `kp --help` is the low-friction static help entrypoint, `kp <name>` prints and copies prompt text with exact clipboard verification, `kp list` owns the prompt-only selector, and `kp scaffold` creates reusable repo support files without creating direct Kit project state. Built-ins now include the issue/branch/pull-request handoff prompt at `kp pr`. The project remains a local macOS utility with no `prompt` namespace and no release packaging in this feature.
- **APPROACH**: 1. Keep command parsing in `cmd/kp` and `internal/cmd`; keep path resolution in `internal/config`; keep prompt parsing, validation, and registry behavior in `internal/prompt`; keep clipboard copy/read/verify behavior in `internal/clipboard`; keep repo support file generation in `internal/scaffold`. 2. Use Cobra for predictable command/flag behavior, help output, command aliases where needed, and future top-level command registration. 3. Model bare prompt commands as root behavior: known command names route to command handlers, while any non-reserved single positional argument routes to prompt lookup. 4. Render root help with a Kit-style grouped layout behind `kp --help`, separating direct prompt commands such as `kp clarify`, `kp handoff`, and `kp pr` from prompt-library commands such as `kp list`, `kp new`, `kp edit`, and `kp rm`, plus utility commands such as `kp scaffold`. 5. Make bare `kp` an `fzf` launcher with Tab/Shift-Tab and arrow navigation, including user prompts plus safe command entries; prompt selections execute copy/print behavior, while side-effecting commands show help instead of writing files. 6. Render the bare `kp` launcher as fixed-width table-like rows so prompt, command, and action columns stay visually aligned. 7. Reserve command names before registry lookup so user prompts cannot shadow `list`, `new`, `edit`, `rm`, `scaffold`, `prompt`, `help`, or `version`. 8. Strip YAML frontmatter once during prompt parsing and store metadata separately from the body that copy, print, and preview consume. 9. Implement exact clipboard verification as string equality after reading `pbpaste`; keep checksums only as optional diagnostics in verbose logs. 10. Use a compact emoji-enhanced `fzf` selector for `kp list`, with `--no-fzf` as the explicit numbered fallback and `--plain`/`--verbose` for non-interactive listing. 11. Make `kp scaffold` skip existing files by default, append missing `.gitignore` patterns, support `--dry-run`, and support `--force` for scaffold files while excluding `.kit.yaml`, `.kit/`, specs, notes, progress summary, and Constitution docs. 12. Keep Darwin-only clipboard behavior behind build-tagged files and provide a non-Darwin unsupported implementation for compile/test feedback outside macOS. 13. Treat README, Makefile, and local install flow as part of the implementation surface because release packaging is out of scope.
- **OPEN ITEMS**: T001-T011 and T013 are complete. T012 remains blocked only for optional performance/RSS evidence; the previous Cmd+V paste validation gap no longer applies after the no-paste default change.
- **POINTERS**: `docs/specs/0001-v0-init-utility/BRAINSTORM.md`, `docs/specs/0001-v0-init-utility/SPEC.md`, `docs/specs/0001-v0-init-utility/PLAN.md`, `docs/specs/0001-v0-init-utility/TASKS.md`

### merge-command

- **STATUS**: deliver
- **PAUSED**: no
- **INTENT**: Provide one low-friction prompt that turns an explicitly authorized PR set into an evidence-backed graphical dependency plan and safe merge waves without overstating merge, deployment, runtime, or production state.
- **APPROACH**: 1. Add one embedded `merge` prompt through the existing registry. 2. Synthesize LoopC's observe-act-remeasure discipline, Merge Controller's PR-forest and downstream-unlock model, and the repository's exact merge gate into six concise steps. 3. Pin exact output and discovery behavior in prompt, registry, list, verbose-list, and help tests. 4. Update README and testing guidance. 5. Validate format, tests, race behavior, vet, build, CLI output, diff hygiene, and source-file size before ready-PR delivery.
- **OPEN ITEMS**: Implementation and local validation are complete; create the ready pull request for issue #16. Hosted pull-request correctness checks remain unavailable because the repository has no validation workflow.
- **POINTERS**: `docs/specs/0003-merge-command/SPEC.md`

## LAST UPDATED

2026-08-18 10:42:00 EDT
