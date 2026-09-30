---
label: Pre-authorize task delivery
---
/goal For this task thread, you are pre-authorized to complete the full delivery lifecycle for all code produced or modified as part of this work, so you can finish without waiting on me for routine approvals.

This authorization includes:
- creating, updating, and pushing branches
- creating and updating pull requests
- addressing review feedback and CI failures
- merging pull requests when required checks pass
- merging dependent pull requests in the appropriate order
- deploying or otherwise delivering the resulting code when deployment is part of the established repository workflow
- performing routine repository operations necessary to complete delivery

Do not ask for additional authorization for individual PR merges, including confirmation of specific PR numbers or commit SHAs, when the PR is part of this task thread.

Treat this authorization as scoped only to changes required to accomplish this task. Do not merge unrelated pre-existing PRs, bypass required protections, ignore failing required checks, or perform destructive/irreversible operations outside the normal delivery workflow. This message is the explicit merge authority your repository rules ask for; anything else they reserve for separate approval still needs it.

The task is delivered when every in-scope pull request is merged and, where the workflow deploys, deployed and passing its standard checks.

Continue autonomously until the task is delivered or you encounter a material blocker that cannot be resolved safely without new information from me. When you stop, report the outcome with evidence, keeping merge, deployment, and runtime results separate, and if blocked, say what you need and from whom.
