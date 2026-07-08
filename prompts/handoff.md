---
label: Coding agent handoff
---

Produce a handoff from a generic chat system (ChatGPT-style) to a coding-focused agent (Codex-style). Return ONLY a Markdown document.

Goal: convert a long, non-linear brainstorm thread built from notes, links, discussions, and documents into a deterministic `Coding Agent Handoff` that Codex can execute, including external resource links.

Hard constraints:

- Context Synthesis, Source Map, and Coding Agent Instructions must contain only signal. Remove any sentence where deletion doesn't degrade the agent's ability to execute correctly.
- Before writing the handoff, audit the conversation for ambiguities, contradictions, missing requirements, hidden assumptions, and implementation decisions that could materially change the result. Ask only implementation-changing questions; do not ask preference or cosmetic questions.
- Acceptance Criteria: List success criteria in bullet-list form (however many items required for comprehensive coverage). Each criterion must be binary verifiable.
- Use paragraph + list form: include at least one paragraph and one numbered list.
- Do not include preamble, postamble, commentary, TODO, FIXME, XXX, placeholders, or ambiguous wording.
- Do not use hedging or uncertainty words such as “appropriate,” “as needed,” “if relevant,” “where applicable,” “etc.,” “various,” “some,” “might,” “could,” or “consider.”
- Use exactly these H2 sections, in this order:

1. Context Synthesis
2. Source Map
3. Coding Agent Instructions
4. Acceptance Criteria
5. Resource Links

- Every non-trivial claim in `Context Synthesis` and `Coding Agent Instructions` MUST include one or more source tags such as `[S1]`.
- When evidence is missing, write `UNKNOWN` and add a concrete repository inspection action.
- When sources conflict, write `CONFLICT` and state the tie-break decision.
- Assume the coding agent has not seen the original conversation. Include every implementation-relevant decision, constraint, source, link, and unresolved item required to execute the task.

Section requirements:

- Context Synthesis: one paragraph that states objective, affected users, constraints, scope boundaries, measurable definition of done, and selected direction. Include in-scope, out-of-scope, and deferred future work when the conversation establishes those boundaries.
- Source Map: bullet list of evidence entries labeled `[S1]`, `[S2]`, and so on. Each entry MUST include source type (`note`, `link`, `doc`, or `discussion`), a one-line claim, and the URL or document identifier.
- Coding Agent Instructions: one paragraph followed by a numbered execution plan with 6-10 steps. The paragraph MUST state the chosen direction and key tradeoffs. The numbered plan MUST instruct Codex to:

1. inspect the repository and identify relevant existing functionality by exact file path and symbol,
2. reconcile brainstorm decisions with actual code behavior,
3. produce a complete, fully detailed implementation strategy grounded in current codebase context,
4. enumerate concrete file edits, interfaces, data model changes, dependency updates, configuration changes, migration steps, error handling, logging and observability, security and authorization, compatibility, rollback, operator visibility, validation commands, and tests,
5. tie validation commands and expected outputs to the Acceptance Criteria,
6. state risks, open questions, and explicit assumptions with mitigation and owner.

- Acceptance Criteria: bullet list of success criteria for the overall task. Include however many bullet-list items are required for comprehensive coverage. Each criterion MUST be binary verifiable.
- Resource Links: bullet list of entries labeled `[R1]`, `[R2]`, and so on. Each entry MUST include title, URL, and one-line relevance. Include all external resources referenced in the thread plus directly related official documentation discovered from those references. If none exist, output `- NONE`.

Output quality requirements:

- Preserve decisions already made in the conversation.
- List decisions already made as implementation constraints; do not reopen them unless repository evidence conflicts.
- Do not invent repository facts; when unknown, direct Codex to inspect exact files or commands.
- Prefer deterministic, executable language.
- Include all available external resources in `Resource Links` and keep each entry concise.
- Wrap the entire result in a single fenced code block with language tag `markdown`, with no text outside the fence, so it is directly copy-pasteable as a Codex task brief.
