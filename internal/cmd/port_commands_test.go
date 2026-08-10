package cmd

import (
	"strings"
	"testing"
)

func TestFindPortBarePromptsForPortAndUsesMenu(t *testing.T) {
	fake := &fakeClipboard{}
	processes := []PortProcess{
		{
			PID:            1234,
			PPID:           1,
			Command:        "/usr/local/bin/node server.js",
			ExecutablePath: "/usr/local/bin/node",
			CWD:            "/Users/test/app",
			Sockets:        []string{"TCP *:4005 (LISTEN)"},
		},
		{
			PID:            5678,
			PPID:           1,
			Command:        "/usr/local/bin/python app.py",
			ExecutablePath: "/usr/local/bin/python",
			CWD:            "/Users/test/other",
			Sockets:        []string{"UDP *:4005"},
		},
	}

	stdout, stderr, err := executeTestCommand(t,
		"find-port",
		withStdin("4005\n3\n2\n"),
		withClipboard(fake),
		withPortLookup(processes, nil),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "Processes found on port 4005:") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "PID 1234") || !strings.Contains(stdout, "PID 5678") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stderr, "Select a process") || !strings.Contains(stderr, "Select an action") {
		t.Fatalf("stderr = %q", stderr)
	}
	if fake.copied != "/usr/local/bin/python app.py" {
		t.Fatalf("clipboard copied=%q", fake.copied)
	}
}

func TestFindPortSelectingAllProcessesDoesNotPanic(t *testing.T) {
	fake := &fakeClipboard{}
	processes := []PortProcess{
		{
			PID:            1234,
			Command:        "/usr/local/bin/node server.js",
			ExecutablePath: "/usr/local/bin/node",
			CWD:            "/Users/test/app",
			Sockets:        []string{"TCP *:4005 (LISTEN)"},
		},
		{
			PID:            5678,
			Command:        "/usr/local/bin/python app.py",
			ExecutablePath: "/usr/local/bin/python",
			CWD:            "/Users/test/other",
			Sockets:        []string{"UDP *:4005"},
		},
	}

	_, _, err := executeTestCommand(t,
		"find-port",
		withStdin("4005\n1\n1\n"),
		withClipboard(fake),
		withPortLookup(processes, nil),
	)
	if err != nil {
		t.Fatal(err)
	}
	if fake.copied != "1234\n5678" {
		t.Fatalf("clipboard copied=%q", fake.copied)
	}
}

func TestFindPortCopyPidCopiesMatchingProcesses(t *testing.T) {
	fakeClipboard := &fakeClipboard{}
	processes := []PortProcess{
		{
			PID:            1234,
			Command:        "/usr/local/bin/node server.js",
			ExecutablePath: "/usr/local/bin/node",
			CWD:            "/Users/test/app",
			Sockets:        []string{"TCP *:4005 (LISTEN)"},
		},
	}

	stdout, stderr, err := executeTestCommand(t,
		"port-find",
		"4005",
		"--copy", "pid",
		withClipboard(fakeClipboard),
		withPortLookup(processes, nil),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "PID 1234") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stderr, "PID copied to clipboard") {
		t.Fatalf("stderr = %q", stderr)
	}
	if fakeClipboard.copied != "1234" {
		t.Fatalf("clipboard copied=%q", fakeClipboard.copied)
	}
}

func TestFindPortStopConfirmsBeforeStopping(t *testing.T) {
	var stoppedPID int
	var stoppedForce bool
	processes := []PortProcess{
		{
			PID:            1234,
			Command:        "/usr/local/bin/node server.js",
			ExecutablePath: "/usr/local/bin/node",
			CWD:            "/Users/test/app",
			Sockets:        []string{"TCP *:4005 (LISTEN)"},
		},
	}

	_, stderr, err := executeTestCommand(t,
		"find-port",
		"4005",
		"--stop",
		withStdin("y\n"),
		func(opts *Options) {
			opts.PortLookup = func(int) ([]PortProcess, error) {
				return processes, nil
			}
			opts.PortStop = func(pid int, force bool) error {
				stoppedPID = pid
				stoppedForce = force
				return nil
			}
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, "Stop these process(es)? [y/N]:") {
		t.Fatalf("stderr = %q", stderr)
	}
	if stoppedPID != 1234 || stoppedForce {
		t.Fatalf("stoppedPID=%d stoppedForce=%v", stoppedPID, stoppedForce)
	}
}

func TestFindPortForceStopsWithoutConfirmation(t *testing.T) {
	var stoppedPID int
	var stoppedForce bool
	processes := []PortProcess{
		{
			PID:            1234,
			Command:        "/usr/local/bin/node server.js",
			ExecutablePath: "/usr/local/bin/node",
			CWD:            "/Users/test/app",
			Sockets:        []string{"TCP *:4005 (LISTEN)"},
		},
	}

	_, stderr, err := executeTestCommand(t,
		"find-port",
		"4005",
		"--stop",
		"--force",
		func(opts *Options) {
			opts.PortLookup = func(int) ([]PortProcess, error) {
				return processes, nil
			}
			opts.PortStop = func(pid int, force bool) error {
				stoppedPID = pid
				stoppedForce = force
				return nil
			}
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stderr, "Stop these process(es)?") {
		t.Fatalf("stderr = %q", stderr)
	}
	if stoppedPID != 1234 || !stoppedForce {
		t.Fatalf("stoppedPID=%d stoppedForce=%v", stoppedPID, stoppedForce)
	}
}
