# Testing Reference

## Purpose

- Record durable repo-wide testing guidance that is broader than one feature
- Keep feature-specific testing details in the current feature's `PLAN.md` or `TASKS.md`

## Current State

- Use `go test ./...` as the default repository-wide automated correctness
  check for Go runtime changes.
- Use `make test` when validating the documented local tooling path.
- Use `make fmt` or `gofmt` before completing Go source edits.
- Use temp config roots for command tests and manual prompt lifecycle checks;
  do not write tests against the real `~/.config/kp/prompts` directory.
- Keep Darwin-only clipboard integration checks build-tagged or skipped so the
  package remains inspectable from non-Darwin machines while runtime clipboard
  support stays Darwin-only.
- Verify clipboard copy behavior with `pbcopy`/`pbpaste` on macOS and compare
  prompt bodies byte-for-byte after frontmatter stripping.
- Treat `kp <prompt>` as print-plus-copy behavior; it must not perform paste
  side effects.
- Record missing optional tools such as `hyperfine` as validation gaps instead
  of marking performance acceptance complete without evidence.
