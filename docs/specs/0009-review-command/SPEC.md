---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0009"
  slug: "review-command"
  dir: "0009-review-command"
references:
  - id: "github-issue"
    name: "Add kp review command"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/58"
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
    used_for: "built-in prompt command pattern, exact-output pinning, and local CLI-only runtime contract"
    status: "active"
  - id: "ship-command"
    name: "kp ship command"
    type: "feature_artifact"
    target: "docs/specs/0008-ship-command/SPEC.md"
    relation: "informs"
    read_policy: "must"
    used_for: "built-in prompt command pattern and distinction from delivery authorization"
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

Add `kp review` as a built-in prompt command that tells a coding agent to
perform an independent correctness review of the current branch or pull
request against the target branch, without making changes.

## CONTEXT

- `kp` already models bare prompt names as embedded Markdown assets. Adding a
  built-in prompt requires no Cobra subcommand or new runtime abstraction.
- `kp pr` establishes a new worklane. It does not review an existing diff.
- `kp punchlist` clusters operational observations. It is not an independent
  code-correctness review.
- `kp ship` pre-authorizes delivery. It must not become a review prompt.
- The command must remain local prompt output. It does not post GitHub reviews
  or mutate the repository itself.
- `kp plan` and `kp ship` established the delivery pattern this feature reuses:
  one embedded asset, existing print/copy/list/help/launcher paths, exact-output
  tests, and README discovery updates.

## REQUIREMENTS

- Ship an embedded prompt named `review` with label
  `Independent branch/PR review`.
- `kp review`, `kp review --print`, and `kp review --copy` must use the existing
  prompt execution and clipboard-verification behavior without special-case
  command code.
- Do not add a review alias. `review` is the canonical name.
- Do not expand `pr`, `punchlist`, `ship`, or `merge`. Worklane creation,
  operational punch lists, delivery authorization, and merge coordination stay
  on those commands.
- The prompt must instruct the receiving agent to:
  - review the current branch or pull request against the target branch;
  - treat the existing implementation as untrusted;
  - focus primarily on correctness;
  - inspect surrounding code, tests, and contracts rather than reviewing only
    the syntactic diff;
  - run the relevant tests, linters, type checks, and other repository
    validation available to it;
  - make no changes yet;
  - return findings ordered by severity with severity, location, problem, risk,
    smallest fix, and proving tests; and
  - avoid manufactured findings, stating correctness explicitly when warranted
    and naming residual validation gaps.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `review` through the existing registry path.
- Automated tests must pin the exact prompt output and updated built-in
  ordering.
- README command and built-in tables must document `kp review`.
- No network dependency, review-host integration, or new third-party dependency
  is in scope.

## ACCEPTED PLAN

1. Add `prompts/review.md` as an embedded built-in with the approved
   independent-review procedure.
2. Pin exact CLI output, the approved body hash, required contract phrases,
   and discovery order in prompt and command tests.
3. Document the command in README and record the feature in
   `docs/PROJECT_PROGRESS_SUMMARY.md`.
4. Run formatting, focused prompt tests, the full Go test suite, race tests,
   vet, both builds, isolated `kp review --print` acceptance, diff hygiene,
   and the affected source-size audit.
5. Self-review, commit with the repository contract, push `GH-58`, and open
   one ready pull request that closes issue #58.

## DECISIONS

- Use a built-in prompt asset instead of a dedicated Cobra command because the
  existing bare-name registry already supplies print, copy, help, launcher, and
  user-override behavior.
- Name the command `review`. The user requested `kp review`; a second alias
  would split discovery without changing behavior.
- Keep `pr`, `punchlist`, `ship`, and `merge` unchanged. Review is a read-only
  correctness pass, not worklane creation, punch-list clustering, delivery
  authorization, or merge orchestration.
- Preserve the supplied body. The review contract is the copied prompt for
  receiving agents.
- Pin the approved body with SHA-256 in addition to comparing CLI output to
  the asset so accidental prompt edits fail tests without duplicating the full
  body in a Go file that would exceed the 300-line source limit.
- Record topology as `single-lane, because tightly coupled and high-overlap:
  one embedded prompt, existing registry paths, and documentation updates in
  a single delivery lane`.

## DISCOVERIES

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; `review` sorts after `punchlist` and before `ship`.
- `kp plan` and `kp ship` are the closest precedents: one asset, no Cobra
  command, hash pin, required-phrase tests, and README/progress-summary
  documentation.
- A separate `review_prompt_test.go` keeps the handwritten test file under the
  300-line limit.
- Help, list, launcher, and override paths discover built-ins dynamically, so
  no command-registration change is required.
- `review` is not a reserved prompt name.

## VALIDATION

- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./internal/prompt ./internal/cmd` — `PASS`.
- `go test ./...` — `PASS` across every package.
- `go test -race ./...` — `PASS` across every package.
- `go vet ./...` — `PASS`.
- `go build ./...` — `PASS`.
- `make build` — `PASS`; produced `bin/kp`.
- Isolated CLI acceptance with an empty temporary config directory — `PASS`:
  `kp review --print` emitted SHA-256
  `ab9233c17c44e46658d68f82d90e229319806d5085bbde7d315e7ae159234999`,
  `list --plain` included `review` between `punchlist` and `ship`, and `--help`
  showed `Independent branch/PR review`.
- `git diff --check` — `PASS`.
- Source-file-size audit of version-eligible Go files — `PASS`: 47 eligible
  handwritten source/test files checked, 0 above 300 physical lines. Changed
  Go files: `review_prompt_test.go` 87, `root_test.go` 179, `builtin_test.go` 84,
  `registry_test.go` 211.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.

## OUTCOME

- `kp review` is an embedded, overrideable built-in prompt exposed through the
  existing print, copy, list, help, and launcher paths.
- The prompt requires an independent, untrusted-implementation correctness
  review of the current branch or pull request, runs available repository
  validation, and returns severity-ordered findings without making changes.
- `pr`, `punchlist`, `ship`, and `merge` are unchanged. Issue #58 tracks
  delivery on `GH-58`.

## REPOSITORY MEMORY

- Created this living specification because the command establishes a durable
  split from `pr`, `punchlist`, `ship`, and `merge`: independent correctness
  review versus worklane creation, operational punch-list clustering, delivery
  authorization, and merge coordination.
- Constitution curation is not required: adding one built-in prompt is
  feature-local and does not change project-wide invariants already stated
  for bare prompt commands.
