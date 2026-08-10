package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

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
