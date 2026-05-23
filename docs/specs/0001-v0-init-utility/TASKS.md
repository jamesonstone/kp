---
kit_metadata_version: 1
artifact: tasks
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
    used_for: project authority order, CLI-first constraints, validation expectations, and progress tracking
    status: active
  - id: brainstorm
    name: v0-init-utility brainstorm
    type: feature_artifact
    target: docs/specs/0001-v0-init-utility/BRAINSTORM.md
    relation: informs
    read_policy: conditional
    used_for: upstream research context where not superseded by SPEC
    status: active
  - id: spec
    name: v0-init-utility spec
    type: feature_artifact
    target: docs/specs/0001-v0-init-utility/SPEC.md
    relation: constrains
    read_policy: must
    used_for: binding requirements, acceptance criteria, prompt inventory, and non-goals
    status: active
  - id: plan
    name: v0-init-utility plan
    type: feature_artifact
    target: docs/specs/0001-v0-init-utility/PLAN.md
    relation: guides
    read_policy: must
    used_for: implementation order, component boundaries, risks, and validation strategy
    status: active
  - id: progress-summary
    name: Project progress summary
    type: repo_doc
    target: docs/PROJECT_PROGRESS_SUMMARY.md
    relation: verifies
    read_policy: evidence
    used_for: highest completed artifact tracking
    status: active
  - id: kit-map
    name: Kit map for v0-init-utility
    type: command
    target: kit map 0001-v0-init-utility
    selector: kit map 0001-v0-init-utility
    selector_type: command
    relation: verifies
    read_policy: evidence
    used_for: verified feature artifacts, phase state, and absence of prior feature relationships
    status: active
  - id: kit-scaffold-source
    name: Kit scaffold source
    type: external_repo
    target: /Users/jamesonstone/go/src/github.com/jamesonstone/kit/pkg/cli/init.go, /Users/jamesonstone/go/src/github.com/jamesonstone/kit/pkg/cli/scaffold_agents.go, /Users/jamesonstone/go/src/github.com/jamesonstone/kit/internal/templates
    relation: informs
    read_policy: evidence
    used_for: kp scaffold task scope, template inventory, write policy, and validation evidence
    status: active
---
# TASKS

## PROGRESS TABLE

| ID | TASK | STATUS | OWNER | DEPENDENCIES |
| -- | ---- | ------ | ----- | ------------ |
| T001 | Scaffold Go module, root command, version metadata, and local build targets [PLAN-APPROACH](PLAN.md#approach) | done | agent | |
| T002 | Implement config path resolution and prompt directory bootstrap [PLAN-COMPONENTS](PLAN.md#components) | done | agent | T001 |
| T003 | Implement prompt types, name validation, reserved names, and frontmatter parsing [PLAN-DATA](PLAN.md#data) | done | agent | T001 |
| T004 | Add embedded built-in prompt assets for `clarify` and `instructions` only [PLAN-DATA](PLAN.md#data) | done | agent | T003 |
| T005 | Implement prompt registry overlay, CRUD, promotion, and sorted listing [PLAN-COMPONENTS](PLAN.md#components) | done | agent | T002, T003, T004 |
| T006 | Implement clipboard copy, read, exact verification, and platform stubs [PLAN-RISKS](PLAN.md#risks) | done | agent | T001 |
| T007 | Wire root help dispatch, prompt print/copy, and list output flows [PLAN-INTERFACES](PLAN.md#interfaces) | done | agent | T002, T005, T006 |
| T008 | Implement `kp list` picker flows with `fzf`, `--no-fzf`, cancellation handling, and preview behavior [PLAN-INTERFACES](PLAN.md#interfaces) | done | agent | T005, T006, T007 |
| T009 | Implement prompt creation, editing, built-in promotion, removal, and editor handling [PLAN-INTERFACES](PLAN.md#interfaces) | done | agent | T005, T007 |
| T010 | Add unit, command, platform, picker, editor, and clipboard tests [PLAN-TESTING](PLAN.md#testing) | done | agent | T002, T003, T004, T005, T006, T007, T008, T009 |
| T011 | Complete README, LICENSE, Makefile polish, and release-scope guardrails [PLAN-DEPENDENCIES](PLAN.md#dependencies) | done | agent | T007, T008, T009 |
| T013 | Implement `kp scaffold` repo support file generation [PLAN-INTERFACES](PLAN.md#interfaces) | done | agent | T007, T010, T011 |
| T012 | Run acceptance, performance, size, memory, and artifact validation evidence [PLAN-TESTING](PLAN.md#testing) | blocked | agent | T010, T011, T013 |

## TASK LIST

- [x] T001: Scaffold Go module, root command, version metadata, and local build targets [PLAN-APPROACH](PLAN.md#approach)
- [x] T002: Implement config path resolution and prompt directory bootstrap [PLAN-COMPONENTS](PLAN.md#components)
- [x] T003: Implement prompt types, name validation, reserved names, and frontmatter parsing [PLAN-DATA](PLAN.md#data)
- [x] T004: Add embedded built-in prompt assets for `clarify` and `instructions` only [PLAN-DATA](PLAN.md#data)
- [x] T005: Implement prompt registry overlay, CRUD, promotion, and sorted listing [PLAN-COMPONENTS](PLAN.md#components)
- [x] T006: Implement clipboard copy, read, exact verification, and platform stubs [PLAN-RISKS](PLAN.md#risks)
- [x] T007: Wire root help dispatch, prompt print/copy, and list output flows [PLAN-INTERFACES](PLAN.md#interfaces)
- [x] T008: Implement `kp list` picker flows with `fzf`, `--no-fzf`, cancellation handling, and preview behavior [PLAN-INTERFACES](PLAN.md#interfaces)
- [x] T009: Implement prompt creation, editing, built-in promotion, removal, and editor handling [PLAN-INTERFACES](PLAN.md#interfaces)
- [x] T010: Add unit, command, platform, picker, editor, and clipboard tests [PLAN-TESTING](PLAN.md#testing)
- [x] T011: Complete README, LICENSE, Makefile polish, and release-scope guardrails [PLAN-DEPENDENCIES](PLAN.md#dependencies)
- [x] T013: Implement `kp scaffold` repo support file generation [PLAN-INTERFACES](PLAN.md#interfaces)
- [ ] T012: Run acceptance, performance, size, memory, and artifact validation evidence [PLAN-TESTING](PLAN.md#testing)

## TASK DETAILS

### T001

- **GOAL**: Establish a compiling Go CLI skeleton named `kp` with version output and local build/test targets.
- **SCOPE**:
  - Create `go.mod` and `go.sum` for module `github.com/jamesonstone/kp`.
  - Add Cobra dependency `github.com/spf13/cobra` v1.8 or newer.
  - Create `cmd/kp/main.go` and the initial root command construction boundary.
  - Add version and commit metadata with deterministic `--version` output.
  - Add initial `Makefile` targets for build, test, install, format, and clean.
- **ACCEPTANCE**:
  - Evidence: `go test ./...` exits `0`.
  - Evidence: `go build -o kp ./cmd/kp` exits `0`.
  - Evidence: `./kp --version` exits `0` and prints version and commit fields.
  - Evidence: no `kp prompt` command is registered.
- **NOTES**: Keep exit-code translation at the process edge so later command tests can exercise command behavior without calling `os.Exit`.

### T002

- **GOAL**: Resolve config roots and create the prompt directory without touching real user state in tests.
- **SCOPE**:
  - Create `internal/config/config.go`.
  - Resolve `--config <dir>` first, then `XDG_CONFIG_HOME`, then `~/.config`.
  - Treat `--config <dir>` as the config root and derive `<dir>/prompts`.
  - Treat XDG/default roots as `<config-root>/kp/prompts`.
  - Resolve relative `--config` paths against the current working directory.
  - Create config and prompt directories with owner-only permissions where supported.
- **ACCEPTANCE**:
  - Evidence: config unit tests pass for XDG override, default home fallback, relative `--config`, directory creation, and filesystem failure mapping.
  - Evidence: tests use temp directories and do not write to `~/.config/kp/prompts`.
  - Evidence: config/IO failures map to exit code `3` through command-level tests after T007.
- **NOTES**: Keep path resolution separate from command parsing and prompt registry loading.

### T003

- **GOAL**: Define prompt data behavior before registry or command wiring depends on it.
- **SCOPE**:
  - Create `internal/prompt/prompt.go`.
  - Define `Prompt`, `Source`, typed errors, and display source values `builtin` and `user`.
  - Validate prompt names with `^[a-z][a-z0-9-]*$`.
  - Centralize reserved prompt names: `help`, `list`, `new`, `edit`, `rm`, `prompt`, and `version`.
  - Create `internal/prompt/frontmatter.go` using `gopkg.in/yaml.v3`.
  - Strip optional YAML frontmatter from copy, print, and preview body output.
- **ACCEPTANCE**:
  - Evidence: `go test ./internal/prompt -run 'Frontmatter|Name|Reserved'` exits `0`.
  - Evidence: tests cover frontmatter label, no frontmatter, malformed YAML, missing label default, empty body, literal delimiter body, invalid names, and reserved names.
  - Evidence: copied, printed, and previewed body test fixtures exclude YAML frontmatter.
- **NOTES**: Frontmatter metadata supports only `label` in this feature.

### T004

- **GOAL**: Add the exact approved built-in prompts and make them embeddable at compile time.
- **SCOPE**:
  - Create `prompts/clarify.md` with the exact source content from `SPEC.md`.
  - Create `prompts/instructions.md` with the exact source content from `SPEC.md`.
  - Exclude `prompts/plan.md`.
  - Add embedded prompt loading support in `internal/prompt`.
- **ACCEPTANCE**:
  - Evidence: built-in prompt tests confirm exactly two built-ins: `clarify` and `instructions`.
  - Evidence: tests compare source content for both prompt files against SPEC-approved fixtures.
  - Evidence: `rg --files prompts | sort` prints only `prompts/clarify.md` and `prompts/instructions.md`.
  - Evidence: `test ! -e prompts/plan.md` exits `0`.
- **NOTES**: `BRAINSTORM.md` contains stale historical `plan` prompt context; `SPEC.md` supersedes it.

### T005

- **GOAL**: Provide prompt registry behavior for listing, lookup, user overrides, creation, removal, and built-in promotion.
- **SCOPE**:
  - Create `internal/prompt/registry.go`.
  - Load built-ins before user prompts.
  - Load user prompt files on every invocation from the resolved prompt directory.
  - Let user prompts shadow built-ins by matching name.
  - Sort visible prompt list by name.
  - Implement add, remove, and promote behavior with typed errors.
  - Reject empty prompt bodies and invalid filename-derived prompt names.
- **ACCEPTANCE**:
  - Evidence: `go test ./internal/prompt` exits `0`.
  - Evidence: tests cover user override precedence, sorted list order, add collision, add empty body, invalid user file name, remove user prompt, reject built-in-only removal, promote built-in to user, and malformed user frontmatter path reporting.
  - Evidence: registry tests use temp prompt directories only.
- **NOTES**: Registry code must not know about Cobra flags or system clipboard commands.

### T006

- **GOAL**: Implement clipboard behavior with exact read-back verification after copy.
- **SCOPE**:
  - Create `internal/clipboard` package.
  - On Darwin, wrap `pbcopy` and `pbpaste`.
  - Verify clipboard contents by exact string equality after copy.
  - Retry verification 5 times over 250 ms total.
  - Do not report success when verification fails.
  - Provide non-Darwin unsupported behavior so builds and tests fail clearly or skip platform-only checks.
  - Keep optional verbose diagnostics non-contractual for checksum algorithm choice.
- **ACCEPTANCE**:
  - Evidence: clipboard unit tests prove match succeeds and mismatch exits the verification path within the 250 ms policy.
  - Evidence: Darwin integration tests prove `Copy` plus `Read` round-trips through `pbcopy` and `pbpaste`.
  - Evidence: non-Darwin tests or compile checks expose unsupported-platform behavior without registering Linux or Windows support.
- **NOTES**: Keep command execution injectable so tests do not depend on real clipboard side effects.

### T007

- **GOAL**: Wire the bare `kp` command surface for help, direct prompt printing/copying, listing, and prompt lifecycle commands.
- **SCOPE**:
  - Implement root command dispatch in `internal/cmd`.
  - Register `list`, global flags, `--copy`, `--print`, `--config`, `--verbose`, `--help`, and `--version`.
  - Treat `kp <name>` as prompt lookup unless `<name>` is a reserved command.
  - Implement `kp list --plain` and `kp list --verbose` output formats.
  - Implement `kp <name> --print`, `kp <name> --copy`, and `kp <name>` side-effect policies.
  - Map typed errors to stable exit codes.
- **ACCEPTANCE**:
  - Evidence: command tests pass for `kp list`, `kp list --verbose`, `kp clarify --print`, `kp clarify --copy`, `kp clarify`, unknown prompt, invalid name, reserved name, verbose logging, and help/version no-side-effect behavior.
  - Evidence: `./kp` prints help when no arguments are passed.
  - Evidence: `./kp list --plain` prints exactly `clarify` then `instructions` when no user prompts exist.
  - Evidence: `./kp clarify --print` prints the prompt body without YAML frontmatter.
  - Evidence: `rg -n 'Use:.*prompt|kp prompt' cmd internal` finds no registered `prompt` command.
- **NOTES**: Normal success output stays quiet except requested stdout data and concise stderr status.

### T008

- **GOAL**: Implement `kp list` interactive prompt selection with `fzf` and numbered fallback behavior.
- **SCOPE**:
  - Make `kp list` launch `fzf` when available.
  - Feed sorted prompt display data and preview body data to `fzf`.
  - Detect missing `fzf` and exit `3` unless `--no-fzf` is passed.
  - Implement `kp list --no-fzf` numbered list on stderr with stdin selection.
  - Treat picker cancellation as exit `130`.
  - Treat non-numeric and out-of-range numbered selections as exit `1`.
- **ACCEPTANCE**:
  - Evidence: picker tests pass for `fzf` selection, preview data, missing `fzf`, cancellation, valid `--no-fzf` selection, invalid `--no-fzf` input, and out-of-range input.
  - Evidence: missing `fzf` stderr contains `brew install fzf` and `--no-fzf`.
  - Evidence: cancellation and invalid selection tests prove clipboard is not modified.
  - Evidence: `kp list` startup instrumentation can measure process start to `fzf` launch for final validation.
- **NOTES**: Use injected picker execution in tests; reserve real `fzf` for manual or integration validation.

### T009

- **GOAL**: Implement user prompt lifecycle commands with editor resolution and safe file mutation.
- **SCOPE**:
  - Register `new`, `edit`, and `rm` commands.
  - Resolve editor command in order: `KP_EDITOR`, `EDITOR`, then `vi`.
  - Create new prompt stubs under the resolved prompt directory.
  - Delete new empty stubs after empty save or cancellation.
  - Promote built-ins to user prompt files before editing.
  - Delete only user prompt files during `rm`.
  - Preserve existing user prompt files on editor cancellation.
- **ACCEPTANCE**:
  - Evidence: command tests pass for successful new, collision with built-in, collision with user prompt, invalid name, missing editor, editor cancellation, empty save cleanup, edit user prompt, promote built-in, missing prompt edit, remove user prompt, and reject built-in-only remove.
  - Evidence: `XDG_CONFIG_HOME=$(mktemp -d) KP_EDITOR=true ./kp edit clarify` creates `kp/prompts/clarify.md` with the built-in clarify source.
  - Evidence: `kp new instructions` exits `1`.
  - Evidence: `kp rm instructions` exits `1` when no user override exists.
- **NOTES**: A no-op editor can validate promotion but cannot validate non-empty `new` without a test helper editor.

### T010

- **GOAL**: Make automated tests cover the public behavior and risky internal seams before final validation.
- **SCOPE**:
  - Add package-level unit tests for `internal/config`, `internal/prompt`, and `internal/clipboard`.
  - Add command tests for every public command path and exit-code category.
  - Add Darwin-tagged integration tests for real clipboard copy/read behavior.
  - Add negative-path tests for each public function and command family.
  - Keep tests deterministic with injected stdout, stderr, stdin, temp config roots, fake editors, fake pickers, fake clipboard, and fake command runners.
- **ACCEPTANCE**:
  - Evidence: `go test ./...` exits `0` on Darwin.
  - Evidence: Darwin-only tests are build-tagged or skipped appropriately on non-Darwin.
  - Evidence: tests cover exit codes `0`, `1`, `2`, `3`, and `130`.
  - Evidence: test fixtures confirm frontmatter never appears in copied, printed, or previewed bodies.
- **NOTES**: This task can add tests for behavior introduced earlier, but it must not introduce new product scope.

### T011

- **GOAL**: Align user-facing documentation and local tooling with the implemented CLI surface and release non-goals.
- **SCOPE**:
  - Update `README.md` with local install, dependency installation, commands, prompt format, config paths, exit codes, performance checks, no-paste scope, and concurrent clipboard limitation.
  - Add `LICENSE` with MIT license text.
  - Finalize `Makefile` targets from T001 against actual source paths.
  - Document that built-ins are `clarify` and `instructions` only.
  - Document that Homebrew, Goreleaser, release archives, GitHub release workflows, Linux, and Windows are out of scope for this feature.
- **ACCEPTANCE**:
  - Evidence: README examples use `kp`, `kp list`, `kp clarify`, `kp clarify --copy`, and `kp clarify --print` with no `kp prompt` examples.
  - Evidence: README documents exactly two built-ins: `clarify` and `instructions`.
  - Evidence: `.goreleaser.yaml` and `.github/workflows/release.yaml` do not exist.
  - Evidence: `make build`, `make test`, `make install`, `make fmt`, and `make clean` have documented behavior and run against local project paths.
- **NOTES**: Keep release packaging instructions out of executable targets for this feature.

### T013

- **GOAL**: Add a top-level `kp scaffold` command that creates approved repo support artifacts without creating direct Kit project state.
- **SCOPE**:
  - Add `internal/scaffold` for scaffold artifact inventory, embedded templates, planning, dry-run, forced writes, and `.gitignore` missing-pattern updates.
  - Register `kp scaffold` under the root command with `--dir <path>`, `--dry-run`, and `--force`.
  - Include `.env`, `.envrc`, `.coderabbit.yaml`, `.github/pull_request_template.md`, `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, `docs/agents/*`, and `docs/references/*`.
  - Exclude `.kit.yaml`, global Kit config, `.kit/`, `docs/specs/**`, `docs/notes/**`, `docs/PROJECT_PROGRESS_SUMMARY.md`, and `docs/CONSTITUTION.md`.
  - Add tests for normal scaffold, dry-run, skip-existing, force overwrite, `.gitignore` append-only behavior, excluded files, and no clipboard/editor side effects.
  - Update README and root help to describe `kp scaffold`.
- **ACCEPTANCE**:
  - Evidence: `go test ./...` exits `0`.
  - Evidence: `kp scaffold --dir <temp>` creates all approved scaffold files and excludes all direct Kit project-state files.
  - Evidence: `kp scaffold --dir <temp> --dry-run` writes no files and reports planned actions.
  - Evidence: `kp scaffold --dir <temp>` skips existing scaffold files by default.
  - Evidence: `kp scaffold --dir <temp> --force` overwrites scaffold files except `.gitignore`, which remains append-only.
  - Evidence: help output and README include `kp scaffold`.
- **VERIFY**:
  - `go test ./...`
  - `go vet ./...`
  - `make build`
  - `./bin/kp scaffold --dir "$(mktemp -d)" --dry-run`
  - `./bin/kp scaffold --dir "$(mktemp -d)"`
- **EXPECTED FILES**:
  - `internal/scaffold/scaffold.go`
  - `internal/scaffold/templates.go`
  - `internal/scaffold/scaffold_test.go`
  - `internal/cmd/root.go`
  - `internal/cmd/root_help.go`
  - `internal/cmd/root_test.go`
  - `internal/prompt/prompt.go`
  - `README.md`
  - `docs/specs/0001-v0-init-utility/SPEC.md`
  - `docs/specs/0001-v0-init-utility/PLAN.md`
  - `docs/specs/0001-v0-init-utility/TASKS.md`
  - `docs/PROJECT_PROGRESS_SUMMARY.md`
- **RISK**: The scaffold command writes many repo-root files, so incorrect target resolution or force behavior could overwrite user content.
- **ROLLBACK**: Remove `internal/scaffold`, unregister `kp scaffold`, remove scaffold tests/docs, and revert `scaffold` from reserved prompt names.
- **NOTES**: Embed templates in `kp`; do not shell out to `kit` or read from a local Kit checkout at runtime.

### T012

- **GOAL**: Produce final implementation evidence mapped to SPEC acceptance criteria and PLAN testing strategy.
- **SCOPE**:
  - Run full automated tests and build commands.
  - Run manual macOS clipboard checks.
  - Run latency benchmarks for copy paths on M-series macOS when available.
  - Measure stripped binary size.
  - Measure picker RSS while interactive `fzf` is active.
  - Confirm no release artifacts or `plan` prompt were created.
  - Update task checkboxes and progress summary only for evidence-backed completion.
- **ACCEPTANCE**:
  - Evidence: `go test ./...` exits `0`.
  - Evidence: `go build -o kp ./cmd/kp` exits `0`.
  - Evidence: `./kp`, `./kp list --plain`, `./kp list --verbose`, `./kp clarify`, `./kp clarify --print`, and `./kp clarify --copy` match SPEC outputs.
  - Evidence: `pbpaste` equals the clarify body byte-for-byte after `./kp clarify --copy`.
  - Evidence: `hyperfine './kp clarify --copy' --warmup 5 --runs 50` reports p99 below 100 ms on M-series macOS.
  - Evidence: `hyperfine './kp clarify' --warmup 5 --runs 50` reports p99 below 200 ms on M-series macOS with a focused lightweight text target.
  - Evidence: `kp list` picker startup to `fzf` launch is below 80 ms.
  - Evidence: stripped binary size is below 10 MB.
  - Evidence: interactive picker RSS is below 30 MB.
  - Evidence: `.goreleaser.yaml`, `.github/workflows/release.yaml`, and `prompts/plan.md` do not exist.
- **NOTES**: If optional tools such as `hyperfine` are unavailable, record the missing tool as a validation gap rather than marking the performance criterion complete.
  Current validation gaps:
  - `hyperfine` is not installed, so the required hyperfine copy latency commands did not run.
  - Interactive picker RSS with real `fzf` active was not measured because the PTY-based measurement attempt did not produce stable process data.
- **VERIFY**:
  - `go test ./...`
  - `go build -o kp ./cmd/kp`
  - `./kp`
  - `./kp list --plain`
  - `./kp list --verbose`
  - `./kp clarify --print`
  - `./kp clarify --copy` followed by byte-for-byte `pbpaste` comparison
  - `hyperfine './kp clarify --copy' --warmup 5 --runs 50`
  - `hyperfine './kp clarify' --warmup 5 --runs 50`
  - `kp list` picker startup instrumentation to `fzf` launch
  - stripped binary size measurement
  - interactive picker RSS measurement while real `fzf` is active
  - `test ! -e prompts/plan.md && test ! -e .goreleaser.yaml && test ! -e .github/workflows/release.yaml`
- **EXPECTED FILES**:
  - `kp`
  - `.kit/runs/<run-id>`
  - no `prompts/plan.md`
  - no `.goreleaser.yaml`
  - no `.github/workflows/release.yaml`
- **RISK**: This task depends on optional benchmark tooling, real clipboard state, and a stable interactive `fzf` process, so validation can be blocked by workstation state rather than code defects.
- **ROLLBACK**: Do not mark T012 done or append reflection completion when any required evidence is missing; preserve the blocked notes and rerun the missing validation after installing tools or stabilizing the local UI target.

## DEPENDENCIES

1. T001 must run first because every implementation task depends on the module and command entrypoint.
2. T002 and T003 can be developed after T001, but T005 depends on both.
3. T004 depends on T003 so prompt parsing and body/frontmatter semantics are established before embedding fixtures.
4. T006 can proceed after T001, but command behavior in T007 depends on its interfaces.
5. T007 must precede T008 and T009 because picker and lifecycle commands share root dispatch, flags, output, and exit-code mapping.
6. T010 depends on component implementation tasks but may add tests for earlier work before final validation.
7. T011 depends on the final command surface from T007 through T009 so README examples match implementation.
8. T013 depends on T007, T010, and T011 because it extends the command surface, requires tests, and changes README/help behavior.
9. T012 depends on T010, T011, and T013 because final evidence must include automated tests, documentation review, and scaffold command validation.
10. External validation dependencies: macOS 13+, `pbcopy`, `pbpaste`, `fzf`, an available editor, and optional `hyperfine`.
11. Blocking product decisions: none. SPEC and PLAN resolve the prompt inventory, command namespace, scaffold file inventory, release exclusions, platform scope, and verification policy.

## NOTES

1. Treat `SPEC.md` and `PLAN.md` as binding inputs; use `BRAINSTORM.md` only as historical research context.
2. Do not add `kp prompt ...`; all command examples and tests must use the bare command surface.
3. Do not create `prompts/plan.md`, `.goreleaser.yaml`, or `.github/workflows/release.yaml` in this feature.
4. Do not split implementation across subagents unless a later execution phase identifies file-disjoint work with low overlap; current ordering assumes serialized execution.
5. Keep `docs/PROJECT_PROGRESS_SUMMARY.md` synchronized with the highest completed artifact and evidence-backed implementation progress.
6. `kp scaffold` must not create `.kit.yaml`, `.kit/`, global Kit config, specs, notes, progress summary, or Constitution docs.
