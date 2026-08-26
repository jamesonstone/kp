# References

## Purpose

- Use `rules/aws-agent-toolkit-guidance.md` before AWS-dependent work to select current Agent Toolkit skills, official documentation, the AWS MCP Server or CLI fallback, verified identity, name-aware material targets, infrastructure approval, and secret-safe handling
- Use `rules/testing-and-environment-validation.md` before implementation and validation to preserve code-level checks and add environment evidence safely
- This directory holds durable repo-local references that are broader than one feature
- Keep long-lived background context here instead of in injected top-level instruction files
- Link these files from feature front matter references when they materially shape work
- Store durable rulesets under `rules/<slug>.md` and link them with `kit rules link` instead of copying rules into agent instruction files
- Use `rules/kit-capabilities-usage.md` in downstream projects for Kit command discovery guidance
- Use `rules/feature-notes.md` when deciding how to load, reference, promote, or ignore source material under `docs/notes/<feature>`
- Use `rules/constitution-curation.md` after implementation and validation to keep the Constitution aligned with demonstrated project-wide truth
- Use `rules/infrastructure-change-approval.md` before mutating public-cloud resources, Kubernetes resources or cluster state, or infrastructure-as-code source, configuration, or state to require one plan-level confirmation per batch, one-pass execution, and explicit confirmation for deletion or removal
- Use `rules/cross-repository-program-coordination.md` when dependent multi-repository delivery, staged activation, or agent/session handoff requires one canonical program ledger and reconciled ready frontier
- Use `rules/testing-and-environment-validation.md` before implementation and validation, including browser automation and browser testing, to preserve code-level checks, browser lifecycle ownership, and environment evidence safely
- Use `rules/source-file-size.md` before editing implementation/source or test files and for whole-project reconcile audits
- Use `rules/codex-thread-initialization.md` to preserve Codex's ordered pre-response rename and pin gate during instruction refresh and reconciliation
- Use `worktrees.md` when present for the canonical native Git worktree hierarchy, naming, shared-state model, safety contract, and optional manual convenience commands
- Use `kit rules add` to import or activate available registry rulesets from the Kit GitHub `main` branch
- Use `kit rules view <slug>` to preview a local or registry ruleset before importing it
- Use `kit init --refresh` to adopt existing registry rules into `.kit.yaml` registry state and pick up safe upstream ruleset updates
- Use `kit rules add --custom` for the interactive `$EDITOR` ruleset builder
- `kit rule` is the singular alias for `kit rules`

- Use `rules/agent-completion-output.md` before substantial terminal completion
- Use `rules/human-authorship.md` before any commit, pull request, issue, review comment, or other attribution text so only the human user is displayed as author
- Use `rules/deadline-mode.md` only when the user explicitly signals a real deadline or time constraint, to narrow testing and implementation scope without weakening required approvals, security, or migration/compatibility invariants
- Use `rules/deletion-safety.md` before designing deletion behavior or deleting persistent project, user, business, or external-system state to require recoverable soft delete by default and exact manual confirmation before hard delete
## Starter Files

- `testing.md` — repo-wide testing norms and evidence expectations
- `tooling.md` — local tooling and command references that are broader than one feature
- `external-systems.md` — durable notes about external systems, APIs, or integrations
- `rules/` — pointer-loaded durable rulesets such as frontend UI rules, testing rules, API conventions, security constraints, or domain rules
- `../notes/<feature>/` — optional feature source material; not canonical truth and private contents remain ignored
