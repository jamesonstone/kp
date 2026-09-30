package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/spf13/cobra"
)

const (
	initInvestigationSentence = "Independently investigate. Do not assume my suspected implementation or root cause is correct."
	initPromptPrefix          = "› "
)

type initAnswers struct {
	objective   string
	context     string
	invariants  string
	constraints string
	done        string
}

type initField struct {
	emoji    string
	title    string
	question string
	set      func(*initAnswers, string)
}

func (f initField) heading() string {
	return f.emoji + "  " + f.title
}

var initFields = []initField{
	{
		emoji:    "🎯",
		title:    "Objective",
		question: "What outcome are you trying to achieve?",
		set:      func(a *initAnswers, value string) { a.objective = value },
	},
	{
		emoji:    "🧭",
		title:    "Known business/domain context",
		question: "What context would the agent be unable to reliably discover itself?",
		set:      func(a *initAnswers, value string) { a.context = value },
	},
	{
		emoji:    "🔒",
		title:    "Invariants",
		question: "What must remain true?",
		set:      func(a *initAnswers, value string) { a.invariants = value },
	},
	{
		emoji:    "🚧",
		title:    "Constraints",
		question: "What hard boundaries must the agent respect?",
		set:      func(a *initAnswers, value string) { a.constraints = value },
	},
	{
		emoji:    "✅",
		title:    "Definition of done",
		question: "What evidence or observable result proves the task is complete?",
		set:      func(a *initAnswers, value string) { a.done = value },
	},
}

// newInitCommand builds `kp task`, which constructs a from-scratch task
// prompt. `init` is its former name and stays as an alias.
func (a *app) newInitCommand() *cobra.Command {
	var outputOnly bool
	cmd := &cobra.Command{
		Use:     "task",
		Aliases: []string{"init"},
		Short:   "Start a task from scratch",
		Long: "Build a prompt for a new conversation, or a new task in an existing thread. " +
			"Open the task template in your editor ($KP_EDITOR, $EDITOR, nvim, or vi), " +
			"then print and copy the finished prompt. Write under each heading, save, and quit; " +
			"an empty file cancels. With piped stdin, read one line per section instead.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runInit(outputOnly)
		},
	}
	cmd.Flags().BoolVar(&outputOnly, "output-only", false, "print the blank prompt template without prompting or copying")
	return cmd
}

func (a *app) runInit(outputOnly bool) error {
	if outputOnly {
		return a.writeInitPrompt(renderInitPrompt(initAnswers{}))
	}

	collect := a.collectInitAnswers
	if a.stdinIsTerminal() {
		collect = a.editInitAnswers
	}
	answers, err := collect()
	if err != nil {
		return err
	}

	body := renderInitPrompt(answers)
	cb := a.clipboardFactory()
	if err := cb.Copy(body); err != nil {
		return NewExitError(ExitSystem, err)
	}
	if err := cb.Verify(body, clipboard.DefaultVerifyTimeout); err != nil {
		return NewExitError(ExitSystem, err)
	}
	fmt.Fprintf(a.stderr, "✅ 📋 Prompt copied to clipboard.\n")
	return a.writeInitPrompt(body)
}

func (a *app) writeInitPrompt(body string) error {
	if _, err := fmt.Fprint(a.stdout, body); err != nil {
		return NewExitError(ExitSystem, err)
	}
	return nil
}

// collectInitAnswers reads one line per section from non-terminal stdin so
// scripts and pipes keep working.
func (a *app) collectInitAnswers() (initAnswers, error) {
	var answers initAnswers
	for i, field := range initFields {
		if i > 0 {
			fmt.Fprintln(a.stderr)
		}
		fmt.Fprintln(a.stderr, field.heading())
		fmt.Fprintln(a.stderr, "    "+field.question)
		value, err := a.readInitAnswer()
		if err != nil {
			return initAnswers{}, err
		}
		field.set(&answers, value)
	}
	return answers, nil
}

func (a *app) readInitAnswer() (string, error) {
	for {
		fmt.Fprint(a.stderr, initPromptPrefix)
		value, err := a.readInitLine()
		if err != nil {
			return "", err
		}
		value = strings.TrimSpace(value)
		if value != "" {
			return value, nil
		}
		fmt.Fprintln(a.stderr, "Enter at least one sentence.")
	}
}

func (a *app) readInitLine() (string, error) {
	line, err := a.inputReader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return "", NewExitError(ExitCancel, errPickerCanceled)
		}
		return "", NewExitError(ExitUser, err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func renderInitPrompt(answers initAnswers) string {
	values := []string{
		answers.objective,
		answers.context,
		answers.invariants,
		answers.constraints,
		answers.done,
	}
	var b strings.Builder
	for i, field := range initFields {
		writeInitSection(&b, field.heading(), values[i])
	}
	b.WriteString(initInvestigationSentence)
	b.WriteByte('\n')
	return b.String()
}

func writeInitSection(b *strings.Builder, heading, body string) {
	b.WriteString(heading)
	b.WriteByte('\n')
	b.WriteByte('\n')
	if body != "" {
		b.WriteString(body)
		b.WriteByte('\n')
		b.WriteByte('\n')
	}
}
