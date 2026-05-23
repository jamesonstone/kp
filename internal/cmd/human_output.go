package cmd

import (
	"io"
	"os"
)

const (
	ansiReset     = "\x1b[0m"
	ansiWhiteBold = "\x1b[1;37m"
)

var terminalWriterCheck = isTerminalWriter

type humanOutputStyle struct {
	enabled bool
}

func styleForWriter(w io.Writer) humanOutputStyle {
	return humanOutputStyle{enabled: terminalWriterCheck(w)}
}

func isTerminalWriter(w io.Writer) bool {
	fileLike, ok := w.(interface {
		Stat() (os.FileInfo, error)
	})
	if !ok {
		return false
	}

	info, err := fileLike.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func (s humanOutputStyle) title(emoji, text string) string {
	if !s.enabled {
		return text
	}
	return ansiWhiteBold + emoji + " " + text + ansiReset
}
