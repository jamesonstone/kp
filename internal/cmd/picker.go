package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jamesonstone/kp/internal/prompt"
)

func (a *app) runPicker() error {
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}

	name := ""
	if a.noFzf {
		name, err = a.pickNumbered(reg.List())
	} else {
		name, err = a.pickFZF(reg.List())
	}
	if err != nil {
		return err
	}

	return a.runPrompt(name)
}

func (a *app) pickFZF(prompts []prompt.Prompt) (string, error) {
	if _, err := a.lookPath("fzf"); err != nil {
		return "", NewExitError(ExitConfig, errors.New("fzf not found; install fzf via 'brew install fzf' or use --no-fzf"))
	}

	if a.fzfRunner != nil {
		name, err := a.fzfRunner(prompts)
		if err != nil {
			return "", mapPickerError(err)
		}
		return name, nil
	}

	name, err := runFZF(prompts)
	if err != nil {
		return "", mapPickerError(err)
	}
	return name, nil
}

func (a *app) pickNumbered(prompts []prompt.Prompt) (string, error) {
	for i, p := range prompts {
		fmt.Fprintf(a.stderr, "%d\t%s\t%s\n", i+1, p.Name, p.Label)
	}

	line, err := a.inputReader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return "", NewExitError(ExitCancel, errors.New("picker cancelled"))
		}
		return "", NewExitError(ExitUser, err)
	}

	choiceText := strings.TrimSpace(line)
	if choiceText == "" {
		return "", NewExitError(ExitCancel, errors.New("picker cancelled"))
	}

	choice, err := strconv.Atoi(choiceText)
	if err != nil {
		return "", NewExitError(ExitUser, fmt.Errorf("invalid selection %q", choiceText))
	}
	if choice < 1 || choice > len(prompts) {
		return "", NewExitError(ExitUser, fmt.Errorf("selection %d out of range", choice))
	}

	return prompts[choice-1].Name, nil
}

var errPickerCanceled = errors.New("picker cancelled")

func mapPickerError(err error) error {
	if errors.Is(err, errPickerCanceled) {
		return NewExitError(ExitCancel, err)
	}
	return NewExitError(ExitUser, err)
}

func runFZF(prompts []prompt.Prompt) (string, error) {
	previewDir, err := os.MkdirTemp("", "kp-preview-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(previewDir)

	var input strings.Builder
	for _, p := range prompts {
		if err := os.WriteFile(filepath.Join(previewDir, p.Name), []byte(p.Body), 0o600); err != nil {
			return "", err
		}
		fmt.Fprintf(&input, "%s\t%s\t%s\n", promptIcon(p), p.Name, p.Label)
	}

	cmd := exec.Command(
		"fzf",
		"--height", "60%",
		"--reverse",
		"--cycle",
		"--prompt", "📋 kp list › ",
		"--pointer", "👉",
		"--header", "enter: copy prompt · tab/shift-tab: cycle · esc: cancel",
		"--delimiter", "\t",
		"--with-nth", "1,2,3",
		"--nth", "2,3",
		"--bind", "tab:down,btab:up",
		"--preview", "cat "+shellQuote(previewDir)+"/{2}",
	)
	cmd.Stdin = strings.NewReader(input.String())
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %v", errPickerCanceled, err)
	}

	selected := strings.TrimSpace(stdout.String())
	if selected == "" {
		return "", errPickerCanceled
	}
	_, rest, ok := strings.Cut(selected, "\t")
	if !ok {
		return "", fmt.Errorf("invalid picker selection")
	}
	name, _, _ := strings.Cut(rest, "\t")
	return name, nil
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func promptIcon(p prompt.Prompt) string {
	if p.Source == prompt.SourceUser {
		return "📝"
	}
	return "📦"
}
