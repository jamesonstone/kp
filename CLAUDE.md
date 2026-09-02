# CLAUDE

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

## Multi-Agent Orchestration Evaluation Hard Gate

- Before finalizing any native implementation plan for a new feature, a substantial architectural or behavioral change, or a multi-file refactor, load `docs/references/rules/agent-team-orchestration.md` and evaluate whether the work benefits from multi-agent or parallel decomposition using that rule's lifecycle and semantic capability profiles.
- A single mechanical edit, a direct question, or read-only research that never forms an implementation plan does not trigger this gate.
- Record the decision before the plan is finalized: either a multi-lane Agent Team Plan, or `single-lane, because <reason>` using that rule's single-lane criteria. Never skip the evaluation silently, even when the recorded answer is single-lane.
- This gate fires during plan formation and precedes the Work Lane Mutation Hard Gate below, which fires later, before the first repository mutation.

## Work Lane Mutation Hard Gate

- Before any coding-agent repository file or delivery mutation, including issue, branch, staging, commit, push, worktree, and pull-request mutations, load `docs/agents/GUARDRAILS.md` and `work-lane-gating` first and complete read-only safety recon.
- Default to a new worklane without asking for the accepted unit of work: create or reuse one human-assigned GitHub issue, exact `GH-<issue-number>` branch, canonical non-primary worktree, and ready pull-request plan. Reuse that recorded lane for subsequent in-scope mutations. A clean or dirty checkout, current feature branch, issue reference, or generic pull-request request does not change this default.
- Continue an existing lane only when the user explicitly directs that outcome for the same unit of work. Prove the non-primary owning worktree, branch, issue scope, protected base, and create-or-update pull-request target.
- Never offer or ask the user to choose between lanes.
- Treat exact existing-PR lifecycle work as continuation: review repair, CI repair, base refresh, conflict resolution, and ordered merge coordination reuse every targeted pull-request head. Never create coordination or corrective pull requests for scope-preserving work. If source repair is not authorized, ask only for bounded in-place-remediation authority; do not allocate a new lane.
- Record a Pull-Request Landing Plan covering the repository, issue, branch, canonical non-primary worktree, protected base, and create-or-update PR target. Verify that plan still matches before every mutation. Ask only when implementation intent or an explicitly named target is materially ambiguous and cannot be resolved from repository evidence.
- Treat the primary/root checkout as read-only. If an ungated or root change exists, preserve it: Do not stage, commit, push, stash, reset, clean, discard, or silently transfer it.


## Testing And Validation Gate

- Before implementation or validation, including browser automation and browser testing, load `docs/references/rules/testing-and-environment-validation.md` and the project's `docs/references/testing.md`
- Preserve language-native code-level tests and pull-request checks; end-to-end and live-integration suites supplement rather than replace them

## GitHub Delivery Hard Gate

- In Kit-managed projects, issue, branch, staging, commit, push, and PR actions are mutation boundaries
- Before any GitHub delivery mutation, load `docs/agents/GUARDRAILS.md` and the relevant `docs/references/rules/*` delivery rules
- Repo-local Kit rules outrank global GitHub/plugin defaults; do not use generic branches, commits, PR bodies, or draft defaults when Kit defines the contract

## GitHub Merge Authorization Hard Gate

- Merge is a distinct mutation boundary. PR-delivery consent, automatic lane allocation, approval, check success, subagent assignment, and a program ledger never imply merge consent.
- Merge only after a direct user request or accepted bounded merge plan names the exact authorized PR set.
- Before any merge or merge-queue mutation, resolve `pull-request-merge` and load `docs/references/rules/github-pr-merge.md`.
- Reconcile the authorization source, authenticated actor, expected head/base, repository merge policy, current reviews/checks, dependencies, and infrastructure or deployment effects before every wave.
- Only exact current `MERGE_READY` nodes may merge. Pending, missing, stale-head, or policy-ineligible skipped checks are not passing.
- Revalidating an authorized target does not require another prompt. Adding a target or materially changing actor, method, environment, infrastructure effect, or recovery requires follow-up authorization.
- Never bypass protection, reviews, required checks, a merge queue, repository policy, or identity safeguards.
- Report merge, hosted workflow, deployment/runtime, and production evidence as separate claims.

## Cross-Repository Program Coordination Gate

- Before implementing or resuming an accepted plan that spans multiple repositories and includes dependent deliverables, staged deployment or activation, or expected agent or session handoff, load `docs/references/rules/cross-repository-program-coordination.md`.
- Designate one coordinator repository and create or adopt one canonical `docs/programs/<program>/PROGRAM.md` ledger before implementation; participant repositories remain authoritative for local specs, delivery state, runbooks, and evidence.
- Dispatch only the reconciled ready frontier, checkpoint every material transition and handoff, and reconcile recorded claims against live repositories, GitHub, runtime, and validation evidence before resume or completion.

## Deletion Safety Hard Gate

- Before designing deletion behavior or deleting persistent project, user, business, or external-system state, load `docs/references/rules/deletion-safety.md`.
- An unqualified delete means soft delete: use a reversible lifecycle state with a supported, authorized, and tested restore path. Task-owned ephemeral scratch that never became authoritative state is outside this retained-state definition; ambiguity remains covered.
- Treat purge, destroy, force deletion, empty-trash operations, destructive replacement, history rewrite, retention expiry, backup or snapshot deletion, cryptographic erasure, and irreversible cascades as hard delete.
- Make the normal product and operational path soft-delete by default. Keep hard delete as a separate privileged, auditable, server-enforced action; a client prompt or `force` flag alone is insufficient.
- Before any hard delete, resolve and present the exact targets, or a bounded selector first resolved to the exact current target set with its current count and materialized target IDs or an immutable snapshot/version token, environment, cascades, why soft delete is insufficient, the loss of restore, backup state, retention or legal impact, and verification plan.
- After that outline, obtain a specific manual confirmation from the human for those exact current targets. Initial requests, general task or plan approval, automation, retention schedules, prior soft-delete approval, and broad cleanup language do not count.
- Bind confirmation to the actor, action, exact targets or immutable snapshot/version, environment, and consequences. Immediately before execution, compare the current target set or version with the confirmed snapshot; any difference requires a new outline and confirmation.
- Preserve stricter repository, legal, privacy, security, infrastructure, and provider controls. One post-outline confirmation may satisfy multiple deletion gates only when the combined outline contains every required field.

## Infrastructure Change Approval Hard Gate

- Before mutating public-cloud resources, Kubernetes resources or cluster state, or infrastructure-as-code source, configuration, or state, load `docs/references/rules/infrastructure-change-approval.md`.
- Routine application operations on already-provisioned workloads, including deployment image updates and ECS or equivalent service interactions that do not create, replace, or delete infrastructure, are not infrastructure-approval batches. Record them; do not stop for a covered-mutation outline.
- Read-only discovery may precede confirmation only when it does not alter cloud resources, Kubernetes objects, remote state, or repository-owned infrastructure source.
- Put one consolidated outline of the target context, resource actions, execution boundary, material impact and risk, rollback or recovery, and validation evidence into the task plan when planning is used; otherwise present it once before the first covered mutation. Obtain one explicit user confirmation for the complete bounded batch.
- Approval of a task plan containing the complete outline counts as confirmation. A sufficiently detailed initial request may also count only when it clearly authorizes the exact bounded batch and the batch does not delete or remove infrastructure.
- Deleting, destroying, or removing infrastructure always requires explicit confirmation after the consolidated outline, even when the initial request asked for it; one confirmation covers every deletion named in that batch. Merge authorization, image deployment, and routine ECS interactions never authorize deletion.
- During merge or release orchestration, do not execute infrastructure deletion, destruction, purge, destructive replacement, or state removal; isolate it as a separate task with its own exact post-outline authorization.
- After confirmation, execute the exact approved batch and continue the rest of the task to completion in one pass without routine command-by-command approval.
- If additional covered infrastructure changes become necessary, collect all then-known changes into one follow-up outline, obtain one confirmation, and execute that follow-up batch in one pass. Do not re-confirm actions already included in an approved batch.
- Treat a material change to target identity, environment, region or cluster, resource set, action type, impact, or recovery as a follow-up batch; compatible tools, commands, and retries inside the approved boundary do not require another prompt.

## AWS Context Hard Gate

- If .kit.yaml defines an enabled aws context, run kit aws verify before the first AWS-dependent command in a task and again immediately before any AWS mutation
- Use the verified configured profile explicitly for every AWS-dependent command, including AWS CLI, SDK, Terraform, CDK, deployment, and project scripts, where supported
- After verification, never use default, another discovered profile, or ambient credentials
- Treat the verified account and ARN as authoritative; on missing credentials, incomplete config, or mismatch, stop and follow docs/agents/GUARDRAILS.md instead of falling back to another profile or default

## Agent Completion Output Contract

- Before a substantial terminal completion or handoff response, load `docs/references/rules/agent-completion-output.md` when present.
- This structured contract does not apply to intermediate progress commentary.
- Answer ordinary conversational requests naturally. Direct questions, definitions, confirmations, rewrites, brief explanations, small read-only lookups, concise recommendations, and focused clarification questions must not receive status tokens, canonical section headings, synthetic None items, task profiles, or repository-memory reporting.
- Use the structured contract when omitting it could hide a blocker, incomplete required scope, required operator action, unresolved failure, repository or external-system mutation, delivery artifact, multiple validation layers, material evidence, owner/dependency handoff, or when the user explicitly requests the canonical report.
- Do not classify by word count, token count, elapsed time, or tool-call count. When uncertain, prefer natural prose unless structure is necessary to preserve operationally important information.
- When the structured contract applies, emit exactly `## What happened`, `## Deviations`, and `## Next steps` in that order, with no other output section.
- Put `**Status: PASS|PARTIAL|BLOCKED|FAIL — <one-sentence outcome>.**` in the first What happened bullet. Do not add a separate status heading.
- Fold task-specific results, validation, delivery, coordination, and repository-memory evidence into concise What happened bullets. Use at most one nested evidence layer and state each fact once.
- Put blockers, incomplete scope, failures, warnings, pending or unknown evidence, skipped validation, and degraded execution under Deviations. Use one `**None.**` bullet when there are no deviations.
- Put independently actionable items under Next steps, required before optional. Name the actor and make every required continuation copy-ready. Use one `**None.**` bullet when no action remains.
- Use PASS only for complete scope and required validation, PARTIAL for usable incomplete work, BLOCKED for a specific external dependency, and FAIL for an unresolved known failure without an external stopping dependency.
- Preserve native evidence states such as PENDING, UNKNOWN, SKIPPED, and NOT_APPLICABLE literally.
- Do not use Markdown pipe tables, additional profile headings, or separate Completed, Validation, Delivery, Feature State, Residual Notes, Coordination, or Repository Memory sections.
- Preserve every field required by active delivery, validation, repository-memory, orchestration, program, and environment contracts inside the three canonical sections without duplication.
- For merge or release orchestration, report only state changes, terminal evidence, and actionable next steps; omit chronological command logs, repeated checks, unchanged polling, and routine tool details.

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

- Keep CLAUDE short and stable so it fits easily into injected context
- Put durable workflow guidance in `docs/agents/*` rather than expanding this file
- Do not add an always-loaded monolithic instruction file
