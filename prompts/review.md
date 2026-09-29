---
label: Independent branch/PR review
---
Perform an independent review of the changes in this branch/PR against the target branch.

Treat the existing implementation as untrusted. Do not assume the approach is correct because another agent wrote it or because tests currently pass.

Focus primarily on correctness.

Review for:
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

Run the relevant tests, linters, type checks, and other repository validation available to you.

Do not make changes yet.

Return findings ordered by severity. For each finding include:
1. severity
2. file/location
3. concrete problem
4. why it is incorrect or risky
5. the smallest appropriate fix
6. test coverage that should prove the fix

Do not manufacture findings. If the implementation is correct, say so explicitly and identify any remaining validation gaps or residual risks.
