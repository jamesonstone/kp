---
label: Punch list control loop
---
# Purpose
Drive a living punch list toward resolution by scanning, clustering, and fixing shared causes. Treat the list as operational evidence about the product, not a queue of independent tickets or an authoritative spec.

# Runtime contract
`kp punchlist` only prints this procedure. It does not locate, fetch, or mutate a punch list. Discover the actual document, tools, worklanes, and validation environment from the current project. Do not assume GitHub, Git, a spreadsheet provider, a CI system, or a deployment platform unless the environment establishes them. If multiple punch-list documents exist and the correct one cannot be determined, ask.

If the invoking user message, project config, or punch-list document supplies an engineering-note author prefix, use it consistently. This command does not bind an author flag.

# Control loop
Scan → Reconcile → Cluster → Research → Plan → Implement → Review → Deploy → Update → Re-test → Repeat.

Optimize for the smallest coherent product or engineering change that resolves the underlying problem, including multiple items when they share a root cause.

Re-read the entire punch list before selecting work, finalizing scope, implementing, opening a change for review, responding to material review discoveries, merging, deploying, updating the punch list, or selecting the next work unit. Never make a state-changing decision from a stale snapshot.

# Invariants
- Fresh-state: the list may change asynchronously. Reconcile new items, edits, status changes, engineering notes, operator notes, priorities, duplicates, related symptoms, previously addressed items, and active work.
- Whole-list: never treat an item as isolated until it has been considered against the entire current list.
- Human-input: operator, tester, user, and validator notes are evidence. Never overwrite, delete, fabricate, or silently dismiss them. If new feedback challenges a previous fix, treat the issue as active again.
- Deployment: implemented is not merged; merged is not deployed; deployed is not validated. Do not request re-testing until the change is available in the validation environment.
- Traceability: keep a recoverable relationship of punch-list item(s) → problem cluster → worklane → change → deployment → validation using existing project systems.

# Interpret items
A punch-list item is an observation, not necessarily a correct diagnosis or implementation specification. Determine what the person is experiencing, what correct product behavior should be, why current behavior occurs, and the smallest coherent fix. Repeated misunderstanding by competent users is itself a potential UX problem. Do not blindly implement the literal requested solution.

# Cluster before fixing
Group related items by shared workflows, screens, components, domain entities, APIs, state transitions, validation rules, backend invariants, failure modes, recovery mechanisms, mental-model mismatches, and root causes.

For each cluster record affected items, priority/severity, symptoms, likely root cause, shared abstraction, current vs desired behavior, existing related work, whether one fix can resolve multiple items, regression risk, and validation strategy.

Prefer one coherent shared fix over repeated item-specific patches when a genuine common cause exists. Do not combine unrelated work merely to reduce the number of changes.

# Research previous work
Before designing a change, determine whether the problem has already been discussed, investigated, partially fixed, implemented elsewhere, reported elsewhere, assigned to active work, included in an open change, merged, deployed, returned for validation, or rejected or reopened. Search the punch list, notes, source, history, worklanes, issues, open and merged changes, review discussions, tests, documentation, and runtime evidence when useful. Do not duplicate existing work.

If a deployed change appears to cover an item, determine whether it predates deployment and only needs validation, is a regression, shows an incomplete fix, or describes a different problem.

# Select work units
Do not mechanically process items top to bottom. After clustering, select the next coherent engineering work unit: one item, several related items, a shared workflow correction, a UX improvement, reusable validation, or a justified refactor.

Prioritize explicit priority, operational blockage, severity, dependencies, number of related problems resolved, regression risk, unblocking subsequent work, and existing active engineering work. Critical work should normally dominate unless dependencies require another ordering.

# Plan before implementation
For each work unit establish coverage, the underlying user/operator problem, evidence, the best-supported root cause, desired observable behavior, the proposed change, validation, and how the change reaches the validation environment. Resolve materially competing hypotheses before implementation.

# Clarification gate
Remain in planning mode until implementation is sufficiently understood. Investigate before asking. Resolve questions from available code, history, tests, logs, documentation, previous changes, or punch-list context whenever possible.

Ask only questions that cannot reasonably be answered through available evidence and could materially change implementation. Number them. For each provide Default, Assumption, Why it matters, and Uncertainty when useful. Accept terse answers such as `1y 2n 3a`. After each clarification batch report `Confidence: NN%`.

Do not implement until confidence is >=95% and there are zero unresolved material questions. Unknown information that does not materially affect implementation is not a reason to block.

# Worklane strategy
Use the project's established isolation mechanism for each coherent work unit. Reuse existing active work when it already represents the problem correctly. Do not create one worklane per punch-list item when several items share a root cause. Do not mix unrelated changes into one worklane. Follow existing naming, issue, commit, and review conventions.

# Implementation standard
Fix the underlying cause rather than only the visible symptom. Prefer domain invariants, reusable workflow logic, explicit state transitions, deterministic handling, clear failure states, recoverable workflows, and consistent behavior across related surfaces. Avoid unnecessary redesign. A broader refactor is justified when punch-list evidence shows the existing abstraction repeatedly causes operational problems. Choose the smallest change that correctly solves the shared problem.

# Testing
Add or update the highest-value tests for reported behavior, intended success, relevant failure, meaningful edge cases, and regressions exposed by related items. When one change addresses several items, test the shared invariant. Run appropriate validation before requesting review. Do not knowingly submit failing work.

# Change review
A review request should explain the operational problem, not merely the code diff. Identify every punch-list item addressed, the cluster, root cause, behavior before and after, approach, tradeoffs, tests, deployment considerations, and validation instructions when useful.

# Review feedback loop
Opening a change is not completion. Monitor checks and human and automated review, evaluate each comment technically, change what is necessary, rerun relevant validation, and repeat until genuinely ready. Do not blindly implement automated reviewer suggestions. If review reveals a broader problem affecting other items, re-scan the punch list and reconsider the cluster before continuing.

# Merge and deployment
Before any issue, branch, worktree, staging, commit, push, PR, merge, or deployment mutation, load repo-local guardrails and work-lane or delivery-gating guidance when they exist, complete read-only safety reconnaissance, and obtain the mandated delivery-consent confirmation. Keep those actions blocked until that gate passes. In Kit-managed repositories this means loading `docs/agents/GUARDRAILS.md` and `work-lane-gating`, asking the required new-lane versus continue-existing question, and recording the Pull-Request Landing Plan before mutation.

Permit merging only after a direct user authorization or an accepted bounded plan explicitly names the exact authorized PR set. Merge is a distinct mutation boundary from PR delivery; approval, passing checks, and punch-list status never imply merge consent.

Before an authorized merge, verify required reviews, passing checks, unresolved feedback, dependencies, safe merge order, and that recent punch-list changes do not invalidate the approach. After merge determine the deployment mechanism, target validation environment, migrations, dependencies, configuration or flags, and deployment success. Verify deployed behavior when possible. If deployment fails or causes a regression, do not request re-testing.

# Engineering notes
Use the punch list's existing engineering-notes mechanism, author prefixes, formatting, status conventions, and re-test terminology. Communicate engineering's understanding, relationships to other items, what changed, whether the change is deployed, expected behavior, exactly what to re-test, and whether clarification is required.

Prefer: `Deployed. Unmatched scans now remain explicitly unmatched instead of presenting a successful scan. Please re-test the unmatched-scan scenario and record any remaining behavior here.` Avoid: `fixed` or `should work now`.

# Status handling
Discover and respect the existing status vocabulary. Do not invent new statuses unless explicitly authorized. Status represents the current handoff state, not merely engineering activity. Where supported, distinguish awaiting engineering, in progress, implemented but unavailable, deployed and awaiting validation, validation complete, and reopened. If the existing model cannot express a necessary state, surface the problem rather than silently changing the workflow.

# Post-deployment handoff
Once a change is available in the validation environment, re-scan the entire punch list, identify every affected item, update engineering notes and statuses using existing conventions, and request re-test where human validation is appropriate. Update every affected item, not only the item that initiated the work.

# Validation loop
Do not assume a requested re-test succeeded. Later scans must inspect validator notes, status changes, new related items, regressions, and additional symptoms. If validation succeeds, follow existing conventions for final resolution. If it fails or produces new concerns, treat the issue as active, inspect the previous reasoning and deployed implementation, and determine whether the feedback is a regression, incomplete implementation, misunderstanding, new requirement, UX deficiency, or another edge case. Then restart the control loop.

# Duplicate handling
Do not silently delete or ignore duplicate reports. Preserve provenance. Associate items that share a root cause, implement the shared fix, update every affected item, and request separate validation when scenarios differ materially.

# Scope heuristic
Before an item-level fix, ask whether another user is likely to encounter essentially the same problem elsewhere; if yes, investigate the shared abstraction. Ask whether this is an implementation defect or a product-model or workflow problem, and whether the product can make correct behavior obvious instead of requiring more training. Do not turn every issue into a redesign.

# Working state
Maintain a concise internal state model of cluster, items, priority, state, worklane, change, deployment, and validation using existing project systems. At meaningful checkpoints report current punch-list state, clusters, root causes, active work, review state, deployment state, items awaiting validation, blockers, and the next action.

# Definition of done
A work unit is not complete until the latest punch list has been re-scanned, affected items are understood, root cause is established, implementation and relevant tests are complete, review feedback is resolved, the change is merged and deployed to the validation environment through the project's delivery and merge gates, affected items contain accurate engineering notes, validators have clear re-test instructions when needed, and subsequent validation feedback is incorporated. Code completion alone is not completion. Do not merge or deploy to satisfy this bar without delivery consent and, for merge, direct user authorization or an accepted bounded plan that names the exact authorized change set.

# Initial execution
When `kp punchlist` begins, do not implement immediately. Locate and read the entire punch list; understand schema, statuses, notes, priorities, and conventions; inventory active items; identify highest-priority unresolved items, duplicates, clusters, notes, feedback, items awaiting re-test, and items potentially covered by existing work; search relevant code, history, issues, changes, tests, and deployments; construct a problem-cluster map; recommend a dependency-aware execution order; identify items that should not result in independent engineering work; ask only material unresolved questions; and report confidence.

Use this initial response structure:

## Punch List state
Concise summary of the current punch list.

## Proposed clusters
For each cluster:

### Cluster N: `<short name>`
- **Items:**
- **Priority:**
- **Shared problem:**
- **Likely root cause:**
- **Existing related work:**
- **Proposed scope:**
- **Expected behavior:**
- **Confidence:**

## Proposed execution order
Explain the dependency-aware order.

## Already addressed / awaiting validation
Identify items that appear to require validation rather than new implementation.

## Clarifications
Ask only material unresolved questions. For each:
- **Default:**
- **Assumption:**
- **Why it matters:**

End with `Confidence: NN%`. Remain in planning mode until confidence is >=95% with zero unresolved material questions.

# Governing principle
A punch list is a continuously updated observational dataset about product behavior. Use the whole list to find patterns. Use previous issues to understand new issues. Use new feedback to challenge previous assumptions. Fix shared causes instead of repeated symptoms. Keep the punch list, engineering work, review state, deployments, and human validation synchronized. Whenever you believe you know what to do next, re-scan the punch list first.
