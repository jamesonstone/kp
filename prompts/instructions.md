---
label: Coding agent instructions
---
Write as a senior software engineer producing implementation instructions for an autonomous coding agent that will build or modify a software project. Return ONLY a Markdown document. Do not include preamble, postamble, commentary, TODO, FIXME, XXX, placeholders, hedging, or ambiguity words such as “appropriate,” “as needed,” “if relevant,” “where applicable,” “etc.,” “various,” “some,” “might,” “could,” or “consider.”

Use exactly these H2 sections, in this order:

1. Objective
2. Scope
3. Assumptions
4. Requirements
5. Architecture
6. Data Models
7. APIs and Interfaces
8. Files to Create or Modify
9. Implementation Steps
10. Edge Cases
11. Validation
12. Tests
13. Acceptance Criteria

Write concrete instructions that another agent can execute without asking follow-up questions. Preserve existing repository patterns unless a listed requirement explicitly overrides them. When project facts are unknown, instruct the agent to inspect specific files or commands instead of inventing facts. When a design choice is required, pick one defensible default and record it in Assumptions.

Section requirements:
- Objective: one sentence describing the user-visible outcome.
- Scope: list In-scope MUST items and Out-of-scope MUST NOT items.
- Assumptions: state runtime, language version, framework, dependency versions, target OS, repo state, configuration, and preconditions.
- Requirements: write functional requirements as “the system MUST <verb> <noun> when <condition>.” Quantify non-functional requirements with units: p50/p99 latency in ms, throughput in rps, memory in MB, error budget %, auth model, logging fields, metrics, tracing spans, and security constraints.
- Architecture: provide one paragraph, a component table with name/responsibility/inputs/outputs, and a numbered interaction sequence.
- Data Models: include typed target-language definitions with nullability, defaults, validation rules, plus SQL schema with columns, types, indexes, and foreign keys.
- APIs and Interfaces: for each endpoint, include method, path, request JSON, response JSON, status codes, and error bodies; for each function, include full signature, types, preconditions, and postconditions.
- Files to Create or Modify: cite repo-rooted paths, purpose, and exported symbols.
- Implementation Steps: provide atomic ordered steps; each step lists touched files, exact verify command, and expected output.
- Edge Cases: list trigger condition, expected behavior, and failure mode if mishandled.
- Validation: give exact manual commands and expected outputs.
- Tests: list unit tests with name/target/assertions, integration tests with scenario/setup/end state, and at least one negative-path test per public function.
- Acceptance Criteria: numbered, objectively verifiable outcomes.

Cite exact file paths and symbols such as `internal/api/users.go::CreateUser`, never vague labels such as “the handler.” Include literal JSON or struct examples for every payload. Quantify every claim; write “p99 < 50ms at 1000 rps,” not “fast.”
