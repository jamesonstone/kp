//go:build darwin

package clipboard

import "testing"

func TestCopyRead_RoundTrip(t *testing.T) {
	c := New()
	const body = "kp clipboard round trip"

	if err := c.Copy(body); err != nil {
		t.Fatal(err)
	}
	got, err := c.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got != body {
		t.Fatalf("Read = %q, want %q", got, body)
	}
}
