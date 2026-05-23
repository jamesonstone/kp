package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreatesApprovedFilesAndExcludesKitState(t *testing.T) {
	target := t.TempDir()

	results, err := Run(Options{Dir: target})
	if err != nil {
		t.Fatal(err)
	}

	requireAction(t, results, ".gitignore", ActionCreate)
	requireAction(t, results, "AGENTS.md", ActionCreate)
	requireAction(t, results, "docs/agents/README.md", ActionCreate)
	requireAction(t, results, "docs/references/external-systems.md", ActionCreate)

	mustExist(t, filepath.Join(target, ".env"))
	mustExist(t, filepath.Join(target, ".envrc"))
	mustExist(t, filepath.Join(target, ".coderabbit.yaml"))
	mustExist(t, filepath.Join(target, ".github", "pull_request_template.md"))
	mustExist(t, filepath.Join(target, ".github", "copilot-instructions.md"))
	mustExist(t, filepath.Join(target, "CLAUDE.md"))

	mustNotExist(t, filepath.Join(target, ".kit.yaml"))
	mustNotExist(t, filepath.Join(target, ".kit"))
	mustNotExist(t, filepath.Join(target, "docs", "CONSTITUTION.md"))
	mustNotExist(t, filepath.Join(target, "docs", "PROJECT_PROGRESS_SUMMARY.md"))
	mustNotExist(t, filepath.Join(target, "docs", "specs"))
	mustNotExist(t, filepath.Join(target, "docs", "notes"))

	env, err := os.ReadFile(filepath.Join(target, ".env"))
	if err != nil {
		t.Fatal(err)
	}
	if string(env) != "" {
		t.Fatalf(".env = %q, want empty", env)
	}
}

func TestRunDryRunDoesNotWriteFiles(t *testing.T) {
	target := t.TempDir()

	results, err := Run(Options{Dir: target, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}

	requireAction(t, results, ".gitignore", ActionWouldCreate)
	requireAction(t, results, "AGENTS.md", ActionWouldCreate)
	mustNotExist(t, filepath.Join(target, ".gitignore"))
	mustNotExist(t, filepath.Join(target, "AGENTS.md"))
}

func TestRunSkipsExistingFilesByDefault(t *testing.T) {
	target := t.TempDir()
	path := filepath.Join(target, "AGENTS.md")
	if err := os.WriteFile(path, []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := Run(Options{Dir: target})
	if err != nil {
		t.Fatal(err)
	}

	requireAction(t, results, "AGENTS.md", ActionSkip)
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "custom" {
		t.Fatalf("AGENTS.md = %q, want custom", got)
	}
}

func TestRunForceOverwritesFilesExceptGitignoreAppendOnly(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "AGENTS.md"), []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, ".gitignore"), []byte("custom.log\n.env\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	results, err := Run(Options{Dir: target, Force: true})
	if err != nil {
		t.Fatal(err)
	}

	requireAction(t, results, "AGENTS.md", ActionUpdate)
	requireAction(t, results, ".gitignore", ActionUpdate)

	agents, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agents), "# AGENTS") || strings.Contains(string(agents), "custom") {
		t.Fatalf("AGENTS.md was not overwritten correctly: %q", agents)
	}

	gitignore, err := os.ReadFile(filepath.Join(target, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gitignore)
	if !strings.Contains(got, "custom.log\n") {
		t.Fatalf(".gitignore lost existing content: %q", got)
	}
	if strings.Count(got, ".env\n") != 1 {
		t.Fatalf(".gitignore duplicated existing .env pattern: %q", got)
	}
	if !strings.Contains(got, ".envrc\n") || !strings.Contains(got, ".kit/runs/\n") {
		t.Fatalf(".gitignore missing appended scaffold patterns: %q", got)
	}
}

func TestRunCreatesMissingTargetDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "nested", "repo")

	if _, err := Run(Options{Dir: target}); err != nil {
		t.Fatal(err)
	}

	mustExist(t, filepath.Join(target, "AGENTS.md"))
}

func requireAction(t *testing.T, results []Result, path string, action Action) {
	t.Helper()
	for _, result := range results {
		if result.Path == path {
			if result.Action != action {
				t.Fatalf("action for %s = %s, want %s", path, result.Action, action)
			}
			return
		}
	}
	t.Fatalf("missing result for %s in %#v", path, results)
}

func mustExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func mustNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be absent, stat err=%v", path, err)
	}
}
