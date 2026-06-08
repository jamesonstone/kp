package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/spf13/cobra"
)

type PortProcess struct {
	PID            int
	PPID           int
	Command        string
	ExecutablePath string
	CWD            string
	Sockets        []string
	Notes          []string
}

func (a *app) newFindPortCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "find-port [port]",
		Aliases: []string{"port-find"},
		Short:   "Inspect processes on a port",
		Long:    "Search TCP and UDP listeners on a port, print process details, and let you copy values or stop the matching process.",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			port := ""
			if len(args) == 1 {
				port = args[0]
			}
			return a.runFindPort(port)
		},
	}
	cmd.Flags().StringVar(&a.portCopy, "copy", "", "copy a field: pid, command, path, cwd, socket, or all")
	cmd.Flags().BoolVar(&a.portStop, "stop", false, "stop matching process(es) after confirmation")
	cmd.Flags().BoolVar(&a.portForce, "force", false, "stop without confirmation")
	return cmd
}

func (a *app) runFindPort(portText string) error {
	portText = strings.TrimSpace(portText)
	if portText == "" {
		var err error
		portText, err = a.promptPort()
		if err != nil {
			return err
		}
	}

	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return NewExitError(ExitUser, fmt.Errorf("invalid port %q", portText))
	}

	processes, err := a.portLookup(port)
	if err != nil {
		return NewExitError(ExitUser, err)
	}
	if len(processes) == 0 {
		return NewExitError(ExitUser, fmt.Errorf("no processes found on port %d", port))
	}

	a.printPortProcesses(port, processes)

	if a.portCopy != "" || a.portStop || a.portForce {
		if err := a.runPortAction(processes, a.portCopy, a.portStop || a.portForce, a.portForce); err != nil {
			return err
		}
		return nil
	}

	return a.runPortMenu(processes)
}

func (a *app) promptPort() (string, error) {
	fmt.Fprint(a.stderr, "Port to inspect: ")
	line, err := a.inputReader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return "", NewExitError(ExitCancel, errors.New("picker cancelled"))
		}
		return "", NewExitError(ExitUser, err)
	}

	portText := strings.TrimSpace(line)
	if portText == "" {
		return "", NewExitError(ExitCancel, errors.New("picker cancelled"))
	}
	return portText, nil
}

func (a *app) printPortProcesses(port int, processes []PortProcess) {
	fmt.Fprintf(a.stdout, "Processes found on port %d:\n", port)
	for _, p := range processes {
		fmt.Fprintf(a.stdout, "%s\n", formatPortProcess(p))
	}
}

func formatPortProcess(p PortProcess) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  PID %d", p.PID)
	if p.PPID > 0 {
		fmt.Fprintf(&b, " (PPID %d)", p.PPID)
	}
	if p.Command != "" {
		fmt.Fprintf(&b, "\n    Command: %s", p.Command)
	}
	if p.ExecutablePath != "" {
		fmt.Fprintf(&b, "\n    Path: %s", p.ExecutablePath)
	}
	if p.CWD != "" {
		fmt.Fprintf(&b, "\n    Cwd: %s", p.CWD)
	}
	if len(p.Sockets) > 0 {
		fmt.Fprintf(&b, "\n    Socket: %s", strings.Join(p.Sockets, "; "))
	}
	for _, note := range p.Notes {
		fmt.Fprintf(&b, "\n    Note: %s", note)
	}
	return b.String()
}

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
	if choice == 0 {
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

func lookupPortProcesses(port int) ([]PortProcess, error) {
	out, err := exec.Command("lsof",
		"-nP",
		"-iTCP:"+strconv.Itoa(port),
		"-sTCP:LISTEN",
		"-iUDP:"+strconv.Itoa(port),
	).CombinedOutput()
	if err != nil && len(bytes.TrimSpace(out)) == 0 {
		return nil, fmt.Errorf("search port %d: %w", port, err)
	}

	processes := map[int]*PortProcess{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "COMMAND") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		proc := processes[pid]
		if proc == nil {
			proc = &PortProcess{PID: pid}
			processes[pid] = proc
		}

		socket := strings.Join(fields[8:], " ")
		proc.Sockets = appendUnique(proc.Sockets, socket)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read lsof output: %w", err)
	}

	results := make([]PortProcess, 0, len(processes))
	for _, proc := range processes {
		detail, note := inspectPortProcess(proc.PID)
		detail.PID = proc.PID
		detail.Sockets = proc.Sockets
		if note != "" {
			detail.Notes = append(detail.Notes, note)
		}
		results = append(results, detail)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no processes found on port %d", port)
	}

	sortPortProcesses(results)
	return results, nil
}

func inspectPortProcess(pid int) (PortProcess, string) {
	proc := PortProcess{PID: pid}
	var notes []string

	if out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "ppid=", "-o", "args=").CombinedOutput(); err == nil || len(bytes.TrimSpace(out)) > 0 {
		if line := firstNonEmptyLine(out); line != "" {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				if ppid, err := strconv.Atoi(fields[0]); err == nil {
					proc.PPID = ppid
				}
				if len(fields) > 1 {
					proc.Command = strings.Join(fields[1:], " ")
				}
			}
		}
		if err != nil && len(bytes.TrimSpace(out)) == 0 {
			notes = append(notes, fmt.Sprintf("ps lookup failed: %v", err))
		}
	} else if err != nil {
		notes = append(notes, fmt.Sprintf("ps lookup failed: %v", err))
	}

	if out, err := exec.Command("lsof", "-nP", "-a", "-p", strconv.Itoa(pid), "-d", "cwd,txt", "-Ffn").CombinedOutput(); err == nil || len(bytes.TrimSpace(out)) > 0 {
		if cwd, exe := parseLsofPaths(out); cwd != "" || exe != "" {
			proc.CWD = cwd
			proc.ExecutablePath = exe
		}
		if err != nil && len(bytes.TrimSpace(out)) == 0 {
			notes = append(notes, fmt.Sprintf("path lookup failed: %v", err))
		}
	} else if err != nil {
		notes = append(notes, fmt.Sprintf("path lookup failed: %v", err))
	}

	return proc, strings.Join(notes, "; ")
}

func parseLsofPaths(out []byte) (string, string) {
	var cwd string
	var exe string
	current := ""

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch line[0] {
		case 'f':
			current = strings.TrimSpace(line[1:])
		case 'n':
			value := strings.TrimSpace(line[1:])
			switch current {
			case "cwd":
				cwd = value
			case "txt":
				exe = value
			}
		}
	}
	return cwd, exe
}

func stopPortProcess(pid int, force bool) error {
	if force {
		return syscall.Kill(pid, syscall.SIGKILL)
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)
	alive, err := processAlive(pid)
	if err != nil {
		return err
	}
	if !alive {
		return nil
	}
	return syscall.Kill(pid, syscall.SIGKILL)
}

func processAlive(pid int) (bool, error) {
	err := syscall.Kill(pid, 0)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, syscall.EPERM):
		return true, nil
	case errors.Is(err, syscall.ESRCH):
		return false, nil
	default:
		return false, err
	}
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func firstNonEmptyLine(out []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			return line
		}
	}
	return ""
}

func sortPortProcesses(processes []PortProcess) {
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if processes[j].PID < processes[i].PID {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
}
