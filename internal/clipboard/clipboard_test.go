package clipboard

import (
	"errors"
	"testing"
	"time"
)

func TestVerify_Match(t *testing.T) {
	fake := &fakeCommands{readValues: []string{"expected"}}
	c := newClipboard(fake)
	c.sleep = func(time.Duration) {}

	if err := c.Verify("expected", DefaultVerifyTimeout); err != nil {
		t.Fatal(err)
	}
	if fake.readCalls != 1 {
		t.Fatalf("readCalls = %d, want 1", fake.readCalls)
	}
}

func TestVerify_TimesOutOnMismatch(t *testing.T) {
	fake := &fakeCommands{readValues: []string{"a", "a", "a", "a", "a"}}
	c := newClipboard(fake)
	sleepCount := 0
	c.sleep = func(time.Duration) { sleepCount++ }

	err := c.Verify("b", DefaultVerifyTimeout)
	if !errors.Is(err, ErrVerifyFailed) {
		t.Fatalf("Verify error = %v, want ErrVerifyFailed", err)
	}
	if fake.readCalls != 5 {
		t.Fatalf("readCalls = %d, want 5", fake.readCalls)
	}
	if sleepCount != 5 {
		t.Fatalf("sleepCount = %d, want 5", sleepCount)
	}
}

func TestCopyThenVerify_ReturnsMismatch(t *testing.T) {
	fake := &fakeCommands{readValues: []string{"wrong", "wrong", "wrong", "wrong", "wrong"}}
	c := newClipboard(fake)
	c.sleep = func(time.Duration) {}

	if err := c.Copy("expected"); err != nil {
		t.Fatal(err)
	}
	err := c.Verify("expected", DefaultVerifyTimeout)
	if !errors.Is(err, ErrVerifyFailed) {
		t.Fatalf("Verify error = %v, want ErrVerifyFailed", err)
	}
	if !fake.copied {
		t.Fatal("Copy was not called")
	}
}

func TestCopyThenVerify_SucceedsAfterVerifiedCopy(t *testing.T) {
	fake := &fakeCommands{readValues: []string{"expected"}}
	c := newClipboard(fake)
	c.sleep = func(time.Duration) {}

	if err := c.Copy("expected"); err != nil {
		t.Fatal(err)
	}
	if err := c.Verify("expected", DefaultVerifyTimeout); err != nil {
		t.Fatal(err)
	}
	if !fake.copied {
		t.Fatal("Copy was not called")
	}
}

type fakeCommands struct {
	readValues []string
	readCalls  int
	copied     bool
}

func (f *fakeCommands) Copy(string) error {
	f.copied = true
	return nil
}

func (f *fakeCommands) Read() (string, error) {
	f.readCalls++
	if len(f.readValues) == 0 {
		return "", nil
	}
	if f.readCalls > len(f.readValues) {
		return f.readValues[len(f.readValues)-1], nil
	}
	return f.readValues[f.readCalls-1], nil
}
