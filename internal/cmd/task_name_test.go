package cmd

import (
	"strings"
	"testing"
)

func TestNewRejectsReservedTaskNames(t *testing.T) {
	for _, name := range []string{"task", "init"} {
		_, _, err := executeTestCommand(t,
			"new", name,
			withEditor(func(string, []string, string) error { return nil }),
		)
		if ExitCode(err) != ExitUser || err == nil || !strings.Contains(err.Error(), "reserved prompt name") {
			t.Fatalf("new %s: ExitCode = %d, err = %v", name, ExitCode(err), err)
		}
	}
}

func TestInitAliasRunsTask(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "init", "--output-only")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != renderInitPrompt(initAnswers{}) {
		t.Fatalf("kp init --output-only = %q", stdout)
	}
}
