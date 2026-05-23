//go:build !darwin

package clipboard

import "time"

func New() Clipboard {
	return unsupportedClipboard{}
}

type unsupportedClipboard struct{}

func (unsupportedClipboard) Copy(string) error {
	return ErrUnsupported
}

func (unsupportedClipboard) Read() (string, error) {
	return "", ErrUnsupported
}

func (unsupportedClipboard) Verify(string, time.Duration) error {
	return ErrUnsupported
}
