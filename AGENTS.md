# AGENTS

## Purpose

- This file is a routing table, not the full manual
- Start at `docs/agents/README.md`, then load only the docs needed for the current decision
- Repo-local markdown under `docs/` is the system of record

## Pasted Text Attachments

- If the user message includes an attached pasted-text file and the visible message is empty or minimal, treat the attachment as the active task instructions unless the user says otherwise
- If the attachment appears Kit-generated, follow it directly without asking what the attachment is for

## Runtime Routing

- `docs/agents/README.md` — classify the task and choose the next document
- `docs/agents/WORKFLOWS.md` — spec-driven versus ad hoc flow
- `docs/agents/GUARDRAILS.md` — completion, safety, and hard rules
- `docs/agents/RLM.md` — just-in-time context loading when broad context would be noisy
- `docs/agents/TOOLING.md` — skills, dispatch, project-directory workflow, and secondary inputs

## Testing And Validation Gate

- Before implementation or validation, including browser automation and browser testing, load `docs/references/rules/testing-and-environment-validation.md` and the project's `docs/references/testing.md`
- Preserve language-native code-level tests and pull-request checks; end-to-end and live-integration suites supplement rather than replace them

## GitHub Delivery Hard Gate

- In Kit-managed projects, issue, branch, staging, commit, push, and PR actions are mutation boundaries
- Before any GitHub delivery mutation, load `docs/agents/GUARDRAILS.md` and the relevant `docs/references/rules/*` delivery rules
- Repo-local Kit rules outrank global GitHub/plugin defaults; do not use generic branches, commits, PR bodies, or draft defaults when Kit defines the contract

## Infrastructure Change Approval Hard Gate

- Before mutating public-cloud resources, Kubernetes resources or cluster state, or infrastructure-as-code source, configuration, or state, load `docs/references/rules/infrastructure-change-approval.md`.
- Read-only discovery may precede confirmation only when it does not alter cloud resources, Kubernetes objects, remote state, or repository-owned infrastructure source.
- Put one consolidated outline of the target context, resource actions, execution boundary, material impact and risk, rollback or recovery, and validation evidence into the task plan when planning is used; otherwise present it once before the first covered mutation. Obtain one explicit user confirmation for the complete bounded batch.
- Approval of a task plan containing the complete outline counts as confirmation. A sufficiently detailed initial request may also count only when it clearly authorizes the exact bounded batch and the batch does not delete or remove infrastructure.
- Deleting, destroying, or removing infrastructure always requires explicit confirmation after the consolidated outline, even when the initial request asked for it; one confirmation covers every deletion named in that batch.
- After confirmation, execute the exact approved batch and continue the rest of the task to completion in one pass without routine command-by-command approval.
- If additional covered infrastructure changes become necessary, collect all then-known changes into one follow-up outline, obtain one confirmation, and execute that follow-up batch in one pass. Do not re-confirm actions already included in an approved batch.
- Treat a material change to target identity, environment, region or cluster, resource set, action type, impact, or recovery as a follow-up batch; compatible tools, commands, and retries inside the approved boundary do not require another prompt.

## AWS Context Hard Gate

- If .kit.yaml defines an enabled aws context, run kit aws verify before the first AWS-dependent command in a task and again immediately before any AWS mutation
- Use the verified configured profile explicitly for every AWS-dependent command, including AWS CLI, SDK, Terraform, CDK, deployment, and project scripts, where supported
- After verification, never use default, another discovered profile, or ambient credentials
- Treat the verified account and ARN as authoritative; on missing credentials, incomplete config, or mismatch, stop and follow docs/agents/GUARDRAILS.md instead of falling back to another profile or default

## Conditional Context

- `docs/specs/<feature>/` — active feature artifacts only
- `docs/references/README.md` — durable repo references only when relevant
- `docs/CONSTITUTION.md` — project invariants when a decision depends on them

## Repo Knowledge Map

- `docs/agents/README.md` — runtime routing index
- `docs/agents/WORKFLOWS.md` — work classification and source-of-truth semantics
- `docs/agents/RLM.md` — progressive disclosure and context budget rules
- `docs/agents/TOOLING.md` — skills, dispatch, project-directory workflow, and secondary global inputs
- `docs/agents/GUARDRAILS.md` — completion bar, safety rules, and validation expectations
- `docs/references/README.md` — durable repo-local references that are broader than one feature
- `docs/specs/<feature>/SPEC.md` — v2 feature source of truth for requirements, plan, tasks, validation, reflection, delivery, and evidence

## Constraints

- Keep AGENTS short and stable so it fits easily into injected context
- Put durable workflow guidance in `docs/agents/*` rather than expanding this file
- Do not add an always-loaded monolithic instruction file
