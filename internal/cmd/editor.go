package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type editorCommand struct {
	name string
	args []string
}

func (a *app) resolveEditor() (editorCommand, error) {
	editorText := a.getenv("KP_EDITOR")
	if strings.TrimSpace(editorText) == "" {
		editorText = a.getenv("EDITOR")
	}
	if strings.TrimSpace(editorText) == "" {
		editorText = "vi"
	}

	fields := strings.Fields(editorText)
	if len(fields) == 0 {
		fields = []string{"vi"}
	}

	name := fields[0]
	if _, err := a.lookPath(name); err != nil {
		return editorCommand{}, NewExitError(ExitConfig, fmt.Errorf("editor %q not found: %w", name, err))
	}
	return editorCommand{name: name, args: fields[1:]}, nil
}

func (a *app) runEditor(editor editorCommand, path string) error {
	if a.editorRunner != nil {
		return mapEditorError(a.editorRunner(editor.name, editor.args, path))
	}

	args := append(append([]string{}, editor.args...), path)
	cmd := exec.Command(editor.name, args...)
	cmd.Stdin = a.stdin
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr
	return mapEditorError(cmd.Run())
}

func mapEditorError(err error) error {
	if err == nil {
		return nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == ExitCancel {
		return NewExitError(ExitCancel, err)
	}
	if errors.Is(err, errPickerCanceled) {
		return NewExitError(ExitCancel, err)
	}
	return NewExitError(ExitConfig, err)
}
