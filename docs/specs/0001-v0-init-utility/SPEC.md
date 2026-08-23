---
kit_metadata_version: 1
artifact: spec
feature:
  id: 0001
  slug: v0-init-utility
  dir: 0001-v0-init-utility
skills:
  - id: rlm
    name: RLM
    target: docs/agents/RLM.md
    trigger: analyze codebase; scan all files; large repository analysis; scan repository; recursive language model
    used_for: progressive discovery, source-attributed synthesis, and downstream planning with parallelization_mode rlm
    status: active
references:
  - id: constitution
    name: Project constitution
    type: repo_doc
    target: docs/CONSTITUTION.md
    relation: constrains
    read_policy: must
    used_for: project authority order, CLI-first rules, dependency policy, validation rules, and progress tracking
    status: active
  - id: brainstorm
    name: v0-init-utility brainstorm
    type: feature_artifact
    target: docs/specs/0001-v0-init-utility/BRAINSTORM.md
    relation: informs
    read_policy: must
    used_for: resolved scope, accepted CLI namespace, built-in prompt bodies, edge cases, and out-of-scope release work
    status: active
  - id: progress-summary
    name: Project progress summary
    type: repo_doc
    target: docs/PROJECT_PROGRESS_SUMMARY.md
    relation: verifies
    read_policy: evidence
    used_for: verified current feature phase and highest completed artifact tracking
    status: active
  - id: kit-map
    name: Kit map for v0-init-utility
    type: command
    target: kit map 0001-v0-init-utility
    selector: kit map 0001-v0-init-utility
    selector_type: command
    relation: verifies
    read_policy: evidence
    used_for: verified no prior feature relationships and confirmed BRAINSTORM/SPEC presence
    status: active
  - id: kit-scaffold-source
    name: Kit scaffold source
    type: external_repo
    target: /Users/jamesonstone/go/src/github.com/jamesonstone/kit/pkg/cli/init.go, /Users/jamesonstone/go/src/github.com/jamesonstone/kit/pkg/cli/scaffold_agents.go, /Users/jamesonstone/go/src/github.com/jamesonstone/kit/internal/templates
    relation: informs
    read_policy: evidence
    used_for: source behavior for the `kp scaffold` file inventory, skip-existing policy, `.gitignore` append policy, and instruction document templates
    status: active
  - id: agents-entrypoint
    name: Agents docs entrypoint
    type: repo_doc
    target: docs/agents/README.md
    relation: guides
    read_policy: must
    used_for: repo-local routing and progressive context loading
    status: active
  - id: workflows
    name: Workflow rules
    type: repo_doc
    target: docs/agents/WORKFLOWS.md
    relation: guides
    read_policy: must
    used_for: spec-driven workflow and source-of-truth order
    status: active
  - id: rlm-guide
    name: RLM guide
    type: repo_doc
    target: docs/agents/RLM.md
    relation: guides
    read_policy: must
    used_for: RLM discovery, prior-work pass, and reference recording requirements
    status: active
  - id: tooling
    name: Tooling rules
    type: repo_doc
    target: docs/agents/TOOLING.md
    relation: guides
    read_policy: must
    used_for: canonical skills discovery, dispatch guidance, worktree rules, and secondary input order
    status: active
  - id: guardrails
    name: Guardrails
    type: repo_doc
    target: docs/agents/GUARDRAILS.md
    relation: constrains
    read_policy: must
    used_for: completion bar, no placeholder sections, docs-first rules, and validation expectations
    status: active
  - id: top-level-agent-routing
    name: Top-level agent routing files
    type: repo_doc
    target: AGENTS.md, CLAUDE.md, .github/copilot-instructions.md
    relation: informs
    read_policy: evidence
    used_for: confirmed top-level files route to repo-local docs and should stay short
    status: active
  - id: repo-inventory
    name: Repository inventory
    type: command
    target: rg --files --hidden -g '!.git'
    selector: rg --files --hidden -g '!.git'
    selector_type: command
    relation: informs
    read_policy: skip
    used_for: historical specification snapshot; verified the pre-implementation repo had docs/config only before v0-init-utility added Go runtime files
    status: stale
  - id: kit-config
    name: Kit configuration
    type: repo_file
    target: .kit.yaml
    relation: constrains
    read_policy: must
    used_for: specs dir, skills dir, constitution path, confidence goal, and feature naming
    status: active
  - id: readme
    name: README
    type: repo_file
    target: README.md
    relation: informs
    read_policy: skip
    used_for: historical specification baseline; README now documents the implemented v0-init-utility command surface
    status: stale
  - id: gitignore
    name: Git ignore rules
    type: repo_file
    target: .gitignore
    relation: constrains
    read_policy: must
    used_for: verified existing Go, env, and Kit local-state ignore coverage
    status: active
  - id: github-pr-template
    name: GitHub pull request template
    type: repo_file
    target: .github/pull_request_template.md
    relation: informs
    read_policy: evidence
    used_for: expected PR description, testing, and ticket fields
    status: active
  - id: coderabbit-config
    name: CodeRabbit config
    type: repo_file
    target: .coderabbit.yaml
    relation: informs
    read_policy: evidence
    used_for: verified implementation files will be reviewed while docs and agent routing are excluded
    status: active
  - id: testing-reference
    name: Testing reference
    type: repo_doc
    target: docs/references/testing.md
    relation: informs
    read_policy: conditional
    used_for: durable repo-wide testing guidance after v0-init-utility implementation
    status: active
  - id: tooling-reference
    name: Tooling reference
    type: repo_doc
    target: docs/references/tooling.md
    relation: informs
    read_policy: conditional
    used_for: durable repo-wide tooling guidance after v0-init-utility implementation
    status: active
  - id: external-systems-reference
    name: External systems reference
    type: repo_doc
    target: docs/references/external-systems.md
    relation: informs
    read_policy: conditional
    used_for: verified no durable external-system notes exist yet
    status: optional
---
# SPEC

## SUMMARY

Builds the initial Darwin-only `kp` Go CLI for help-first command discovery, prompt printing, exact clipboard verification, an emoji-enhanced interactive `kp list` selector, and a local `kp scaffold` command without a `prompt` subcommand. The feature must embed prompt assets, support user prompt overrides, scaffold repo support files, and stay limited to local build/install/test scope.

## PROBLEM

Prompt use needs a low-friction CLI path that makes available commands obvious, prints the selected prompt for inspection, writes the complete prompt body to the clipboard, and verifies the clipboard content exactly. The no-argument command should show help, while interactive selection belongs under `kp list`.

## GOALS

1. Create a Go 1.22+ module at `github.com/jamesonstone/kp`.
2. Provide a `kp` binary for macOS 13+ on Apple Silicon and Intel Macs.
3. Support bare prompt commands with no `prompt` subcommand: `kp`, `kp <name>`, `kp list`, `kp new <name>`, `kp edit <name>`, `kp rm <name>`, and `kp scaffold`.
4. Embed built-in prompt files for `clarify` and `instructions`.
5. Load user prompts from `$XDG_CONFIG_HOME/kp/prompts/`, defaulting to `~/.config/kp/prompts/`, on every invocation.
6. Let user prompt files shadow built-in prompts by matching name.
7. Copy prompt body text to the clipboard and verify exact clipboard equality.
8. Print direct prompt selections to stdout by default.
9. Support `--copy`, `--print`, `--no-fzf`, `--config <dir>`, `--verbose`, `--help`, and `--version`.
10. Provide deterministic command output and stable documented exit codes.
11. Keep p99 `kp <name> --copy` latency below 100 ms over 50 invocations on M-series macOS.
12. Keep p99 `kp <name>` print-and-copy latency below 100 ms over 50 invocations on M-series macOS.
13. Keep `kp list` interactive startup to `fzf` process launch below 80 ms on M-series macOS.
14. Keep stripped binary size below 10 MB.
15. Keep RSS below 30 MB while the interactive picker is active.
16. Provide `kp scaffold` to create repo support files modeled on Kit-managed scaffold artifacts while excluding direct Kit project state.

## NON-GOALS

1. Do not implement Linux or Windows support.
2. Do not implement a `kp prompt ...` subcommand.
3. Do not publish Homebrew formulae or use `github.com/jamesonstone/homebrew-tap` in this feature.
4. Do not add Goreleaser archives or GitHub release workflows in this feature.
5. Do not add cloud sync, network calls, telemetry, analytics, or external services.
6. Do not implement a full TUI beyond shelling out to `fzf`.
7. Do not add an always-loaded monolithic instruction file.
8. Do not create a plugin framework or generic extension system before a second top-level utility command exists.
9. Do not make prompt frontmatter part of copied, printed, or previewed prompt content.
10. Do not preserve empty prompt files created by a cancelled or empty `kp new <name>` edit flow.
11. Do not embed a built-in `plan` prompt in this feature.
12. Do not make `kp scaffold` create `.kit.yaml`, `.kit/`, global Kit config, `docs/specs/**`, `docs/notes/**`, `docs/PROJECT_PROGRESS_SUMMARY.md`, or `docs/CONSTITUTION.md`.

## USERS

1. Primary user: a local macOS developer who repeatedly pastes long prompt templates into focused applications and needs reliable behavior under macOS Secure Input.
2. Secondary user: a future coding agent or maintainer adding additional local utility commands without rewriting existing prompt command internals.

## SKILLS

Skills are tracked in front matter. `rlm` is selected because this feature started from broad repository analysis and must preserve progressive-disclosure discovery in downstream planning. `parallelization_mode: "rlm"` should be recorded in `PLAN.md` or execution metadata for later implementation planning.

## RELATIONSHIPS

No prior feature relationships apply. `kit map 0001-v0-init-utility` reports no incoming relationships, no outgoing relationships, and no prior feature artifacts beyond this feature's brainstorm and spec.

## DEPENDENCIES

1. Runtime language: Go 1.22 or newer.
2. Module path: `github.com/jamesonstone/kp`.
3. CLI framework: `github.com/spf13/cobra` v1.8 or newer.
4. YAML parser: `gopkg.in/yaml.v3` for optional prompt frontmatter metadata.
5. macOS tools expected on PATH or available through the system: `pbcopy` and `pbpaste`.
6. User-installed interactive dependency: `fzf`, installable with `brew install fzf`.
7. Editor resolution order: `KP_EDITOR`, then `EDITOR`, then `vi`.
8. Local workflow tooling: `Makefile` targets for build, test, install, format, and clean.
9. Release tooling: not required for this feature.
10. Existing repo docs and configuration: `docs/CONSTITUTION.md`, `docs/agents/*`, `.kit.yaml`, `.gitignore`, `.github/pull_request_template.md`, `.coderabbit.yaml`, and `docs/specs/0001-v0-init-utility/BRAINSTORM.md`.

## REQUIREMENTS

1. The system MUST create and maintain a Go module with module path `github.com/jamesonstone/kp` when the feature is implemented.
2. The system MUST expose one binary named `kp`.
3. The system MUST support macOS 13+ on `darwin/arm64` and `darwin/amd64`.
4. The system MUST reject non-Darwin clipboard execution with a clear unsupported-platform error if the binary is built or tested outside Darwin.
5. The system MUST keep command parsing at the edge and keep prompt registry, config resolution, and clipboard behavior in separate internal packages.
6. The system MUST NOT create or register a `prompt` subcommand.
7. The system MUST treat `kp` with no arguments as the help command and MUST render a grouped help page that separates direct prompt commands from prompt-library commands.
8. The system MUST treat `kp <name>` as direct prompt selection for the prompt named `<name>`.
9. The system MUST treat `kp list` as the interactive prompt selector command.
10. The system MUST treat `kp new <name>` as the user prompt creation command.
11. The system MUST treat `kp edit <name>` as the user prompt edit command.
12. The system MUST treat `kp rm <name>` as the user prompt deletion command.
13. The system MUST reserve command names `help`, `list`, `new`, `edit`, `rm`, `scaffold`, `prompt`, and `version` so prompt names cannot shadow commands.
14. The system MUST validate prompt names with `^[a-z][a-z0-9-]*$`.
15. The system MUST reject invalid prompt names with exit code `1`.
16. The system MUST embed `prompts/clarify.md` and `prompts/instructions.md` at compile time.
17. The system MUST NOT embed or ship `prompts/plan.md` in this feature.
18. The system MUST ship `prompts/clarify.md` with this exact source content:

    ```markdown
    ---
    label: Clarify before implementing
    ---
    Clarify before implementing. Stay in planning mode. Ask numbered questions with defaults, assumptions, and uncertainty until >=95% confidence and 0 unresolved. Accept y/n shorthand. Report confidence each batch.
    ```

19. The system MUST ship `prompts/instructions.md` with this exact source content:

    ```markdown
    ---
    label: Coding agent instructions
    ---
    Write as a senior software engineer producing implementation instructions for an autonomous coding agent that will build or modify a software project. Return ONLY a Markdown document. Do not include preamble, postamble, commentary, TODO, FIXME, XXX, placeholders, hedging, or ambiguity words such as “appropriate,” “as needed,” “if relevant,” “where applicable,” “etc.,” “various,” “some,” “might,” “could,” or “consider.”
    
    Use exactly these H2 sections, in this order:
    
    1. Objective
    2. Scope
    3. Assumptions
    4. Requirements
    5. Architecture
    6. Data Models
    7. APIs and Interfaces
    8. Files to Create or Modify
    9. Implementation Steps
    10. Edge Cases
    11. Validation
    12. Tests
    13. Acceptance Criteria
    
    Write concrete instructions that another agent can execute without asking follow-up questions. Preserve existing repository patterns unless a listed requirement explicitly overrides them. When project facts are unknown, instruct the agent to inspect specific files or commands instead of inventing facts. When a design choice is required, pick one defensible default and record it in Assumptions.
    
    Section requirements:
    - Objective: one sentence describing the user-visible outcome.
    - Scope: list In-scope MUST items and Out-of-scope MUST NOT items.
    - Assumptions: state runtime, language version, framework, dependency versions, target OS, repo state, configuration, and preconditions.
    - Requirements: write functional requirements as “the system MUST <verb> <noun> when <condition>.” Quantify non-functional requirements with units: p50/p99 latency in ms, throughput in rps, memory in MB, error budget %, auth model, logging fields, metrics, tracing spans, and security constraints.
    - Architecture: provide one paragraph, a component table with name/responsibility/inputs/outputs, and a numbered interaction sequence.
    - Data Models: include typed target-language definitions with nullability, defaults, validation rules, plus SQL schema with columns, types, indexes, and foreign keys.
    - APIs and Interfaces: for each endpoint, include method, path, request JSON, response JSON, status codes, and error bodies; for each function, include full signature, types, preconditions, and postconditions.
    - Files to Create or Modify: cite repo-rooted paths, purpose, and exported symbols.
    - Implementation Steps: provide atomic ordered steps; each step lists touched files, exact verify command, and expected output.
    - Edge Cases: list trigger condition, expected behavior, and failure mode if mishandled.
    - Validation: give exact manual commands and expected outputs.
    - Tests: list unit tests with name/target/assertions, integration tests with scenario/setup/end state, and at least one negative-path test per public function.
    - Acceptance Criteria: numbered, objectively verifiable outcomes.
    
    Cite exact file paths and symbols such as `internal/api/users.go::CreateUser`, never vague labels such as “the handler.” Include literal JSON or struct examples for every payload. Quantify every claim; write “p99 < 50ms at 1000 rps,” not “fast.”
    ```

20. The system MUST parse optional YAML frontmatter from prompt files.
21. The system MUST use frontmatter `label` only as metadata for list and picker display.
22. The system MUST exclude YAML frontmatter from copied, printed, and previewed prompt bodies.
23. The system MUST default a prompt label from the prompt name when frontmatter is absent.
24. The system MUST load built-in prompts before user prompts.
25. The system MUST load user prompts from `$XDG_CONFIG_HOME/kp/prompts/` when `XDG_CONFIG_HOME` is set.
26. The system MUST load user prompts from `~/.config/kp/prompts/` when `XDG_CONFIG_HOME` is unset.
27. The system MUST support `--config <dir>` to override the config root for the current invocation.
28. The system MUST treat the prompt directory under `--config <dir>` as `<dir>/prompts`.
29. The system MUST load user prompt files on every invocation and MUST NOT cache prompt registry state across process runs.
30. The system MUST let user prompt files shadow built-in prompts by matching prompt name.
31. The system MUST list prompts sorted by name in ascending lexical order.
32. The system MUST make `kp list --plain` print one prompt name per line to stdout.
33. The system MUST make `kp list --verbose` print tab-separated `name`, `label`, and `source` columns to stdout.
34. The system MUST make `source` values exactly `builtin` or `user`.
35. The system MUST make `kp <name> --print` write the prompt body to stdout and perform no clipboard or paste side effects.
36. The system MUST make `kp <name> --copy` write the prompt body to the clipboard, verify clipboard equality, and perform no paste side effects.
37. The system MUST make `kp <name>` write the prompt body to stdout, write the prompt body to the clipboard, verify clipboard equality, and perform no paste side effects.
38. The system MUST read the clipboard back after writing to it and verify exact equality with the expected prompt body before reporting success.
39. The system MUST retry clipboard verification 5 times over 250 ms total before treating verification as failed.
40. The system MUST exit `2` without reporting success when clipboard verification fails.
41. The system MAY include short checksums or byte counts in verbose diagnostic logs, but checksum algorithm is not contractual.
42. The system MUST NOT invoke `osascript` or issue Cmd+V in this feature.
43. The system MUST make `kp list` launch `fzf` over the prompt list with a preview of the prompt body when `fzf` is available.
44. The system MUST decorate the `kp list` selector and root help page with concise emojis where they improve scanability.
45. The system MUST exit `3` with install and fallback instructions when `fzf` is missing and `--no-fzf` is not passed.
46. The system MUST make `kp list --no-fzf` display a numbered list to stderr, read one line from stdin, and run the selected prompt through the default print-and-copy flow.
47. The system MUST exit `1` for non-numeric or out-of-range `--no-fzf` selection input.
48. The system MUST exit `130` when the user cancels `fzf`, the numbered picker, or the editor.
49. The system MUST make `kp new <name>` create a user prompt file at the resolved prompt directory and open it in the resolved editor.
50. The system MUST make `kp new <name>` fail with exit `1` when a prompt with that name already exists as either a built-in or user prompt.
51. The system MUST make `kp new <name>` delete the stub file and exit `1` when the saved prompt body is empty.
52. The system MUST make `kp edit <name>` open the user prompt when a user prompt exists.
53. The system MUST make `kp edit <name>` promote a built-in prompt to a user prompt by copying its source file before opening the editor when the name resolves only to a built-in prompt.
54. The system MUST make `kp edit <name>` print a promotion message to stderr when it promotes a built-in prompt.
55. The system MUST make `kp rm <name>` delete only user prompt files.
56. The system MUST make `kp rm <name>` reject built-in-only prompts with exit `1`.
57. The system MUST resolve editor command in order: `KP_EDITOR`, `EDITOR`, `vi`.
58. The system MUST exit `3` when the resolved editor executable is unavailable.
59. The system MUST create config and prompt directories with owner-only permissions where the platform supports them.
60. The system MUST make `kp scaffold` create missing repo support artifacts in the current working directory by default.
61. The system MUST make `kp scaffold --dir <path>` create missing repo support artifacts under `<path>`.
62. The system MUST make `kp scaffold --dry-run` print planned create, skip, and update actions without writing files.
63. The system MUST make `kp scaffold --force` overwrite scaffold file artifacts except `.gitignore`, which MUST still use append-only missing-pattern updates.
64. The system MUST make `kp scaffold` skip existing scaffold files by default.
65. The system MUST make `kp scaffold` append missing scaffold ignore patterns to `.gitignore` without deleting existing ignore rules.
66. The system MUST make `kp scaffold` include `.env`, `.envrc`, `.coderabbit.yaml`, `.github/pull_request_template.md`, `AGENTS.md`, `CLAUDE.md`, `.github/copilot-instructions.md`, `docs/agents/README.md`, `docs/agents/WORKFLOWS.md`, `docs/agents/RLM.md`, `docs/agents/TOOLING.md`, `docs/agents/GUARDRAILS.md`, `docs/references/README.md`, `docs/references/testing.md`, `docs/references/tooling.md`, and `docs/references/external-systems.md`.
67. The system MUST make `kp scaffold` exclude `.kit.yaml`, global Kit config, `.kit/`, `docs/specs/**`, `docs/notes/**`, `docs/PROJECT_PROGRESS_SUMMARY.md`, and `docs/CONSTITUTION.md`.
68. The system MUST make `kp scaffold` use embedded templates from the `kp` binary and MUST NOT shell out to `kit` or read templates from a local Kit checkout at runtime.
69. The system MUST make `kp scaffold` perform no clipboard, paste, or editor side effects.
70. The system MUST make `kp scaffold` print a concise status table or line set that shows created, updated, and skipped files.
71. The system MUST document stable exit codes in README.
72. The system MUST make `--verbose` emit structured key=value logs to stderr.
73. The system MUST keep normal success output quiet except for requested stdout data and concise stderr status.
74. The system MUST keep `kp <name> --copy` p99 latency below 100 ms over 50 invocations on M-series macOS.
75. The system MUST keep `kp <name>` p99 latency below 100 ms over 50 invocations on M-series macOS.
76. The system MUST keep `kp list` interactive startup to `fzf` process launch below 80 ms on M-series macOS.
77. The system MUST keep stripped binary size below 10 MB.
78. The system MUST keep RSS below 30 MB while the interactive picker is active.
79. The system MUST document the known clipboard race limitation for concurrent invocations.
80. The system MUST document that automatic paste is out of scope for this feature.
81. The system MUST provide local build, test, install, format, and clean commands.
82. The system MUST update README with local installation, command usage, prompt format, scaffold behavior, exit codes, dependencies, and known limitations.
83. The system MUST keep release archive generation, Homebrew publishing, and GitHub release workflows out of this feature.
84. The system MUST render genuine interactive picker cancellation as one whimsical farewell from a finite in-process rotation with a random starting offset, and MUST NOT expose the internal picker-cancellation text or child-process exit status to the user.
85. The system MUST rotate picker farewells without an immediate repeat inside one process and MUST NOT add persistent state, network behavior, or a third-party dependency to vary them across invocations.
86. The system MUST preserve exit `130`, empty stdout, and unchanged clipboard state for picker cancellation while retaining the original diagnostic text and non-cancellation exit classification for operational picker failures.

## ACCEPTANCE

1. `go test ./...` exits `0` on macOS arm64.
2. `go test ./...` exits `0` on macOS amd64 or an amd64-equivalent CI runner.
3. `go build ./cmd/kp` exits `0` on macOS.
4. `./kp --version` exits `0` and prints a version string containing version and commit fields.
5. `./kp` exits `0`, prints help, and performs no prompt registry, clipboard, or paste side effects.
6. `./kp list --verbose` exits `0` and prints two tab-separated rows with `source` equal to `builtin` when no user prompts exist.
7. `./kp clarify --print` exits `0`, writes only the clarify body to stdout, and does not include YAML frontmatter.
8. `./kp clarify --copy` exits `0`, and `pbpaste` equals the clarify body byte-for-byte.
9. `./kp clarify` exits `0`, writes the clarify body to stdout, copies it to the clipboard, verifies the clipboard, and performs no paste side effects.
10. `./kp list --no-fzf` lets the user choose a numbered built-in prompt and runs the selected prompt through the default print-and-copy flow.
11. `PATH` without `fzf` plus `./kp list` exits `3` and prints instructions containing `brew install fzf` and `--no-fzf`.
12. `XDG_CONFIG_HOME=$(mktemp -d) ./kp new test1` opens the editor, saves a non-empty prompt, and creates `$XDG_CONFIG_HOME/kp/prompts/test1.md`.
13. `XDG_CONFIG_HOME=<same-dir> ./kp list --plain` includes `test1`.
14. `XDG_CONFIG_HOME=<same-dir> ./kp rm test1` exits `0` and removes `test1.md`.
15. `XDG_CONFIG_HOME=<same-dir> ./kp edit clarify` with a no-op editor creates `$XDG_CONFIG_HOME/kp/prompts/clarify.md` containing the built-in clarify prompt source.
16. `XDG_CONFIG_HOME=<same-dir> ./kp list --verbose` shows `clarify` with `source` equal to `user` after promotion.
17. `./kp new instructions` exits `1` because built-in prompt names are reserved from `new`.
18. `./kp rm instructions` exits `1` when no user override exists.
19. `./kp "Bad Name"` exits `1`.
20. `./kp nonexistent` exits `1`.
21. `./kp scaffold --dir "$(mktemp -d)"` exits `0`, creates the approved scaffold files, and does not create `.kit.yaml`, `.kit/`, `docs/CONSTITUTION.md`, `docs/specs`, `docs/notes`, or `docs/PROJECT_PROGRESS_SUMMARY.md`.
22. `./kp scaffold --dir "$(mktemp -d)" --dry-run` exits `0`, prints planned scaffold actions, and writes no files.
23. `./kp scaffold --dir <dir>` preserves existing files by default and appends only missing `.gitignore` patterns.
24. `./kp scaffold --dir <dir> --force` overwrites scaffold file artifacts while preserving `.gitignore` append-only behavior.
25. Clipboard verification mismatch exits `2` and does not invoke `osascript`.
26. Missing editor executable exits `3`.
27. User cancellation in `fzf`, numbered picker, or editor exits `130`.
28. `hyperfine './kp clarify --copy' --warmup 5 --runs 50` reports p99 below 100 ms on M-series macOS.
29. `hyperfine './kp clarify' --warmup 5 --runs 50` reports p99 below 100 ms on M-series macOS.
30. `kp list` picker startup instrumentation shows less than 80 ms from process start to `fzf` process launch on M-series macOS.
31. Stripped `kp` binary size is below 10 MB.
32. Interactive picker RSS stays below 30 MB while `fzf` is active.
33. README documents local installation, dependencies, command examples, prompt format, scaffold behavior, exit codes, performance targets, no-paste scope, and concurrent clipboard limitation.
34. `.goreleaser.yaml` and `.github/workflows/release.yaml` are not created by this feature.
35. Cancelling `kp`, `kp list`, or another interactive picker exits `130`, writes one approved whimsical farewell to stderr, leaves stdout empty, and does not expose `picker cancelled` or `exit status 130`.
36. Calling the picker-farewell selector once for every configured style visits every style without an immediate repeat before the rotation wraps.
37. An operational picker failure retains its original diagnostic message and non-cancellation exit classification.

## EDGE-CASES

1. Missing `fzf`: exit `3` unless `--no-fzf` is passed; print `brew install fzf` and `--no-fzf` instructions.
2. Invalid `--no-fzf` input: exit `1`; do not modify the clipboard.
3. User cancels picker: exit `130`; do not modify the clipboard or stdout; render one whimsical rotating farewell instead of internal cancellation or child-process status text.
4. User cancels editor: exit `130`; preserve existing files and delete only a new empty stub created by `kp new`.
5. Empty body after `kp new`: delete the stub file and exit `1`.
6. Built-in-only edit: copy the built-in source file to the user prompt directory, print the promotion path to stderr, then open the editor.
7. Built-in-only remove: exit `1`; do not delete embedded prompts or user files with different names.
8. User override: return `source=user` and use the user prompt body for copy, print, and picker preview.
9. Malformed YAML frontmatter: exit `1` for user prompt operations that load the malformed file and include the file path in stderr.
10. Frontmatter without `label`: use the default label derived from the prompt name.
11. Prompt file with metadata-only content and empty body: treat as empty body for `new` and reject it.
12. Clipboard mismatch after 250 ms: exit `2`; do not report success.
13. Clipboard tools unavailable: exit `2` for copy/read failures.
14. Automatic paste requested by a future command: out of scope for this feature.
15. Concurrent invocations: document that the global clipboard is shared process state and do not implement locking in this feature.
16. Config directory cannot be created or read: exit `3` and include the path in stderr.
17. Config override path is relative: resolve it against the current working directory before deriving `<dir>/prompts`.
18. Non-Darwin runtime: return an unsupported-platform error for clipboard flows and skip Darwin-only integration tests.
19. Reserved command name used as a prompt name: exit `1`.
20. User prompt with uppercase, leading digit, whitespace, slash, dot, or underscore in filename-derived name: reject the prompt name with exit `1`.
21. Existing scaffold file: `kp scaffold` skips it unless `--force` is passed.
22. Existing `.gitignore`: `kp scaffold` appends only missing scaffold patterns and preserves existing rules.
23. Scaffold dry run: `kp scaffold --dry-run` reports planned actions and writes no files.
24. Scaffold target path does not exist: `kp scaffold --dir <path>` creates parent directories as needed for scaffold files.

## OPEN-QUESTIONS

No unresolved questions remain. Confidence is at least 95%, and unresolved assumptions are 0.

## VALIDATION EVIDENCE

- `GH-24` whimsical picker farewells — `PASS`: Go 1.23.4 preflight,
  formatting, focused picker and launcher tests, `go test ./...`,
  `go test -race ./...`, `go vet ./...`, `go build ./...`, `make build`, and
  `git diff --check` passed.
- Built-binary cancellation acceptance — `PASS`: eight real `fzf` no-match
  cancellations produced five distinct approved farewells; every run exited
  `130`, kept stdout empty, and omitted `picker cancelled` and
  `exit status 130`. Numbered-picker EOF produced the same typed exit and one
  approved farewell.
- `kit reconcile --all --dry-run` source-file-size audit — `PASS`: 40 eligible
  handwritten source/test files checked, zero above 300 physical lines. Its 10
  warnings in seven untouched managed instruction files remain outside issue
  #24.
- Hosted pull-request correctness checks — `UNAVAILABLE`: this repository has
  no hosted format, test, race, vet, or build workflow.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI with no
  deployed service or production environment.

## OUTCOME

- Genuine picker cancellation retains its typed sentinel and process exit
  `130`, while the executable renders a finite whimsical farewell set from a
  random starting offset and rotates without immediate in-process repeats.
- Operational picker failures and non-picker cancellations retain their
  original diagnostic output and exit classification.
- Issue #24 carries this refinement on `GH-24`; pull-request delivery does not
  authorize merge.
