# CONSTITUTION

## Authority

This document is the canonical project contract for `kp`.

Authority order:

1. safety, filesystem, privacy, and permission constraints
2. the current user request
3. this constitution
4. feature `SPEC.md`
5. feature `PLAN.md`
6. feature `TASKS.md`
7. feature `BRAINSTORM.md`
8. existing repo conventions

Repo-local Markdown under `docs/` is the system of record. `AGENTS.md`,
`CLAUDE.md`, and `.github/copilot-instructions.md` are routing tables only.
Keep them short, aligned with `docs/agents/*`, and free of durable project
manual content.

## Current Project State

`kp` is a Darwin-first Go CLI for local prompt utilities. The repository
contains the agent instruction system, Kit configuration, feature docs, a Go
module, runtime source, embedded prompt assets, tests, local Make targets, and
user-facing README documentation.

The current implementation ships the initial `v0-init-utility` surface: bare
prompt commands such as `kp clarify`, an interactive `kp list` selector,
prompt CRUD commands, user prompt overrides under the resolved config
directory, and exact clipboard verification on macOS. Release packaging,
Homebrew publishing, Linux support, Windows support, and automatic paste remain
out of scope until a future feature spec changes that contract.

## Product Vision

Build `kp` as a small, reliable command-line utility set for prompt-oriented
developer workflows.

Long-term direction:

- Favor local, scriptable CLI behavior over services, dashboards, or hidden
  background state.
- Keep commands easy to compose with shell pipelines and other developer tools.
- Make output deterministic unless a feature explicitly requires nondeterminism.
- Keep startup and execution fast enough for repeated interactive use.
- Preserve a compact repo instruction layer so agents and humans can work from
  the same source of truth.
- Grow through explicit feature specs instead of speculative framework work.

## PRINCIPLES

- Correctness comes before speed.
- Minimal, production-ready code is preferred over speculative framework work.
- Explicit control flow, clear names, and useful errors are preferred over
  cleverness.
- Runtime behavior should be local, deterministic, and scriptable unless a
  feature spec explicitly requires otherwise.
- Documentation and implementation must stay aligned through feature artifacts
  and the project progress summary.

## Goals

- Provide minimal, production-ready CLI prompt utilities.
- Optimize for correctness, clarity, and performance in that order.
- Keep implementation code explicit, idiomatic, and easy to inspect.
- Keep dependencies small, justified, and documented by purpose.
- Make behavior changes traceable through feature docs and progress summaries.
- Keep validation evidence close to the change that required it.

## Non-Goals

- Do not build a web app, hosted service, daemon, or API server unless a future
  feature spec explicitly changes the product direction.
- Do not add an always-loaded monolithic instruction file.
- Do not create generic frameworks, plugin systems, or abstraction layers before
  a real repeated need exists.
- Do not use `.claude/skills` as canonical discovery input.
- Do not treat uncommitted local secrets, generated state, or machine-specific
  cache files as project state.
- Do not add broad dependencies for convenience when the standard library or a
  small local function is sufficient.

## Architectural Patterns

### CLI First

The primary runtime surface is expected to be command-line execution.

Rules for future implementation:

- Keep command parsing and process I/O at the edge.
- Put reusable behavior in small internal functions with explicit inputs and
  outputs.
- Prefer pure transformations for prompt/text manipulation.
- Avoid global mutable state except for process-level configuration that is
  explicitly initialized at startup.
- Use stdin, stdout, stderr, files, and exit codes in conventional CLI ways.
- Make failure modes observable with clear errors and non-zero exits.

### Minimal Core

Prefer a small core with boring control flow.

- Add modules only when they clarify ownership or reduce real duplication.
- Keep public APIs narrow. If a symbol is only used locally, do not export it.
- Remove dead code, unused exports, and unused public surface as part of the
  change that makes them obsolete.
- Keep data structures plain unless stronger modeling is needed for correctness.
- Prefer explicit branching and error handling over clever control flow.

### Documentation-Governed Features

Feature work lives under `docs/specs/<feature>/` when the work is new,
substantial, cross-cutting, or already has feature docs.

Canonical feature flow:

1. `BRAINSTORM.md`
2. `SPEC.md`
3. `PLAN.md`
4. `TASKS.md`
5. implementation
6. validation and documentation reflection

Do not mix multiple features in one `docs/specs/<feature>/` directory.

### Progressive Context Loading

Use RLM-style discovery for broad or uncertain work:

1. identify the immediate decision
2. load the smallest relevant artifact
3. extract only facts needed for that decision
4. recurse only while uncertainty remains
5. stop loading once the decision is supported

Prefer repo-local docs before global or vendor instructions. Treat secondary
global inputs as fallback context only after repo-local docs are exhausted.

## Code Style And Naming

Project-wide rules:

- Prefer explicit over clever.
- Avoid premature generalization.
- Use idiomatic patterns for the language in use.
- Keep code minimal and production-ready.
- Use clear names that describe behavior or ownership.
- Use comments sparingly, only where they explain non-obvious intent or
  constraints.
- Format code with the language-standard formatter before completion.

Naming conventions:

- CLI command names and feature slugs use lowercase kebab case.
- Environment variables use uppercase snake case.
- Markdown feature artifacts use canonical uppercase names:
  `BRAINSTORM.md`, `SPEC.md`, `PLAN.md`, and `TASKS.md`.
- Durable reference files under `docs/references/` use lowercase kebab case.
- Feature directories use the Kit numbering convention from `.kit.yaml`:
  four-digit numeric prefix, hyphen separator, descriptive slug.

Go-specific conventions:

- Use `gofmt` and idiomatic package names.
- Keep packages small and purpose-specific.
- Return errors explicitly with useful context.
- Use the standard library unless a third-party dependency has a clear,
  documented purpose.
- Keep command parsing and process exit handling at the CLI edge.
- Keep reusable behavior in `internal/*` packages with narrow exported
  contracts.

## Dependencies And Tooling

Current verified dependencies and tools:

- Kit: governs feature docs, instruction scaffolding, `docs/specs`, skills
  routing, and the constitution path via `.kit.yaml`.
- Go module `github.com/jamesonstone/kp`: defines the runtime dependency model
  and builds the `kp` binary from `cmd/kp`.
- Cobra (`github.com/spf13/cobra`): provides CLI command parsing, help/version
  behavior, and future top-level command registration.
- YAML v3 (`gopkg.in/yaml.v3`): parses optional prompt frontmatter metadata.
- macOS clipboard tools: `pbcopy` and `pbpaste` implement local copy and exact
  read-back verification on Darwin.
- `fzf`: optional user-installed interactive picker dependency. Missing `fzf`
  must produce a clear fallback instruction unless `--no-fzf` is used.
- `$KP_EDITOR`, `$EDITOR`, or `vi`: editor resolution order for user prompt
  creation and editing.
- Make: local convenience targets for format, test, build, install, and clean.
- GitHub repository metadata: `.github/copilot-instructions.md` routes Copilot
  to repo-local docs, and `.github/pull_request_template.md` defines minimal PR
  evidence fields.
- CodeRabbit: `.coderabbit.yaml` configures review path filters and currently
  excludes docs and top-level agent routing files.
- direnv-style local environment loading: `.envrc` calls `dotenv_if_exists` and
  `.env` is ignored. Do not commit secrets.
- Go-oriented ignore patterns: `.gitignore` excludes common Go build, test, and
  coverage artifacts.

Rules for adding dependencies:

- Add a dependency only when it materially improves correctness,
  maintainability, portability, or performance.
- Document the dependency's purpose near the change that introduces it.
- Prefer standard tooling and standard libraries for CLI behavior.
- Avoid runtime network dependencies unless a feature spec requires them.

## CHANGE CLASSIFICATION

Classify work before acting:

- Use the spec-driven flow for new features, substantial behavioral changes,
  cross-component changes, or work that already has feature docs.
- Use the ad hoc flow for contained bug fixes, small docs updates, dependency
  updates, config changes, and narrow refinements.
- Promote ad hoc work to the spec-driven flow when scope grows enough to affect
  requirements, interfaces, validation strategy, or multiple feature artifacts.
- Update existing feature docs when ad hoc work changes behavior covered by
  those docs.

## Process Rules

### Work Classification

Classify before acting.

Spec-driven work applies to new features, substantial behavioral changes,
cross-component changes, or work with existing feature docs.

Ad hoc work applies to contained bug fixes, reviews, dependency updates, config
changes, and small refinements.

Ad hoc work that touches behavior covered by existing feature docs must update
those docs unless the change is purely mechanical.

### Spec-Driven Execution

For feature work, execution order is:

1. read the relevant `TASKS.md` entry
2. read the linked `PLAN.md` section
3. read the linked `SPEC.md` requirement
4. read `BRAINSTORM.md` only for unresolved rationale
5. inspect implementation files
6. update canonical docs first when implementation changes behavior,
   requirements, or approach
7. implement
8. validate
9. update progress tracking

Run the readiness gate before implementation. Challenge the docs for
contradictions, ambiguity, hidden assumptions, missing failure modes, task gaps,
and scope creep. If the gate fails, fix the canonical docs first.

### Ad Hoc Execution

For ad hoc work:

1. inspect relevant files before editing
2. use existing repo patterns
3. make the smallest production-ready change
4. verify with the smallest relevant check
5. update practical docs when behavior, commands, setup, or expectations change

Do not create feature docs for ad hoc work unless scope grows enough to require
the spec-driven track.

### Progress Tracking

`docs/PROJECT_PROGRESS_SUMMARY.md` must reflect the highest completed artifact
per feature at all times.

Update it whenever a feature artifact is created, completed, renamed, removed,
or advanced. If no feature docs exist, the summary must explicitly say so.

Highest completed artifact order:

1. none
2. `BRAINSTORM.md`
3. `SPEC.md`
4. `PLAN.md`
5. `TASKS.md`
6. implementation validated
7. reflected complete

### Skills And References

For feature-scoped work, read canonical front matter `skills` first. Fall back
to the legacy `SPEC.md` `## SKILLS` table only when front matter is absent.

Open each referenced `SKILL.md` and use only the skills needed for the current
decision. Record materially used docs, skills, and references in canonical
feature front matter when feature docs are touched.

Use `docs/references/*` for durable repo-wide context. Use
`docs/references/rules/*` only when linked or directly relevant.

### Dispatch And Subagents

Keep broad or noisy discovery in RLM first. Use `kit dispatch` or subagents only
after workstreams are narrow enough to predict overlap.

Rules:

- cluster work conservatively by likely touched files and interfaces
- parallelize only independent, low-overlap areas
- serialize dependent or cross-cutting work
- keep the main agent responsible for synthesis, integration, validation, and
  communication
- keep isolated worktrees flat under `~/worktrees/`

## Validation Rules

- Never claim tests passed unless they ran.
- Never claim files were inspected unless they were inspected.
- Never guess file contents, APIs, dependencies, or behavior.
- If validation cannot run, state why.
- Fix relevant lint and test failures before calling work complete.
- For docs-only changes, validate by reviewing formatting and source-of-truth
  consistency.
- For code changes, run the smallest checks that exercise the changed behavior.

## CONSTRAINTS


### Kit-Managed Baseline Rules

<!-- BEGIN KIT-MANAGED BASELINE RULES -->
- Treat `docs/CONSTITUTION.md` as the canonical project contract.
- Keep `AGENTS.md`, `CLAUDE.md`, and `.github/copilot-instructions.md` aligned with the repo-local docs tree.
- Use native agent planning for research, clarification, design, and implementation planning.
- Before implementation, inspect code and repository memory; create or adopt `SPEC.md` when material rationale exists.
- After validation, curate feature rationale, project invariants, reusable practices, and domain knowledge into their scope-appropriate canonical documents.
- Allow a justified `not required` repository-memory decision when code and tests preserve the complete durable truth.
- Keep every version-control-eligible handwritten implementation/source and test file at 300 physical lines or less.
- Before delivery, audit the complete affected source/test scope; whole-project reconcile and scheduled maintenance audit the entire repository.
- Exclude documentation files, all `docs/**`, all `.kit/**`, `.kit.yaml`, ignored files, vendored dependencies, and proven generated files.
- Split oversized files by semantic responsibility while preserving stable public entry points and behavior; never use minification or arbitrary numbered chunks to claim compliance.
<!-- END KIT-MANAGED BASELINE RULES -->
## Non-Negotiable Constraints

- Keep changes minimal and reversible.
- Prefer explicit error handling over silent failure.
- Keep `AGENTS.md`, `CLAUDE.md`, and `.github/copilot-instructions.md` aligned
  with the repo-local docs tree.
- Do not run `git add` or `git commit` without explicit approval.
- Do not run `coderabbit --prompt-only` unless explicitly requested or approved.
- Do not commit `.env`, `.envrc`, Kit cache/state, build artifacts, coverage
  files, or machine-local scratch files.
- Do not expand top-level injected instruction files into full manuals.
- Do not treat generated `.kit/state.json` or task bundles as canonical
  behavior documentation.

## Definitions

- `kp`: the CLI prompt utilities project in this repository.
- `Kit`: the local workflow/spec tooling configured by `.kit.yaml`.
- `RLM`: the repo's just-in-time context-routing pattern for progressive
  disclosure.
- `Feature`: a scoped unit of product or behavior work under
  `docs/specs/<feature>/`.
- `Artifact`: one of the canonical feature documents or tracked implementation
  completion states used by `docs/PROJECT_PROGRESS_SUMMARY.md`.
- `Reference`: durable repo-local context under `docs/references/`.
- `Secondary global input`: a non-repo instruction or skill source used only
  after repo-local docs are exhausted.
