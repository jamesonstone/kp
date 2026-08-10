```text
██╗  ██╗██████╗
██║ ██╔╝██╔══██╗
█████╔╝ ██████╔╝
██╔═██╗ ██╔═══╝
██║  ██╗██║
╚═╝  ╚═╝╚═╝   prompt pockets
```

**`kp` is a tiny macOS CLI for keeping reusable coding-agent prompts one
keystroke away.** It prints the prompt you ask for, copies the exact same body
to your clipboard after a read-back verification, and stays out of the way.

<!-- BEGIN KIT-MANAGED README BADGES -->
[![Last commit](https://img.shields.io/github/last-commit/jamesonstone/kp)](https://github.com/jamesonstone/kp/commits) [![Open issues](https://img.shields.io/github/issues/jamesonstone/kp)](https://github.com/jamesonstone/kp/issues) [![Pull requests](https://img.shields.io/github/issues-pr/jamesonstone/kp)](https://github.com/jamesonstone/kp/pulls) [![Release](https://img.shields.io/github/v/release/jamesonstone/kp)](https://github.com/jamesonstone/kp/releases)
<!-- END KIT-MANAGED README BADGES -->

No accounts. No sync service. No paste injection. Just local prompt utilities
that feel good in a shell.

## Why Use It?

- 🧠 Keep high-signal prompts in one predictable place.
- 📋 Copy prompt bodies only after `pbpaste` matches the expected text.
- 🪶 Use bare commands such as `kp clarify`, not a heavy command tree.
- 🧰 Create and edit your own Markdown prompts.
- 🔍 Inspect and stop processes by port when you need to clean up a dev server.
- 🏗️ Scaffold lightweight agent/review docs for other repos.
- 🔒 Stay local-only: no network calls, telemetry, or cloud storage.

`kp` is early, practical, and intentionally small. The current release targets
macOS because it depends on `pbcopy` and `pbpaste`.

## Install From Source

```sh
git clone https://github.com/jamesonstone/kp.git
cd kp
make build
./bin/kp
```

For a local `go install`:

```sh
make install
kp --version
```

## Quick Start

```sh
# open the interactive launcher
kp

# see the grouped help page
kp --help

# print and copy a built-in prompt
kp clarify

# browse prompts interactively
kp list

# inspect listeners on a port
kp find-port 4005

# add your own prompt
kp new rubber-duck

# preview repo support-file scaffolding
kp scaffold --dry-run
```

The launcher keeps its first view focused on prompts, Find port, and Help. Use
`j`/`k`, the arrow keys, or Tab/Shift-Tab to move; the preview pane wraps long
text. Select Help for the complete prompt-management, scaffolding, version, and
utility command reference.

## Command Guide

| Command                    | What It Does                                                      | Best Used In                      | Status          |
| -------------------------- | ----------------------------------------------------------------- | --------------------------------- | --------------- |
| `kp`                       | Open the focused launcher for prompts, Find port, and Help.       | Interactive terminals             | ✅ ready        |
| `kp --help`                | Show the grouped help page without opening the launcher.          | Terminal discovery, docs, scripts | ✅ ready        |
| `kp clarify`               | Print and copy the clarify-before-implementing prompt.            | Coding-agent chats                | ✅ ready        |
| `kp continue`              | Print and copy the autonomous-continuation prompt.                | Coding-agent chats                | ✅ ready        |
| `kp handoff`               | Print and copy the coding-agent handoff prompt.                   | Coding-agent chats                | ✅ ready        |
| `kp parentthread`          | Print and copy the parent-thread response prompt.                 | Coding-agent chats                | ✅ ready        |
| `kp pr`                    | Print and copy the issue/branch/pull-request workflow prompt.     | Coding-agent chats, repo handoff  | ✅ ready        |
| `kp <name> --copy`         | Copy a prompt without printing it.                                | Shell aliases, scripts            | ✅ ready        |
| `kp <name> --print`        | Print a prompt without touching the clipboard.                    | Pipes, inspection                 | ✅ ready        |
| `kp find-port <port>`      | Inspect processes listening on a port and choose an action.       | Dev-server cleanup                | ✅ ready        |
| `kp port-find <port>`      | Alias for `kp find-port <port>`.                                  | Dev-server cleanup                | ✅ ready        |
| `kp list`                  | Open the emoji-enhanced `fzf` prompt picker.                      | Interactive terminals             | ✅ ready        |
| `kp list --no-fzf`         | Use the numbered fallback picker.                                 | Minimal terminals                 | ✅ ready        |
| `kp list --plain`          | Print prompt names, one per line.                                 | Scripts, completion experiments   | ✅ ready        |
| `kp list --verbose`        | Print `name`, `label`, and `source`.                              | Prompt audits                     | ✅ ready        |
| `kp new <name>`            | Create a user prompt in your editor.                              | Prompt authoring                  | ✅ ready        |
| `kp edit <name>`           | Edit a user prompt or promote a built-in first.                   | Prompt tuning                     | ✅ ready        |
| `kp rm <name>`             | Remove a user prompt; built-ins stay protected.                   | Prompt cleanup                    | ✅ ready        |
| `kp scaffold`              | Create reusable repo support files.                               | New or lightly configured repos   | ✅ ready        |
| `kp scaffold --dry-run`    | Preview scaffold actions without writing files.                   | Before touching a repo            | ✅ ready        |
| `kp scaffold --dir <path>` | Scaffold a specific directory.                                    | Scripts, test checkouts           | ✅ ready        |
| `kp scaffold --force`      | Overwrite scaffold files except `.gitignore`.                     | Refreshing support files          | ✅ ready        |
| `kp --version`             | Print version and commit metadata.                                | Diagnostics                       | ✅ ready        |
| Release packaging          | Homebrew, Goreleaser, GitHub releases.                            | Release automation                | 🚫 out of scope |
| Automatic paste            | Cmd+V injection into GUI apps.                                    | Focused apps                      | 🚫 out of scope |

## Built-In Prompts

| Prompt         | Label                       | Status      |
| -------------- | --------------------------- | ----------- |
| `clarify`      | Clarify before implementing | ✅ embedded |
| `continue`     | Continue autonomously       | ✅ embedded |
| `handoff`      | Coding agent handoff        | ✅ embedded |
| `parentthread` | Parent thread response      | ✅ embedded |
| `pr`           | Pull request workflow       | ✅ embedded |

Prompt names are bare commands. There is no `kp prompt ...` namespace.

Reserved names cannot be used for prompts: `help`, `list`, `new`, `edit`,
`rm`, `scaffold`, `prompt`, and `version`.

## Your Prompt Library

User prompts live here:

- `$XDG_CONFIG_HOME/kp/prompts/` when `XDG_CONFIG_HOME` is set
- `~/.config/kp/prompts/` otherwise

Use `--config <dir>` to read and write prompts under `<dir>/prompts` for one
invocation.

Prompt files are Markdown with optional YAML frontmatter:

```markdown
---
label: Clarify before implementing
---

Clarify before implementing. Stay in planning mode.
```

Only `label` is metadata. Frontmatter is not copied, printed, or previewed.
User prompts shadow built-ins with the same name.

## Repo Scaffolding

`kp scaffold` is a pocket-sized version of Kit's repo setup helpers. It writes
agent and review support files without turning the target into a Kit project.

Created when missing:

- `.env`
- `.envrc`
- `.coderabbit.yaml`
- `.github/pull_request_template.md`
- `.github/copilot-instructions.md`
- `AGENTS.md`
- `CLAUDE.md`
- `docs/agents/README.md`
- `docs/agents/WORKFLOWS.md`
- `docs/agents/RLM.md`
- `docs/agents/TOOLING.md`
- `docs/agents/GUARDRAILS.md`
- `docs/references/README.md`
- `docs/references/testing.md`
- `docs/references/tooling.md`
- `docs/references/external-systems.md`

`.gitignore` is append-only: existing lines are preserved, and missing local
environment or scratch patterns are added once. Existing files are skipped
unless `--force` is passed; `--dry-run` prints planned actions without writing
files.

Not created by `kp scaffold`: `.kit.yaml`, `.kit/`, global Kit config,
`docs/specs/`, `docs/notes/`, `docs/PROJECT_PROGRESS_SUMMARY.md`, and
`docs/CONSTITUTION.md`.

## Requirements

| Requirement                                   | Why                                                         |
| --------------------------------------------- | ----------------------------------------------------------- |
| macOS 13+ on `darwin/arm64` or `darwin/amd64` | Clipboard support uses macOS tools.                         |
| Go 1.22+                                      | Builds the local binary.                                    |
| `pbcopy` and `pbpaste`                        | Copy and verify prompt bodies.                              |
| `fzf`                                         | Powers `kp` and `kp list`; install with `brew install fzf`. |
| `KP_EDITOR`, `EDITOR`, or `vi`                | Opens prompts for `new` and `edit`.                         |

## Development

```sh
make build
make test
go vet ./...
```

The project favors direct, boring Go: small packages, explicit errors, stable
CLI output, and tests around filesystem and clipboard boundaries.

## Contributing

Issues and pull requests are welcome while the project is young. Good changes
keep `kp` small, local-first, and easy to reason about.

Before opening a PR:

1. Run `make test`.
2. Run `go vet ./...`.
3. Update this README when command behavior changes.
4. Keep new dependencies rare and justified.

## Exit Codes

| Code | Meaning                                                               |
| ---- | --------------------------------------------------------------------- |
| 0    | Success                                                               |
| 1    | User input or prompt data error                                       |
| 2    | Clipboard, platform, or system-command failure                        |
| 3    | Config, filesystem bootstrap, editor lookup, or missing `fzf` failure |
| 130  | User cancellation                                                     |

## License

MIT. See [LICENSE](LICENSE).

## Maintainers

Maintained with 🪖 and ❤️ by [Jameson](https://github.com/jamesonstone) (`jamesonstone`).
