package cmd

import (
	"unicode/utf8"
)

type initKeyAction int

const (
	initKeyNeedMore initKeyAction = iota
	initKeyIgnore
	initKeyChar
	initKeyNewline
	initKeySubmit
	initKeyBackspace
	initKeyCancel
	initKeyEOF
)

func consumeInitKey(buf []byte) (action initKeyAction, r rune, consumed int) {
	if len(buf) == 0 {
		return initKeyNeedMore, 0, 0
	}

	switch buf[0] {
	case 0x03:
		return initKeyCancel, 0, 1
	case 0x04:
		return initKeyEOF, 0, 1
	case 0x08, 0x7f:
		return initKeyBackspace, 0, 1
	case '\r':
		n := 1
		if len(buf) > 1 && buf[1] == '\n' {
			n = 2
		}
		return initKeySubmit, 0, n
	case '\n':
		return initKeyNewline, 0, 1
	case 0x1b:
		return consumeInitEscape(buf)
	}

	r, size := utf8.DecodeRune(buf)
	if r == utf8.RuneError && size == 1 && !utf8.FullRune(buf) {
		return initKeyNeedMore, 0, 0
	}
	if r < 0x20 {
		return initKeyIgnore, 0, size
	}
	return initKeyChar, r, size
}

func consumeInitEscape(buf []byte) (initKeyAction, rune, int) {
	if len(buf) == 1 {
		return initKeyNeedMore, 0, 0
	}

	if buf[1] == '\r' {
		n := 2
		if len(buf) > 2 && buf[2] == '\n' {
			n = 3
		}
		return initKeyNewline, 0, n
	}
	if buf[1] == '\n' {
		return initKeyNewline, 0, 2
	}
	if buf[1] != '[' {
		return initKeyIgnore, 0, 2
	}

	for i := 2; i < len(buf); i++ {
		c := buf[i]
		if c == '~' || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			seq := buf[:i+1]
			if isShiftEnterSeq(seq) {
				return initKeyNewline, 0, i + 1
			}
			return initKeyIgnore, 0, i + 1
		}
	}
	if len(buf) >= 32 {
		return initKeyIgnore, 0, 1
	}
	return initKeyNeedMore, 0, 0
}

func isShiftEnterSeq(seq []byte) bool {
	switch string(seq) {
	case "\x1b[27;2;13~", "\x1b[27;2;10~", "\x1b[13;2u", "\x1b[13;2~":
		return true
	default:
		return false
	}
}
