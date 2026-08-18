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
  - operate only on an exact user-authorized PR set;
  - build and output a Mermaid dependency DAG plus topological wave table;
  - distinguish ordinary merge dependencies from edges that require deployed
    or production-accepted predecessors;
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
  - preserve repository safeguards and required merge queues; and
  - report merge, hosted workflow, deployment, runtime, production acceptance,
    recovery, and rollback as separate claims.
- Prompt listing, verbose listing, grouped help, launcher discovery, and user
  override behavior must include `merge` through the existing registry path.
- Automated tests must pin the exact prompt output and updated built-in ordering.
- README command and built-in tables must document `kp merge`.
- No network dependency, background controller, GitHub API integration, merge
  executor, or new third-party dependency is in scope.

## ACCEPTED PLAN

1. Add one embedded `prompts/merge.md` asset containing the concise synthesized
   orchestration policy.
2. Update built-in registry, command-output, list, verbose-list, and grouped-help
   tests for the new alphabetically sorted prompt.
3. Update README command discovery and built-in prompt documentation.
4. Run formatting, focused prompt tests, the full Go test suite, vet, build,
   exact `kp merge --print` inspection, and the affected source-size audit.
5. Self-review, commit with the repository contract, push `GH-16`, and open one
   ready pull request that closes issue #16.

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

## DISCOVERIES

- The prompt registry loads every embedded `prompts/*.md` file dynamically and
  sorts by name; the implementation change is one asset plus tests and docs.
- `kit status` reported an unrelated managed refresh. Its dry-run would change
  12 managed instruction and ruleset files, so it is intentionally excluded
  from issue #16. The same dry-run reported 40 eligible handwritten source/test
  files and zero files above the 300-line limit before implementation.

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

## OUTCOME

- `kp merge` is an embedded, overrideable built-in prompt exposed through the
  existing print, copy, list, help, and launcher paths.
- Its six-step output builds a Mermaid PR DAG and wave table, binds exact-current
  readiness, maximizes safe independent concurrency and downstream unlocks,
  revalidates every wave, isolates failures, and separates merge from deployment
  and production evidence.
- Implementation and local validation are complete. Delivery remains one ready
  pull request from `GH-16` to `main`; merge is not authorized by this work.

## REPOSITORY MEMORY

- Created this living specification because the command establishes durable
  merge-orchestration behavior synthesized from multiple projects.
- Updated `README.md` with command discovery and built-in prompt documentation.
- Replaced placeholder testing guidance in `docs/references/testing.md` with
  the project's actual local commands, non-applicable environments, and hosted
  correctness-check gap.
- Updated `docs/PROJECT_PROGRESS_SUMMARY.md` with this feature's intent,
  approach, status, and pointer.
- The existing v0 feature artifacts remain unchanged because this is a separate
  public command and policy surface.
