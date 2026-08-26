---
label: Context-aware PR merge and deployment
---
Coordinate the already-scoped pull requests using the conversation and current repository context. Derive the exact PR set, heads, bases, dependencies, permitted merge methods, deployment targets, acceptance gates, and authorization boundaries; do not ask the user to restate discoverable facts or expand scope.

## Analysis and approval

1. Before any repository, delivery, or infrastructure mutation, load `docs/agents/GUARDRAILS.md` plus the applicable `work-lane-gating` and infrastructure rules, then complete read-only safety reconnaissance. For repository or delivery work, verify the lane rule's exact consent question was explicitly answered for the current lane; if not, ask it verbatim and wait. Generic approval is not lane consent.
2. Load the remaining applicable repository-local merge, orchestration, and completion rules. Build the dependency and deployment graph from authoritative evidence. Classify each node as `MERGE_READY`, `BLOCKED`, or `UNKNOWN`; missing, stale, pending, conflicted, or unattributable evidence never passes.
3. Record each node's repository, PR, exact head/base, method, dependencies, infrastructure effects, recovery, and acceptance signal. Keep merge, CI, deployment, runtime, and production acceptance distinct.
4. Before the first merge or infrastructure mutation, present one consolidated approval request for the exact current frontier and all known infrastructure effects unless equivalent exact approval from the conversation remains valid. Never delete, destroy, purge, destructively replace, or remove infrastructure; stop and isolate that work for separate explicit authorization.
5. Run one complete preflight immediately before each consequential mutation. Repeat only when a material fact changes: head/base, reviews/checks, policy, actor, dependency, deployment target/effect, approval, or acceptance window.

## Execution

- One primary coordinator owns the graph, authority, wave selection, recovery, and final acceptance.
- When the host supports explicit model selection, use lower-cost or lower-capability agents only for exact, bounded `MERGE_READY` merges and deployment monitoring. Keep graph changes, repair decisions, recovery, and acceptance with the coordinator.
- Parallelize nodes only when both source and deployment effects are independent. Serialize shared bases, services, environments, databases, migrations, queues, and acceptance gates.
- Merge only the authorized ready frontier with repository-permitted methods and required queues. Never bypass policy, switch identity, force-push, weaken a gate, or explicitly delete PR branches.
- Monitor with event-driven waits or bounded backoff; do not emit or repeat unchanged polling.
- Reconcile failures autonomously within the approved scope: diagnose once, apply only authorized in-lane repair, rerun affected evidence, and refresh that node and its dependents. Do not retry blindly or introduce a new PR, infrastructure effect, target, method, or authority boundary.
- A changed head returns to `UNKNOWN` and requires fresh current-head evidence and authorization. Continue independent valid nodes; stop when recovery cannot make progress safely.

## Final response

Use repository-required status vocabulary; otherwise emit:

`Status: SUCCESS | PARTIAL | BLOCKED | FAILURE`

`Result:` one to three sentences stating the exact outcome and evidence boundary.

`Next steps:` `none` or at most three specific, copy-ready sentences.

Do not include a chronological work log, repeated checks, unchanged polling, or routine command details.
