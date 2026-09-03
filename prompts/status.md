---
label: Thread completion audit
---
# Purpose
Perform a comprehensive completion audit of all work discussed in this thread. Determine whether everything implied, planned, requested, discovered, or required by this thread has actually been completed.

# Runtime contract
`kp status` only prints this procedure. It does not inspect repositories, GitHub, CI, deployments, or conversation history by itself. Discover the thread context, project systems, worklanes, issues, pull requests, branches, tests, and validation evidence from the current environment. Do not assume GitHub, Git, a CI system, or a deployment platform unless the environment establishes them.

Be skeptical. Search for omissions rather than trying to justify completion. Do not infer completion merely because implementation work occurred. Verify the actual current repository and GitHub state wherever possible.

# Audit scope
Review the entire conversation, including early planning, later implementation, tangential discussions, side quests, follow-up discoveries, and requirements that emerged during execution.

Check for all of the following:

- Explicit tasks, requirements, acceptance criteria, and deliverables.
- Tasks identified during planning or analysis that may not have been revisited.
- Tangential or side-quest work that was discussed, suggested, discovered, or implicitly accepted as necessary.
- Work required for correctness, completeness, production readiness, integration, cleanup, migration, documentation, testing, or validation.
- TODOs, deferred items, placeholders, temporary implementations, follow-ups, or unresolved decisions.
- Open GitHub issues related to this work.
- Open pull requests, including draft PRs, stacked or dependent PRs, or PRs awaiting review, fixes, merge, or deployment.
- Branches or worklanes containing relevant unmerged work.
- Failed, skipped, incomplete, or missing tests and validation.
- Changes that were implemented but not merged, deployed, documented, or verified.
- Review feedback or CI failures that remain unresolved.
- Dependencies discovered during the work that were never completed.
- Any mismatch between what we said we would deliver and the repository's current state.

# Classification
Classify every identified item as:

- **DONE** — fully implemented, merged, and validated as required.
- **OPEN** — known work remains.
- **PARTIAL** — started but incomplete.
- **BLOCKED** — cannot currently be completed; explain why.
- **NOT NEEDED** — discussed but ultimately determined unnecessary; explain the decision.

# Required answers
Then answer:

1. **Is all work from this thread complete?** Yes or no.
2. If no, what exactly remains?
3. Are there any open or unmerged PRs related to the thread?
4. Are there any unresolved GitHub issues or tasks?
5. Are there any requirements or side quests from earlier in the conversation that were forgotten?
6. Is there anything else required before this thread can legitimately be considered complete?

# Completion statement
If everything is complete, explicitly state:

> **THREAD COMPLETE:** I found no remaining implementation, integration, testing, review, merge, deployment, documentation, issue, PR, or follow-up work attributable to this thread.

Otherwise, provide the remaining work as a prioritized completion checklist.

# Initial execution
When `kp status` begins, do not implement immediately. Inventory the thread's implied work, verify current repository and GitHub state where possible, classify every item, answer the required questions, and end with either the thread-complete statement or a prioritized completion checklist.

Use this initial response structure:

## Thread work inventory
Concise summary of everything this thread implied, planned, requested, discovered, or required.

## Item classification
For each identified item:

### Item N: `<short name>`
- **Status:** DONE | OPEN | PARTIAL | BLOCKED | NOT NEEDED
- **Evidence:**
- **Notes:**

## Required answers
Answer all six required questions directly.

## Completion outcome
Either the exact thread-complete statement or a prioritized completion checklist.

# Governing principle
A thread is not complete because code was written or a plan was discussed. Completion requires reconciling the full conversation against current evidence for implementation, integration, validation, review, merge, deployment, documentation, issues, pull requests, and follow-up work.
