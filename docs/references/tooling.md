# Tooling Reference

## Purpose

- Record durable repo-wide tooling notes, command references, and local development expectations
- Keep short-lived implementation notes in feature docs instead of here

## Current State

- `make fmt` formats Go source with `gofmt`.
- `make test` runs `go test ./...`.
- `make build` builds the local binary at `./kp` from `./cmd/kp`.
- `make install` installs `./cmd/kp` with `go install`.
- `make clean` removes the local `./kp` build artifact.
- `kit map <feature>` is the preferred feature index command before broad
  repository reads.
- `kit verify <feature> --task <task-id>` records task verification evidence
  when a task declares checks for Kit to run.
- macOS runtime copy behavior depends on `pbcopy` and `pbpaste`.
- Interactive picker behavior depends on user-installed `fzf` unless
  `--no-fzf` is used.
- `hyperfine` is optional local tooling for latency acceptance evidence and is
  not a build dependency.
