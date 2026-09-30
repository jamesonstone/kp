---
label: Independent branch/PR review
---
Perform an independent review of the changes in this branch/PR against the target branch. The goal is to catch defects before merge, so report the problems the author would want to fix.

Establish the scope first: find the target branch (from the PR, or the repository's default branch), review the merge-base diff against its remote copy (for example `origin/main`) plus any uncommitted changes, and read the change's stated intent (PR description, linked issue, or spec).

Treat the existing implementation as untrusted. Do not assume the approach is correct because another agent wrote it or because tests currently pass.

Focus primarily on correctness.

Review for:
- behavior that differs from the change's stated intent
- incorrect behavior or logic
- regressions
- edge cases and boundary conditions
- error handling and failure modes
- concurrency, ordering, idempotency, and retry issues where applicable
- state consistency and transactional correctness
- security or data-integrity problems
- violations of existing contracts, invariants, or repository conventions
- missing or insufficient tests
- unnecessary complexity that creates correctness risk

Trace important execution paths rather than reviewing only the diff syntactically. Inspect surrounding code and existing tests/contracts when necessary to determine intended behavior.

Run the relevant tests, linters, type checks, and other repository validation available to you, without modifying tracked files or external state: no auto-fixers, snapshot updates, migrations, or deployments. If a check cannot run, say which one and why.

Do not make changes yet.

Search broadly, then confirm each finding before reporting it: trace the failing path in the code, or demonstrate it with a command or a scratch test in a temporary copy of the repository. List concerns you cannot confirm as open questions rather than dropping them. Report problems this change introduces or makes reachable, and list pre-existing issues you notice separately as follow-ups. Skip style or preference nits unless a repository rule requires them or they create correctness risk.

Return findings ordered by severity: P0 for data loss, an exploitable security flaw, or an outage; P1 for incorrect behavior in normal use; P2 for edge-case or robustness risk; P3 for low-impact correctness issues. For each finding include:
1. severity
2. file and line
3. concrete problem
4. why it is incorrect or risky
5. the smallest appropriate fix
6. test coverage that should prove the fix

Do not manufacture findings. If the implementation is correct, say so explicitly and identify any remaining validation gaps or residual risks.
