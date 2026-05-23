# kp

**Tiny prompt pockets for your terminal.**

`kp` keeps your reusable coding-agent prompts close at hand. Ask for a prompt by
name, and `kp` prints the clean body while copying the same text to the macOS
clipboard after an exact read-back check. No paste magic, no cloud cupboard, no
extra ceremony.

## Quick Start

```sh
make build
./bin/kp
./bin/kp clarify
./bin/kp list
./bin/kp scaffold --dry-run
```

`make build` writes `./bin/kp`. `make install` installs the local binary with
version metadata through `go install ./cmd/kp`.

## Command Guide

| Command | Description | Use It In | Status |
| ------- | ----------- | --------- | ------ |
| `kp` | Show the grouped help page with prompt commands separated from library commands. | Terminal | ✅ ready |
| `kp clarify` | Print the `clarify` prompt body and copy it to the clipboard after exact verification. | Terminal, coding-agent chat input | ✅ ready |
| `kp instructions` | Print the coding-agent instruction prompt body and copy it to the clipboard after exact verification. | Terminal, coding-agent chat input | ✅ ready |
| `kp parentthread` | Print the parent-thread response prompt body and copy it to the clipboard after exact verification. | Terminal, coding-agent chat input | ✅ ready |
| `kp <name> --copy` | Copy a prompt without printing it. | Terminal, shell aliases, scripts | ✅ ready |
| `kp <name> --print` | Print a prompt without touching the clipboard. | Terminal, pipes, inspection | ✅ ready |
| `kp list` | Open the emoji-enhanced `fzf` picker with prompt previews. | Interactive terminal | ✅ ready |
| `kp list --no-fzf` | Pick from a numbered fallback list when `fzf` is unavailable. | Minimal terminal sessions | ✅ ready |
| `kp list --plain` | Print prompt names, one per line. | Scripts, shell completion experiments | ✅ ready |
| `kp list --verbose` | Print `name`, `label`, and `source` columns. | Auditing prompt overrides | ✅ ready |
| `kp new <name>` | Create a user prompt and open it in the configured editor. | Terminal editor workflow | ✅ ready |
| `kp edit <name>` | Edit a user prompt or promote a built-in prompt before editing. | Terminal editor workflow | ✅ ready |
| `kp rm <name>` | Remove a user prompt; built-ins are protected. | Terminal cleanup | ✅ ready |
| `kp scaffold` | Create repo support files such as agent instructions, CodeRabbit config, PR template, `.env`, and `.envrc`. | New or lightly configured repos | ✅ ready |
| `kp scaffold --dry-run` | Show what scaffold would create, update, or skip without writing files. | Before touching a repo | ✅ ready |
| `kp scaffold --dir <path>` | Scaffold a specific directory instead of the current directory. | Scripts, tests, alternate checkouts | ✅ ready |
| `kp scaffold --force` | Overwrite scaffold files except `.gitignore`, which remains append-only. | Refreshing generated support files | ✅ ready |
| `kp --version` | Print version and commit metadata. | Terminal, diagnostics | ✅ ready |
| Release packaging | Homebrew, Goreleaser archives, and GitHub release workflows. | Release automation | 🚫 out of scope |
| Automatic paste | Cmd+V injection through `osascript`. | Focused GUI apps | 🚫 out of scope |

## Built-In Prompts

| Prompt | Label | Status |
| ------ | ----- | ------ |
| `clarify` | Clarify before implementing | ✅ embedded |
| `instructions` | Coding agent instructions | ✅ embedded |
| `parentthread` | Parent thread response | ✅ embedded |

Prompt names are bare commands. There is no `kp prompt ...` namespace, and
prompt names cannot be `help`, `list`, `new`, `edit`, `rm`, `prompt`, or
`version`.

## Prompt Files

User prompts live in `$XDG_CONFIG_HOME/kp/prompts/` when `XDG_CONFIG_HOME` is
set, and `~/.config/kp/prompts/` otherwise. Use `--config <dir>` to read and
write prompts under `<dir>/prompts` for one invocation.

Prompt files are Markdown with optional YAML frontmatter:

```markdown
---
label: Clarify before implementing
---
Clarify before implementing. Stay in planning mode.
```

Only `label` is metadata. Frontmatter is not copied, printed, or previewed.
User prompts shadow built-ins with the same name.

## Scaffold Files

`kp scaffold` is a small pocket version of Kit's repo setup helpers. It writes
the reusable support files without turning the target into a Kit project.

Created when missing:

- `.env`
- `.envrc`
- `.coderabbit.yaml`
- `.github/pull_request_template.md`
- `AGENTS.md`
- `CLAUDE.md`
- `.github/copilot-instructions.md`
- `docs/agents/README.md`
- `docs/agents/WORKFLOWS.md`
- `docs/agents/RLM.md`
- `docs/agents/TOOLING.md`
- `docs/agents/GUARDRAILS.md`
- `docs/references/README.md`
- `docs/references/testing.md`
- `docs/references/tooling.md`
- `docs/references/external-systems.md`

`.gitignore` is updated by appending missing local-environment and scratch
patterns. Existing files are skipped unless `--force` is passed; `--dry-run`
prints planned actions without writing files.

Not created by `kp scaffold`: `.kit.yaml`, `.kit/`, global Kit config,
`docs/specs/`, `docs/notes/`, `docs/PROJECT_PROGRESS_SUMMARY.md`, and
`docs/CONSTITUTION.md`.

## Requirements

| Requirement | Why |
| ----------- | --- |
| macOS 13+ on `darwin/arm64` or `darwin/amd64` | Clipboard support uses macOS tools. |
| Go 1.22+ | Builds the local binary. |
| `pbcopy` and `pbpaste` | Copy and verify prompt bodies. |
| `fzf` | Powers `kp list`; install with `brew install fzf`. |
| `KP_EDITOR`, `EDITOR`, or `vi` | Opens prompts for `new` and `edit`. |

## Exit Codes

| Code | Meaning |
| ---- | ------- |
| 0 | Success |
| 1 | User input or prompt data error |
| 2 | Clipboard, platform, or system-command failure |
| 3 | Config, filesystem bootstrap, editor lookup, or missing `fzf` failure |
| 130 | User cancellation |

## 👤 Maintainer

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/jamesonstone">
        <img src="https://github.com/jamesonstone.png" width="100px;" alt="Jameson Stone"/>
        <br />
        <sub><b>Jameson Stone</b></sub>
      </a>
      <br />
      <sub>Lead Maintainer</sub>
    </td>
  </tr>
</table>
