package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLookupPortProcessesDoesNotStartAfterCancellation(t *testing.T) {
	tempDir := t.TempDir()
	markerPath := filepath.Join(tempDir, "started")
	lsofPath := filepath.Join(tempDir, "lsof")
	script := "#!/bin/sh\nprintf started > \"$KP_MARKER\"\n"
	if err := os.WriteFile(lsofPath, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KP_MARKER", markerPath)
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := lookupPortProcesses(ctx, 4005); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context cancellation", err)
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("lsof started after cancellation: %v", err)
	}
}

func TestStopPortProcessRejectsMissingIdentity(t *testing.T) {
	err := stopPortProcess(context.Background(), PortProcess{PID: os.Getpid()}, true)
	if err == nil || !strings.Contains(err.Error(), "identity unavailable") {
		t.Fatalf("err = %v, want unavailable identity error", err)
	}
}

func TestStopPortProcessRejectsChangedIdentity(t *testing.T) {
	ctx := context.Background()
	identity, err := readProcessIdentity(ctx, os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	process := PortProcess{PID: os.Getpid(), identity: identity + " stale"}
	err = stopPortProcess(ctx, process, true)
	if err == nil || !strings.Contains(err.Error(), "identity changed") {
		t.Fatalf("err = %v, want changed identity error", err)
	}
}

func TestWaitForPortProcessHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := waitForPortProcess(ctx, time.Minute); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context cancellation", err)
	}
}

func TestRunFindPortMapsLookupCancellation(t *testing.T) {
	app := newApp(Options{})
	app.portLookup = func(context.Context, int) ([]PortProcess, error) {
		return nil, context.Canceled
	}

	if err := app.runFindPort(context.Background(), "4005"); ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestStopPortProcessesMapsCancellation(t *testing.T) {
	app := newApp(Options{})
	app.portStopper = func(context.Context, PortProcess, bool) error {
		return context.Canceled
	}

	err := app.stopPortProcesses(context.Background(), []PortProcess{{PID: 1234}}, true)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}
