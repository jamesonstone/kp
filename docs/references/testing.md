# Testing Reference

## Purpose

- Record the project's durable commands, suites, environments, automation, and evidence expectations
- Follow `rules/testing-and-environment-validation.md` for the mandatory cross-project testing and production-safety contract
- Keep feature-specific testing details in the current feature's `SPEC.md` Validation Map and Evidence sections; legacy staged flows may still use `PLAN.md` or `TASKS.md`

## Code-Level Validation

| Layer | Command | PR workflow or check | Required | Notes |
| --- | --- | --- | --- | --- |
| Formatting | `test -z "$(gofmt -l prompts.go cmd internal)"` | none | yes | Go source and tests must already be formatted. |
| Unit and component | `go test ./...` | none | yes | Covers prompt parsing, registry, command behavior, clipboard boundaries, ports, config, and scaffolding. |
| Race detection | `go test -race ./...` | none | yes | Run before pull-request handoff. |
| Static analysis | `go vet ./...` | none | yes | Run across the complete module. |
| Build | `go build ./...` and `make build` | none | yes | `make build` produces the local Darwin CLI at `bin/kp`. |

## High-Level Suites

| Suite | Type | Environment | Command | Automation | Evidence |
| --- | --- | --- | --- | --- | --- |
| CLI prompt acceptance | component | local | Build, then run `./bin/kp --config <empty-temp-dir> <prompt> --print`, `list --plain`, and `--help` | manual fallback | terminal output recorded in the feature spec |

Production end-to-end and live-integration suites are `NOT_APPLICABLE`; `kp`
is a local CLI with no deployed service or production environment.

## Environment Preflights

- Use Go 1.22 or newer.
- Use an empty temporary `--config` directory for CLI prompt acceptance so a
  user override cannot shadow the built-in under test.
- macOS is required for real clipboard acceptance; unit tests inject clipboard
  fakes, and unsupported-platform tests cover the non-Darwin boundary.

## Credentials And Test Data

- No credentials, external services, production data, or synthetic records are
  required.

## Evidence And Retention

- Record feature validation in the active `SPEC.md`.
- No high-level run directory or `tests/RUN_STATUS.md` is required while the
  project has no end-to-end or live-integration suite.

## Automation And Fallbacks

- Run the code-level commands locally in the table's order before pull-request
  handoff.
- The repository currently has no hosted pull-request correctness workflow.
  Local results do not substitute for hosted checks; report hosted validation
  as unavailable until such a workflow exists.

## Known Gaps

- GitHub-hosted format, test, race, vet, and build checks are unavailable.
- Real `pbcopy`/`pbpaste` acceptance is manual and is not required for a
  `--print`-only prompt addition covered by injected clipboard tests.
