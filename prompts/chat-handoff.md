---
label: Chat-to-agent handoff
---
Create a zero-context handoff from the current chat or brainstorm into a filesystem-capable coding agent. Preserve every implementation-relevant decision, source, constraint, acceptance condition, and unresolved fact. Remove repetition, not necessary detail.

## Origin phase 1: clarify

- Read the full conversation and every supplied note, document, attachment, and link.
- Inventory the objective, users, scope, non-goals, decisions, constraints, sources, acceptance conditions, risks, dependencies, and requested authority.
- Resolve apparent ambiguity from the conversation before asking anything.
- If an unresolved answer would materially change implementation, output only the smallest numbered set of questions and stop. Do not produce the handoff in the same response.
- Do not ask about cosmetic preferences or repository facts that the destination agent can determine safely by inspection.

## Origin phase 2: emit

Begin only after phase 1 has zero unresolved implementation-changing questions. Return only one fenced code block with language tag `markdown`; place no text outside it. The document must be titled `# Coding Agent Handoff` and use these H2 sections in order:

1. `## Objective and Scope`
2. `## Decisions and Constraints`
3. `## Source Map`
4. `## Required Work`
5. `## Acceptance and Validation`
6. `## Risks, Conflicts, and Unknowns`
7. `## Resource Links`
8. `## Destination Protocol`

### Evidence rules

- Assume the destination has no conversation history.
- Tag every non-trivial claim in Objective and Scope, Decisions and Constraints, and Required Work with one or more source IDs such as `[S1]`.
- In Source Map, define each source ID with source type (`discussion`, `note`, `document`, `attachment`, or `link`), exact identifier or URL, and the claim it supports.
- Preserve explicit user decisions as constraints. Do not reopen them unless live repository evidence conflicts.
- Write `UNKNOWN` for missing evidence and give the destination one exact read-only inspection action.
- Write `CONFLICT` for contradictory sources, identify both sources, and apply any tie-break decision already made in the conversation. Otherwise require destination clarification.
- Never invent repository paths, symbols, commands, outputs, owners, dates, issue state, or links.
- Never include passwords, tokens, cookies, secret values, signed URLs, private keys, PHI, or customer data. Include safe credential names and retrieval boundaries only when required.

### Section requirements

- Objective and Scope: state the outcome, affected users, measurable definition of done, selected direction, in-scope work, non-goals, and deferred work.
- Decisions and Constraints: list settled product, technical, compatibility, safety, authorization, and delivery decisions without duplicating background.
- Required Work: provide an ordered, implementation-ready plan. Name known components and interfaces; use exact destination inspection actions for unknown repository facts. Include data and state transitions, failure behavior, observability, security, compatibility, migration, rollback, operator visibility, documentation, and tests when the task establishes those concerns.
- Acceptance and Validation: use binary-verifiable bullets. Tie each validation command or inspection to the criterion it proves and state the expected result; use `UNKNOWN` when the chat cannot establish a repository command.
- Risks, Conflicts, and Unknowns: list only material remaining items with mitigation, owner, and resolution action.
- Resource Links: include every external URL from the conversation with title and one-line relevance. Output `- NONE` when there are no external links.

### Destination protocol content

The handoff's Destination Protocol must instruct the receiving coding agent to use two phases:

1. Clarify: load repository-local instructions, inspect the live repository and delivery state read-only, reconcile the handoff against actual files and behavior, and ask only implementation-changing questions that neither the handoff nor safe inspection can answer. Output only those questions and stop.
2. Hydrate: after all answers, output `## Context Hydration` summarizing the confirmed objective, decisions, live state, remaining work, authority boundaries, validation contract, divergences, and exact next safe action. Then ask `Proceed with the hydrated task?` and do not implement or mutate anything until the user explicitly agrees.
