---
label: Dependency-ordered PR merge
---
Create and execute an evidence-backed merge plan only for the exact pull-request set directly authorized by the user or an accepted bounded plan.

1. Observe current GitHub, repository, workflow, dependency, deployment, runtime, and rollback evidence. Build a Mermaid DAG where A --> B means A must be merged—or, when the edge requires it, deployed and accepted—before B.
2. Bind each node to its repository, expected head OID, base, actor, allowed merge method, review policy, required current-head checks, dependency closure, and infrastructure effects. Classify it exactly as MERGE_READY, BLOCKED, or UNKNOWN; missing, stale, pending, skipped-without-proven-eligibility, cyclic, provisional, conflicted, or ambiguous evidence never passes.
3. Output the graph and a wave table. Repeatedly select only the zero-unmet-dependency MERGE_READY frontier. Maximize safe concurrency among independent nodes and prioritize nodes that unlock the most downstream work or shorten the critical path; serialize dependency chains and same-base, deployment-coupled, or otherwise interacting operations.
4. Immediately before every wave, revalidate authorization, identity, head/base, policy, checks, approvals, dependencies, deployment effects, and rollback ownership. Use the required merge queue and a repository-permitted method; never bypass safeguards or add an unapproved PR.
5. After every merge or queue transition, re-observe and recompute the graph. Stop the failed node and its dependents, but continue proven-independent authorized nodes.
6. Record merge, hosted workflow, deployment, runtime, production acceptance, recovery, and rollback as separate claims. Finish only when the authorized dependency closure reaches its required terminal state; otherwise report exact blockers, unknowns, completed waves, and the next safe action.
