package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/picker"
)

func TestListPickerSelection(t *testing.T) {
	fake := &fakeClipboard{}
	var got []picker.Item
	stdout, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		withPicker(t, func(items []picker.Item) (string, error) {
			got = items
			return "ship", nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(itemIDs(got), ",") != "merge,review,ship" {
		t.Fatalf("list items = %v", itemIDs(got))
	}
	if got[2].Detail != "kp ship" || got[2].Preview != stdout {
		t.Fatalf("ship item = %+v", got[2])
	}
	if fake.copied != stdout || fake.pasted {
		t.Fatalf("clipboard copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestListPickerCancel(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		withPicker(t, func([]picker.Item) (string, error) { return "", errPickerCanceled }),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
	assertPickerFarewell(t, ExitMessage(err))
}

func TestExitMessageRotatesWhimsicalPickerFarewells(t *testing.T) {
	seen := make(map[string]bool, len(pickerFarewells))
	first := ""
	previous := ""
	for range pickerFarewells {
		message := ExitMessage(NewExitError(ExitCancel, errPickerCanceled))
		assertPickerFarewell(t, message)
		if message == previous {
			t.Fatalf("farewell repeated immediately: %q", message)
		}
		if first == "" {
			first = message
		}
		seen[message] = true
		previous = message
	}
	if len(seen) != len(pickerFarewells) {
		t.Fatalf("saw %d farewells, want %d", len(seen), len(pickerFarewells))
	}
	if got := ExitMessage(NewExitError(ExitCancel, errPickerCanceled)); got != first {
		t.Fatalf("wrapped farewell = %q, want %q", got, first)
	}
}

func TestExitMessagePreservesOperationalFailure(t *testing.T) {
	err := NewExitError(ExitUser, errors.New("picker exploded"))
	if got := ExitMessage(err); got != "picker exploded" {
		t.Fatalf("ExitMessage = %q, want original diagnostic", got)
	}
	if got := ExitCode(err); got != ExitUser {
		t.Fatalf("ExitCode = %d, want %d", got, ExitUser)
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
	if !strings.Contains(stderr, "1\tmerge\tContext-aware PR merge and deployment\n") {
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
	assertPickerFarewell(t, ExitMessage(err))
}

func assertPickerFarewell(t *testing.T, message string) {
	t.Helper()
	for _, farewell := range pickerFarewells {
		if message == farewell {
			if strings.Contains(strings.ToLower(message), "cancel") || strings.Contains(message, "exit status") {
				t.Fatalf("farewell exposes internal cancellation text: %q", message)
			}
			return
		}
	}
	t.Fatalf("unrecognized picker farewell: %q", message)
}
