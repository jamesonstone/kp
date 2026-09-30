package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/clipboard"
)

func TestRenderInitPromptPinsExactBody(t *testing.T) {
	got := renderInitPrompt(initAnswers{
		objective:   "Ship the init command.",
		context:     "kp is a local macOS prompt CLI.",
		invariants:  "Clipboard verification stays exact.",
		constraints: "Do not add network calls.",
		done:        "kp init prints and copies the prompt.",
	})

	want := "" +
		"🎯  Objective\n" +
		"\n" +
		"Ship the init command.\n" +
		"\n" +
		"🧭  Known business/domain context\n" +
		"\n" +
		"kp is a local macOS prompt CLI.\n" +
		"\n" +
		"🔒  Invariants\n" +
		"\n" +
		"Clipboard verification stays exact.\n" +
		"\n" +
		"🚧  Constraints\n" +
		"\n" +
		"Do not add network calls.\n" +
		"\n" +
		"✅  Definition of done\n" +
		"\n" +
		"kp init prints and copies the prompt.\n" +
		"\n" +
		"Independently investigate. Do not assume my suspected implementation or root cause is correct.\nWhen you finish, show evidence that the definition of done is met, flagging anything unmet or unverified, and list work the objective doesn't need as follow-ups. If the definition of done conflicts with an invariant or constraint, stop and ask me.\n"
	if got != want {
		t.Fatalf("renderInitPrompt = %q, want %q", got, want)
	}

	order := []string{
		"🎯  Objective",
		"🧭  Known business/domain context",
		"🔒  Invariants",
		"🚧  Constraints",
		"✅  Definition of done",
		initInvestigationSentence,
	}
	previous := -1
	for _, section := range order {
		index := strings.Index(got, section)
		if index < 0 {
			t.Fatalf("missing section %q", section)
		}
		if index <= previous {
			t.Fatalf("section %q out of order", section)
		}
		previous = index
	}
}

func TestRenderInitPromptPreservesMultilineAnswers(t *testing.T) {
	got := renderInitPrompt(initAnswers{
		objective:   "Ship it.\nCover the happy path.",
		context:     "Local CLI.",
		invariants:  "Exact clipboard bytes.",
		constraints: "No network.",
		done:        "Prompt prints and copies.",
	})
	if !strings.Contains(got, "Ship it.\nCover the happy path.") {
		t.Fatalf("multiline objective not preserved:\n%s", got)
	}
}

func TestRenderInitPromptPreservesUserText(t *testing.T) {
	long := strings.Repeat("Keep this exact sentence. ", 80) + "Done."
	got := renderInitPrompt(initAnswers{
		objective:   "Preserve punctuation, numbers 123, and symbols @#$%.",
		context:     long,
		invariants:  "Two sentences here. Both stay.",
		constraints: "Do not rewrite this.",
		done:        "The original text is still present.",
	})
	for _, want := range []string{
		"Preserve punctuation, numbers 123, and symbols @#$%.",
		long,
		"Two sentences here. Both stay.",
		"Do not rewrite this.",
		"The original text is still present.",
		initInvestigationSentence,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("rendered prompt missing %q", want)
		}
	}
}

func TestInitOutputOnlyPrintsBlankTemplate(t *testing.T) {
	var factoryCalled bool
	stdout, stderr, err := executeTestCommand(t, "init", "--output-only", func(opts *Options) {
		opts.ClipboardFactory = func() clipboard.Clipboard {
			factoryCalled = true
			return &fakeClipboard{}
		}
		opts.Stdin = strings.NewReader("this must not be read\n")
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "" +
		"🎯  Objective\n" +
		"\n" +
		"🧭  Known business/domain context\n" +
		"\n" +
		"🔒  Invariants\n" +
		"\n" +
		"🚧  Constraints\n" +
		"\n" +
		"✅  Definition of done\n" +
		"\n" +
		"Independently investigate. Do not assume my suspected implementation or root cause is correct.\nWhen you finish, show evidence that the definition of done is met, flagging anything unmet or unverified, and list work the objective doesn't need as follow-ups. If the definition of done conflicts with an invariant or constraint, stop and ask me.\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if factoryCalled {
		t.Fatal("clipboard factory was called")
	}
}

func TestInitCopiesStdoutToClipboard(t *testing.T) {
	fake := &fakeClipboard{}
	input := strings.Join([]string{
		"Ship the init command.",
		"kp is a local macOS prompt CLI.",
		"Clipboard verification stays exact.",
		"Do not add network calls.",
		"kp init prints and copies the prompt.",
	}, "\n") + "\n"
	stdout, stderr, err := executeTestCommand(t, "init", withClipboard(fake), withStdin(input))
	if err != nil {
		t.Fatal(err)
	}

	want := renderInitPrompt(initAnswers{
		objective:   "Ship the init command.",
		context:     "kp is a local macOS prompt CLI.",
		invariants:  "Clipboard verification stays exact.",
		constraints: "Do not add network calls.",
		done:        "kp init prints and copies the prompt.",
	})
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if fake.copied != stdout || fake.verified != stdout {
		t.Fatalf("clipboard copied=%q verified=%q stdout=%q", fake.copied, fake.verified, stdout)
	}
	previous := -1
	for _, title := range []string{
		"🎯  Objective",
		"🧭  Known business/domain context",
		"🔒  Invariants",
		"🚧  Constraints",
		"✅  Definition of done",
	} {
		index := strings.Index(stderr, title+"\n")
		if index < 0 || index <= previous {
			t.Fatalf("stderr missing or reordered %q:\n%s", title, stderr)
		}
		previous = index
	}
	if !strings.Contains(stderr, "    What outcome do you want, and why does it matter?") {
		t.Fatalf("stderr missing indented question:\n%s", stderr)
	}
	if !strings.Contains(stderr, initPromptPrefix) {
		t.Fatalf("stderr missing %q:\n%s", initPromptPrefix, stderr)
	}
	if !strings.Contains(stderr, "✅ 📋 Prompt copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
	if strings.Contains(stdout, "What outcome do you want, and why does it matter?") {
		t.Fatalf("questionnaire leaked into stdout: %q", stdout)
	}
	if strings.Contains(stderr, "Shift+Enter") {
		t.Fatalf("TTY hint leaked into piped stdin:\n%s", stderr)
	}
}

func TestInitRetriesEmptyAnswers(t *testing.T) {
	fake := &fakeClipboard{}
	input := "\n   \nShip it.\nContext stays.\nInvariants stay.\nConstraints stay.\nDone stays.\n"
	stdout, stderr, err := executeTestCommand(t, "init", withClipboard(fake), withStdin(input))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "Ship it.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if strings.Count(stderr, "Enter at least one sentence.") != 2 {
		t.Fatalf("stderr = %q", stderr)
	}
	if fake.copied != stdout {
		t.Fatalf("clipboard copied=%q stdout=%q", fake.copied, stdout)
	}
}

func TestInitCancelDoesNotTouchClipboard(t *testing.T) {
	assertInitCancelLeavesClipboardUntouched(t, "")
}

func TestInitCancelAfterPartialAnswersDoesNotTouchClipboard(t *testing.T) {
	assertInitCancelLeavesClipboardUntouched(t, "First answer.\nSecond answer.\n")
}

func TestInitClipboardFailureDoesNotPrintPrompt(t *testing.T) {
	fake := &fakeClipboard{verifyErr: errors.New("verify failed")}
	input := "Obj.\nCtx.\nInv.\nCon.\nDone.\n"
	stdout, stderr, err := executeTestCommand(t, "init", withClipboard(fake), withStdin(input))
	if ExitCode(err) != ExitSystem {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if strings.Contains(stderr, "Prompt copied to clipboard") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestNewRejectsReservedInitName(t *testing.T) {
	_, _, err := executeTestCommand(t,
		"new", "init",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if err == nil || !strings.Contains(err.Error(), "reserved prompt name") {
		t.Fatalf("err = %v, want reserved prompt name", err)
	}
}

func TestInitRejectsExtraArgs(t *testing.T) {
	_, _, err := executeTestCommand(t, "init", "extra")
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestInitStdoutWriteFailure(t *testing.T) {
	cmd := NewRoot(Options{
		Stdout: errWriter{err: errors.New("write failed")},
		Stdin:  strings.NewReader(""),
		ClipboardFactory: func() clipboard.Clipboard {
			t.Fatal("clipboard factory was called")
			return nil
		},
	})
	cmd.SetArgs([]string{"init", "--output-only"})
	err := cmd.Execute()
	if ExitCode(err) != ExitSystem {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

type errWriter struct {
	err error
}

func (w errWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func assertInitCancelLeavesClipboardUntouched(t *testing.T, input string) {
	t.Helper()
	var factoryCalled bool
	stdout, _, err := executeTestCommand(t, "init", func(opts *Options) {
		opts.ClipboardFactory = func() clipboard.Clipboard {
			factoryCalled = true
			return &fakeClipboard{}
		}
		opts.Stdin = strings.NewReader(input)
	})
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if factoryCalled {
		t.Fatal("clipboard factory was called")
	}
}
