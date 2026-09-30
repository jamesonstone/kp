# Coding-Agent Prompt Practice

## Purpose

- Review built-in prompts (`prompts/*.md` and the `kp task` template) against
  current vendor guidance before changing them.
- Last reviewed 2026-09-30 for issue #68. The guidance covers GPT-6 and
  GPT-5.6 (OpenAI, Codex) and current Claude models (Anthropic, Claude Code).
  Re-check the sources when a new model generation ships.
- Each prompt's spec still wins over this rubric; record any deliberate
  deviation in that spec.

## The Core Test

Keep a sentence only if it tells the model something it cannot work out on
its own: the outcome and why it matters, the scope and authority boundaries,
the user's preferences, the check that defines done, the shape of the report,
or a guard against a failure that has actually happened. Leave investigation,
sequencing, tool choice, and strategy to the model. A prompt that grows when
it is "improved" is a signal to re-run this test; issue #68 trimmed its own
first revision (review had grown from 188 to 405 words) after applying it.

## Rubric

1. Lead with the outcome, a check that defines done, and when to stop or hand
   back. Let the agent choose the path. (O2, O4, O3, A6)
2. Give the reason behind a constraint; models generalize from explained
   intent. (A1, A3, O6)
3. Keep one source of truth. Contradictions cost more than missing detail.
   Point to repository rules instead of copying them; where a prompt must
   restate policy, say which source wins and ask the agent to name
   conflicts. (O2, C2, O4, U2, O1)
4. Reserve absolutes (NEVER, MUST) for real invariants and say what to do
   instead; broad emphasis now causes over-triggering. (A1, C1, O2, O3)
5. Investigate before acting and ground claims in files that were actually
   opened. (A1, O2, A4, C1)
6. Name the check that proves done and ask for evidence; report failed,
   skipped, or unverified checks as exactly that. Avoid blanket "verify
   again" instructions: the newest models over-verify. (C1, A3, A5, O2, A2,
   O1, O3)
7. Deliver what was asked. Report pre-existing problems and adjacent ideas as
   follow-ups instead of widening scope. (A1, A4, A2, A5, O1)
8. Place each state-changing action in a tier: proceed, proceed and report, or
   stop and ask. Make approval the last step, not a recurring interruption.
   Prompt text is advisory; permissions and hooks enforce. (O2, O1, A1, A6,
   C2)
9. Code review: search broadly, then report only verified problems the change
   introduces, with severity and file and line; list pre-existing issues and
   unconfirmed concerns separately; skip nits. "Only report high severity"
   makes literal models under-report. (O8, O7, C3, C5, C1, A2)
10. Final reports lead with the outcome, then evidence, caveats, and the next
    action. Do not ask the agent to write out its reasoning. (O2, A2, A3, A4,
    C1, O5)
11. Prefer the smallest set of high-signal instructions; heuristics over
    scripts. Trimming stale text is most of the work when models change. (E1,
    A1, A3, O3, C2, C3)

## Sources

- A1 https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices
- A2 https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-opus-5
- A3 https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-fable-5
- A4 https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-fable-5-1
- A5 https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-sonnet-5-5
- A6 https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-opus-5-5
- C1 https://code.claude.com/docs/en/best-practices
- C2 https://code.claude.com/docs/en/memory
- C3 https://code.claude.com/docs/en/code-review
- C5 https://github.com/anthropics/claude-code/blob/main/plugins/code-review/README.md
- E1 https://www.anthropic.com/engineering/effective-context-engineering-for-ai-agents
- O1 https://developers.openai.com/api/docs/guides/latest-model
- O2 https://developers.openai.com/api/docs/guides/prompt-guidance-gpt-5p6
- O3 https://developers.openai.com/blog/rethinking-skills-and-prompts-for-gpt-6-astra
- O4 https://learn.chatgpt.com/guides/best-practices
- O5 https://developers.openai.com/cookbook/examples/gpt-5/codex_prompting_guide
- O6 https://developers.openai.com/blog/custom-code-review-rules-for-codex
- O7 https://alignment.openai.com/scaling-code-verification/
- O8 https://github.com/openai/codex/blob/main/codex-rs/prompts/templates/review/rubric.md
- U2 https://cursor.com/docs/context/rules
