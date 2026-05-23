//go:build darwin

package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func New() Clipboard {
	return newClipboard(darwinCommands{})
}

type darwinCommands struct{}

func (darwinCommands) Copy(body string) error {
	_, err := runCommand("pbcopy", nil, body)
	return err
}

func (darwinCommands) Read() (string, error) {
	return runCommand("pbpaste", nil, "")
}

func runCommand(name string, args []string, input string) (string, error) {
	cmd := exec.Command(name, args...)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message != "" {
			return "", fmt.Errorf("%s: %s: %w", name, message, err)
		}
		return "", fmt.Errorf("%s: %w", name, err)
	}

	return stdout.String(), nil
}
