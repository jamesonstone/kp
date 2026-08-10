package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func lookupPortProcesses(ctx context.Context, port int) ([]PortProcess, error) {
	out, err := exec.CommandContext(ctx, "lsof",
		"-nP",
		"-iTCP:"+strconv.Itoa(port),
		"-sTCP:LISTEN",
		"-iUDP:"+strconv.Itoa(port),
	).CombinedOutput()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}
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
		detail, note := inspectPortProcess(ctx, proc.PID)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
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

func inspectPortProcess(ctx context.Context, pid int) (PortProcess, string) {
	proc := PortProcess{PID: pid}
	var notes []string

	if identity, err := readProcessIdentity(ctx, pid); err != nil {
		notes = append(notes, fmt.Sprintf("process identity lookup failed: %v", err))
	} else {
		proc.identity = identity
	}

	if out, err := exec.CommandContext(ctx, "ps", "-p", strconv.Itoa(pid), "-o", "ppid=", "-o", "args=").CombinedOutput(); err == nil || len(bytes.TrimSpace(out)) > 0 {
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

	if out, err := exec.CommandContext(ctx, "lsof", "-nP", "-a", "-p", strconv.Itoa(pid), "-d", "cwd,txt", "-Ffn").CombinedOutput(); err == nil || len(bytes.TrimSpace(out)) > 0 {
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

var errProcessNotFound = errors.New("process not found")

func stopPortProcess(ctx context.Context, process PortProcess, force bool) error {
	if force {
		return signalPortProcess(ctx, process, syscall.SIGKILL)
	}

	if err := signalPortProcess(ctx, process, syscall.SIGTERM); err != nil {
		return err
	}

	if err := waitForPortProcess(ctx, 250*time.Millisecond); err != nil {
		return err
	}
	return signalPortProcess(ctx, process, syscall.SIGKILL)
}

func readProcessIdentity(ctx context.Context, pid int) (string, error) {
	out, err := exec.CommandContext(ctx, "ps", "-p", strconv.Itoa(pid), "-o", "lstart=", "-o", "command=").CombinedOutput()
	identity := strings.TrimSpace(string(out))
	if err == nil && identity != "" {
		return identity, nil
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", ctxErr
	}
	if alive, aliveErr := processAlive(pid); aliveErr != nil {
		return "", fmt.Errorf("check process %d after identity lookup: %w", pid, aliveErr)
	} else if !alive {
		return "", errProcessNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read process %d identity: %w", pid, err)
	}
	return "", fmt.Errorf("read process %d identity: empty result", pid)
}

func portProcessIdentityMatches(ctx context.Context, process PortProcess) (bool, error) {
	if process.identity == "" {
		return false, fmt.Errorf("process %d identity unavailable; refusing to signal", process.PID)
	}

	identity, err := readProcessIdentity(ctx, process.PID)
	if errors.Is(err, errProcessNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if identity != process.identity {
		return false, fmt.Errorf("process %d identity changed; refusing to signal", process.PID)
	}
	return true, nil
}

func signalPortProcess(ctx context.Context, process PortProcess, signal syscall.Signal) error {
	matches, err := portProcessIdentityMatches(ctx, process)
	if err != nil {
		return err
	}
	if !matches {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	err = syscall.Kill(process.PID, signal)
	if errors.Is(err, syscall.ESRCH) {
		return nil
	}
	return err
}

func waitForPortProcess(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
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
