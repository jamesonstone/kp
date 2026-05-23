# Testing Reference

## Purpose

- Record durable repo-wide testing guidance that is broader than one feature
- Keep feature-specific testing details in the current feature's `PLAN.md` or `TASKS.md`

## Current State

- Add the default test, lint, and verification commands for this repository.
- Use temp directories for tests that write config or generated files.
- Never claim tests passed unless the commands ran in the current checkout.
