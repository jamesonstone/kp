---
label: Construct an executable goal
---

Construct an executable `/goal` from the current request and accumulated conversation before producing an implementation plan.

Stay in planning mode. Perform read-only research when useful, but do not edit files, create delivery artifacts, deploy, perform state-changing setup, or execute the goal.

Treat the initial request as a thesis to investigate, not automatically as a complete specification. Preserve its intent while reconciling it against current repository, runtime, external-system, and user evidence when available.

Maintain one accumulated goal model. At minimum, keep track of the intended outcome, success criteria, scope boundaries, material assumptions, unresolved material questions, and relevant constraints or dependencies. Add dimensions such as authority, state sequence, failure behavior, recovery, evidence, or acceptance boundaries only when they are material to the objective.

Research progressively. Resolve discoverable facts from repository instructions, current code, tests, specifications, configuration, deployment/runtime evidence, documentation, and owning-system contracts before asking the user. Treat historical conclusions as hypotheses until reconciled against current evidence. Stop retrieving when additional context would improve wording rather than resolve a material question.

Ask only questions whose answers could materially change the intended outcome, implementation, scope, authority, architecture, state sequence, risk, dependency ownership, success criteria, required evidence, or recovery behavior. Ask the smallest useful numbered batch. Recommend a default when one is justified, state material assumptions or uncertainty, and accept concise answers such as `1y 2n`.

After every material discovery or answer, synthesize it into the accumulated goal model. Accepted defaults become decisions. Later answers supersede earlier assumptions. Surface contradictions explicitly and reopen only the affected dimensions. Do not treat responses as unrelated notes or repeatedly restate the full conversation.

Challenge the emerging goal before convergence. Look for unsupported assumptions, contradictory requirements, ambiguous success criteria, circular dependencies, false PASS conditions, unclear authority, and unsafe or undefined failure behavior when applicable.

If execution can safely create or resolve a prerequisite through a supported interface, treat that prerequisite as bootstrap work rather than a blocking precondition.

Distinguish implementation, deployment, configuration, runtime readiness, operator or integration acceptance, and broader business or physical acceptance when those distinctions materially affect what success means. State what PASS proves and what it must not claim.

For consequential workflows, define enough failure behavior to make the goal executable: the last safe checkpoint, retry or reconciliation requirements, repair behavior, and resume semantics when applicable. Do not force these dimensions onto trivial tasks where they cannot change implementation or acceptance.

Urgency may change prioritization and the order in which blockers are investigated. It must not weaken truth, safety, authorization, or acceptance criteria.

Converge when:

- no known unresolved material question remains;
- no known contradiction remains;
- PASS and its relevant claim boundaries are unambiguous;
- scope is sufficiently defined for implementation planning;
- applicable dependencies, authority boundaries, and failure concerns are resolved; and
- one adversarial review of the accumulated goal exposes no new material ambiguity.

Residual uncertainty may remain documented when it cannot materially change implementation or acceptance.

Then return one cohesive `/goal` contract, not an implementation plan or transcript. Adapt its structure to the task rather than forcing irrelevant sections. It should communicate, at minimum:

1. the intended outcome;
2. binary or otherwise objectively verifiable PASS criteria;
3. relevant scope and non-goals;
4. material constraints, dependencies, assumptions, and claim boundaries;
5. any additional state, authority, failure, recovery, evidence, or acceptance semantics required to make the goal executable; and
6. any intentionally deferred non-material uncertainty.

Ask the user to confirm the complete `/goal` before handing it to implementation planning.

If convergence is blocked, report the exact missing decision or evidence and the smallest action needed. Do not manufacture a complete goal.
