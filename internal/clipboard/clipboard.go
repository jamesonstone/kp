package clipboard

import (
	"errors"
	"fmt"
	"time"
)

const (
	DefaultVerifyTimeout = 250 * time.Millisecond
	defaultPollInterval  = 50 * time.Millisecond
)

var (
	ErrVerifyFailed = errors.New("clipboard verification failed")
	ErrUnsupported  = errors.New("clipboard is unsupported on this platform")
)

type Clipboard interface {
	Copy(body string) error
	Read() (string, error)
	Verify(expected string, timeout time.Duration) error
}

type commands interface {
	Copy(body string) error
	Read() (string, error)
}

type systemClipboard struct {
	commands     commands
	pollInterval time.Duration
	sleep        func(time.Duration)
}

func newClipboard(commands commands) *systemClipboard {
	return &systemClipboard{
		commands:     commands,
		pollInterval: defaultPollInterval,
		sleep:        time.Sleep,
	}
}

func (c *systemClipboard) Copy(body string) error {
	return c.commands.Copy(body)
}

func (c *systemClipboard) Read() (string, error) {
	return c.commands.Read()
}

func (c *systemClipboard) Verify(expected string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultVerifyTimeout
	}

	retries := int(timeout / c.pollInterval)
	if retries < 1 {
		retries = 1
	}

	actual := ""
	for i := 0; i < retries; i++ {
		value, err := c.Read()
		if err != nil {
			return err
		}
		actual = value
		if actual == expected {
			return nil
		}
		c.sleep(c.pollInterval)
	}

	return fmt.Errorf("%w: expected_bytes=%d actual_bytes=%d", ErrVerifyFailed, len(expected), len(actual))
}
