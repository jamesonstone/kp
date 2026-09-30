package cmd

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// withInitEditor simulates a terminal session whose editor applies edit to
// the draft on each call.
func withInitEditor(t *testing.T, calls *int, edit func(call int, draft string) string) func(*Options) {
	return func(opts *Options) {
		withEditor(func(_ string, _ []string, path string) error {
			*calls++
			draft, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			return os.WriteFile(path, []byte(edit(*calls, string(draft))), 0o600)
		})(opts)
		opts.StdinIsTerminal = func() bool { return true }
	}
}

// fillInitDraft writes answers under each heading of the preloaded template.
func fillInitDraft(draft string, answers map[string]string) string {
	for _, field := range initFields {
		value, ok := answers[field.title]
		if !ok {
			continue
		}
		marker := "<!-- " + field.question + " -->\n"
		draft = strings.Replace(draft, marker, marker+value+"\n", 1)
	}
	return draft
}

var fullInitAnswers = map[string]string{
	"Objective":                     "Ship the editor flow.\n\n  - keep spacing\n  - keep bullets",
	"Known business/domain context": "kp is a local CLI.",
	"Invariants":                    "Clipboard verification stays exact.",
	"Constraints":                   "No network calls.",
	"Definition of done":            "The prompt prints and copies.",
}

func TestInitEditorPrintsAndCopiesWithSpacingPreserved(t *testing.T) {
	fake := &fakeClipboard{}
	calls := 0
	stdout, _, err := executeTestCommand(t, "task", withClipboard(fake),
		withInitEditor(t, &calls, func(_ int, draft string) string {
			for _, field := range initFields {
				if !strings.Contains(draft, field.heading()) || !strings.Contains(draft, field.question) {
					t.Fatalf("draft missing %q section:\n%s", field.title, draft)
				}
			}
			return fillInitDraft(draft, fullInitAnswers)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := renderInitPrompt(initAnswers{
		objective:   fullInitAnswers["Objective"],
		context:     fullInitAnswers["Known business/domain context"],
		invariants:  fullInitAnswers["Invariants"],
		constraints: fullInitAnswers["Constraints"],
		done:        fullInitAnswers["Definition of done"],
	})
	if calls != 1 || stdout != want || fake.copied != want || fake.verified != want {
		t.Fatalf("calls=%d stdout=%q copied=%q\nwant %q", calls, stdout, fake.copied, want)
	}
	if strings.Contains(stdout, "<!--") {
		t.Fatalf("comments leaked into the prompt: %q", stdout)
	}
}

func TestInitEditorEmptyOrUntouchedFileCancels(t *testing.T) {
	for name, edit := range map[string]func(int, string) string{
		"deleted everything": func(int, string) string { return "" },
		"saved untouched":    func(_ int, draft string) string { return draft },
	} {
		t.Run(name, func(t *testing.T) {
			fake := &fakeClipboard{}
			calls := 0
			stdout, _, err := executeTestCommand(t, "task", withClipboard(fake), withInitEditor(t, &calls, edit))
			if ExitCode(err) != ExitCancel || calls != 1 {
				t.Fatalf("ExitCode = %d, calls = %d, err = %v", ExitCode(err), calls, err)
			}
			if stdout != "" || fake.copied != "" {
				t.Fatalf("stdout=%q copied=%q, want no side effects", stdout, fake.copied)
			}
		})
	}
}

func TestInitEditorReopensWithNoticeUntilComplete(t *testing.T) {
	fake := &fakeClipboard{}
	calls := 0
	partial := map[string]string{"Objective": "Do it.", "Invariants": "Stay exact."}
	_, _, err := executeTestCommand(t, "task", withClipboard(fake),
		withInitEditor(t, &calls, func(call int, draft string) string {
			if call == 1 {
				return "stray note\n" + fillInitDraft(draft, partial)
			}
			if !strings.HasPrefix(draft, initNoticePrefix) {
				t.Fatalf("reopened draft lacks a notice:\n%s", draft)
			}
			for _, want := range []string{"Known business/domain context is empty.", "Constraints is empty.", "Definition of done is empty.", "above the first heading", "Do it."} {
				if !strings.Contains(draft, want) {
					t.Fatalf("reopened draft missing %q:\n%s", want, draft)
				}
			}
			draft = strings.Replace(draft, "stray note\n", "", 1)
			return fillInitDraft(draft, map[string]string{
				"Known business/domain context": "Ctx.",
				"Constraints":                   "None.",
				"Definition of done":            "Done.",
			})
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || !strings.Contains(fake.copied, "Do it.") || !strings.Contains(fake.copied, "Done.") {
		t.Fatalf("calls=%d copied=%q", calls, fake.copied)
	}
}

func TestInitEditorAbortCancels(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t, "task", withClipboard(fake), func(opts *Options) {
		withEditor(func(string, []string, string) error {
			return exec.Command("sh", "-c", "exit 1").Run() // like vim's :cq
		})(opts)
		opts.StdinIsTerminal = func() bool { return true }
	})
	if ExitCode(err) != ExitCancel || fake.copied != "" {
		t.Fatalf("ExitCode = %d, copied = %q, err = %v", ExitCode(err), fake.copied, err)
	}
}

func TestInitEditorMissingEditorExitsConfig(t *testing.T) {
	_, _, err := executeTestCommand(t, "task", func(opts *Options) {
		opts.StdinIsTerminal = func() bool { return true }
		opts.Getenv = func(string) string { return "" }
		opts.LookPath = func(string) (string, error) { return "", errors.New("not found") }
	})
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestParseInitDraftAcceptsPlainAndMarkdownHeadings(t *testing.T) {
	content := "## objective\nA\n\n# Known business/domain context\nB\nInvariants\nC\n🚧 Constraints\nD\n✅  Definition of done\n\n  E\n\n"
	answers, problems, empty := parseInitDraft(content)
	if empty || len(problems) != 0 {
		t.Fatalf("empty=%v problems=%v", empty, problems)
	}
	if answers.objective != "A" || answers.context != "B" || answers.done != "  E" {
		t.Fatalf("answers = %+v", answers)
	}
}

func TestParseInitDraftKeepsListAndQuoteLinesAsAnswerText(t *testing.T) {
	content := "Objective\nBuild it.\n- Constraints\n  - no network\n> Invariants\n" +
		"Known business/domain context\nB\nInvariants\nC\nConstraints\nD\nDefinition of done\nE\n"
	answers, problems, _ := parseInitDraft(content)
	if len(problems) != 0 {
		t.Fatalf("problems = %v", problems)
	}
	if want := "Build it.\n- Constraints\n  - no network\n> Invariants"; answers.objective != want {
		t.Fatalf("objective = %q, want %q", answers.objective, want)
	}
	if answers.constraints != "D" || answers.invariants != "C" {
		t.Fatalf("answers = %+v", answers)
	}
}

func TestParseInitDraftKeepsFencedExamplesAsAnswerText(t *testing.T) {
	fenced := "Build it.\n```md\n# Constraints\n<!-- kept -->\n```"
	content := "Objective\n" + fenced + "\nKnown business/domain context\nB\nInvariants\nC\n" +
		"Constraints\nD\nDefinition of done\nE\n"
	answers, problems, _ := parseInitDraft(content)
	if len(problems) != 0 {
		t.Fatalf("problems = %v", problems)
	}
	if answers.objective != fenced || answers.constraints != "D" {
		t.Fatalf("objective = %q, constraints = %q", answers.objective, answers.constraints)
	}
}

func TestParseInitDraftKeepsTrailingSpaces(t *testing.T) {
	content := "Objective\nline one  \r\nline two\nKnown business/domain context\nB\nInvariants\nC\nConstraints\nD\nDefinition of done\nE\n"
	answers, _, _ := parseInitDraft(content)
	if answers.objective != "line one  \nline two" {
		t.Fatalf("objective = %q", answers.objective)
	}
}

func TestStripInitCommentsRemovesSpansAndCommentOnlyLines(t *testing.T) {
	got := stripInitComments("keep <!-- a --> this\n<!-- only -->\n<!-- multi\nline -->\nafter\n")
	if got != "keep  this\nafter\n" {
		t.Fatalf("stripInitComments = %q", got)
	}
}

func TestRenderInitDraftPointsCursorAtFirstAnswer(t *testing.T) {
	draft, line := renderInitDraft()
	lines := strings.Split(draft, "\n")
	if lines[line-2] != "<!-- "+initFields[0].question+" -->" || lines[line-1] != "" {
		t.Fatalf("cursor line %d lands on %q after %q", line, lines[line-1], lines[line-2])
	}
}

func TestResolveEditorPrefersNvimFallback(t *testing.T) {
	a := newApp(Options{
		Getenv:   func(string) string { return "" },
		LookPath: func(name string) (string, error) { return "/usr/bin/" + name, nil },
	})
	editor, err := a.resolveEditor()
	if err != nil || editor.name != "nvim" {
		t.Fatalf("editor = %+v, err = %v", editor, err)
	}
}
