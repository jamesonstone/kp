---
label: Pre-authorize task delivery
---
/goal For this task thread, you are pre-authorized to complete the full delivery lifecycle for all code produced or modified as part of this work. The point is to let you finish delivery without waiting on me for routine approvals.

This authorization includes:
- creating, updating, and pushing branches
- creating and updating pull requests
- addressing review feedback and CI failures
- merging pull requests when required checks pass
- merging dependent pull requests in the appropriate order
- deploying or otherwise delivering the resulting code when deployment is part of the established repository workflow
- performing routine repository operations necessary to complete delivery

Do not ask for additional authorization for individual PR merges, including confirmation of specific PR numbers or commit SHAs, when the PR is part of this task thread.

Treat this authorization as scoped only to changes required to accomplish this task. Do not merge unrelated pre-existing PRs, bypass required protections, ignore failing required checks, or perform destructive/irreversible operations outside the normal delivery workflow. For in-scope merges, this message is the explicit authority your repository rules ask for. Any other action a repository rule reserves for separate approval still needs it: ask when you reach one, and keep doing independent in-scope work meanwhile.

The task is delivered when every in-scope pull request is merged and, where the repository's workflow deploys, the deployment has completed and passed its standard checks.

Continue autonomously until the task is delivered or you encounter a material blocker that cannot be resolved safely without new information from me. When you stop, including to wait on someone else's review or access, lead with the outcome, then report merge, CI, deployment, and runtime evidence as separate claims, mark anything unverified, and if blocked, name exactly what is needed and from whom.
