---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0004"
  slug: "punchlist-command"
  dir: "0004-punchlist-command"
references:
  - id: "github-issue"
    name: "Add kp punchlist command"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/28"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "merge-command"
    name: "kp merge command"
    type: "feature_artifact"
    target: "docs/specs/0003-merge-command/SPEC.md"
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
delivery_intent: "issue_branch_pr_in_progress"
---
# SPEC

## PURPOSE

Add `kp punchlist` as a built-in prompt command that tells a coding agent how
to drive a living punch list toward resolution by scanning, clustering, and
fixing shared causes rather than treating items as an independent ticket queue.

## CONTEXT

- `kp` already models bare prompt names as embedded Markdown assets. Adding a
  built-in prompt requires no Cobra subcommand or new runtime abstraction.
- Operators, testers, and users maintain punch lists as observational evidence:
  defects, confusing UX, missing functionality, workflow friction, and
  regressions. Those notes are not engineering specifications.
- The command must remain local prompt output. It does not locate a punch-list
  document, call a spreadsheet or issue tracker, mutate statuses, or implement
  product changes itself.
- `kp merge` established the delivery pattern this feature reuses: one embedded
  asset, existing print/copy/list/help/launcher paths, exact-output tests, and
  README discovery updates.

## REQUIREMENTS

- Ship an embedded prompt named `punchlist` with label
  `Punch list control loop`.
- `kp punchlist`, `kp punchlist --print`, and `kp punchlist --copy` must use
  the existing prompt execution and clipboard-verification behavior without
  special-case command code.
- The prompt must instruct the receiving agent to:
  - discover the punch-list document, source repositories, worklanes, review,
    CI, deployment, runtime, logs, and project-specific agent instructions
    from the current environment without assuming GitHub, Git, a spreadsheet
    provider, a CI system, or a deployment platform;
  - re-read the entire punch list before every state-changing decision;
  - treat items as observations, cluster related symptoms, and prefer one
    coherent shared fix over repeated item-specific patches;
  - preserve human notes, never invent statuses, and use existing
    engineering-note conventions including an author prefix when the user
    message, project config, or punch-list document supplies one;
  - remain in planning mode until confidence is at least 95% with zero
    unresolved material questions;
  - distinguish implemented, merged, deployed, and validated states and not
    request re-testing until the change is available in the validation
    environment; prohibit re-testing of a failed or regressed deployment, then
    request re-testing after a successful corrective deployment;
  - load repo-local guardrails and work-lane or delivery-gating guidance when
    they exist, complete read-only safety recon, and obtain delivery-consent
    confirmation before issue, branch, worktree, staging, commit, push, PR,
    merge, or deployment mutation; permit merge only after direct user
    authorization or an accepted bounded plan that names the exact authorized
    PR set;
  - use the project's established worklane mechanism, reuse active work that
    already covers the cluster, and keep recoverable traceability from items
    through change, deployment, and validation; and
  - begin with the required planning report covering punch-list state,
    proposed clusters, execution order, items awaiting validation, and
    numbered clarifications.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `punchlist` through the existing registry
  path.
- Automated tests must pin the exact prompt output and updated built-in
  ordering.
- README command and built-in tables must document `kp punchlist`.
- No network dependency, punch-list host integration, status writer, or new
  third-party dependency is in scope.

## ACCEPTED PLAN

1. Add `prompts/punchlist.md` as an embedded built-in with the approved
   control-loop procedure.
2. Pin exact CLI output, the approved body hash, required contract phrases,
   and discovery order in prompt and command tests.
3. Document the command in README and record the feature in
   `docs/PROJECT_PROGRESS_SUMMARY.md`.
4. Run formatting, focused prompt tests, the full Go test suite, race tests,
   vet, both builds, isolated `kp punchlist --print` acceptance, diff hygiene,
   and the affected source-size audit.
5. Self-review, commit with the repository contract, push `GH-28`, and open
   one ready pull request that closes issue #28.

## DECISIONS

- Use a built-in prompt asset instead of a dedicated Cobra command because the
  existing bare-name registry already supplies print, copy, help, launcher, and
  user-override behavior.
- Keep author-prefix handling inside the prompt. A root `--author` flag would
  be special-case command code and would apply to unrelated prompts.
- Keep the procedure environment-agnostic. The receiving agent must discover
  the punch-list document and local conventions rather than assuming GitHub or
  a spreadsheet. When the environment is Kit-managed, name `GUARDRAILS.md` and
  `work-lane-gating` as the concrete delivery-consent gate rather than inventing
  a generic mutation path.
- Pin the approved body with SHA-256 in addition to comparing CLI output to
  the asset so accidental prompt edits fail tests without duplicating the full
  body in a Go file that would exceed the 300-line source limit.
- Leave extra positional context out of scope. Operators can paste additional
  notes in the agent chat after copying the prompt.

## DISCOVERIES

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; `punchlist` sorts after `pr`.
- `kp merge` is the closest precedent: one asset, no Cobra command, exact
  output tests, and README/progress-summary documentation.
- `internal/cmd/prompt_execution_test.go` already pins the large `merge` body.
  A separate `punchlist_prompt_test.go` keeps both files under the 300-line
  handwritten source/test limit.
- PR #29 review found that merge and the definition of done could authorize
  delivery mutations without work-lane gating or merge consent. The prompt now
  requires read-only recon and delivery consent before those mutations, and
  merge only after a named authorized PR set.
- A later PR #29 review found that "do not request re-testing" after any
  deployment failure blocked validation of a successful corrective deployment.
  The prompt now prohibits re-testing only of the failed or regressed
  deployment and requires re-testing after a successful corrective deployment.

## VALIDATION

- `go test ./internal/prompt ./internal/cmd` — `PASS`.
- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./...` — `PASS` across every package.
- `go test -race ./...` — `PASS` across every package.
- `go vet ./...` — `PASS`.
- `go build ./...` — `PASS`.
- `make build` — `PASS`; produced `bin/kp`.
- Isolated CLI acceptance with an empty temporary config directory — `PASS`:
  `kp punchlist --print` emitted SHA-256
  `7efd70f586a3f959f1365d8bc9af95ff4789f99f19d9466266b9d3dee6335423`,
  `list --plain` included `punchlist` in sorted order, and `--help` showed
  `Punch list control loop`.
- `git diff --check` — `PASS`.
- Source-file-size audit of tracked Go files — `PASS`: 40 eligible handwritten
  source/test files checked, 0 above 300 physical lines.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.
- PR #29 review repair — `PASS`: focused prompt and command tests, full tests,
  race tests, vet, both builds, isolated `kp punchlist --print` hash
  `7efd70f586a3f959f1365d8bc9af95ff4789f99f19d9466266b9d3dee6335423`,
  `git diff --check`, and the 40-file source-size audit were rerun after
  scoping re-test prohibition to failed or regressed deployments and requiring
  re-testing after a successful corrective deployment.

## OUTCOME

- `kp punchlist` is an embedded, overrideable built-in prompt exposed through
  the existing print, copy, list, help, and launcher paths.
- The prompt requires environment discovery without assuming a tracker or
  platform, whole-list clustering, a 95% clarification gate, worklane reuse,
  preservation of human notes, delivery-consent before repository mutation,
  merge only of an exact authorized PR set, separate implemented, merged,
  deployed, and validated states, and re-testing after a successful corrective
  deployment rather than after a failed or regressed one.
- Issue #28 tracks delivery on `GH-28`; merge is not authorized by this work.

## REPOSITORY MEMORY

- Created this living specification because the command establishes durable
  punch-list operating behavior for coding agents.
- Updated `README.md` with command discovery and built-in prompt documentation.
- Updated `docs/PROJECT_PROGRESS_SUMMARY.md` with this feature's intent,
  approach, status, and pointer.
- The existing v0 and merge feature artifacts remain unchanged because this is
  a separate public command surface.
