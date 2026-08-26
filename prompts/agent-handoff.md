---
label: Agent-to-agent handoff
---
Create a lossless, zero-context handoff from the current coding-agent session to a different coding-agent system. The destination may use different models, tools, shells, repository integrations, and instruction formats. Preserve facts and authority, not provider-specific assumptions. Completeness outranks brevity; remove redundancy before removing detail.

## Origin phase 1: clarify

- Read the full conversation, current task plan, repository instructions, specifications, ledgers, attachments, links, and tool results.
- Inspect the live workspace read-only. Resolve the current repository, working directory, primary checkout, worktree, branch, base, `HEAD`, upstream, staged/unstaged/untracked paths, relevant diff, commits, worktrees, issue and pull-request state, reviews/checks, deployment/runtime state, validation evidence, and exact next safe action.
- Separate observed facts, user decisions, agent inference, and stale or unverified claims.
- Resolve questions from current conversation, files, Git, repository host, and other authorized read-only sources before asking the user.
- If an unresolved answer would materially change scope, architecture, safety, authority, or the next action, output only the smallest numbered set of questions and stop. Do not produce the handoff in the same response.
- Do not mutate files, Git state, issues, pull requests, deployments, or external systems solely to prepare the handoff.

## Origin phase 2: emit

Begin only after phase 1 has zero unresolved implementation-changing questions. Refresh volatile Git, review, check, deployment, and runtime evidence once, then return only one fenced code block with language tag `markdown`; place no text outside it. The document must be titled `# Coding Agent Handoff` and use these H2 sections in order:

1. `## Executive Context`
2. `## Authority and Safety Boundaries`
3. `## Repository and Workspace State`
4. `## Delivery and External State`
5. `## Decisions and Constraints`
6. `## Completed Work`
7. `## Remaining Work and Next Safe Action`
8. `## Source Map`
9. `## Acceptance and Validation`
10. `## Risks, Blockers, and Unknowns`
11. `## Resource Links`
12. `## Destination Protocol`

### Evidence and safety rules

- Assume the destination has no prior messages, memory, filesystem state, tool output, or implicit authority.
- Begin with the UTC snapshot time and identify the originating agent system/model only when known.
- Tag every non-trivial claim in Executive Context, Decisions and Constraints, Completed Work, and Remaining Work with one or more source IDs such as `[S1]`.
- In Source Map, define each source ID with source type (`conversation`, `file`, `command`, `test`, `issue`, `pull request`, `workflow`, `runtime`, or `external document`), exact path/URL/command or stable identifier, observation time when volatile, and supported claim.
- Preserve native evidence states such as `PENDING`, `UNKNOWN`, `SKIPPED`, `BLOCKED`, `NOT_RUN`, and `NOT_APPLICABLE`. Never convert them into success.
- Write `CONFLICT` when sources disagree; state the competing evidence and the repository or user tie-break rule. Do not guess.
- Never include passwords, tokens, cookies, authorization headers, secret values, signed URLs, private keys, PHI, customer data, or credential-bearing file contents. Include safe variable/secret names and retrieval boundaries only when required.

### State requirements

- Authority and Safety Boundaries: distinguish read/plan, implementation, delivery, review-thread, merge, deployment, infrastructure, production, and destructive-action authority. Quote or cite the exact approval source. State every prohibited or separately gated action.
- Repository and Workspace State: record repository and remotes, primary checkout, active worktree, branch/base/HEAD/upstream, divergence, dirty paths by category, relevant diff intent, commits introduced, unpushed work, user-owned or concurrent changes, and environment/process dependencies. State `none` explicitly when verified absent.
- Delivery and External State: record issue/branch/PR URLs and exact heads, review/check status, merge state, artifacts, deployment/runtime identity, migrations, feature flags, external coordination, and durable ledger/checkpoint state. Keep source, CI, merge, deployment, runtime, and business acceptance separate.
- Decisions and Constraints: preserve settled requirements, architecture, interfaces, data ownership, compatibility, security, rollout, rollback, testing, documentation, and rejected alternatives that constrain future work.
- Completed Work: list exact files/symbols or external resources changed, behavior delivered, commits, validations, and evidence. Mark work the destination must not repeat.
- Remaining Work and Next Safe Action: order the unfinished work by dependency. Name exact files/symbols when known, expected edits, validation, blockers, owner, and the single next safe action.
- Acceptance and Validation: use binary-verifiable criteria. Map each criterion to exact commands, suites, hosted checks, environment, expected literal result, and current evidence state tied to a commit or artifact.
- Risks, Blockers, and Unknowns: include only material items with impact, mitigation, owner, and exact resolution action. For missing facts, write `UNKNOWN` and give one bounded read-only inspection.
- Resource Links: include all relevant repository, issue, pull-request, workflow, evidence, specification, design, and official external documentation URLs with one-line relevance. Output `- NONE` only when no URL exists.

### Destination protocol content

The handoff's Destination Protocol must instruct the receiving coding agent to use two phases:

1. Clarify: treat the handoff as an attributable snapshot, load destination and repository-local instructions, inspect live state read-only, verify paths/symbols/heads/checks/artifacts, detect drift, and ask only implementation-changing questions not answerable from the handoff or inspection. Output only those questions and stop. Do not repeat completed work.
2. Hydrate: after all answers, output `## Context Hydration` with the confirmed goal, source snapshot, live-state differences, preserved decisions, completed work, remaining dependency graph, authority boundaries, validation contract, blockers, and exact next safe action. Then ask `Proceed with the hydrated task?` and do not implement or mutate anything until the user explicitly agrees.
