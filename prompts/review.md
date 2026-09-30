---
label: Independent branch/PR review
---
Perform an independent review of the changes in this branch/PR against the target branch. The goal is to catch defects before merge, so report the problems the author would want to fix.

Treat the existing implementation as untrusted. Do not assume the approach is correct because another agent wrote it or because tests currently pass. Check it against its stated intent (PR description, issue, or spec), not only its code.

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

Do not make changes yet, to the checkout or to external state.

Confirm each finding before you report it, and list what you could not confirm as open questions. Report issues that predate this change separately, and skip style nits that no repository rule requires.

Return findings ordered by severity. For each finding include:
1. severity
2. file and line
3. concrete problem
4. why it is incorrect or risky
5. the smallest appropriate fix
6. test coverage that should prove the fix

Do not manufacture findings. If the implementation is correct, say so explicitly and identify any remaining validation gaps or residual risks.
