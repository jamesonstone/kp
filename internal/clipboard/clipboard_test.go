//go:build darwin

package clipboard

import (
	"testing"
	"time"
)

func TestCopyVerify_Match(t *testing.T) {
	c := New()
	if err := c.Copy("abc"); err != nil {
		t.Fatal(err)
	}
	if err := c.Verify("abc", DefaultVerifyTimeout); err != nil {
		t.Fatal(err)
	}
}

func TestVerify_TimesOutOnMismatch(t *testing.T) {
	c := New()
	if err := c.Copy("a"); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	if err := c.Verify("b", DefaultVerifyTimeout); err == nil {
		t.Fatal("expected error")
	}
	if elapsed := time.Since(start); elapsed > 300*time.Millisecond {
		t.Fatalf("verify took too long: %s", elapsed)
	}
}
