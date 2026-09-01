---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0007"
  slug: "goal-command"
  dir: "0007-goal-command"
references:
  - id: "github-issue"
    name: "Add kp goal command"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/46"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "plan-command"
    name: "kp plan command"
    type: "feature_artifact"
    target: "docs/specs/0006-plan-command/SPEC.md"
    relation: "informs"
    read_policy: "must"
    used_for: "built-in prompt command pattern, exact-output pinning, and downstream plan-convergence contract"
    status: "active"
  - id: "v0-init-utility"
    name: "v0 init utility"
    type: "feature_artifact"
    target: "docs/specs/0001-v0-init-utility/SPEC.md"
    relation: "constrains"
    read_policy: "must"
    used_for: "bare prompt-name execution, clipboard verification, list/help/launcher discovery"
    status: "active"
delivery_intent: "issue_branch_pr_ready"
---
# SPEC

## PURPOSE

Add `kp goal` as a built-in prompt command that tells a coding agent how to
turn incomplete or ambiguous engineering intent into a user-confirmed
executable `/goal` before implementation planning.

## CONTEXT

- `kp` already models bare prompt names as embedded Markdown assets. Adding a
  built-in prompt requires no Cobra subcommand or new runtime abstraction.
- `kp clarify` remains a short ambiguity check when the objective is already
  mostly understood.
- `kp plan` remains downstream plan convergence. This feature teaches `plan`
  to consume an accepted `/goal` as its intent contract.
- The command must remain local prompt output. It does not research
  repositories, write goals, or execute implementation itself.
- `kp plan` established the delivery pattern this feature reuses: one
  embedded asset, existing print/copy/list/help/launcher paths, exact-output
  tests, and README discovery updates.

## REQUIREMENTS

- Ship an embedded prompt named `goal` with label
  `Construct an executable goal`.
- `kp goal`, `kp goal --print`, and `kp goal --copy` must use the existing
  prompt execution and clipboard-verification behavior without special-case
  command code.
- Keep `clarify` and `plan` as distinct phases. Do not expand `clarify` into
  evidence-backed goal convergence, and do not let `plan` replace goal
  construction.
- The prompt must instruct the receiving agent to:
  - stay in planning mode and refuse file edits, delivery artifacts,
    deployment, state-changing setup, and goal execution;
  - treat the initial request as a thesis, not automatically as a complete
    specification;
  - maintain one accumulated goal model;
  - resolve discoverable facts before asking the user;
  - ask only material questions in the smallest useful numbered batch, with
    defaults and `1y 2n` shorthand;
  - synthesize every discovery and answer into the same model;
  - challenge circular dependencies, false PASS conditions, and unsafe
    failure behavior when material;
  - treat safely creatable prerequisites as bootstrap work rather than
    blocking preconditions;
  - distinguish implementation, deployment, runtime readiness, and
    acceptance only when those distinctions are material;
  - converge after one adversarial review exposes no new material
    ambiguity; and
  - return one cohesive `/goal` for user confirmation, not an
    implementation plan.
- Amend `prompts/plan.md` so an accepted `/goal` is the intent contract and
  is consumed unchanged. Research must not redefine the objective; a material
  contradiction is reported and requires an explicit revised user decision
  before the intent contract may change.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `goal` through the existing registry path.
- Automated tests must pin the exact prompt output and updated built-in
  ordering.
- README command and built-in tables must document `kp goal`.
- No network dependency or new third-party dependency is in scope.

## ACCEPTED PLAN

1. Add `prompts/goal.md` as an embedded built-in with the approved
   Evidence-Backed Goal Convergence procedure.
2. Add the accepted `/goal` consumer paragraph to `prompts/plan.md`.
3. Pin exact CLI output, approved body hashes, required contract phrases,
   and discovery order in prompt and command tests.
4. Document the command in README and record the feature in
   `docs/PROJECT_PROGRESS_SUMMARY.md`.
5. Run formatting, focused prompt tests, the full Go test suite, race tests,
   vet, both builds, isolated `kp goal --print` acceptance, diff hygiene,
   and the affected source-size audit.
6. Self-review, commit with the repository contract, push `GH-46`, and open
   one ready pull request that closes issue #46.

## DECISIONS

- Use a built-in prompt asset instead of a dedicated Cobra command because the
  existing bare-name registry already supplies print, copy, help, launcher, and
  user-override behavior.
- Keep `clarify` unchanged. Quick ambiguity checks and evidence-backed goal
  construction are complementary commands, not one expanded prompt.
- Encode the `goal` → accepted `/goal` → `plan` lifecycle in `plan.md` rather
  than leaving it as README-only guidance. `kp plan` consumes the confirmed
  `/goal` unchanged and may not redefine it from research; a material
  contradiction requires an explicit revised user decision.
- Pin approved bodies with SHA-256 in addition to comparing CLI output to the
  assets so accidental prompt edits fail tests without duplicating the full
  body in a Go file that would exceed the 300-line source limit.
- Leave extra positional context out of scope. Operators can paste a request
  in the agent chat after copying the prompt.

## DISCOVERIES

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; `goal` sorts after `continue` and before `merge`.
- `kp plan` is the closest precedent: one asset, no Cobra command, hash pin,
  required-phrase tests, and README/progress-summary documentation.
- A separate `goal_prompt_test.go` keeps the handwritten test file under the
  300-line limit.
- Help, list, launcher, and override paths discover built-ins dynamically, so
  no command-registration change is required.

## VALIDATION

- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./internal/prompt ./internal/cmd` — `PASS`.
- `go test ./...` — `PASS` across every package.
- `go test -race ./...` — `PASS` across every package.
- `go vet ./...` — `PASS`.
- `go build ./...` — `PASS`.
- `make build` — `PASS`; produced `bin/kp`.
- Isolated CLI acceptance with an empty temporary config directory — `PASS`:
  `kp goal --print` emitted SHA-256
  `ac318360ff5f848bf8fc9d673ebb22515f4aa13d6bfe3ec1aeda1e2043a8e308`,
  `list --plain` included `goal` between `continue` and `merge`, and `--help`
  showed `Construct an executable goal`.
- `git diff --check` — `PASS`.
- Source-file-size audit of tracked Go files — `PASS`: 45 eligible handwritten
  source/test files checked, 0 above 300 physical lines. Changed Go files:
  `goal_prompt_test.go` 77, `plan_prompt_test.go` 84, `root_test.go` 167,
  `builtin_test.go` 81, `registry_test.go` 211.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.

## OUTCOME

- `kp goal` is an embedded, overrideable built-in prompt exposed through the
  existing print, copy, list, help, and launcher paths.
- The prompt keeps the agent in planning mode, maintains one accumulated goal
  model, researches discoverable facts before asking, asks only material
  questions, challenges circular dependencies and false PASS conditions, and
  returns one user-confirmed executable `/goal`.
- `kp plan` now consumes an accepted `/goal` as its intent contract and does
  not redefine that objective from research alone.
- `clarify` is unchanged. Issue #46 tracks delivery on `GH-46`; merge is not
  authorized by this work.

## REPOSITORY MEMORY

- Created this living specification because the command establishes durable
  intent-to-goal convergence behavior and a durable split from `clarify` and
  `plan`.
- Constitution curation is not required: adding one built-in prompt is
  feature-local and does not change project-wide invariants already stated
  for bare prompt commands.
