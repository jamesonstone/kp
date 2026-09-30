package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	initModifyOtherKeysOn  = "\x1b[>4;2m"
	initModifyOtherKeysOff = "\x1b[>4;0m"
)

func (a *app) useInitTTY() bool {
	f, ok := a.stdin.(*os.File)
	return ok && term.IsTerminal(int(f.Fd()))
}

func (a *app) readInitField() (string, error) {
	if a.useInitTTY() {
		return a.readInitFieldTTY()
	}
	return a.readInitLine()
}

func (a *app) readInitFieldTTY() (string, error) {
	f := a.stdin.(*os.File)
	fd := int(f.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return a.readInitLine()
	}
	defer func() {
		_ = term.Restore(fd, state)
		_, _ = f.WriteString(initModifyOtherKeysOff)
	}()
	if _, err := f.WriteString(initModifyOtherKeysOn); err != nil {
		return "", NewExitError(ExitSystem, err)
	}

	var field strings.Builder
	pending := make([]byte, 0, 32)
	chunk := make([]byte, 64)
	for {
		n, readErr := f.Read(chunk)
		if n > 0 {
			pending = append(pending, chunk[:n]...)
		}
		for len(pending) > 0 {
			action, r, consumed := consumeInitKey(pending)
			if action == initKeyNeedMore {
				break
			}
			pending = pending[consumed:]
			value, done, handleErr := a.applyInitKey(&field, action, r)
			if handleErr != nil {
				return "", handleErr
			}
			if done {
				return value, nil
			}
		}
		if readErr != nil {
			text := strings.TrimSpace(field.String())
			if text != "" {
				return text, nil
			}
			if errors.Is(readErr, io.EOF) {
				return "", NewExitError(ExitCancel, errPickerCanceled)
			}
			return "", NewExitError(ExitUser, readErr)
		}
	}
}

func (a *app) applyInitKey(field *strings.Builder, action initKeyAction, r rune) (string, bool, error) {
	switch action {
	case initKeyCancel:
		return "", false, NewExitError(ExitCancel, errPickerCanceled)
	case initKeyEOF:
		text := strings.TrimSpace(field.String())
		if text == "" {
			return "", false, NewExitError(ExitCancel, errPickerCanceled)
		}
		fmt.Fprintln(a.stderr)
		return text, true, nil
	case initKeySubmit:
		fmt.Fprintln(a.stderr)
		return strings.TrimSpace(field.String()), true, nil
	case initKeyNewline:
		field.WriteByte('\n')
		fmt.Fprint(a.stderr, "\n"+initPromptPrefix)
	case initKeyBackspace:
		a.backspaceInitField(field)
	case initKeyChar:
		field.WriteRune(r)
		fmt.Fprint(a.stderr, string(r))
	}
	return "", false, nil
}

func (a *app) backspaceInitField(field *strings.Builder) {
	s := field.String()
	if s == "" {
		return
	}
	runes := []rune(s)
	last := runes[len(runes)-1]
	field.Reset()
	field.WriteString(string(runes[:len(runes)-1]))
	if last != '\n' {
		fmt.Fprint(a.stderr, "\b \b")
	}
}
