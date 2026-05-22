package clipboard

import "time"

const DefaultVerifyTimeout = 250 * time.Millisecond

type Clipboard interface {
	Copy(body string) error
	Read() (string, error)
	Verify(expected string, timeout time.Duration) error
	Paste() error
	CopyAndPaste(body string) error
}

func New() Clipboard {
	return newClipboard()
}
