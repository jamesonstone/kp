package cmd

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jamesonstone/kp/internal/clipboard"
)

func (a *app) runPortMenu(processes []PortProcess) error {
	if len(processes) == 1 {
		return a.runPortActionMenu(processes)
	}

	target, err := a.choosePortTarget(processes)
	if err != nil {
		return err
	}
	return a.runPortActionMenu(target)
}

func (a *app) choosePortTarget(processes []PortProcess) ([]PortProcess, error) {
	choices := make([]string, 0, len(processes)+1)
	choices = append(choices, "All matching processes")
	for _, p := range processes {
		choices = append(choices, summarizePortProcessChoice(p))
	}

	choice, err := a.chooseMenu("Select a process", choices)
	if err != nil {
		return nil, err
	}
	if choice == 1 {
		return processes, nil
	}
	return []PortProcess{processes[choice-2]}, nil
}

func (a *app) runPortActionMenu(processes []PortProcess) error {
	action, err := a.chooseMenu("Select an action", []string{
		"Copy PID",
		"Copy command",
		"Copy executable path",
		"Copy working directory",
		"Copy socket details",
		"Stop process(es)",
	})
	if err != nil {
		return err
	}

	switch action {
	case 1:
		return a.copyPortField("pid", processes)
	case 2:
		return a.copyPortField("command", processes)
	case 3:
		return a.copyPortField("path", processes)
	case 4:
		return a.copyPortField("cwd", processes)
	case 5:
		return a.copyPortField("socket", processes)
	case 6:
		return a.stopPortProcesses(processes, false)
	default:
		return NewExitError(ExitUser, fmt.Errorf("invalid action selection"))
	}
}

func (a *app) runPortAction(processes []PortProcess, copyField string, shouldStop bool, force bool) error {
	if copyField != "" {
		if err := a.copyPortField(copyField, processes); err != nil {
			return err
		}
	}
	if shouldStop {
		return a.stopPortProcesses(processes, force)
	}
	return nil
}

func (a *app) chooseMenu(title string, choices []string) (int, error) {
	if len(choices) == 0 {
		return 0, NewExitError(ExitUser, fmt.Errorf("no choices available"))
	}

	fmt.Fprintln(a.stderr, title)
	for i, choice := range choices {
		fmt.Fprintf(a.stderr, "%d\t%s\n", i+1, choice)
	}
	fmt.Fprint(a.stderr, "Selection: ")

	line, err := a.inputReader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return 0, NewExitError(ExitCancel, errors.New("picker cancelled"))
		}
		return 0, NewExitError(ExitUser, err)
	}

	choiceText := strings.TrimSpace(line)
	if choiceText == "" {
		return 0, NewExitError(ExitCancel, errors.New("picker cancelled"))
	}

	choice, err := strconv.Atoi(choiceText)
	if err != nil {
		return 0, NewExitError(ExitUser, fmt.Errorf("invalid selection %q", choiceText))
	}
	if choice < 1 || choice > len(choices) {
		return 0, NewExitError(ExitUser, fmt.Errorf("selection %d out of range", choice))
	}
	return choice, nil
}

func summarizePortProcessChoice(p PortProcess) string {
	summary := fmt.Sprintf("PID %d", p.PID)
	if p.Command != "" {
		summary += " — " + truncateForMenu(p.Command, 64)
	}
	if len(p.Sockets) > 0 {
		summary += " [" + truncateForMenu(strings.Join(p.Sockets, "; "), 48) + "]"
	}
	return summary
}

func truncateForMenu(value string, max int) string {
	if len(value) <= max {
		return value
	}
	if max <= 1 {
		return value[:max]
	}
	return value[:max-1] + "…"
}

func (a *app) copyPortField(field string, processes []PortProcess) error {
	text, label, err := portCopyContent(field, processes)
	if err != nil {
		return NewExitError(ExitUser, err)
	}

	cb := a.clipboardFactory()
	if err := cb.Copy(text); err != nil {
		return NewExitError(ExitSystem, err)
	}
	if err := cb.Verify(text, clipboard.DefaultVerifyTimeout); err != nil {
		return NewExitError(ExitSystem, err)
	}
	fmt.Fprintf(a.stderr, "✅ 📋 %s copied to clipboard.\n", label)
	return nil
}

func portCopyContent(field string, processes []PortProcess) (string, string, error) {
	var text string
	var label string

	switch field {
	case "pid":
		text, label = joinPortValues(processes, func(p PortProcess) string { return strconv.Itoa(p.PID) }), "PID"
	case "command":
		text, label = joinPortValues(processes, func(p PortProcess) string { return p.Command }), "command"
	case "path":
		text, label = joinPortValues(processes, func(p PortProcess) string { return p.ExecutablePath }), "executable path"
	case "cwd":
		text, label = joinPortValues(processes, func(p PortProcess) string { return p.CWD }), "working directory"
	case "socket":
		text, label = joinPortValues(processes, func(p PortProcess) string { return strings.Join(p.Sockets, "; ") }), "socket details"
	case "all":
		text, label = renderPortSummary(processes), "process summary"
	default:
		return "", "", fmt.Errorf("invalid copy field %q", field)
	}

	if strings.TrimSpace(text) == "" {
		return "", "", fmt.Errorf("%s unavailable for selected process(es)", label)
	}
	return text, label, nil
}

func joinPortValues(processes []PortProcess, pick func(PortProcess) string) string {
	values := make([]string, 0, len(processes))
	for _, p := range processes {
		if value := strings.TrimSpace(pick(p)); value != "" {
			values = append(values, value)
		}
	}
	return strings.Join(values, "\n")
}

func renderPortSummary(processes []PortProcess) string {
	var b strings.Builder
	for i, p := range processes {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(formatPortProcess(p))
	}
	return b.String()
}

func (a *app) stopPortProcesses(processes []PortProcess, force bool) error {
	if !force {
		if err := a.confirmPortStop(processes); err != nil {
			return err
		}
	}

	for _, p := range processes {
		if err := a.portStopper(p.PID, force); err != nil {
			return NewExitError(ExitSystem, err)
		}
	}
	fmt.Fprintf(a.stderr, "stopped %d process(es)\n", len(processes))
	return nil
}

func (a *app) confirmPortStop(processes []PortProcess) error {
	fmt.Fprintln(a.stderr, "About to stop:")
	for _, p := range processes {
		fmt.Fprintf(a.stderr, "%s\n", formatPortProcess(p))
	}
	fmt.Fprint(a.stderr, "Stop these process(es)? [y/N]: ")

	line, err := a.inputReader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return NewExitError(ExitCancel, errors.New("picker cancelled"))
		}
		return NewExitError(ExitUser, err)
	}

	answer := strings.ToLower(strings.TrimSpace(line))
	switch answer {
	case "y", "yes":
		return nil
	default:
		return NewExitError(ExitCancel, errors.New("stop cancelled"))
	}
}
