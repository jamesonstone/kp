---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0003"
  slug: "merge-command"
  dir: "0003-merge-command"
references:
  - id: "github-issue"
    name: "Add kp merge command"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/16"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "loopc-control-model"
    name: "LoopC control model"
    type: "external"
    target: "https://github.com/jamesonstone/loopc/blob/main/docs/CONSTITUTION.md"
    relation: "informs"
    read_policy: "must"
    used_for: "observe-act-remeasure discipline, evidence binding, and fail-closed brakes"
    status: "active"
  - id: "merge-controller-pr-forest"
    name: "Merge Controller PR-forest control loop"
    type: "external"
    target: "https://github.com/jamesonstone/merge-controller/blob/main/docs/specs/0015-pr-forest-control-loop/SPEC.md"
    relation: "informs"
    read_policy: "must"
    used_for: "dependency graph, readiness, frontier selection, deployment, and runtime semantics"
    status: "active"
  - id: "github-pr-merge-rule"
    name: "Repository pull-request merge rule"
    type: "internal"
    target: "docs/references/rules/github-pr-merge.md"
    relation: "constrains"
    read_policy: "must"
    used_for: "authorization, readiness, concurrency, and evidence boundaries"
    status: "active"
  - id: "github-issue-refinement"
    name: "Refine merge safety instructions"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/22"
    relation: "supports"
    read_policy: "must"
    used_for: "refinement issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "github-issue-in-place-remediation"
    name: "Prefer in-place PR remediation during merge orchestration"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/26"
    relation: "supports"
    read_policy: "must"
    used_for: "in-place remediation requirements and delivery traceability"
    status: "active"
  - id: "github-issue-context-aware-merge"
    name: "Make the merge prompt concise and context-aware"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/30"
    relation: "supports"
    read_policy: "must"
    used_for: "context derivation, concise execution, recovery, and delivery traceability"
    status: "active"
delivery_intent: "issue_branch_pr_in_progress"
---
# SPEC

## PURPOSE

Add `kp merge` as a built-in prompt command that tells a coding agent how to
create and execute an evidence-backed, graphical pull-request merge plan. The
prompt must preserve dependency order while maximizing safe concurrency and
prioritizing work that unlocks the greatest downstream dependency closure.

## CONTEXT

- `kp` already models bare prompt names as embedded Markdown assets. Adding a
  built-in prompt requires no Cobra subcommand or new runtime abstraction.
- The referenced ChatGPT conversation identifies Merge Controller as the
  concrete PR-orchestration service and LoopC as its control-loop foundation.
- LoopC establishes observe, type, select, act, and remeasure; missing or stale
  evidence fails closed, and ineffective repetition must stop or escalate.
- Merge Controller establishes an evidence-backed PR forest, exact-head
  readiness, dependency confidence, quiet revalidation, downstream-unlock
  priority, and separate merge, deployment, runtime, production, and rollback
  evidence.
- The repository's current merge rule allows concurrent merge waves only for
  independent `MERGE_READY` nodes. Dependency chains, same-base operations
  whose ordering matters, and coupled deployment or release state remain
  serialized.
- Running `kp merge` remains local prompt output. It does not discover PRs,
  call GitHub, grant merge authority, or perform a merge itself.
- Issue #30 replaces the detailed standalone procedure with a concise overlay
  that derives known PR state from conversation and repository context while
  delegating durable policy to Kit-managed repository rules.

## REQUIREMENTS

- Keep `merge` as an embedded, overrideable built-in prompt using the existing
  print, copy, list, help, launcher, and clipboard paths.
- Replace only the prompt contract; do not add a Cobra command, network access,
  GitHub integration, background controller, or dependency.
- Derive the exact PR set, current heads/bases, dependencies, permitted merge
  methods, deployment targets, acceptance gates, and authorization boundaries
  from conversation and repository context without asking the user to restate
  discoverable facts or expanding scope.
- Require repository-local merge, infrastructure, orchestration, and completion
  rules to remain authoritative.
- Name `docs/agents/GUARDRAILS.md` and `work-lane-gating` before mutation,
  require read-only safety recon, and verify the exact lane-choice question was
  explicitly answered for the current lane. Generic approval cannot substitute
  for lane consent, while an already-recorded exact choice is not repeated.
- Name `testing-and-environment-validation.md` and the project testing reference
  before implementation or validation so prompt brevity cannot bypass the
  repository's required environment procedure.
- Build the dependency/deployment graph from authoritative evidence and fail
  closed as `BLOCKED` or `UNKNOWN` when readiness is incomplete.
- Present one consolidated approval request before the first merge or covered
  infrastructure mutation unless equivalent exact approval remains valid.
- Prohibit infrastructure deletion, destruction, purge, destructive
  replacement, and state removal inside the merge workflow; isolate them for a
  separate explicitly authorized task.
- Run one complete preflight immediately before a consequential mutation and
  refresh only after material state change or freshness expiry.
- Keep one primary coordinator responsible for graph changes, authority,
  recovery, waves, and acceptance. When supported, permit lower-cost or
  lower-capability agents only for exact bounded ready-node merges and
  deployment monitoring.
- Parallelize only nodes independent in both source and deployment effects;
  serialize shared bases, services, environments, databases, migrations,
  queues, and acceptance gates.
- Attempt bounded autonomous recovery within approved scope, rerun only affected
  evidence, and stop blind or non-progressing retries.
- Prefer event-driven waits or bounded backoff and omit unchanged polling.
- Preserve merge, hosted CI, deployment, runtime, and production acceptance as
  separate claims.
- Use repository-required completion vocabulary when present; otherwise emit
  `SUCCESS|PARTIAL|BLOCKED|FAILURE`, a one-to-three-sentence result, and
  `none` or at most three copy-ready next steps.
- Pin exact prompt output in tests and keep README discovery aligned.
- Keep every changed handwritten source and test file at or below 300 lines.

## ACCEPTED PLAN

1. Use issue #30, branch `GH-30`, and its canonical non-primary worktree.
2. Replace only `prompts/merge.md` with the concise context-aware contract.
3. Update the exact-output test, README, and this living specification without
   changing command execution behavior.
4. Run focused and complete local validation, source-size audit, diff hygiene,
   and isolated prompt acceptance before ready-PR delivery.

## DECISIONS

- Accepted: Kit-managed repository rules own durable merge and infrastructure
  policy; `kp merge` is a concise conversational invocation layer.
- Accepted: current conversation and repository evidence define the PR graph,
  so the prompt asks only for materially missing authority or state.
- Accepted: one meaningful preflight replaces redundant unchanged-state
  rechecking, while material drift still invalidates affected evidence.
- Accepted: lower-capability agents may perform only exact bounded mechanical
  merges or deployment monitoring; the primary coordinator retains judgment.
- Accepted: infrastructure deletion is never reconciled inside a release wave;
  it becomes a separate explicit task and approval boundary.
- Accepted: terminal output is status-first and concise, with repository-local
  completion vocabulary taking precedence.

- Use a built-in prompt asset instead of a dedicated Cobra command because the
  existing bare-name registry already supplies print, copy, help, launcher, and
  user-override behavior.
- Require a Mermaid DAG and wave table so dependency order and parallel frontiers
  are both visible, rather than emitting an unstructured ordered list.
- Interpret "maximized for complexity" as maximizing safe concurrency and
  prioritizing the frontier node that unlocks the most downstream work or
  shortens the critical path.
- Permit concurrent merges only for proven-independent ready nodes. This keeps
  the prompt aligned with repository policy while retaining LoopC's requirement
  to remeasure after material state changes.
- Keep the prompt concise by expressing operational invariants rather than
  copying the full LoopC, Merge Controller, or Kit rulesets.
- Bind accepted-plan authority to the exact repository, pull requests, and
  expected heads, and carry repo-local Kit and AWS identity gates into the
  execution instructions rather than assuming generic platform defaults.
- Treat shared artifacts as coupling evidence only. Dependency direction needs
  an explicit base relationship, producer/consumer contract, prerequisite,
  canonical graph, or user or repository requirement.
- Classify deployment recovery from what can actually be restored or retained;
  schema migration is not itself proof of irreversibility.
- Preserve protected existing behavior through exact-artifact compatibility
  gates before activation rather than inferring safety from additive or
  default-off source.
- Define post-merge terminal states and their audit fields so merge, artifact
  publication, deployment, activation, and acceptance cannot collapse into one
  overstated claim.
- Freeze exact heads only for merge authorization. Routine corrective work may
  update the existing PR under separate repair authority, but the new head must
  lose prior readiness and receive fresh checks, review, revalidation, and
  exact-head merge authorization before merging.
- Keep replacement pull requests exceptional. Material scope or architecture
  change, an unsafe or inaccessible original head, or explicit policy or user
  direction justifies replacement; minor in-scope fixes do not.

## DISCOVERIES

- Issue #30 showed that the prior prompt duplicated the complete Kit ruleset,
  encouraged unnecessary rechecking, and forced callers to restate context the
  receiving agent already possessed. A short overlay preserves behavior while
  reducing prompt noise and policy drift.

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; the implementation change is one asset plus tests and docs.
- `kit status` reported an unrelated managed refresh. Its dry-run would change
  12 managed instruction and ruleset files, so it is intentionally excluded
  from issue #16. The same dry-run reported 40 eligible handwritten source/test
  files and zero files above the 300-line limit before implementation.
- PR #17 review identified two valid safety gaps in the prompt: accepted-plan
  authority was not explicitly bound to repository and expected heads, and the
  execution step omitted repo-local Kit and verified-profile AWS precedence.
  Both were corrected in the prompt and its exact-output test.
- Issue #18 restructured the prompt into named sections but left the pinned
  exact-output test unchanged, so `TestMergePromptPrintsApprovedInstructions`
  failed on `GH-18` until this revision. The prompt asset and its pin must
  always change together.
- The `GH-18` revision closed further correctness gaps: sibling nodes sharing a
  just-merged base kept stale check evidence; the dependency graph had no
  completeness or edge-evidence requirement; deferred merges such as auto-merge
  escaped wave revalidation; deployment edges defined no acceptance signal,
  reversibility class, expand-then-contract ordering, deployed-identity check,
  or proven rollback mechanism; an authorized node depending on an unauthorized
  node had no defined outcome; and the observe/act loop had no no-progress
  brake despite LoopC's stop-or-escalate requirement.
- The prompt body stores backticks as `~` inside the pinned Go raw string
  literal because a Mermaid fence cannot appear in one directly.
- The issue #22 comparison found four improvements worth adopting from a longer
  generic procedure: directional-edge precision, recovery classification by
  actual behavior, protected-workload gates, and explicit terminal audit
  states. Its optional graph, deferred-merge, and weaker freshness behavior
  remain excluded.
- Issue #26 showed that the prior in-place-repair exclusion overextended the
  exact-head merge freeze. Preserving the original PR for routine repairs avoids
  recursive corrective PRs while still invalidating every prior readiness,
  review, check, and merge-authorization claim tied to the old head.

## VALIDATION

- `go test ./internal/prompt ./internal/cmd` — `PASS`.
- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./...` — `PASS` across every package.
- `go test -race ./...` — `PASS` across every package.
- `go vet ./...` — `PASS`.
- `go build ./...` — `PASS`.
- `make build` — `PASS`; produced `bin/kp`.
- Isolated CLI acceptance with an empty temporary config directory — `PASS`:
  `kp merge --print` emitted the exact approved body, `list --plain` included
  `merge` in sorted order, and `--help` showed its label.
- `git diff --check` — `PASS`.
- `kit reconcile --all --dry-run` source-file-size audit — `PASS`: 40 eligible
  handwritten source/test files checked, 0 above 300 physical lines.
- Hosted pull-request correctness checks — `UNAVAILABLE`: the repository has no
  hosted format, test, race, vet, or build workflow. This pre-existing gap is
  recorded in `docs/references/testing.md`; local results are not represented
  as hosted evidence.
- Production validation — `NOT_APPLICABLE`: `kp` is a local CLI and this change
  adds no deployed service or external integration.
- PR #17 review repair — `PASS`: focused prompt output, prompt and command
  packages, full tests, race tests, vet, both builds, isolated CLI acceptance,
  feature validation, diff hygiene, and the source-file-size audit were rerun
  after the authorization and credential-boundary fixes.
- `GH-18` prompt revision — `PASS`: `gofmt`, `go vet ./...`, `go test ./...`,
  `go test -race ./...`, `go build ./...`, `make build`, `git diff --check`, and
  isolated CLI acceptance (`merge --print`, `list --plain`, `--help`) were rerun
  against an empty temporary config directory after the rewrite and its updated
  exact-output pin.
- `GH-22` prompt refinement — `PASS`: Go 1.23.4 preflight, `gofmt`, focused
  prompt and command tests, `go test ./...`, `go test -race ./...`,
  `go vet ./...`, `go build ./...`, `make build`, `git diff --check`, and
  isolated CLI acceptance (`merge --print`, `list --plain`, and `--help`) all
  passed. The exact printed prompt SHA-256 was
  `e00544329def274721ac52c801a750bbe257a87d6e9f693070314ba9e488b560`.
  `kit reconcile --all --dry-run` checked 40 eligible handwritten source/test
  files with zero above 300 lines; its 10 warnings in seven untouched managed
  instruction files remain outside issue #22. Hosted pull-request correctness
  checks remain `UNAVAILABLE`, and production validation is `NOT_APPLICABLE`
  for this local prompt-only CLI refinement.
- `GH-26` in-place remediation refinement — `PASS`: Go 1.23.4 preflight,
  focused exact-output coverage, `gofmt`, `go test ./...`,
  `go test -race ./...`, `go vet ./...`, `go build ./...`, `make build`,
  `kit check --all`, `git diff --check`, and isolated empty-config CLI
  acceptance (`merge --print`, `list --plain`, and `--help`) all passed. The
  exact printed prompt SHA-256 was
  `46284e8be6d7735c12ba3b80fafbc439bf03dd9c297e91e67bbd20c3e6c57666`.
  `kit reconcile --all --dry-run` checked 40 eligible handwritten source/test
  files with zero above 300 lines and reported 10 pre-existing managed-refresh
  warnings in seven untouched files. `kit check --project` failed with the
  same nine pre-existing blocking findings on both `GH-26` and clean `main`, so
  that repo-wide scaffold refresh remains outside issue #26. Hosted
  pull-request correctness checks remain `UNAVAILABLE`, and production
  validation is `NOT_APPLICABLE` for this local prompt-only CLI refinement.
- `GH-30` context-aware prompt replacement — `PASS`: focused prompt/command
  tests, formatting, full tests, race tests, vet, both builds, isolated
  empty-config `merge --print`, `list --plain`, and help acceptance, plus
  `git diff --check` all pass. Kit's audit checks 41 eligible handwritten
  source/test files with zero above 300 physical lines; its ten pre-existing
  managed-refresh warnings remain outside issue #30. Gitleaks scans 47 commits
  and 1.39 MB with no leaks. Hosted correctness checks remain `UNAVAILABLE`,
  and production validation is `NOT_APPLICABLE` for this local prompt-only
  replacement.
- PR #31 review repair — `PASS`: the mutation preflight now names Guardrails,
  work-lane gating, read-only safety recon, and exact current-lane consent while
  preserving the no-redundant-recheck requirement. A follow-up review finding
  also added the mandatory testing and environment-validation references.

## OUTCOME

- `kp merge` remains an embedded, overrideable prompt with unchanged command,
  printing, copying, listing, help, launcher, and override behavior.
- The prompt now derives scoped PR and deployment state from available context,
  loads repository rules, asks once for any missing consolidated authority,
  parallelizes only fully independent nodes, confines recovery to approved
  lanes, and suppresses redundant rechecks and polling.
- A primary coordinator owns decisions while explicitly supported
  lower-capability agents may handle bounded ready-node merges and deployment
  monitoring.
- Terminal responses are concise and keep merge, CI, deployment, runtime, and
  production acceptance separate.

## REPOSITORY MEMORY

Issue #30 records the durable split between Kit's canonical policy registry and
KP's concise informal invocation layer. The prompt replacement intentionally
changes no command behavior.

- Created this living specification because the command establishes durable
  merge-orchestration behavior synthesized from multiple projects.
- Updated `README.md` with command discovery and built-in prompt documentation.
- Replaced placeholder testing guidance in `docs/references/testing.md` with
  the project's actual local commands, non-applicable environments, and hosted
  correctness-check gap.
- Updated `docs/PROJECT_PROGRESS_SUMMARY.md` with this feature's intent,
  approach, status, and pointer.
- Issue #26 records the durable distinction between exact-head merge authority
  and separately authorized in-place remediation across the prompt, merge rule,
  merge workflow, exact-output coverage, and feature rationale.
- The existing v0 feature artifacts remain unchanged because this is a separate
  public command and policy surface.
