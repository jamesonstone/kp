---
kit_metadata_version: 1
artifact: brainstorm
feature:
  id: 0001
  slug: v0-init-utility
  dir: 0001-v0-init-utility
references:
  - id: constitution
    name: Project constitution
    type: repo_doc
    target: docs/CONSTITUTION.md
    relation: constrains
    read_policy: must
    used_for: project authority order, CLI-first architecture rules, dependency policy, validation rules, progress tracking
    status: active
  - id: progress-summary
    name: Project progress summary
    type: repo_doc
    target: docs/PROJECT_PROGRESS_SUMMARY.md
    relation: informs
    read_policy: evidence
    used_for: verified feature phase and absence of prior completed feature artifacts
    status: active
  - id: kit-map
    name: Kit map for v0-init-utility
    type: command
    target: kit map 0001-v0-init-utility
    selector: kit map 0001-v0-init-utility
    selector_type: command
    relation: informs
    read_policy: evidence
    used_for: verified current feature docs, missing SPEC/PLAN/TASKS, no incoming or outgoing relationships, and optional feature notes link
    status: active
  - id: agents-entrypoint
    name: Agents docs entrypoint
    type: repo_doc
    target: docs/agents/README.md
    relation: guides
    read_policy: must
    used_for: confirmed progressive loading rule and repo-local docs as system of record
    status: active
  - id: workflows
    name: Workflow rules
    type: repo_doc
    target: docs/agents/WORKFLOWS.md
    relation: guides
    read_policy: must
    used_for: spec-driven workflow, source-of-truth order, and readiness gate expectations
    status: active
  - id: guardrails
    name: Guardrails
    type: repo_doc
    target: docs/agents/GUARDRAILS.md
    relation: constrains
    read_policy: must
    used_for: completion bar, no placeholder sections, and docs-first safety rules
    status: active
  - id: rlm
    name: RLM
    type: repo_doc
    target: docs/agents/RLM.md
    relation: guides
    read_policy: must
    used_for: prior-work pass, context loading limits, and reference recording requirements
    status: active
  - id: tooling
    name: Tooling rules
    type: repo_doc
    target: docs/agents/TOOLING.md
    relation: guides
    read_policy: must
    used_for: canonical skills lookup, dispatch constraints, and worktree rules
    status: active
  - id: top-level-agent-routing
    name: Top-level agent routing files
    type: repo_doc
    target: AGENTS.md, CLAUDE.md, .github/copilot-instructions.md
    relation: informs
    read_policy: evidence
    used_for: verified these files are routing tables and should remain aligned with docs/agents
    status: active
  - id: repo-inventory
    name: Repository inventory
    type: command
    target: rg --files --hidden -g '!.git'
    selector: rg --files --hidden -g '!.git'
    selector_type: command
    relation: informs
    read_policy: skip
    used_for: historical planning snapshot; verified the repository had docs/config only before v0-init-utility implementation added Go source, module files, prompts, Makefile, and tests
    status: stale
  - id: kit-config
    name: Kit configuration
    type: repo_file
    target: .kit.yaml
    relation: constrains
    read_policy: must
    used_for: feature directory naming, specs dir, skills dir, and constitution path
    status: active
  - id: readme
    name: README
    type: repo_file
    target: README.md
    relation: informs
    read_policy: skip
    used_for: historical brainstorm baseline; README now documents the implemented v0-init-utility command surface
    status: stale
  - id: gitignore
    name: Git ignore rules
    type: repo_file
    target: .gitignore
    relation: constrains
    read_policy: must
    used_for: verified Go artifact ignores, env ignores, and Kit local-state ignores already exist
    status: active
  - id: github-pr-template
    name: GitHub pull request template
    type: repo_file
    target: .github/pull_request_template.md
    relation: informs
    read_policy: evidence
    used_for: verified expected PR evidence fields are description, how to test, and ticket
    status: active
  - id: coderabbit-config
    name: CodeRabbit config
    type: repo_file
    target: .coderabbit.yaml
    relation: informs
    read_policy: evidence
    used_for: verified docs and top-level agent routing files are excluded from CodeRabbit review
    status: active
  - id: testing-reference
    name: Testing reference
    type: repo_doc
    target: docs/references/testing.md
    relation: informs
    read_policy: conditional
    used_for: durable repo-wide testing guidance added during implementation reflection
    status: active
  - id: tooling-reference
    name: Tooling reference
    type: repo_doc
    target: docs/references/tooling.md
    relation: informs
    read_policy: conditional
    used_for: durable repo-wide tooling guidance added during implementation reflection
    status: active
  - id: external-systems-reference
    name: External systems reference
    type: repo_doc
    target: docs/references/external-systems.md
    relation: informs
    read_policy: conditional
    used_for: verified no durable external-system notes exist yet; release tap assumptions need confirmation before becoming durable
    status: optional
  - id: feature-notes
    name: Feature notes
    type: notes
    target: docs/notes/0001-v0-init-utility
    relation: informs
    read_policy: conditional
    used_for: optional pre-brainstorm research input
    status: optional
---
# BRAINSTORM

## SUMMARY

`v0-init-utility` should turn the current docs-only `kp` scaffold into a
Darwin-first Go CLI for fast prompt selection, clipboard verification, and
Cmd+V paste. The likely direction is a small Cobra-based command tree with a
prompt registry, embedded built-ins, user overrides, macOS clipboard/paste
integration, and local build/install/test tooling. Prompt operations should use
the bare command surface rather than a `prompt` subcommand. The approved initial
built-ins are `clarify` and `instructions`; `plan` is explicitly deferred from
this feature.

## USER THESIS

# SPEC: `kp` — utility CLI

## 1. Objective
Build a Go CLI named `kp` that provides fast, reliable prompt-paste functionality (replacing Espanso's keystroke injection, which fails under macOS Secure Input) and serves as an extensible host for future utility subcommands. v1 ships a single command tree with built-in and user-defined prompts, fuzzy selection, clipboard-with-verification, and Cmd+V paste.

## 2. Scope

**In-scope (v1):**
- Subcommand framework supporting future utility commands (extensible architecture)
- `kp` subcommand: list, pick (interactive), get/paste (direct), new, edit, rm
- Built-in prompts: `instructions`, `clarify`, `plan` (embedded via `//go:embed`)
- User prompt overrides at `~/.config/kp/prompts/*.md`
- Clipboard copy with MD5 verification before paste
- Cmd+V paste via `osascript` (survives macOS Secure Input)
- macOS arm64 + amd64 distribution via goreleaser + homebrew tap

**Out-of-scope (v1):**
- Linux/Windows support
- Cloud sync, network calls, telemetry
- TUI beyond `fzf` shell-out

## 3. Assumptions
- runtime: Go 1.22+
- platform: macOS 13+ (Darwin); CI matrix: macos-14 (arm64), macos-13 (amd64)
- framework: `github.com/spf13/cobra` v1.8+
- external tools available on user PATH: `fzf` (brew), `pbcopy`, `pbpaste`, `osascript` (macOS built-in)
- distribution: `goreleaser` v2+, homebrew tap at `github.com/jamesonstone/homebrew-tap`
- license: MIT
- module path: `github.com/jamesonstone/kp`
- editor: `$EDITOR` env var, default `vi`; override via `KP_EDITOR`

## 4. Requirements

**Functional:**
1. The system MUST embed built-in prompts at compile time via `//go:embed prompts/*.md`.
2. The system MUST load user prompts from `$XDG_CONFIG_HOME/kp/prompts/` (default `~/.config/kp/prompts/`) at every invocation.
3. User prompts with the same name as a built-in MUST shadow the built-in.
4. `kp <name>` MUST copy the body to the clipboard and issue Cmd+V by default.
5. The system MUST verify clipboard contents match the expected MD5 before issuing Cmd+V; on mismatch within 250ms (5 retries × 50ms), exit 2 without pasting.
6. `kp <name> --copy` MUST copy without pasting.
7. `kp <name> --print` MUST write body to stdout and exit; no clipboard or paste side effects.
8. `kp ` (no name) MUST launch `fzf` over the prompt list with a preview pane; selection invokes paste flow.
9. `kp new <name>` MUST create `~/.config/kp/prompts/<name>.md` and open in `$EDITOR`.
10. `kp edit <name>` MUST open the user prompt in `$EDITOR`; if the name resolves to a built-in only, copy the built-in body to the user dir first, then open.
11. `kp rm <name>` MUST delete only user prompts; built-ins are not deletable.

**Non-functional:**
- latency: `kp prompt <name> --copy` p99 < 100ms over 50 invocations on M-series macOS.
- latency: `kp prompt <name>` (with paste) p99 < 200ms.
- latency: `kp prompt` (interactive) startup to fzf-visible < 80ms.
- binary size: < 10MB stripped.
- memory: < 30MB RSS during interactive picker.
- observability: `--verbose` flag emits structured logs (key=value) to stderr; exit codes documented and stable.

## 5. Architecture

Single statically-linked Go binary. Cobra subcommand tree rooted at `kp`. Each top-level subcommand is a self-contained package under `internal/cmd/`, designed so new utility subcommands can be added without touching existing code. The `prompt` subcommand composes from `internal/prompt/` (registry + types) and `internal/clipboard/` (copy/verify/paste). Built-in prompts are markdown files in `prompts/` embedded via `embed.FS`. The registry overlays user files on top of embedded built-ins at load time.

| Component | Responsibility | Inputs | Outputs |
|---|---|---|---|
| `cmd/kp` | binary entrypoint; cobra root; `--version` | `os.Args`, env | exit code |
| `internal/cmd/prompt` | `prompt` subcommand tree | flags, args | stdout/stderr |
| `internal/prompt` | registry: load, list, get, add, remove | `embed.FS`, user dir path | `[]Prompt`, errors |
| `internal/clipboard` | macOS clipboard + paste with verification | string body | error |
| `internal/config` | XDG path resolution, dir creation | env, defaults | paths, errors |

Interaction sequence (for `kp prompt <name>`):
1. `main.go` calls cobra `Execute()`.
2. `prompt` subcommand parses positional + flags.
3. Subcommand instantiates registry: `prompt.NewRegistry(embeddedFS, userDir)`.
4. Registry returns `Prompt`, source-tagged.
5. Subcommand calls `clipboard.CopyAndPaste(body)`.
6. `CopyAndPaste` invokes `pbcopy`, polls `pbpaste` until MD5 matches, then runs `osascript -e 'tell application "System Events" to keystroke "v" using command down'`.
7. On success: log `✅ <name>` to stderr; exit 0.

## 6. Data Models

```go
// internal/prompt/prompt.go
type Source int

const (
    SourceBuiltIn Source = iota
    SourceUser
)

type Prompt struct {
    Name     string // unique key, kebab-case; e.g. "instructions"
    Label    string // human label; defaults to titlecased Name
    Source   Source
    FilePath string // empty for built-in; absolute path for user
    Body     string // full text including any leading "---" separator
}

type Registry interface {
    List() []Prompt                              // sorted by Name
    Get(name string) (Prompt, error)             // ErrNotFound if missing
    Add(name, body string) (Prompt, error)       // writes to user dir; ErrEmpty if body=="" 
    Remove(name string) error                    // user only; ErrBuiltIn if name is built-in only
    PromoteToUser(name string) (Prompt, error)   // copies built-in body to user dir
}

var (
    ErrNotFound = errors.New("prompt not found")
    ErrEmpty    = errors.New("prompt body is empty")
    ErrBuiltIn  = errors.New("cannot modify built-in prompt; promote first")
    ErrExists   = errors.New("prompt already exists")
)
```

**Prompt file format** (markdown with optional YAML frontmatter):
```markdown
---
label: Coding agent instructions
---
You are a senior software engineer producing implementation instructions...
```
Frontmatter is optional. If absent, `Label` is `strings.Title(strings.ReplaceAll(name, "-", " "))`.

## 7. APIs and Interfaces

### CLI surface

```
kp                                  # root help
kp --version                        # version + commit
kp                            # interactive picker, paste on select
kp <name>                    # direct paste
kp <name> --copy             # copy only, no paste
kp <name> --print            # body → stdout
kp list                      # name per line
kp list --verbose            # name<TAB>label<TAB>source
kp new <name>                # opens $EDITOR; writes user file
kp edit <name>               # opens $EDITOR; promotes built-in if needed
kp rm <name>                 # removes user file
kp --no-fzf                  # fallback numbered-list picker (for environments without fzf)
```

Global flags: `--verbose`, `--config <dir>` (override XDG), `--help`, `--version`.

### Exit codes (stable, documented)

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | user error (prompt not found, invalid name, ErrExists on `new`) |
| 2 | system error (clipboard failure, paste failure, Secure Input) |
| 3 | config/IO error (cannot read or write config dir; editor not found) |
| 130 | user cancelled (Ctrl+C in fzf or editor) |

### Internal Go API

```go
// internal/clipboard/clipboard.go
type Clipboard interface {
    Copy(body string) error
    Read() (string, error)
    Verify(expected string, timeout time.Duration) error
    Paste() error
    CopyAndPaste(body string) error
}

func New() Clipboard // returns darwin impl; panics on non-darwin in v1
```

```go
// internal/prompt/registry.go
func NewRegistry(builtin embed.FS, userDir string) (Registry, error)
```

## 8. Files to Create or Modify

```
kp/
├── cmd/kp/main.go                       # entrypoint; cobra root; version
├── internal/cmd/
│   ├── prompt.go                        # prompt subcommand root
│   ├── prompt_pick.go                   # default (no-args) interactive + direct paste
│   ├── prompt_list.go                   # list
│   ├── prompt_new.go                    # new
│   ├── prompt_edit.go                   # edit (with built-in promotion)
│   └── prompt_rm.go                     # rm
├── internal/prompt/
│   ├── prompt.go                        # types, errors
│   ├── registry.go                      # registry impl
│   ├── frontmatter.go                   # YAML frontmatter parser
│   └── *_test.go
├── internal/clipboard/
│   ├── clipboard.go                     # interface + darwin impl (build tag)
│   └── clipboard_test.go                # darwin build tag
├── internal/config/
│   ├── config.go                        # XDG paths, dir bootstrap
│   └── config_test.go
├── prompts/                             # //go:embed source
│   ├── instructions.md
│   ├── clarify.md
│   └── plan.md
├── .goreleaser.yaml
├── .github/workflows/release.yaml       # goreleaser on tag push
├── Makefile                             # build, test, install, release
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

## 9. Implementation Steps

Each step MUST compile and pass tests before the next. Each step lists touched files and the verification command.

1. **Scaffold module.**
   - touches: `go.mod`, `cmd/kp/main.go`, `Makefile`
   - actions: `go mod init github.com/jamesonstone/kp`; add cobra; root command with `--version`.
   - verify: `go build ./... && ./kp --version` → prints semver placeholder.

2. **Config package.**
   - touches: `internal/config/config.go`, `internal/config/config_test.go`
   - actions: resolve `$XDG_CONFIG_HOME`, default `~/.config`; expose `PromptsDir()`; create dir with 0700 on demand.
   - verify: `go test ./internal/config` passes; runtime check `ls ~/.config/kp/prompts/` after first call.

3. **Prompt types + frontmatter.**
   - touches: `internal/prompt/prompt.go`, `internal/prompt/frontmatter.go`, `internal/prompt/frontmatter_test.go`
   - actions: define `Prompt`, `Source`, errors; implement frontmatter parser using `gopkg.in/yaml.v3`.
   - verify: `go test ./internal/prompt -run Frontmatter` passes table-driven cases.

4. **Built-in prompts + embed.**
   - touches: `prompts/instructions.md`, `prompts/clarify.md`, `prompts/plan.md`, `internal/prompt/registry.go`
   - actions: paste full bodies from prior conversation into the three files; declare `//go:embed prompts/*.md` in registry.
   - verify: `go test ./internal/prompt -run BuiltIn` confirms 3 prompts loaded.

5. **Registry implementation.**
   - touches: `internal/prompt/registry.go`, `internal/prompt/registry_test.go`
   - actions: implement `List`, `Get`, `Add`, `Remove`, `PromoteToUser`; user overrides built-ins by name.
   - verify: `go test ./internal/prompt` passes including override test, rm-builtin-rejection, add-empty-rejection.

6. **Clipboard package (darwin).**
   - touches: `internal/clipboard/clipboard.go` (`//go:build darwin`), `clipboard_test.go`
   - actions: `Copy` via `exec.Command("pbcopy")` with stdin pipe; `Read` via `pbpaste`; `Verify` polls MD5 every 50ms up to 5x; `Paste` via `osascript`; `CopyAndPaste` composes them.
   - verify: `go test ./internal/clipboard` passes on darwin; manual check that `Copy` populates pbpaste.

7. **`list` subcommand.**
   - touches: `internal/cmd/prompt.go`, `internal/cmd/prompt_list.go`
   - actions: register `list` under `prompt`; `--verbose` prints tab-separated name/label/source.
   - verify: `./kp prompt list` shows `instructions\nclarify\nplan`.

8. **`pick` subcommand (core).**
   - touches: `internal/cmd/prompt_pick.go`
   - actions: no positional → fzf shell-out with preview; with positional → direct. Flags: `--copy`, `--print`, `--no-fzf`. Default action is paste.
   - verify: `./kp prompt clarify --print` writes body to stdout; `./kp prompt clarify --copy && pbpaste | head -1` shows `---`.

9. **`new` / `edit` / `rm` subcommands.**
   - touches: `internal/cmd/prompt_new.go`, `prompt_edit.go`, `prompt_rm.go`
   - actions: `new` writes empty body, opens editor, then reloads from disk; `edit` calls `PromoteToUser` if name is built-in-only; `rm` removes user file; all reject built-in modification.
   - verify: round-trip `kp prompt new test; kp prompt list | grep test; kp prompt rm test`.

10. **CI + release.**
    - touches: `.github/workflows/release.yaml`, `.goreleaser.yaml`
    - actions: release workflow runs `goreleaser release` on tag.
    - verify: `goreleaser release --snapshot --clean` produces `dist/kp_*.tar.gz` for both architectures.

11. **README.**
    - touches: `README.md`
    - actions: install via `brew install jamesonstone/tap/kp`; usage; skhd hotkey example; Secure Input note.
    - verify: human read-through; commands in README execute as documented.

## 10. Edge Cases

- **fzf not on PATH:** detect via `exec.LookPath("fzf")`; if missing, print `install fzf via 'brew install fzf' or use --no-fzf` and exit 3.
- **`--no-fzf` mode:** print numbered list to stderr, read single line from stdin; non-numeric or out-of-range → exit 1.
- **Built-in prompt edit:** first edit copies embedded body to user dir, then opens editor. User sees this in stderr: `promoted built-in 'instructions' to ~/.config/kp/prompts/instructions.md`.
- **Clipboard verification timeout:** after 5×50ms retries, MD5 still mismatched → exit 2 with `clipboard write failed: expected <hash> got <hash>`. Do NOT issue Cmd+V; would paste wrong content.
- **Secure Input held:** osascript paste may fail or silently no-op. Detect via paste exit code; if non-zero, exit 2 with hint: `paste failed; check Secure Input holder via 'ioreg -l -w 0 | grep SecureInputPID'`.
- **Empty prompt body on `new`:** if editor exits with empty file, exit 1 and delete the stub file.
- **`$EDITOR` unset, vi missing:** exit 3 with `set $EDITOR or install vi`.
- **Concurrent invocations:** two `kp prompt <name>` in parallel race on clipboard. Documented as known limitation; first to verify wins.
- **Invalid prompt name:** names MUST match `^[a-z][a-z0-9-]*$`; reject `Add`/`new` with exit 1.
- **Prompt name collision on `new`:** `ErrExists` → exit 1 with `prompt 'x' already exists; use 'edit' instead`.
- **Ctrl+C in fzf or editor:** detect non-zero exit; do not modify clipboard; exit 130.

## 11. Validation

Manual verification sequence on macOS 14 arm64:

1. `brew install fzf && go install ./cmd/kp` → `which kp` returns path.
2. `kp --version` → prints version string.
3. `kp list` → outputs `clarify\ninstructions\nplan` (sorted).
4. `kp list --verbose` → 3 rows, all `source=builtin`.
5. Open TextEdit, focus a doc; run `kp prompt clarify` → body pastes; stderr shows `✅ clarify`.
6. `kp clarify --copy && pbpaste | head -c 3` → outputs `---`.
7. `kp clarify --print | wc -c` → byte count matches embedded file.
8. `kp new test1` → editor opens; save body `hello`; exit; `kp prompt list` includes `test1`.
9. `kp test1` in TextEdit → `hello` pasted.
10. `kp edit clarify` → editor opens with built-in body; save; `kp prompt list --verbose` shows `clarify` with `source=user`.
11. `kp rm test1; kp prompt list | grep test1` → no match.
12. `time kp clarify --copy` → real < 100ms on M-series.

## 12. Tests

**Unit:**
- `prompt.TestFrontmatter_WithLabel`: `"---\nlabel: X\n---\nbody"` → `Prompt{Label:"X", Body:"body"}`.
- `prompt.TestFrontmatter_None`: `"body"` → `Prompt{Label:"<titlecased name>", Body:"body"}`.
- `prompt.TestFrontmatter_MalformedYAML`: returns parse error.
- `prompt.TestRegistry_UserOverrides`: user `instructions.md` shadows built-in; `Get` returns user; `Source==SourceUser`.
- `prompt.TestRegistry_AddRejectsEmpty`: `Add("x", "")` returns `ErrEmpty`.
- `prompt.TestRegistry_RemoveBuiltinFails`: `Remove("instructions")` with no user override returns `ErrBuiltIn`.
- `prompt.TestRegistry_PromoteToUser_CopiesBody`: built-in body becomes user file; subsequent `Get` returns `SourceUser` with identical body.
- `prompt.TestNameValidation_RejectsInvalid`: names `Foo`, `1abc`, `foo bar` return validation error.
- `clipboard.TestCopyVerify_Match` (darwin): `Copy(s); Verify(s, 250ms)` returns nil.
- `clipboard.TestVerify_TimesOutOnMismatch` (darwin): after `Copy("a")`, `Verify("b", 250ms)` returns timeout error in < 300ms.
- `config.TestPromptsDir_XDGOverride`: `XDG_CONFIG_HOME=/tmp/x` yields `/tmp/x/kp/prompts`.
- `config.TestPromptsDir_DefaultsToHome`: env unset → `~/.config/kp/prompts`.

**Integration (darwin build tag):**
- `IntegrationTest_ListEmbedded`: shell out to compiled `kp prompt list`; assert stdout contains 3 names.
- `IntegrationTest_CopyClipboard`: `kp prompt clarify --copy`; `pbpaste` returns clarify body.
- `IntegrationTest_RoundTripCRUD`: in temp `XDG_CONFIG_HOME`, run `new` → `list` → `edit` → `rm` → `list`; assert state at each step.
- `IntegrationTest_PromotionOnEdit`: `kp prompt edit clarify` with `EDITOR=true` (no-op editor); assert user file exists with built-in body.

**Negative-path (at least one per public function):**
- `TestPick_PromptNotFound`: `kp prompt nonexistent` exits 1.
- `TestNew_NameCollision`: `kp prompt new instructions` exits 1 with `ErrExists`.
- `TestNew_InvalidName`: `kp prompt new "Bad Name"` exits 1.
- `TestRm_BuiltinNotRemovable`: `kp prompt rm instructions` exits 1 unless promoted.
- `TestPick_NoFzfAvailable`: with mocked PATH, `kp prompt` exits 3.

## 13. Acceptance Criteria

1. `go test ./...` exits 0 on macOS arm64 with all unit + integration tests.
2. `kp list` exits 0 and prints `clarify`, `instructions`, `plan` (one per line, sorted).
3. `kp clarify` in a focused text field results in the body appearing within 200ms p99 over 50 runs.
4. `kp clarify --copy` exits 0 and `pbpaste` matches the embedded body byte-for-byte.
5. `kp clarify --print` writes the body to stdout and exits 0.
6. `kp new test; kp prompt list | grep -q test; kp prompt rm test` exits 0 end to end.
7. `kp edit clarify` (with no-op `EDITOR`) creates `~/.config/kp/prompts/clarify.md` and `kp prompt list --verbose | grep clarify` shows `source=user`.
8. `goreleaser release --snapshot --clean` succeeds and produces `dist/kp_*_darwin_arm64.tar.gz` and `dist/kp_*_darwin_amd64.tar.gz`, each containing a binary < 10MB.
9. `kp clarify --copy` p99 latency < 100ms over 50 invocations (`hyperfine 'kp prompt clarify --copy'`).
10. Adding a new top-level subcommand (e.g. `kp foo`) requires only a new package under `internal/cmd/foo` and one line in the root command registration — no edits to existing subcommand packages.

## RELATIONSHIPS

No prior feature relationships are in scope. `kit map 0001-v0-init-utility`
reports no incoming relationships, no outgoing relationships, and only this
feature's `BRAINSTORM.md` is present under `docs/specs/0001-v0-init-utility/`.

## CODEBASE FINDINGS

1. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/CONSTITUTION.md`
   constrains this feature to a CLI-first, local, scriptable utility with
   deterministic behavior unless explicitly required otherwise. It also says
   source code and package manifests do not yet exist and must not be inferred
   from absent files.
2. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/CONSTITUTION.md`
   requires small explicit internals, command parsing at the edge, conventional
   stdin/stdout/stderr/exit-code behavior, clear errors, narrow public APIs,
   and standard-library preference unless a dependency has a documented purpose.
3. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/CONSTITUTION.md`
   defines the progress invariant: `docs/PROJECT_PROGRESS_SUMMARY.md` must
   reflect the highest completed artifact per feature. The current summary
   already lists feature `0001` in `brainstorm` phase.
4. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/agents/WORKFLOWS.md`
   classifies this as spec-driven work because it creates a new feature and
   substantial runtime architecture. The next durable artifact after this phase
   should be `SPEC.md`, not code.
5. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/agents/RLM.md`
   requires index-first prior-work discovery. `kit map 0001-v0-init-utility`
   and `docs/PROJECT_PROGRESS_SUMMARY.md` show no prior feature docs to inspect.
6. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/notes/0001-v0-init-utility/`
   contains only `.gitkeep`, so there are no usable note files to import. The
   feature notes reference remains optional.
7. Repository inventory via `rg --files --hidden -g '!.git'` shows no `go.mod`,
   no `go.sum`, no `cmd/`, no `internal/`, no `prompts/`, no `Makefile`, no
   `.goreleaser.yaml`, and no workflow under `.github/workflows/`. The feature
   will create the runtime and release surface from scratch.
8. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.kit.yaml` sets
   `docs/specs` as the specs dir, `.agents/skills` as the skills dir,
   `docs/CONSTITUTION.md` as the constitution path, `goal_percentage: 95`, and
   feature naming with a four-digit numeric prefix plus hyphen separator.
9. There is no `.agents/skills` directory in the current repo. The current
   brainstorm front matter has no `skills` list, and no legacy `SPEC.md`
   `## SKILLS` table exists because `SPEC.md` has not been created.
10. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/README.md` currently
    contains only the title and "CLI Prompt Utilities" tagline. Usage,
    local installation, Secure Input rationale, command reference, and known
    limitations will need to be authored as part of this feature.
11. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.gitignore` already
    excludes common Go binaries/test artifacts/coverage files, `.env`, Go
    workspace files, and Kit cache/state/temp artifacts. No release `dist/`
    ignore is required while Goreleaser is deferred.
12. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.github/pull_request_template.md`
    expects PRs to include description, testing steps, and ticket closure. The
    future implementation should include concrete macOS test evidence there.
13. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.coderabbit.yaml`
    excludes `docs/**`, `AGENTS.md`, and `CLAUDE.md` from CodeRabbit review, so
    implementation files will be reviewed while this brainstorm update will not.
14. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/references/testing.md`,
    `docs/references/tooling.md`, and
    `docs/references/external-systems.md` are placeholders. No additional
    durable repo-wide testing, tooling, or external-system guidance exists yet.
15. The original user thesis had a namespace conflict: several requirements and
    validation examples use root-level commands such as `kp <name>`, `kp list`,
    `kp new`, `kp edit`, and `kp rm`, while architecture sections describe
    `internal/cmd/prompt` and examples such as `kp prompt list` and
    `kp prompt clarify`. The first clarification response rejects the `prompt`
    subcommand and requests bare commands. The second clarification response
    confirms `kit clarify` was a typo; the canonical command is `kp clarify`.
16. The original user thesis had a prompt-body/frontmatter conflict: the data model says
    `Body` includes any leading `---` separator and validation expects
    `pbpaste | head -c 3` to output `---`, but the proposed unit test expects
    frontmatter input to parse into `Body:"body"`. The first clarification
    response accepts stripped-frontmatter semantics: YAML frontmatter is
    metadata and is excluded from copied, printed, and pasted prompt bodies.
17. The user thesis names built-in prompts `instructions`, `clarify`, and
    `plan`. The first clarification response provides exact bodies for
    `instructions` and `clarify`; the third clarification response approves a
    concise default `plan` body for now, with the expectation that it may be
    updated later.
18. The user thesis originally required MD5 verification before Cmd+V. The
    first clarification response makes exact clipboard byte equality
    contractual, while checksum algorithm is not contractual. Error logs may
    include short checksums or byte counts for diagnostics only.
19. The user thesis assumes `fzf`, `pbcopy`, `pbpaste`, and `osascript` on PATH.
    `pbcopy`, `pbpaste`, and `osascript` are macOS built-ins, while `fzf` is a
    user-installed dependency. The first clarification response accepts the
    explicit-fallback behavior: missing `fzf` exits `3` unless `--no-fzf` is
    passed.
20. The user thesis uses both "v0-init-utility" as the feature name and "v1" as
    the shipped scope. SPEC should make clear whether this feature builds the
    initial production v1 CLI surface despite the feature slug beginning with
    `v0`; the first clarification response accepts this framing.
21. The first clarification response removes Homebrew from this release. The
    second clarification response also defers Goreleaser release archives and
    release workflows to a later feature. This feature should cover local
    build/install and tests only.

## AFFECTED FILES

1. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/go.mod` - create the
   Go module `github.com/jamesonstone/kp`; declare Go version and dependencies.
2. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/go.sum` - generated
   dependency checksum file once dependencies are added.
3. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/cmd/kp/main.go` -
   binary entrypoint, version variables, and root command execution.
4. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/internal/cmd/` -
   command packages for root command registration and bare prompt command
   behavior; no `prompt` subcommand should be implemented unless a later answer
   changes this.
5. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/internal/prompt/` -
   prompt model, name validation, frontmatter parsing, embedded built-in load,
   user overlay registry, CRUD operations, and tests.
6. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/internal/clipboard/` -
   Darwin clipboard copy/read/verify/paste implementation, interface seams for
   tests, and macOS build tags.
7. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/internal/config/` -
   XDG config resolution, config override handling, prompt directory creation,
   and tests.
8. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/prompts/` - embedded
   markdown built-ins for `instructions` and `clarify`; `plan` is out of scope
   for this feature.
9. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/Makefile` - build,
   test, local install, and formatting targets if accepted.
10. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/README.md` - local
    install, command reference, Secure Input explanation, examples, and known
    limitations.
11. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/LICENSE` - MIT
    license if confirmed.
12. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.gitignore` -
    update only if local implementation-generated files need ignoring; release
    `dist/` output is not in scope while Goreleaser is deferred.
13. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/specs/0001-v0-init-utility/SPEC.md`
    - next canonical artifact after this brainstorm.
14. `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/docs/PROJECT_PROGRESS_SUMMARY.md`
    - update when the feature advances beyond brainstorm.
15. Deferred out of this feature:
    `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.goreleaser.yaml`
    and
    `/Users/jamesonstone/go/src/github.com/jamesonstone/kp/.github/workflows/release.yaml`.

## DEPENDENCIES

1. Current verified repo tooling:
   - Kit via `.kit.yaml` for feature docs, instruction scaffolding, and progress
     tracking.
   - GitHub PR template for human review evidence.
   - CodeRabbit config for implementation review path filters.
2. Proposed Go/runtime dependencies from the user thesis:
   - Go 1.22+ for the CLI implementation.
   - `github.com/spf13/cobra` v1.8+ for command parsing and future subcommand
     registration.
   - `gopkg.in/yaml.v3` if YAML frontmatter remains in scope.
   - Standard library packages for `embed`, `os/exec`, filesystem access,
     checksums, sorting, timing, and tests.
3. Proposed system/runtime tools:
   - `pbcopy` for clipboard writes on macOS.
   - `pbpaste` for clipboard reads and verification on macOS.
   - `osascript` for issuing Cmd+V through System Events.
   - `fzf` for interactive prompt selection.
   - `$KP_EDITOR`, `$EDITOR`, or `vi` for prompt editing.
4. Release dependencies:
   - Homebrew tap publishing is explicitly out of scope for this release.
   - Goreleaser release archives and GitHub release workflows are deferred to a
     later feature.
   - This feature should require local build/install and test tooling only.
5. Dependency policy from `docs/CONSTITUTION.md`:
   - Each dependency must have a documented purpose.
   - Prefer the standard library where it is sufficient.
   - Avoid runtime network dependencies for v1.
6. Accepted built-in prompt body for `prompts/instructions.md`:

   ```markdown
   ---
   label: Coding agent instructions
   ---
   You are a senior software engineer producing implementation instructions for another coding agent. Output ONLY a markdown document (no preamble, no postamble, no commentary) containing these 13 sections in this exact order: 1. Objective (one sentence, user-visible outcome); 2. Scope (In-scope MUST list, Out-of-scope MUST NOT list); 3. Assumptions (runtime, language version, framework, dependencies with versions, target OS, preconditions); 4. Requirements (functional as "the system MUST <verb> <noun> when <condition>"; non-functional quantified with units: latency p50/p99 ms, throughput rps, memory MB, error budget %, auth model, observability hooks); 5. Architecture (one paragraph + component table with name/responsibility/inputs/outputs + numbered interaction sequence); 6. Data Models (typed definitions in target language with nullability/defaults/validation + SQL schema with columns, types, indexes, foreign keys); 7. APIs and Interfaces (per endpoint: method, path, literal JSON request/response examples, status codes, error bodies; per function: full signature with types, pre/postconditions); 8. Files to Create or Modify (repo-rooted paths, purpose, exported symbols); 9. Implementation Steps (atomic, ordered, each step compiles and passes tests before the next, lists touched files and exact verify command + expected output); 10. Edge Cases (trigger condition, expected behavior, failure mode if mishandled); 11. Validation (exact manual commands + expected outputs); 12. Tests (unit with name/target/assertions, integration with scenario/setup/end state, at least one negative-path test per public function); 13. Acceptance Criteria (numbered, objectively verifiable: command exits 0, response matches schema, metric below threshold). Rules: cite exact file paths and symbol names like `internal/api/users.go::CreateUser`, never "the handler"; quantify every non-functional claim ("p99 < 50ms at 1000 rps", never "fast"); include literal JSON or struct examples for all payloads; do NOT emit TODO, FIXME, XXX, or placeholder text—pick a defensible default and record it in Assumptions; do NOT hedge ("might", "could", "perhaps", "consider"); use imperative verbs at sentence start; replace every ambiguity word ("appropriate", "as needed", "if relevant", "etc.", "where applicable", "various", "some") with a concrete value, list, or number.
   ```

7. Accepted built-in prompt body for `prompts/clarify.md`:

   ```markdown
   ---
   label: Clarify before implementing
   ---
   Clarify before implementing. Stay in planning mode. Ask numbered questions with defaults, assumptions, and uncertainty until >=95% confidence and 0 unresolved. Accept y/n shorthand. Report confidence each batch.
   ```

8. Deferred candidate built-in prompt body for `prompts/plan.md`:

   ```markdown
   ---
   label: Plan before implementing
   ---
   Stay in planning and information-gathering mode. Do not implement code, write production changes, or move into execution. Read the repository instructions and relevant docs first. Produce a numbered implementation plan with assumptions, risks, affected files, validation steps, and unresolved questions. Ask for approval before execution.
   ```

   `SPEC.md` supersedes this brainstorm note and excludes `prompts/plan.md`
   from v0-init-utility.

## QUESTIONS

1. Resolved: do not implement a `prompt` subcommand. Prompt operations should
   use the bare `kp` command surface, such as `kp clarify`, `kp list`, `kp new`,
   `kp edit`, and `kp rm`.
2. Resolved: reserve command names such as `list`, `new`, `edit`, `rm`,
   `prompt`, `help`, and `version`; prompt names must not shadow commands.
3. Resolved: YAML frontmatter is metadata and is excluded from copied, printed,
   and pasted prompt body.
4. Resolved: exact initial bodies are approved for `instructions` and
   `clarify`. `plan` is removed from this release.
5. Resolved: no-arg interactive mode fails with exit `3` if `fzf` is missing
   unless `--no-fzf` is explicitly passed.
6. Resolved: after writing to the clipboard, `kp` must read the clipboard back
   and verify exact equality with the expected prompt body before Cmd+V. On
   mismatch after 5 retries over 250ms, exit `2` without pasting. Error logs may
   include short checksums or byte counts, but checksum algorithm is not
   contractual.
7. Resolved: prompt CRUD user mistakes such as invalid names, duplicate names,
   and missing prompts use exit `1`; system clipboard/paste failures use exit
   `2`; config/editor failures use exit `3`; user cancellation uses exit `130`.
8. Resolved: `new` creates a stub file, opens the editor, and deletes the stub
   if the saved body is empty.
9. Resolved: Homebrew is not part of this release.
10. Resolved: keep feature slug `v0-init-utility` and describe the product
    scope as the initial v1 CLI surface.
11. Resolved: defer Goreleaser release archives and GitHub release workflow to a
    later feature; keep this feature focused on local build/install and tests.

## OPTIONS

1. Accepted CLI namespace option - bare prompt UX:
   - Shape: `kp <name>`, `kp list`, `kp new`, `kp edit`, `kp rm`, and no
     `kp prompt ...` subcommand.
   - Benefits: fastest daily workflow and matches the paste utility's core use.
   - Costs: future utility subcommands can collide with prompt names unless
     reserved names and collision rules are explicit.
2. Rejected CLI namespace option - prompt subcommand:
   - Shape: `kp prompt <name>`, `kp prompt list`, `kp prompt new`.
   - Benefit: clean namespacing.
   - Rejection reason: user explicitly rejected the `prompt` subcommand.
3. Rejected CLI namespace option - root shortcuts plus canonical `prompt`
   namespace:
   - Shape: root shortcuts plus duplicate `kp prompt ...`.
   - Benefit: future namespace clarity.
   - Rejection reason: user explicitly rejected the `prompt` subcommand.
4. Prompt body option A - frontmatter stripped:
   - Shape: YAML frontmatter provides metadata only; copy/print/paste use body
     after frontmatter.
   - Benefits: conventional markdown metadata behavior and aligns with proposed
     frontmatter tests.
   - Costs: conflicts with the validation example expecting copied output to
     start with `---` unless built-in prompt content intentionally starts with a
     non-frontmatter separator.
5. Prompt body option B - frontmatter included:
   - Shape: full file bytes are copied/printed/pasted.
   - Benefits: simplest byte model and aligns with one data model comment.
   - Costs: metadata labels would be pasted into target apps, which is usually
     surprising.
6. Picker option A - require `fzf` unless `--no-fzf`:
   - Shape: missing `fzf` exits `3` with install/fallback instructions.
   - Benefits: failure mode is explicit and matches the thesis edge case.
   - Costs: less graceful for first-time users.
7. Picker option B - automatic fallback:
   - Shape: missing `fzf` automatically uses numbered-list picker.
   - Benefits: more robust out of the box.
   - Costs: behavior differs from the thesis edge case and can surprise users
     expecting `fzf`.
8. Accepted release option - no Homebrew in this release:
   - Shape: local build/install and tests are required; release archive config
     and Homebrew tap publishing are out of scope.
   - Benefits: removes release credential/tap coupling from the initial feature.
   - Costs: users cannot install via Homebrew or download release archives until
     a later release feature.

## RECOMMENDED STRATEGY

1. Treat this as a spec-driven greenfield implementation of the first production
   CLI surface for `kp`, because no runtime files currently exist.
2. Use the bare command surface and do not implement a `prompt` subcommand.
   Canonical examples should use `kp clarify`, `kp list`, `kp new`, `kp edit`,
   and `kp rm`.
3. Reserve built-in command names (`help`, `list`, `new`, `edit`, `rm`,
   `prompt`, `version`) so prompt names cannot shadow commands. User prompt
   names should also pass `^[a-z][a-z0-9-]*$`.
4. Use stripped-frontmatter semantics: YAML frontmatter is metadata for `Label`;
   copied, printed, and pasted text is the body after frontmatter.
5. Keep the implementation small: Cobra at the edge, `internal/prompt` for
   registry and parsing, `internal/config` for paths, and `internal/clipboard`
   for Darwin integration.
6. Make macOS support explicit with Darwin build tags for clipboard/paste code
   and a clear non-Darwin unsupported message or compile exclusion for v1.
7. Keep `fzf` as the only TUI path and implement `--no-fzf` as an explicit
   fallback. If `fzf` is missing and `--no-fzf` is not passed, exit `3` with
   install/fallback instructions.
8. Put feature-specific tests and performance checks in future `PLAN.md` and
   `TASKS.md`; promote only stable cross-feature testing conventions to
   `docs/references/testing.md` later.
9. Exclude Homebrew tap publishing, Goreleaser archives, and GitHub release
   workflows from this feature.
10. Do not ship `prompts/plan.md` in this feature; introduce it only through a
    later feature or explicit ad hoc change.
11. Update `docs/PROJECT_PROGRESS_SUMMARY.md` when `SPEC.md` is created so the
   highest completed artifact remains accurate.

## NEXT STEP

The brainstorm phase is sufficiently resolved for the next workflow step:
`kit spec v0-init-utility`. Current understanding is at least 95%, with no
unresolved questions blocking SPEC.
