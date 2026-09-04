---
label: Thread completion audit
---
Audit whether this thread is actually complete. Be skeptical. Verify current repository and GitHub state. Do not infer completion from implementation, merge, or a successful deploy workflow.

Internally check everything the thread implied, planned, requested, discovered, or required: forgotten planning items, side quests, TODOs, open issues and PRs, unmerged branches, failed or skipped validation, review or CI leftovers, undeployed or unverified production changes, and mismatches between promised delivery and current state.

Then output only:

## Summary
Two to four short sentences: what this thread did and whether it is complete. No chronology, no evidence dumps, no process narration.

## Remaining
If incomplete: a numbered list of remaining actions only. One line each. Include only work that still needs a decision or an action.

If complete:

> **THREAD COMPLETE:** nothing left from this thread.

Do not emit inventories, item classifications, required-answer sections, DONE items, awareness asides, or preambles.
