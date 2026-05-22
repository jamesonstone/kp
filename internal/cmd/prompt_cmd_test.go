package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kp "github.com/jamesonstone/kp"
)

func testOptions() Options {
	return Options{Version: "test", Commit: "sha", BuiltinFS: kp.BuiltinPromptsFS}
}

func TestPick_PromptNotFound(t *testing.T) {
	var out, errb bytes.Buffer
	code := executeWithArgs(testOptions(), []string{"--config", t.TempDir(), "prompt", "nonexistent", "--print"}, &out, &errb)
	if code != 1 {
		t.Fatalf("got code %d stderr=%s", code, errb.String())
	}
}

func TestNew_InvalidName(t *testing.T) {
	var out, errb bytes.Buffer
	code := executeWithArgs(testOptions(), []string{"--config", t.TempDir(), "prompt", "new", "Bad Name"}, &out, &errb)
	if code != 1 {
		t.Fatalf("got code %d stderr=%s", code, errb.String())
	}
}

func TestNew_NameCollision(t *testing.T) {
	conf := t.TempDir()
	d := filepath.Join(conf, "kp", "prompts")
	if err := os.MkdirAll(d, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "instructions.md"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errb bytes.Buffer
	code := executeWithArgs(testOptions(), []string{"--config", conf, "prompt", "new", "instructions"}, &out, &errb)
	if code != 1 || !strings.Contains(errb.String(), "already exists") {
		t.Fatalf("got code %d stderr=%s", code, errb.String())
	}
}

func TestRM_BuiltinNotRemovable(t *testing.T) {
	var out, errb bytes.Buffer
	code := executeWithArgs(testOptions(), []string{"--config", t.TempDir(), "prompt", "rm", "instructions"}, &out, &errb)
	if code != 1 {
		t.Fatalf("got code %d stderr=%s", code, errb.String())
	}
}

func TestPick_NoFzfAvailable(t *testing.T) {
	t.Setenv("PATH", "")
	var out, errb bytes.Buffer
	code := executeWithArgs(testOptions(), []string{"--config", t.TempDir(), "prompt"}, &out, &errb)
	if code != 3 || !strings.Contains(errb.String(), "install fzf") {
		t.Fatalf("got code %d stderr=%s", code, errb.String())
	}
}
