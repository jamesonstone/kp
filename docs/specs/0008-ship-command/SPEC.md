---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0008"
  slug: "ship-command"
  dir: "0008-ship-command"
references:
  - id: "github-issue"
    name: "Add kp ship command"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/48"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "goal-command"
    name: "kp goal command"
    type: "feature_artifact"
    target: "docs/specs/0007-goal-command/SPEC.md"
    relation: "informs"
    read_policy: "must"
    used_for: "built-in prompt command pattern, exact-output pinning, and distinction from executable /goal construction"
    status: "active"
  - id: "merge-command"
    name: "kp merge command"
    type: "feature_artifact"
    target: "docs/specs/0003-merge-command/SPEC.md"
    relation: "informs"
    read_policy: "must"
    used_for: "existing merge-coordination prompt that still requires an authorized PR frontier"
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

Add `kp ship` as a built-in prompt command that pre-authorizes a coding agent
to complete the current task thread through the full delivery lifecycle,
including in-scope pull-request merge, without expanding `continue`, `goal`,
or `merge`.

## CONTEXT

- `kp` already models bare prompt names as embedded Markdown assets. Adding a
  built-in prompt requires no Cobra subcommand or new runtime abstraction.
- `kp continue` tells an agent to keep working autonomously, but it does not
  grant branch, pull-request, merge, or deployment authority.
- `kp merge` coordinates already-scoped pull requests and still requires an
  authorized frontier before merge.
- `kp goal` constructs a user-confirmed executable `/goal` in planning mode.
  `kp ship` is a different command: it copies a thread-scoped delivery
  authorization whose body begins with `/goal` so a receiving Cursor agent can
  arm a long-running goal while shipping.
- The command must remain local prompt output. It does not create branches,
  merge pull requests, or deploy.

## REQUIREMENTS

- Ship an embedded prompt named `ship` with label
  `Pre-authorize task delivery`.
- `kp ship`, `kp ship --print`, and `kp ship --copy` must use the existing
  prompt execution and clipboard-verification behavior without special-case
  command code.
- Do not add `deliver`, `land`, `autopilot`, or another alias. `ship` is the
  canonical name chosen by the user.
- Keep `continue`, `goal`, and `merge` unchanged.
- The prompt body must preserve the supplied authorization contract,
  including the leading `/goal`, and must instruct the receiving agent to:
  - treat the current task thread as pre-authorized for the full delivery
    lifecycle of code produced or modified as part of that work;
  - create, update, and push branches; create and update pull requests;
    address review feedback and CI failures; merge in-scope pull requests
    when required checks pass; merge dependent pull requests in order;
    deploy when that is part of the established repository workflow; and
    perform routine repository operations needed to complete delivery;
  - not ask for additional authorization for individual in-thread PR
    merges, including confirmation of specific PR numbers or commit SHAs;
  - scope that authorization only to changes required for this task;
  - refuse unrelated pre-existing PR merges, protection bypass, ignored
    failing required checks, and destructive or irreversible operations
    outside the normal delivery workflow; and
  - continue autonomously until the task is delivered or a material
    blocker cannot be resolved safely without new information.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `ship` through the existing registry path.
- Automated tests must pin the exact prompt output and updated built-in
  ordering.
- README command and built-in tables must document `kp ship`.
- No network dependency or new third-party dependency is in scope.

## ACCEPTED PLAN

1. Add `prompts/ship.md` as an embedded built-in with the approved
   thread-scoped delivery authorization.
2. Pin exact CLI output, the approved body hash, required contract phrases,
   and discovery order in prompt and command tests.
3. Document the command in README and record the feature in
   `docs/PROJECT_PROGRESS_SUMMARY.md`.
4. Run formatting, focused prompt tests, the full Go test suite, race tests,
   vet, both builds, isolated `kp ship --print` acceptance, diff hygiene,
   and the affected source-size audit.
5. Self-review, commit with the repository contract, push `GH-48`, and open
   one ready pull request that closes issue #48.

## DECISIONS

- Use a built-in prompt asset instead of a dedicated Cobra command because the
  existing bare-name registry already supplies print, copy, help, launcher, and
  user-override behavior.
- Name the command `ship`. The user chose that name over `deliver`, `land`,
  `autopilot`, `greenlight`, `full-send`, and `finish`.
- Keep `continue`, `goal`, and `merge` unchanged. Autonomy without delivery
  authority, executable-goal construction, and authorized merge coordination
  remain separate prompts.
- Preserve the supplied body, including the leading `/goal`. That token is
  part of the copied prompt for receiving agents; it does not make `kp ship`
  an alias of `kp goal`.
- Pin the approved body with SHA-256 in addition to comparing CLI output to
  the asset so accidental prompt edits fail tests without duplicating the full
  body in a Go file that would exceed the 300-line source limit.
- Record topology as `single-lane, because tightly coupled and high-overlap:
  one embedded prompt, existing registry paths, and documentation updates in
  a single delivery lane`.

## DISCOVERIES

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; `ship` sorts after `punchlist`.
- `kp goal` landed on `main` as feature `0007` while this work was scoped, so
  this feature is `0008` and the built-in inventory is eleven prompts after
  `ship` is added.
- `kp punchlist` and `kp goal` are the closest precedents: one asset, no Cobra
  command, hash pin, required-phrase tests, and README/progress-summary
  documentation.
- A separate `ship_prompt_test.go` keeps the handwritten test file under the
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
  `kp ship --print` emitted SHA-256
  `ddd0fecaa4a333e69785d4d838cca092c5dd6806e2b9fa9295e6cdc3894268b8`,
  `list --plain` included `ship` after `punchlist`, and `--help` showed
  `Pre-authorize task delivery`.
- `git diff --check` — `PASS`.
- Source-file-size audit of version-eligible Go files — `PASS`: 45 eligible
  handwritten source/test files checked, 0 above 300 physical lines. Changed
  Go files: `ship_prompt_test.go` 80, `root_test.go` 171, `builtin_test.go` 82,
  `registry_test.go` 211.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.

## OUTCOME

- `kp ship` is an embedded, overrideable built-in prompt exposed through the
  existing print, copy, list, help, and launcher paths.
- The prompt pre-authorizes the current task thread for the full delivery
  lifecycle, scoped only to that task, and keeps going until delivery or a
  material blocker.
- `continue`, `goal`, and `merge` are unchanged. Issue #48 tracks delivery on
  `GH-48`; merge is not authorized by this work.

## REPOSITORY MEMORY

- Created this living specification because the command establishes a durable
  split from `continue`, `goal`, and `merge`: thread-scoped delivery
  pre-authorization versus autonomy, goal construction, and merge
  coordination.
- Constitution curation is not required: adding one built-in prompt is
  feature-local and does not change project-wide invariants already stated
  for bare prompt commands.
