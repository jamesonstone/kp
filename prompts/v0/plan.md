---
label: Drive plan to implementation-ready
---
Drive the current implementation plan to a practical, evidence-backed local maximum. Stay in planning mode: perform read-only research, but do not edit files, create delivery artifacts, deploy, or execute the plan.

If a user-confirmed `/goal` is present, treat it as the intent contract. Consume it unchanged rather than rediscovering or redefining the objective. If implementation research exposes a material contradiction, report it and obtain an explicit revised user decision before changing the objective.

The target is the simplest decision-complete plan that another coding agent can implement without making material design decisions. Do not maximize plan length or complexity. Recommend a change only when it closes a concrete requirement gap, ambiguity, contradiction, unsupported assumption, failure mode, authority gap, or validation gap.

Research progressively:

- Start with the directly owning repository, current plan/specification, relevant code, tests, instructions, and deployed evidence.
- Inspect another repository or external system only when a discovered interface, dependency, authority, deployment, or validation question requires it.
- Resolve discoverable facts through evidence; do not ask the user to restate them.
- Treat historical plans and prior conclusions as hypotheses until reconciled against current evidence.
- Stop researching when additional retrieval would improve wording rather than resolve a material question.

If user input is required, ask the smallest consolidated batch of numbered questions. Include a recommended default, assumptions, uncertainty, and current confidence. Accept concise or y/n answers. Resume the convergence loop after the response.

Review and improve the plan across these lenses:

- goals, scope, and binary success criteria;
- authority, ownership, and external dependencies;
- interfaces, schemas, data flow, and state chronology;
- concurrency, replay, idempotency, and ambiguous starts;
- crash boundaries, partial commits, retry, recovery, and rollback;
- security, privacy, tenancy, secrets, and cost controls;
- compatibility, migrations, deployment, activation, and kill switches;
- observability, operator recovery, and acceptance evidence;
- requirement-to-implementation and requirement-to-test traceability;
- unnecessary complexity or speculative scope.

Silently iterate: research, identify material gaps, integrate the smallest complete corrections, and run another adversarial review. Do not narrate every pass or reproduce the plan after each pass.

Finalize only when all are true:

- zero unresolved material questions;
- at least 95% evidence-backed goal coverage;
- every requirement maps to an implementation change and validation;
- every state-changing or cost-bearing boundary has explicit authority, idempotency, retry, recovery, and safe-stop behavior;
- every external dependency has an exact owner, contract, and readiness gate;
- crash-boundary analysis proves retry cannot duplicate irreversible or cost-bearing effects;
- source, CI, merge, deployment, runtime, activation, operator acceptance, and business acceptance remain distinct;
- one complete adversarial pass produces no new material recommendation;
- a final verification pass also produces no new material recommendation.

If evidence or authority prevents convergence, stop with the exact blocker and the smallest action needed. Do not inflate confidence or hide the gap with assumptions.

Return one complete replacement plan containing:

- outcome and success criteria;
- interfaces, data flow, and state transitions;
- ordered implementation workstreams and dependencies;
- failure, concurrency, security, and recovery behavior;
- deployment, activation, rollback, and monitoring boundaries;
- test and acceptance plan;
- explicit assumptions, defaults, and deferred scope.

End with a concise local-maximum audit stating confidence, evidence-backed coverage, unresolved questions, the adversarial lenses checked, and why further iteration would be editorial or speculative rather than materially useful.
