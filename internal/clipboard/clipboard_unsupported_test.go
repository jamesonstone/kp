//go:build !darwin

package clipboard

import (
	"errors"
	"testing"
)

func TestNewUnsupported(t *testing.T) {
	c := New()
	if err := c.Copy("body"); !errors.Is(err, ErrUnsupported) {
		t.Fatalf("Copy error = %v, want ErrUnsupported", err)
	}
}
