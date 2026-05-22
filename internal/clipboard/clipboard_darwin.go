//go:build darwin

package clipboard

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type darwinClipboard struct{}

func newClipboard() Clipboard { return darwinClipboard{} }

func (darwinClipboard) Copy(body string) error {
	c := exec.Command("pbcopy")
	c.Stdin = strings.NewReader(body)
	return c.Run()
}

func (darwinClipboard) Read() (string, error) {
	out, err := exec.Command("pbpaste").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (d darwinClipboard) Verify(expected string, timeout time.Duration) error {
	expectedHash := fmt.Sprintf("%x", md5.Sum([]byte(expected)))
	deadline := time.Now().Add(timeout)
	for {
		got, err := d.Read()
		if err != nil {
			return err
		}
		gotHash := fmt.Sprintf("%x", md5.Sum([]byte(got)))
		if gotHash == expectedHash {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("clipboard write failed: expected %s got %s", expectedHash, gotHash)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (darwinClipboard) Paste() error {
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("paste failed; check Secure Input holder via 'ioreg -l -w 0 | grep SecureInputPID': %w", err)
	}
	return nil
}

func (d darwinClipboard) CopyAndPaste(body string) error {
	if err := d.Copy(body); err != nil {
		return err
	}
	if err := d.Verify(body, DefaultVerifyTimeout); err != nil {
		return err
	}
	if err := d.Paste(); err != nil {
		return err
	}
	return nil
}

var _ Clipboard = darwinClipboard{}
var _ = errors.New
