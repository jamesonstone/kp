//go:build !darwin

package clipboard

func newClipboard() Clipboard {
	panic("clipboard is only supported on darwin")
}
