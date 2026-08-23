package cmd

import (
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/clipboard"
)

func TestPromptPrint(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "clarify", "--print")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter: %q", stdout)
	}
	if strings.Contains(stderr, "copied to clipboard") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPRPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "pr", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := "Use the Kit workflow implementation path and create a new worklane for this work: a new issue, canonical worktree, `GH-<issue-number>` branch, and ready pull request in every affected project repository."
	if strings.TrimSpace(stdout) != want {
		t.Fatalf("stdout = %q, want %q", strings.TrimSpace(stdout), want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestMergePromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "merge", "--print")
	if err != nil {
		t.Fatal(err)
	}

	// The approved body stores "~" in place of a backtick so the Mermaid fence
	// and inline code spans survive a Go raw string literal.
	want := strings.ReplaceAll(`# Purpose
Provide a concise, auditable procedure to merge dependent pull requests in dependency-ordered waves so every merge is authorized, evidence-backed, and recoverable.

# Scope & Authorization
Apply only to the exact repository, pull-request set, and expected heads directly authorized by the user or covered by a bounded plan the user explicitly accepted. Adding a pull request, repository, base, deployment environment, infrastructure effect, merge method, or actor is scope expansion and requires follow-up authorization. If an authorized node's dependency closure contains an unauthorized node, that node is BLOCKED: report the exact missing authorization instead of merging the predecessor or silently dropping the dependent.

# Definitions
- MERGE_READY: exact-current head, base, checks, reviews, approvals, policy, actor, and dependencies are all satisfied and attributable.
- BLOCKED: a required gate failed, or an explicit dependency or approval is unmet.
- UNKNOWN: evidence is missing, stale, ambiguous, or not attributable to the expected head, target, policy, or actor. UNKNOWN is not a soft pass.
- MERGED: the repository host confirms the PR at the expected head merged and the resulting base commit is known.
- DEPLOYED: the expected artifact from the merged source is observed in the declared environment.
- ACCEPTED: the exact deployed artifact passed every applicable pre-declared runtime, compatibility, and protected-workload gate.
- Wave: a set of proven-independent MERGE_READY nodes merged together.
- Dependency closure: the transitive set of nodes required before a target node can be MERGE_READY.

# Pre-wave checklist
1. Revalidate authorization, actor identity, and exact expected head OIDs for every node. Record the UTC observation time and a stated freshness bound; evidence older than that bound is UNKNOWN.
2. Bind each node to: repository, base branch, expected head OID, actor, allowed merge method, review policy, required current-head checks, dependency closure, reversibility class, infrastructure effects, and protected-workload invariants.
3. Confirm required status checks against the current head, review approvals, and merge-policy eligibility. Pending, missing, earlier-head, locally-substituted, and skipped-without-proven-eligibility checks never pass. An empty required-check set is not passing evidence; state it explicitly and classify from actual repository policy.
4. Confirm each node is conflict-free against its current base. Exact-head merge authorization never authorizes source mutation. For routine remediation that stays within the pull request's issue and declared scope, prefer updating its existing head under bounded repair authority from the user or an accepted plan: ordinarily merge the current base into the head branch, apply or regenerate the repair, commit, and push without rebasing, force-pushing, or retargeting. Mark the node UNKNOWN, remove it from the frontier, rerun current-head checks and review, and require later exact-head merge authorization. Use a replacement pull request only when the repair materially changes scope or architecture, the original head cannot be updated safely, or repository policy or the user requires replacement.
5. Verify infrastructure and deployment credentials before AWS-dependent evidence or action. In Kit-managed repositories run ~kit aws verify~, use only the verified configured profile, and stop on missing, incomplete, or mismatched credentials, configuration, account, or ARN; never fall back to ambient credentials.
6. Load and obey repo-local rules before any merge or queue transition. In Kit-managed repositories ~docs/agents/GUARDRAILS.md~ and related local rulesets override generic GitHub or plugin defaults.
7. Produce evidence artifacts: Mermaid dependency DAG, wave table, per-node merge method, acceptance signal, and named rollback owner and mechanism.

# Execution algorithm
1. Build a dependency DAG where A --> B means A must be merged, or when the edge requires it deployed and accepted, before B. Derive direction only from authoritative evidence: branch stacking or an explicit base relationship, a declared producer/consumer API, schema, or artifact contract, a required migration or infrastructure prerequisite, a canonical program DAG, or an explicit user or repository requirement. Shared files, symbols, configuration, environments, or databases show coupling to examine but do not establish direction. An asserted edge without authoritative directional evidence is UNKNOWN.
2. Prove graph completeness before allowing concurrency. Compare the diffs of any nodes proposed for the same wave and rule out undeclared coupling; unexamined overlap is a missing edge, not independence.
3. Classify every node exactly as MERGE_READY, BLOCKED, or UNKNOWN. Missing, stale, pending, skipped-without-proven-eligibility, cyclic, provisional, conflicted, or ambiguous evidence never yields MERGE_READY.
4. Select only the zero-unmet-dependency MERGE_READY frontier. Maximize safe concurrency among proven-independent nodes and prioritize nodes that unlock the most downstream work or shorten the critical path. Serialize dependency chains and same-base, deployment-coupled, or otherwise interacting operations whose order affects conflict, queue, release, or deployment state.
5. Immediately before every wave, revalidate authorization, identity, head/base, policy, checks, approvals, dependencies, deployment effects, and rollback readiness. Never bypass protection, reviews, required checks, a merge queue, repository policy, or identity safeguards.
6. Merge only assigned wave members, using a repository-permitted method and the merge queue when policy requires it. Do not enable auto-merge or any deferred merge that executes outside a revalidated wave. Do not pass branch-deletion flags; leave head-branch cleanup to repository policy. Queue admission is not a merge: dependents stay BLOCKED until the queue produces the merge commit on the base.
7. After every merge or queue transition, re-observe and recompute the graph. A completed merge supersedes the base, so return every remaining node sharing that base to UNKNOWN until its required checks re-run against the new base head, or until the repository's require-up-to-date-branch policy is verified to enforce it. The base branch's own post-merge workflow must reach a terminal successful state before the next wave targets that base.
8. Contain failures to the failed node and its dependents, and continue only proven-independent authorized nodes whose readiness remains valid. Do not mutate a head during an authorized wave; perform any authorized in-place remediation between waves, preserve the original pull request, invalidate its old head evidence, and recompute the graph. Never force, bypass, substitute identity, or broaden scope to recover a wave.
9. Stop and escalate when a wave produces no frontier advance, or when the same node returns the same blocker on consecutive revalidations. Repetition without progress is a blocker, not a retry.
10. Repeat until every node in the authorized closure reaches its declared terminal state: merged for merge-only nodes, deployed and accepted for deployment-gated nodes.

# Deployment gates
- Classify every deployment-affecting node as reversible-by-redeploy, forward-fix or schema-retaining cutback, or irreversible/destructive. Use actual recovery behavior: a migration or backfill label alone is insufficient; irreversible means restoration or compensation is required or prior state cannot be fully restored.
- For expand/contract changes, merge, deploy, and accept the backward-compatible expansion before the contract or removal becomes eligible; wait for every required consumer to migrate and never place expansion and contraction in the same wave.
- Declare acceptance before the wave, not after: triggering workflow, target account, environment, region or cluster, expected deployed identity, expected migration or resource action, health or SLO signal, and observation window.
- After deployment, verify that the observed deployed identity, such as image digest, release tag, or commit SHA, equals the merged head OID. A successful pipeline that shipped a different artifact is a failed deployment.
- Before merging a deployment-triggering node, state the exact rollback mechanism and confirm it is currently executable: prior artifact retained and redeployable, revert path open, migration reversal or compensating action defined. Named ownership without a proven mechanism is UNKNOWN.
- Confirm that no change-freeze or release window blocks the wave and that the rollback owner is available for its duration.

# Protected-workload gates
- Before merging or deploying, inventory every existing workload, customer lane, simulation or replay, API, schema, hash, event stream, or data boundary that must remain unchanged.
- For each protected lane record the invariant, configuration and activation state, exact validation command or workflow, required literal result, execution timing, exact artifact or deployment identity, and failure recovery.
- Additive or default-off source is not runtime proof. Keep flags, routing, allowlists, writers, effects, and migrations in their declared safe state; after an actual deployment run every required compatibility gate against the exact artifact and do not activate new behavior until its separate prerequisites pass.
- Task-specific user or canonical-program safety gates augment and outrank this generic prompt.

# Example
~~~mermaid
graph LR
  A --> B
  B --> C
  D --> C
~~~
Wave table (example):
- Wave 1: A, D
- Wave 2: B
- Wave 3: C

# Durable state
When the authorized set spans repositories with dependent deliverables, staged deployment, or expected session handoff, record the graph, wave state, evidence, and next safe action in one canonical durable ledger and checkpoint after every material transition. A chat transcript is not program state.

# Evidence & Recording
For every merge or queue completion record the repository and PR, expected head, resulting merge or squash commit, new base commit, method, actor, and UTC observation time; verify that the resulting base contains the intended change, accounting for synthesized squash or merge-queue commits. Record MERGED, artifact publication, DEPLOYED, activation, ACCEPTED, recovery, and rollback as separate claims attributable to exact identities. Merge success is never deployment, runtime, or production proof. Report exact completed waves, blockers, unknowns, unobserved claims, and the next safe action whenever the closure is not satisfied.

# References
- docs/agents/GUARDRAILS.md
- docs/references/rules/github-pr-merge.md
- docs/references/workflows/pull-request-merge.md
`, "~", "`")
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestContinuePromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "continue", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := "Resolve all issues autonomously and continue until the goal is fully complete; ask permission only before large-scale deletion or deleting sensitive files.\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestHandoffPromptIncludesAcceptanceCriteriaWithoutWordLimit(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "handoff", "--print")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"4. Acceptance Criteria",
		"5. Resource Links",
		"Acceptance Criteria: List success criteria in bullet-list form",
		"Each criterion must be binary verifiable.",
		"Remove any sentence where deletion doesn't degrade the agent's ability to execute correctly.",
		"audit the conversation for ambiguities, contradictions, missing requirements, hidden assumptions",
		"Assume the coding agent has not seen the original conversation.",
		"Include in-scope, out-of-scope, and deferred future work",
		"error handling, logging and observability, security and authorization, compatibility, rollback, operator visibility",
		"List decisions already made as implementation constraints",
	}
	for _, text := range expected {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
	if strings.Contains(stdout, "650") || strings.Contains(stdout, "word limit") {
		t.Fatalf("stdout includes removed word-limit language:\n%s", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPromptCopy(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t, "clarify", "--copy", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called")
	}
	if fake.verified == "" {
		t.Fatal("Verify was not called")
	}
	if fake.pasted {
		t.Fatal("Paste was called for --copy")
	}
	if !strings.Contains(stderr, "✅ 📋 Prompt \"clarify\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPromptDefaultShowsClipboardInstructionsWithSpacing(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t, "clarify", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stderr, "✅ 📋 Prompt \"clarify\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
	if !strings.Contains(stderr, "🧾 Full prompt content is printed to stdout below.\n\n") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPromptDefaultPrintsAndCopiesWithoutPaste(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, _, err := executeTestCommand(t, "clarify", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called")
	}
	if fake.verified == "" {
		t.Fatal("Verify was not called")
	}
	if fake.pasted {
		t.Fatal("Paste was called")
	}
}

func TestPromptCopyVerifyFailureExitsSystem(t *testing.T) {
	fake := &fakeClipboard{verifyErr: clipboard.ErrVerifyFailed}
	_, _, err := executeTestCommand(t, "clarify", "--copy", withClipboard(fake))
	if ExitCode(err) != ExitSystem {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.pasted {
		t.Fatal("Paste was called after verification failure")
	}
}

func TestPromptDefaultVerifyFailureExitsSystem(t *testing.T) {
	fake := &fakeClipboard{verifyErr: clipboard.ErrVerifyFailed}
	_, _, err := executeTestCommand(t, "clarify", withClipboard(fake))
	if ExitCode(err) != ExitSystem {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.pasted {
		t.Fatal("Paste was called after verification failure")
	}
}

func TestPromptVerboseLogsToStderr(t *testing.T) {
	fake := &fakeClipboard{}
	_, stderr, err := executeTestCommand(t, "clarify", "--copy", "--verbose", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(stderr, "event=copy name=clarify bytes=") {
		t.Fatalf("stderr = %q", stderr)
	}
}
