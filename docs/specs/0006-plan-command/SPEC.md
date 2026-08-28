---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0006"
  slug: "plan-command"
  dir: "0006-plan-command"
references:
  - id: "github-issue"
    name: "Add kp plan command"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/44"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "punchlist-command"
    name: "kp punchlist command"
    type: "feature_artifact"
    target: "docs/specs/0004-punchlist-command/SPEC.md"
    relation: "informs"
    read_policy: "must"
    used_for: "built-in prompt command pattern, exact-output pinning, and local CLI-only runtime contract"
    status: "active"
  - id: "v0-init-utility"
    name: "v0 init utility"
    type: "feature_artifact"
    target: "docs/specs/0001-v0-init-utility/SPEC.md"
    relation: "constrains"
    read_policy: "must"
    used_for: "bare prompt-name execution, clipboard verification, list/help/launcher discovery"
    status: "active"
  - id: "openai-prompting-guidance"
    name: "OpenAI model prompting guidance"
    type: "external"
    target: "https://developers.openai.com/api/docs/guides/latest-model"
    relation: "informs"
    read_policy: "must"
    used_for: "outcome-first prompting, autonomy boundaries, and stopping conditions"
    status: "active"
delivery_intent: "issue_branch_pr_ready"
---
# SPEC

## PURPOSE

Add `kp plan` as a built-in prompt command that tells a coding agent how to
drive the current implementation plan to a practical, evidence-backed local
maximum while remaining in planning mode.

## CONTEXT

- `kp` already models bare prompt names as embedded Markdown assets. Adding a
  built-in prompt requires no Cobra subcommand or new runtime abstraction.
- `kp clarify` establishes intent with a short clarification loop. It does not
  perform plan convergence, adversarial review, or produce a replacement plan.
- The original v0 feature deferred `prompts/plan.md`. This feature introduces
  that command as a medium-length outcome-first prompt rather than a short
  stub or a full process-prescriptive master prompt.
- The command must remain local prompt output. It does not research
  repositories, write plans, or execute implementation itself.
- `kp punchlist` established the delivery pattern this feature reuses: one
  embedded asset, existing print/copy/list/help/launcher paths, exact-output
  tests, and README discovery updates.

## REQUIREMENTS

- Ship an embedded prompt named `plan` with label
  `Drive plan to implementation-ready`.
- `kp plan`, `kp plan --print`, and `kp plan --copy` must use the existing
  prompt execution and clipboard-verification behavior without special-case
  command code.
- Do not add `plan-max` or another alias. `plan` is the canonical name.
- Do not expand `clarify`. Intent-setting stays on `clarify`; convergence and
  adversarial verification stay on `plan`.
- The prompt must instruct the receiving agent to:
  - stay in planning mode and refuse file edits, delivery artifacts,
    deployment, and plan execution;
  - target the simplest decision-complete plan another coding agent can
    implement without making material design decisions;
  - penalize plan bloat and recommend a change only when it closes a concrete
    requirement, ambiguity, contradiction, assumption, failure-mode,
    authority, or validation gap;
  - research progressively from the owning repository outward, resolve
    discoverable facts from evidence, treat historical conclusions as
    hypotheses, and stop when further retrieval is editorial;
  - ask the smallest numbered question batch only when evidence cannot
    decide, with defaults, assumptions, uncertainty, and confidence;
  - review across the listed adversarial lenses and iterate silently until
    one complete adversarial pass and a final verification pass produce no
    new material recommendation;
  - finalize only at the stated stop conditions, including zero unresolved
    material questions and at least 95% evidence-backed goal coverage;
  - stop with the exact blocker when evidence or authority prevents
    convergence rather than hiding the gap with assumptions; and
  - return one complete replacement plan plus a concise local-maximum audit.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `plan` through the existing registry path.
- Automated tests must pin the exact prompt output and updated built-in
  ordering.
- README command and built-in tables must document `kp plan`.
- No network dependency, plan-host integration, or new third-party dependency
  is in scope.

## ACCEPTED PLAN

1. Add `prompts/plan.md` as an embedded built-in with the approved
   medium-length convergence procedure.
2. Pin exact CLI output, the approved body hash, required contract phrases,
   and discovery order in prompt and command tests.
3. Document the command in README and record the feature in
   `docs/PROJECT_PROGRESS_SUMMARY.md`.
4. Run formatting, focused prompt tests, the full Go test suite, race tests,
   vet, both builds, isolated `kp plan --print` acceptance, diff hygiene,
   and the affected source-size audit.
5. Self-review, commit with the repository contract, push `GH-44`, and open
   one ready pull request that closes issue #44.

## DECISIONS

- Use a built-in prompt asset instead of a dedicated Cobra command because the
  existing bare-name registry already supplies print, copy, help, launcher, and
  user-override behavior.
- Name the command `plan`, not `plan-max`. The user requested `kp plan`; the
  original deferred built-in name is `plan`; a second alias would split
  discovery without changing behavior.
- Keep `clarify` unchanged. Short intent-setting and medium plan-convergence
  are complementary commands, not one expanded prompt.
- Use the medium-length body rather than the deferred v0 stub. The short form
  lacks convergence discipline; a longer master prompt over-prescribes process.
- Pin the approved body with SHA-256 in addition to comparing CLI output to
  the asset so accidental prompt edits fail tests without duplicating the full
  body in a Go file that would exceed the 300-line source limit.
- Leave extra positional context out of scope. Operators can paste a plan or
  specification in the agent chat after copying the prompt.

## DISCOVERIES

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; `plan` sorts after `parentthread` and before `pr`.
- `kp punchlist` is the closest precedent: one asset, no Cobra command, hash
  pin, required-phrase tests, and README/progress-summary documentation.
- A separate `plan_prompt_test.go` keeps the handwritten test file under the
  300-line limit.
- Help, list, launcher, and override paths discover built-ins dynamically, so
  no command-registration change is required.

## VALIDATION

- `go test ./internal/prompt ./internal/cmd` — `PASS`.
- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./...` — `PASS` across every package.
- `go test -race ./...` — `PASS` across every package.
- `go vet ./...` — `PASS`.
- `go build ./...` — `PASS`.
- `make build` — `PASS`; produced `bin/kp`.
- Isolated CLI acceptance with an empty temporary config directory — `PASS`:
  `kp plan --print` emitted SHA-256
  `5d36e5ce4e46b70c46df111cac6a3bfe5e7af3f2b6bba81ec72504c974e00d1a`,
  `list --plain` included `plan` between `parentthread` and `pr`, and `--help`
  showed `Drive plan to implementation-ready`.
- `git diff --check` — `PASS`.
- Source-file-size audit of tracked Go files — `PASS`: 43 eligible handwritten
  source/test files checked, 0 above 300 physical lines.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.

## OUTCOME

- `kp plan` is an embedded, overrideable built-in prompt exposed through the
  existing print, copy, list, help, and launcher paths.
- The prompt keeps the agent in planning mode, researches progressively, asks
  only when evidence cannot decide, iterates silently across adversarial
  lenses, and returns one implementation-ready replacement plan plus a
  local-maximum audit.
- `clarify` is unchanged. Issue #44 tracks delivery on `GH-44`; merge is not
  authorized by this work.

## REPOSITORY MEMORY

- Created this living specification because the command establishes durable
  planning-mode convergence behavior for coding agents and a durable split
  from `clarify`.
- Constitution curation is not required: adding one built-in prompt is
  feature-local and does not change project-wide invariants already stated
  for bare prompt commands.
