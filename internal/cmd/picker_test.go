package cmd

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestPickerFZFSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		withFZF("clarify", nil),
	)
	if err != nil {
		t.Fatal(err)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called after fzf selection")
	}
	if fake.pasted {
		t.Fatal("Paste was called after fzf selection")
	}
}

func TestPickerFZFCancel(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		withFZF("", errPickerCanceled),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestFZFRunErrorClassifiesOnlyUserOutcomesAsCancellation(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		wantCancel bool
		wantExit   int
	}{
		{name: "no match", command: "exit 1", wantCancel: true, wantExit: ExitCancel},
		{name: "interrupted", command: "exit 130", wantCancel: true, wantExit: ExitCancel},
		{name: "operational error", command: "exit 2", wantCancel: false, wantExit: ExitUser},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runErr := exec.Command("sh", "-c", tt.command).Run()
			if runErr == nil {
				t.Fatalf("%q returned nil error", tt.command)
			}
			mappedErr := fzfRunError(context.Background(), runErr)
			if got := errors.Is(mappedErr, errPickerCanceled); got != tt.wantCancel {
				t.Fatalf("cancel classification = %v, want %v", got, tt.wantCancel)
			}
			if got := ExitCode(mapPickerError(mappedErr)); got != tt.wantExit {
				t.Fatalf("exit code = %d, want %d", got, tt.wantExit)
			}
		})
	}
}

func TestFZFRunErrorMapsCallerCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := fzfRunError(ctx, errors.New("process killed")); !errors.Is(err, errPickerCanceled) {
		t.Fatalf("err = %v, want picker cancellation", err)
	}
}

func TestPickerFZFMissing(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		func(opts *Options) {
			opts.LookPath = func(string) (string, error) {
				return "", errors.New("not found")
			}
		},
	)
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if !strings.Contains(err.Error(), "brew install fzf") || !strings.Contains(err.Error(), "--no-fzf") {
		t.Fatalf("err = %v", err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerNoFZFValidSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, stderr, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin("1\n"),
		withClipboard(fake),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, "1\tclarify\tClarify before implementing\n") {
		t.Fatalf("stderr = %q", stderr)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called after numbered selection")
	}
	if fake.pasted {
		t.Fatal("Paste was called after numbered selection")
	}
}

func TestPickerNoFZFInvalidSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin("abc\n"),
		withClipboard(fake),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerNoFZFOutOfRange(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin("99\n"),
		withClipboard(fake),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerNoFZFCancel(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin(""),
		withClipboard(fake),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}
