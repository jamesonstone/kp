package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

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
	identity       string
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
			return a.runFindPort(cmd.Context(), port)
		},
	}
	cmd.Flags().StringVar(&a.portCopy, "copy", "", "copy a field: pid, command, path, cwd, socket, or all")
	cmd.Flags().BoolVar(&a.portStop, "stop", false, "stop matching process(es) after confirmation")
	cmd.Flags().BoolVar(&a.portForce, "force", false, "stop without confirmation")
	return cmd
}

func (a *app) runFindPort(ctx context.Context, portText string) error {
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

	processes, err := a.portLookup(ctx, port)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return NewExitError(ExitCancel, err)
		}
		return NewExitError(ExitUser, err)
	}
	if len(processes) == 0 {
		return NewExitError(ExitUser, fmt.Errorf("no processes found on port %d", port))
	}

	a.printPortProcesses(port, processes)

	if a.portCopy != "" || a.portStop || a.portForce {
		if err := a.runPortAction(ctx, processes, a.portCopy, a.portStop || a.portForce, a.portForce); err != nil {
			return err
		}
		return nil
	}

	return a.runPortMenu(ctx, processes)
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
