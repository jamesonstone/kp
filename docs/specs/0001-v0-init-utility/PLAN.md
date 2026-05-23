---
kit_metadata_version: 1
artifact: plan
feature:
  id: 0001
  slug: v0-init-utility
  dir: 0001-v0-init-utility
parallelization_mode: rlm
references:
  - id: constitution
    name: Project constitution
    type: repo_doc
    target: docs/CONSTITUTION.md
    relation: constrains
    read_policy: must
    used_for: CLI-first architecture, minimal dependency policy, validation expectations, and progress tracking
    status: active
  - id: spec
    name: v0-init-utility spec
    type: feature_artifact
    target: docs/specs/0001-v0-init-utility/SPEC.md
    relation: constrains
    read_policy: must
    used_for: binding requirements, accepted prompt set, acceptance evidence, and non-goals
    status: active
  - id: brainstorm
    name: v0-init-utility brainstorm
    type: feature_artifact
    target: docs/specs/0001-v0-init-utility/BRAINSTORM.md
    relation: informs
    read_policy: conditional
    used_for: historical research, codebase findings, and resolved command-surface tradeoffs; SPEC supersedes prompt inventory
    status: active
  - id: kit-map
    name: Kit map for v0-init-utility
    type: command
    target: kit map 0001-v0-init-utility
    selector: kit map 0001-v0-init-utility
    selector_type: command
    relation: verifies
    read_policy: evidence
    used_for: verified no prior feature relationships and confirmed current phase artifacts
    status: active
  - id: progress-summary
    name: Project progress summary
    type: repo_doc
    target: docs/PROJECT_PROGRESS_SUMMARY.md
    relation: verifies
    read_policy: evidence
    used_for: highest completed artifact tracking
    status: active
  - id: agents-entrypoint
    name: Agents docs entrypoint
    type: repo_doc
    target: docs/agents/README.md
    relation: guides
    read_policy: must
    used_for: repo-local routing and context-loading order
    status: active
  - id: workflows
    name: Workflow rules
    type: repo_doc
    target: docs/agents/WORKFLOWS.md
    relation: guides
    read_policy: must
    used_for: spec-driven source-of-truth order and readiness gate
    status: active
  - id: rlm
    name: RLM guide
    type: repo_doc
    target: docs/agents/RLM.md
    relation: guides
    read_policy: must
    used_for: progressive discovery and deterministic synthesis
    status: active
  - id: guardrails
    name: Guardrails
    type: repo_doc
    target: docs/agents/GUARDRAILS.md
    relation: constrains
    read_policy: must
    used_for: completion bar, placeholder cleanup, docs-first rules, and validation claims
    status: active
  - id: tooling
    name: Tooling rules
    type: repo_doc
    target: docs/agents/TOOLING.md
    relation: guides
    read_policy: must
    used_for: dispatch constraints, worktree rules, and secondary input order
    status: active
  - id: repo-inventory
    name: Repository inventory
    type: command
    target: rg --files --hidden -g '!.git'
    selector: rg --files --hidden -g '!.git'
    selector_type: command
    relation: informs
    read_policy: skip
    used_for: historical planning snapshot; verified the pre-implementation repo had no Go runtime surface before v0-init-utility added it
    status: stale
  - id: kit-config
    name: Kit configuration
    type: repo_file
    target: .kit.yaml
    relation: constrains
    read_policy: must
    used_for: specs directory, constitution path, and feature naming
    status: active
  - id: readme
    name: README
    type: repo_file
    target: README.md
    relation: informs
    read_policy: skip
    used_for: historical planning baseline; README now documents the implemented v0-init-utility command surface
    status: stale
  - id: gitignore
    name: Git ignore rules
    type: repo_file
    target: .gitignore
    relation: constrains
    read_policy: must
    used_for: existing Go, env, and Kit ignore coverage
    status: active
  - id: github-pr-template
    name: GitHub pull request template
    type: repo_file
    target: .github/pull_request_template.md
    relation: informs
    read_policy: evidence
    used_for: expected implementation validation evidence
    status: active
  - id: coderabbit-config
    name: CodeRabbit config
    type: repo_file
    target: .coderabbit.yaml
    relation: informs
    read_policy: evidence
    used_for: implementation review scope
    status: active
  - id: kit-scaffold-source
    name: Kit scaffold source
    type: external_repo
    target: /Users/jamesonstone/go/src/github.com/jamesonstone/kit/pkg/cli/init.go, /Users/jamesonstone/go/src/github.com/jamesonstone/kit/pkg/cli/scaffold_agents.go, /Users/jamesonstone/go/src/github.com/jamesonstone/kit/internal/templates
    relation: informs
    read_policy: evidence
    used_for: implementation strategy for `kp scaffold`, approved file inventory, skip-existing behavior, `.gitignore` append behavior, and embedded template content
    status: active
---
# PLAN

## SUMMARY

Implement `kp` as a small Darwin-first Go CLI with Cobra at the command boundary, a prompt registry that overlays user files onto embedded built-ins, a Darwin clipboard adapter that verifies exact clipboard contents after copy, and a local scaffold writer for repo support artifacts. Keep the first release local-only: create module, source tree, tests, prompts, scaffold templates, README, and Make targets, while excluding Homebrew, Goreleaser, release workflows, automatic paste, direct Kit project state, and a `prompt` subcommand.

## APPROACH

1. Start with the module and command skeleton so every later package can be compiled through `go test ./...`.
2. Keep command parsing in `cmd/kp` and `internal/cmd`; keep path resolution in `internal/config`; keep prompt parsing, validation, and registry behavior in `internal/prompt`; keep system clipboard copy/read/verify operations in `internal/clipboard`.
3. Use Cobra for predictable command/flag behavior, help output, command aliases where needed, and future top-level command registration.
4. Model bare prompt commands as root behavior: known command names route to command handlers, while any non-reserved single positional argument routes to prompt selection.
5. Reserve command names before registry lookup so user prompts cannot shadow `list`, `new`, `edit`, `rm`, `prompt`, `help`, or `version`.
6. Strip YAML frontmatter once during prompt parsing and store metadata separately from the body that copy, print, and preview consume.
7. Implement exact clipboard verification as string equality after reading `pbpaste`; keep checksums only as optional diagnostics in verbose logs.
8. Use dependency injection at package boundaries for command execution, filesystem roots, stdin/stdout/stderr, clock/polling, and editor command lookup so error paths can be tested without mutating real user state.
9. Keep Darwin-only system behavior behind build-tagged files and provide a non-Darwin unsupported implementation for compile/test feedback outside macOS.
10. Add `internal/scaffold` as a small filesystem package that owns the approved scaffold file inventory, embedded templates, default skip-existing behavior, `--force`, `--dry-run`, and `.gitignore` append-only updates.
11. Treat README, Makefile, and local install flow as part of the implementation surface because release packaging is out of scope.
12. Do not split implementation across subagents yet; overlap is high across command wiring, registry behavior, scaffold behavior, and acceptance tests. Preserve `parallelization_mode: rlm` for later task planning, but serialize implementation unless TASKS identifies file-disjoint work.

Tradeoffs:

1. Cobra adds one dependency but prevents ad hoc command parsing and supports future utility commands.
2. `gopkg.in/yaml.v3` adds one dependency but keeps frontmatter parsing structured instead of string-only heuristics.
3. The root-command prompt shortcut is more collision-prone than a `prompt` namespace; reserved command names and strict prompt-name validation are the containment strategy.
4. A process-local lock is intentionally omitted because the SPEC documents concurrent clipboard races as a known limitation.
5. Goreleaser and Homebrew are deferred to avoid release credential and archive concerns in the first implementation.
6. `kp scaffold` embeds adapted templates instead of importing or shelling out to Kit, which duplicates template text but keeps `kp` self-contained and usable outside a Kit checkout.

## COMPONENTS

1. `cmd/kp`
   - Owns `main.go`, process exit handling, injected version metadata, and root Cobra execution.
   - Keeps OS exit-code translation at the process edge instead of scattering `os.Exit`.
2. `internal/cmd`
   - Owns Cobra command construction, global flags, stdout/stderr wiring, command-name reservation, and the root prompt-dispatch rule.
   - Delegates prompt storage, config resolution, and clipboard behavior to internal packages.
3. `internal/config`
   - Resolves config roots from `--config`, `XDG_CONFIG_HOME`, and home directory fallback.
   - Derives `<config-root>/prompts`, creates needed directories, and reports config/IO failures as exit-code `3` candidates.
4. `internal/prompt`
   - Owns `Prompt`, `Source`, prompt-name validation, reserved-name checks, frontmatter parsing, built-in loading, user overlay behavior, CRUD, promotion, and sorted listing.
   - Treats parsed body as the only text eligible for copy, print, and preview.
5. `internal/clipboard`
   - Owns clipboard copy/read/verify orchestration and Darwin system command integration.
   - Separates command execution from verification policy so mismatch paths are testable.
6. `internal/scaffold`
   - Owns approved scaffold artifacts, template content, file-writing policy, dry-run planning, forced overwrite behavior, and `.gitignore` missing-pattern detection.
   - Excludes direct Kit project state such as `.kit.yaml`, `.kit/`, specs, notes, progress summary, and Constitution docs.
7. `prompts/`
   - Contains only `clarify.md` and `instructions.md` for this feature.
   - Excludes `plan.md`.
8. `README.md`
   - Documents local install, dependencies, commands, prompt file format, exit codes, performance checks, no-paste scope, and concurrent clipboard limitation.
9. `Makefile`
   - Provides local build, test, install, format, clean, and benchmark helper targets.
   - Does not include release, tap, or Goreleaser targets.

## DATA

1. Prompt source files
   - Built-ins live at `prompts/clarify.md` and `prompts/instructions.md`.
   - User prompts live at `<config-root>/prompts/<name>.md`.
   - File names define prompt names by removing `.md`; accepted names match `^[a-z][a-z0-9-]*$`.
2. Prompt parsing result
   - `Name`: unique prompt key, lowercase kebab-like format.
   - `Label`: frontmatter `label` or default label derived from `Name`.
   - `Source`: enum with `builtin` and `user` display values.
   - `FilePath`: absolute path for user prompts and empty string for embedded prompts.
   - `Body`: markdown text after optional YAML frontmatter.
3. Registry behavior
   - Built-ins load first.
   - User prompts overlay built-ins by name.
   - Listing sorts final visible prompts by `Name`.
4. Prompt frontmatter
   - Optional YAML block begins at byte 0 with `---` and ends at the next `---` delimiter line.
   - Only `label` is meaningful for this feature.
   - Malformed YAML is a user-facing load error.
5. Exit-code categories
   - `0`: success.
   - `1`: user input or prompt data error.
   - `2`: clipboard, platform, or system-command failure.
   - `3`: config, filesystem bootstrap, editor lookup, or missing `fzf` failure.
   - `130`: user cancellation.
6. Runtime state
   - No daemon, database, cache, telemetry, network state, or cross-process lock exists in this feature.
7. Scaffold artifacts
   - Files: `.env`, `.envrc`, `.coderabbit.yaml`, `.github/pull_request_template.md`, `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, `docs/agents/README.md`, `docs/agents/WORKFLOWS.md`, `docs/agents/RLM.md`, `docs/agents/TOOLING.md`, `docs/agents/GUARDRAILS.md`, `docs/references/README.md`, `docs/references/testing.md`, `docs/references/tooling.md`, and `docs/references/external-systems.md`.
   - `.gitignore` is handled as an append-only pattern set rather than a normal overwrite file.
   - Excluded direct Kit state: `.kit.yaml`, global Kit config, `.kit/`, `docs/specs/**`, `docs/notes/**`, `docs/PROJECT_PROGRESS_SUMMARY.md`, and `docs/CONSTITUTION.md`.

## INTERFACES

1. `kp`
   - Input: no positional arguments, global flags.
   - Output: command help.
   - Side effects: none.
2. `kp <name>`
   - Input: one prompt name and optional `--copy` or `--print`.
   - Output: body to stdout by default and for `--print`; quiet success for `--copy`; diagnostics to stderr on failure.
   - Side effects: default print-and-copy flow or copy-only flow; no side effects for print.
3. `kp list`
   - Input: optional `--no-fzf`, `--plain`, or `--verbose`.
   - Output: interactive picker through `fzf` by default; numbered fallback with `--no-fzf`; names only with `--plain`; tab-separated `name`, `label`, `source` with `--verbose`.
   - Side effects: default selector path runs the selected prompt through print-and-copy behavior; `--plain` and `--verbose` have no clipboard side effects.
4. `kp new <name>`
   - Input: valid non-reserved prompt name.
   - Output: editor launch and stderr diagnostics.
   - Side effects: creates `<config-root>/prompts/<name>.md`, deletes empty stub on empty save or cancellation.
5. `kp edit <name>`
   - Input: existing built-in or user prompt name.
   - Output: editor launch and optional promotion message.
   - Side effects: opens user prompt or promotes built-in source to user prompt before editing.
6. `kp rm <name>`
   - Input: user prompt name.
   - Output: concise stderr status or failure diagnostic.
   - Side effects: deletes only user prompt files.
7. `kp scaffold`
   - Input: optional `--dir <path>`, `--dry-run`, and `--force`.
   - Output: concise stdout status lines for created, updated, skipped, and dry-run planned artifacts.
   - Side effects: creates missing approved scaffold files in the target directory; appends missing `.gitignore` patterns; does not create direct Kit project state; does not touch clipboard or editor.
8. Global flags
   - `--config <dir>` changes config root for the current process.
   - `--verbose` emits key=value diagnostics to stderr.
   - `--copy` and `--print` apply only to prompt selection.
   - `--no-fzf` applies only to `kp list` picker behavior.
   - `--help` and `--version` must not initialize clipboard or prompt-edit side effects.
9. Files and artifacts touched by implementation
   - Create: `go.mod`, `go.sum`, `cmd/kp/main.go`, `internal/cmd/*`, `internal/config/*`, `internal/prompt/*`, `internal/clipboard/*`, `internal/scaffold/*`, `prompts/clarify.md`, `prompts/instructions.md`, `Makefile`, `LICENSE`.
   - Modify: `README.md`.
   - Do not create: `prompts/plan.md`, `.goreleaser.yaml`, `.github/workflows/release.yaml`.

## DEPENDENCIES

1. Binding contract: `docs/specs/0001-v0-init-utility/SPEC.md`.
2. Research context: `docs/specs/0001-v0-init-utility/BRAINSTORM.md`, only where not superseded by SPEC.
3. Project constraints: `docs/CONSTITUTION.md`.
4. Workflow constraints: `docs/agents/README.md`, `docs/agents/WORKFLOWS.md`, `docs/agents/RLM.md`, `docs/agents/GUARDRAILS.md`, and `docs/agents/TOOLING.md`.
5. Current repo baseline: `README.md`, `.gitignore`, `.github/pull_request_template.md`, `.coderabbit.yaml`, `.kit.yaml`, and `docs/PROJECT_PROGRESS_SUMMARY.md`.
6. Go dependencies: Go 1.22+, `github.com/spf13/cobra` v1.8+, and `gopkg.in/yaml.v3`.
7. macOS tools: `pbcopy`, `pbpaste`; user-installed `fzf`; editor from `KP_EDITOR`, `EDITOR`, or `vi`.
8. Optional validation tools: `hyperfine` for latency measurements and standard macOS process/memory inspection for RSS checks.
9. No external APIs, design assets, datasets, MCP tools, release services, network calls, or cloud resources shape this plan.
10. Kit scaffold source in `/Users/jamesonstone/go/src/github.com/jamesonstone/kit` informs template inventory only; `kp` must embed its own templates and must not depend on that checkout at runtime.

## RISKS

1. Root command ambiguity
   - Risk: future utility commands or current reserved words collide with prompt names.
   - Mitigation: centralize reserved-name validation and test each reserved word through `new`, direct selection, and list behavior.
2. Frontmatter/body boundary bugs
   - Risk: metadata leaks into printed/copied prompts or body bytes are trimmed incorrectly.
   - Mitigation: make parser tests cover frontmatter, no frontmatter, malformed frontmatter, empty body, and body beginning with literal delimiters.
3. Clipboard race
   - Risk: concurrent invocations share the global clipboard and can interfere.
   - Mitigation: keep exact read-back verification and document the known limitation; do not add locking outside SPEC.
4. System-command flakiness
   - Risk: `pbcopy`, `pbpaste`, `fzf`, or editor behavior differs across user environments.
   - Mitigation: wrap command execution behind narrow adapters and test missing command, non-zero exit, timeout/mismatch, and cancellation paths.
5. macOS-only behavior in developer environments
   - Risk: non-Darwin machines cannot run clipboard integration tests.
   - Mitigation: use build tags and unsupported stubs; keep Darwin integration tests tagged and make non-Darwin feedback explicit.
6. Performance regression from process spawning
   - Risk: repeated `pbcopy`/`pbpaste` invocations miss latency targets.
   - Mitigation: keep registry loading linear over small prompt sets, avoid unnecessary filesystem reads after resolution, and benchmark copy paths with `hyperfine`.
7. User data mutation in tests
   - Risk: tests or manual checks write to the real `~/.config/kp/prompts`.
   - Mitigation: default tests to temp `XDG_CONFIG_HOME` or `--config`; reserve real user config for explicit manual checks only.
8. Docs drift
   - Risk: README examples, exit codes, or prompt bodies diverge from implementation.
   - Mitigation: include README verification in acceptance evidence and keep prompt source text copied from SPEC exactly.
9. Scaffold overwrite risk
   - Risk: scaffold writes clobber user-maintained repo instruction files or config.
   - Mitigation: skip existing files by default, require `--force` for overwrites, keep `.gitignore` append-only, and test dry-run behavior with a temp target directory.

## TESTING

1. Unit tests
   - Cover prompt-name validation, reserved names, label derivation, frontmatter parsing, malformed YAML, empty bodies, sorted list order, user override precedence, built-in promotion, user-only removal, and config path resolution.
2. Command tests
   - Exercise Cobra command construction with injected stdout/stderr, temp config roots, fake registry, fake clipboard, fake editor, and fake picker dependencies.
   - Cover `kp`, `kp <name>`, `kp <name> --copy`, `kp <name> --print`, `kp list`, `kp list --verbose`, `kp new`, `kp edit`, and `kp rm`.
3. Clipboard tests
   - Unit-test verification policy with fake read values and fake time.
   - Darwin integration-test real `pbcopy`/`pbpaste` copy and read-back behavior.
   - Test mismatch path proves success output is not produced.
4. Picker/editor tests
   - Test missing `fzf`, `--no-fzf` valid selection, invalid selection, cancellation, missing editor, editor cancellation, empty save, and successful save.
5. Scaffold tests
   - Test full scaffold into a temp directory, dry-run no-write behavior, skip-existing behavior, force overwrite behavior, `.gitignore` append-only missing-pattern updates, excluded direct Kit state, and command no-side-effect behavior for clipboard/editor dependencies.
6. Platform tests
   - Build and run unit tests on Darwin.
   - Compile non-Darwin unsupported stubs where feasible; skip Darwin-only integration tests outside Darwin.
7. Acceptance evidence
   - Map each SPEC acceptance criterion to one of: automated unit test, automated command/integration test, performance benchmark, binary-size check, RSS check, or README review.
   - Record commands and observed outputs in the PR testing section.
8. Performance validation
   - Use `hyperfine './kp clarify --copy' --warmup 5 --runs 50` for copy latency.
   - Use `hyperfine './kp clarify' --warmup 5 --runs 50` for print-and-copy latency.
   - Instrument `kp list` startup to record process start to `fzf` process launch under 80 ms.
9. Documentation validation
   - Verify README command examples use `kp` without a `prompt` subcommand.
   - Verify README documents two built-ins only: `clarify` and `instructions`.
   - Verify README documents `kp scaffold` status, included artifact categories, and direct Kit-state exclusions.
   - Verify README states `prompts/plan.md`, Homebrew, Goreleaser, and release workflows are out of scope for this feature.
