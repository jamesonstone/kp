package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolve_ConfigOverride(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "xdg"))

	cwd := t.TempDir()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}

	paths, err := Resolve(Options{ConfigDir: "local-kp"})
	if err != nil {
		t.Fatal(err)
	}

	actualCwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	wantRoot, err := filepath.Abs(filepath.Join(actualCwd, "local-kp"))
	if err != nil {
		t.Fatal(err)
	}
	if paths.ConfigRoot != wantRoot {
		t.Fatalf("ConfigRoot = %q, want %q", paths.ConfigRoot, wantRoot)
	}
	if paths.PromptsDir != filepath.Join(wantRoot, promptsDirName) {
		t.Fatalf("PromptsDir = %q", paths.PromptsDir)
	}
}

func TestResolve_XDGOverride(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)

	paths, err := Resolve(Options{})
	if err != nil {
		t.Fatal(err)
	}

	wantRoot := filepath.Join(xdg, appDirName)
	if paths.ConfigRoot != wantRoot {
		t.Fatalf("ConfigRoot = %q, want %q", paths.ConfigRoot, wantRoot)
	}
	if paths.PromptsDir != filepath.Join(wantRoot, promptsDirName) {
		t.Fatalf("PromptsDir = %q", paths.PromptsDir)
	}
}

func TestResolve_DefaultsToHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", home)

	paths, err := Resolve(Options{})
	if err != nil {
		t.Fatal(err)
	}

	wantRoot := filepath.Join(home, ".config", appDirName)
	if paths.ConfigRoot != wantRoot {
		t.Fatalf("ConfigRoot = %q, want %q", paths.ConfigRoot, wantRoot)
	}
	if paths.PromptsDir != filepath.Join(wantRoot, promptsDirName) {
		t.Fatalf("PromptsDir = %q", paths.PromptsDir)
	}
}

func TestEnsure_CreatesConfigAndPromptsDirs(t *testing.T) {
	root := filepath.Join(t.TempDir(), "kp-config")

	paths, err := Ensure(Options{ConfigDir: root})
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{paths.ConfigRoot, paths.PromptsDir} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %q: %v", path, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", path)
		}
		if got := info.Mode().Perm(); got != dirMode {
			t.Fatalf("%q mode = %v, want %v", path, got, os.FileMode(dirMode))
		}
	}
}

func TestEnsure_ReportsCreateFailure(t *testing.T) {
	base := t.TempDir()
	conflict := filepath.Join(base, "kp-config")
	if err := os.WriteFile(conflict, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := Ensure(Options{ConfigDir: conflict})
	if err == nil {
		t.Fatal("Ensure error = nil, want create failure")
	}
}
