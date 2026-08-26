---
kit_metadata_version: 1
artifact: "spec"
workflow_version: 3
phase: "deliver"
feature:
  id: "0005"
  slug: "handoff-prompts"
  dir: "0005-handoff-prompts"
references:
  - id: "github-issue"
    name: "Add explicit cross-agent handoff prompts"
    type: "external"
    target: "https://github.com/jamesonstone/kp/issues/32"
    relation: "supports"
    read_policy: "must"
    used_for: "issue, branch, commit, and pull-request traceability"
    status: "active"
  - id: "v0-init-utility"
    name: "v0 init utility"
    type: "feature_artifact"
    target: "docs/specs/0001-v0-init-utility/SPEC.md"
    relation: "constrains"
    read_policy: "must"
    used_for: "embedded prompt execution, clipboard verification, discovery, and overrides"
    status: "active"
  - id: "openai-prompting-guidance"
    name: "OpenAI model prompting guidance"
    type: "external"
    target: "https://developers.openai.com/api/docs/guides/latest-model"
    relation: "informs"
    read_policy: "must"
    used_for: "outcome-first prompting, evidence preservation, autonomy boundaries, and stopping conditions"
    status: "active"
delivery_intent: "issue_branch_pr_ready"
---
# SPEC

## PURPOSE

Replace ambiguous `kp handoff` with two provider-neutral prompt commands:
`kp chat-handoff` for chat or brainstorm context entering a coding agent, and
`kp agent-handoff` for transferring a live filesystem-aware coding task between
agent systems without losing evidence, authority, or execution state.

## CONTEXT

- `kp` already exposes embedded prompt assets as bare commands with print,
  copy, help, listing, launcher, and user-override behavior. No command-specific
  runtime is required.
- The existing handoff prompt assumes a generic chat source and Codex
  destination. It cannot represent live Git/worktree/delivery/runtime state or
  provider-neutral agent-to-agent transfer.
- The user selected two explicit command names, prompt-only behavior, removal
  of legacy `handoff`, clarification at both origin and destination, context
  hydration after destination clarification, and an explicit proceed request
  before destination implementation.
- OpenAI guidance favors outcome-first prompts that state each rule once while
  preserving goal, hard constraints, evidence, success criteria, authority,
  and stopping conditions.

## REQUIREMENTS

- Add embedded prompts named `chat-handoff` and `agent-handoff` with labels
  `Chat-to-agent handoff` and `Agent-to-agent handoff`.
- Remove the built-in `handoff` prompt without adding a compatibility alias.
- Preserve the existing prompt registry and command behavior; do not add target
  flags, provider integrations, network calls, or a new runtime abstraction.
- Both prompts must use an origin clarification phase that resolves available
  evidence first, asks only implementation-changing questions, outputs only
  those questions when needed, and emits no handoff until none remain.
- Both generated handoffs must embed a destination clarification phase that
  reconciles the snapshot with repository-local instructions and live state
  before mutation.
- After destination questions are answered, require a `Context Hydration`
  summary and the exact question `Proceed with the hydrated task?`; prohibit
  implementation or mutation until the user explicitly agrees.
- The chat handoff must preserve objective, scope, decisions, sources, required
  work, binary acceptance criteria, validation, risks, conflicts, unknowns,
  and external links. Missing repository facts use `UNKNOWN` plus an exact
  read-only destination inspection.
- The agent handoff must additionally preserve UTC snapshot time, repository,
  primary checkout/worktree, branch/base/HEAD/upstream, dirty paths, diff and
  commit intent, issue/PR/check/review state, deployment/runtime state,
  completed work not to repeat, remaining dependency order, exact next safe
  action, authority boundaries, and validation tied to source or artifact
  identity.
- Both prompts must distinguish facts, user decisions, inference, conflict,
  stale evidence, and native states such as `PENDING`, `UNKNOWN`, `BLOCKED`,
  `NOT_RUN`, and `NOT_APPLICABLE`.
- Both prompts must prohibit credentials, secrets, signed URLs, PHI, customer
  data, and other sensitive values while permitting safe variable/secret names
  and retrieval boundaries.
- Exact command output, required contract phrases, prompt ordering, labels,
  help, list, override, collision, edit, and removal behavior must be tested.
- README and project progress documentation must describe both commands and the
  removal of legacy `handoff`.
- Every changed handwritten source/test file must remain at or below 300 lines.

## ACCEPTED PLAN

1. Replace `prompts/handoff.md` with `prompts/chat-handoff.md` and
   `prompts/agent-handoff.md`.
2. Keep command execution dynamic through the existing embedded registry.
3. Pin each printed prompt by SHA-256 and required behavioral phrases in a
   focused test file; update discovery, help, registry, and library tests.
4. Update README, progress summary, and this living specification.
5. Run focused and complete tests, race detection, vet, both builds, isolated
   CLI acceptance, source-size audit, diff hygiene, and secret scan.
6. Deliver issue #32 through `GH-32` and one ready pull request to `main`.

## DECISIONS

- Accepted: `chat-handoff` and `agent-handoff` are shorter and clearer than the
  brainstorm/coding alternatives while preserving the `<source>-handoff`
  naming pattern.
- Accepted: legacy `handoff` is removed. A duplicate alias would retain the
  ambiguity and duplicate prompt contracts.
- Accepted: destination implementation requires a new explicit proceed answer
  after live reconciliation and context hydration; transferred context alone
  does not invent mutation authority.
- Accepted: provider-specific model names are evidence only. The handoff
  describes capabilities, state, sources, and authority so another agent system
  can map them to its own tools.
- Accepted: prompt completeness has no arbitrary word limit. Redundancy is
  removed before implementation-relevant detail.

## DISCOVERIES

- Built-in prompts are discovered from `prompts/*.md`, sorted by name, and
  automatically exposed through print, copy, list, help, launcher, and override
  paths.
- The old handoff test checked phrases but did not pin exact output. SHA-256
  tests preserve the complete prompt without placing two long bodies in Go
  source or violating the 300-line rule.
- The agent-to-agent transfer needs a volatile evidence refresh immediately
  before emission; otherwise the destination can inherit stale heads, checks,
  deployment identity, or runtime state.

## VALIDATION

- `go test ./internal/prompt ./internal/cmd` — `PASS`.
- `test -z "$(gofmt -l prompts.go cmd internal)"` — `PASS`.
- `go test ./...` and `go test -race ./...` — `PASS` across every package.
- `go vet ./...`, `go build ./...`, and `make build` — `PASS`.
- Isolated empty-config CLI acceptance — `PASS`: `agent-handoff --print`
  emitted SHA-256
  `74a5c9e682397a603f09b71b503fc96451da5f0b56f195058a5811331e880aed`;
  `chat-handoff --print` emitted
  `9e350823259e3856c63a987a92428772f69872c11a1c640b91d53e91b8cfd063`;
  list/help showed both commands; legacy `handoff` failed as absent.
- `kit check --all` — `PASS` for all four feature contracts.
- `kit reconcile --all --dry-run` — `PASS` source-size audit: 42 eligible
  handwritten source/test files checked, zero above 300 physical lines. Ten
  pre-existing managed-refresh warnings in seven untouched files remain
  outside issue #32.
- `git diff --check` — `PASS`.
- `gitleaks git --redact --no-banner` — `PASS`: 51 commits and 1.42 MB scanned
  with no leaks.
- Hosted correctness checks — `UNAVAILABLE`: this repository has no hosted
  format, test, race, vet, or build workflow.
- Production validation — `NOT_APPLICABLE`: `kp` is a local prompt CLI.

## OUTCOME

- Source and local validation are complete on issue #32 and `GH-32`.
- `kp chat-handoff` preserves clarified brainstorm context for a zero-context
  coding agent; `kp agent-handoff` additionally preserves live repository,
  delivery, runtime, authority, completed-work, validation, and next-action
  state.
- Both generated handoffs require destination reconciliation, destination
  clarification, context hydration, and explicit permission before mutation.
- Ready-PR delivery is authorized; merge is not authorized.

## REPOSITORY MEMORY

- Created this specification because the two-phase transfer protocol,
  zero-context state contract, provider-neutral boundary, legacy-name removal,
  and destination authority gate are durable behavior not recoverable from the
  prompt filenames alone.
