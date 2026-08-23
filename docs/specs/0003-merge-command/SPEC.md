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

## REQUIREMENTS

- Ship an embedded prompt named `merge` with label
  `Dependency-ordered PR merge`.
- `kp merge`, `kp merge --print`, and `kp merge --copy` must use the existing
  prompt execution and clipboard-verification behavior without special-case
  command code.
- The prompt must instruct the receiving agent to:
  - operate only on the exact repository, pull-request, and expected-head set
    directly authorized by the user or covered by a bounded plan explicitly
    accepted by the user;
  - build and output a Mermaid dependency DAG plus topological wave table;
  - distinguish ordinary merge dependencies from edges that require deployed
    or production-accepted predecessors;
  - require authoritative directional evidence for dependency edges and treat
    shared files, symbols, configuration, environments, or databases as
    coupling to inspect rather than proof of direction;
  - bind readiness to exact repository, head, base, actor, merge method,
    policy, reviews, checks, dependencies, infrastructure effects, and recovery;
  - classify every node as `MERGE_READY`, `BLOCKED`, or `UNKNOWN` and fail
    closed on missing, stale, pending, unproven-skipped, cyclic, provisional,
    conflicted, or ambiguous evidence;
  - select only the zero-unmet-dependency ready frontier, maximize concurrency
    among independent nodes, and prioritize downstream unlocks and critical-path
    reduction;
  - serialize dependency chains and same-base, deployment-coupled, or otherwise
    interacting operations;
  - revalidate immediately before every wave, re-observe after every mutation,
    isolate failures to the failed node and its dependents, and recompute the
    graph;
  - preserve repository safeguards and required merge queues;
  - distinguish exact-head merge authority from source-repair authority,
    prefer separately authorized in-place updates for routine remediation that
    remains inside the existing pull request's issue and scope, and invalidate
    prior readiness and merge authority whenever the head changes;
  - reserve replacement pull requests for material scope or architecture
    changes, original heads that cannot be updated safely, or explicit
    repository-policy or user requirements;
  - classify migration and deployment recovery from actual behavior, including
    reversible redeploy, schema-retaining cutback or forward-fix, and genuinely
    irreversible or destructive changes;
  - inventory protected workloads and bind their invariants, safe activation
    state, exact validation, literal result, artifact identity, and recovery;
  - load repo-local rules before merge actions, give Kit-managed rules priority
    over generic GitHub or plugin defaults, and require `kit aws verify` plus the
    verified configured profile for AWS-dependent evidence or actions without
    ambient-credential fallback; and
  - report merge, hosted workflow, deployment, runtime, production acceptance,
    recovery, and rollback as separate claims; and
  - define `MERGED`, `DEPLOYED`, and `ACCEPTED` separately and record the exact
    post-merge commit, base, method, actor, observation time, and containment
    evidence needed to support those claims.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `merge` through the existing registry path.
- Automated tests must pin the exact prompt output and updated built-in ordering.
- README command and built-in tables must document `kp merge`.
- No network dependency, background controller, GitHub API integration, merge
  executor, or new third-party dependency is in scope.

## ACCEPTED PLAN

1. Replace the recursive corrective-PR default with separately authorized,
   scope-preserving repair on the existing pull-request head between waves.
2. Preserve the exact-head merge freeze by returning any changed head to
   `UNKNOWN` until fresh checks, review, revalidation, and later exact-head
   merge authorization are complete.
3. Update the durable merge rule, workflow, exact-output test, and canonical
   feature documentation together so the instruction surfaces remain aligned.
4. Run formatting, focused prompt tests, the full Go test suite, race tests,
   vet, both builds, isolated `kp merge --print` acceptance, diff hygiene, and
   the affected source-size audit.
5. Self-review, commit with the repository contract, push `GH-26`, and open one
   ready pull request that closes issue #26.

## DECISIONS

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

## OUTCOME

- `kp merge` is an embedded, overrideable built-in prompt exposed through the
  existing print, copy, list, help, and launcher paths.
- Its sectioned output (purpose, scope and authorization, definitions, pre-wave
  checklist, execution algorithm, deployment gates, example, durable state, and
  evidence) builds a Mermaid PR DAG and wave table, binds exact-current
  readiness, requires evidence-cited edges and proven graph completeness before
  concurrency, maximizes safe independent concurrency and downstream unlocks,
  revalidates every wave, enforces explicit user acceptance plus repo-local Kit
  and verified-profile AWS gates, isolates failures, brakes on no-progress
  repetition, gates deployment on reversibility class and pre-declared
  acceptance, and separates merge from deployment and production evidence.
- The issue #22 refinement additionally requires authoritative directional
  evidence, classifies recovery from actual behavior, protects existing lanes
  through exact-artifact compatibility gates, and records `MERGED`, `DEPLOYED`,
  and `ACCEPTED` without relaxing the established graph, authorization,
  revalidation, credential, rollback, or no-progress gates.
- The issue #26 refinement keeps routine, scope-preserving remediation on the
  original pull request under separate repair authority, makes all changed
  heads re-enter as `UNKNOWN`, and reserves replacement PRs for material or
  otherwise unsafe changes.
- PR #17 delivered the original command and PR #19 delivered the first prompt
  correctness refinement; PR #23 delivered the issue #22 safety refinement.
  Issue #26 tracks the current in-place remediation correction on `GH-26`;
  merge is not authorized by this work.

## REPOSITORY MEMORY

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
