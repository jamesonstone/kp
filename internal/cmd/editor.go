package cmd

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func openEditor(path string) error {
	editor := strings.TrimSpace(os.Getenv("KP_EDITOR"))
	if editor == "" {
		editor = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	if editor == "" {
		if _, err := exec.LookPath("vi"); err != nil {
			return configErr(errors.New("set $EDITOR or install vi"))
		}
		editor = "vi"
	}
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		return configErr(errors.New("set $EDITOR or install vi"))
	}
	parts = append(parts, path)
	c := exec.Command(parts[0], parts[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		var ex *exec.ExitError
		if errors.As(err, &ex) && ex.ExitCode() == 130 {
			return &exitError{Code: 130, Err: errors.New("editor canceled")}
		}
		return configErr(err)
	}
	return nil
}
