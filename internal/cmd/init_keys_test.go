package cmd

import (
	"strings"
	"testing"
)

func TestConsumeInitKeySubmitAndNewline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     initKeyAction
		consumed int
	}{
		{name: "enter", input: "\r", want: initKeySubmit, consumed: 1},
		{name: "enter crlf", input: "\r\n", want: initKeySubmit, consumed: 2},
		{name: "linefeed newline", input: "\n", want: initKeyNewline, consumed: 1},
		{name: "shift enter modify other keys", input: "\x1b[27;2;13~", want: initKeyNewline, consumed: 10},
		{name: "shift enter csi u", input: "\x1b[13;2u", want: initKeyNewline, consumed: 7},
		{name: "esc enter", input: "\x1b\r", want: initKeyNewline, consumed: 2},
		{name: "backspace", input: "\x7f", want: initKeyBackspace, consumed: 1},
		{name: "ctrl c", input: "\x03", want: initKeyCancel, consumed: 1},
		{name: "ctrl d", input: "\x04", want: initKeyEOF, consumed: 1},
		{name: "incomplete csi", input: "\x1b[27;2", want: initKeyNeedMore, consumed: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, _, consumed := consumeInitKey([]byte(tt.input))
			if action != tt.want || consumed != tt.consumed {
				t.Fatalf("consumeInitKey(%q) = %d consumed=%d, want %d consumed=%d", tt.input, action, consumed, tt.want, tt.consumed)
			}
		})
	}
}

func TestConsumeInitKeyChar(t *testing.T) {
	action, r, consumed := consumeInitKey([]byte("A"))
	if action != initKeyChar || r != 'A' || consumed != 1 {
		t.Fatalf("got action=%d r=%q consumed=%d", action, r, consumed)
	}
}

func TestApplyInitKeyBuildsMultilineField(t *testing.T) {
	a := &app{stderr: &strings.Builder{}}
	var field strings.Builder
	if _, done, err := a.applyInitKey(&field, initKeyChar, 'A'); err != nil || done {
		t.Fatalf("char err=%v done=%v", err, done)
	}
	if _, done, err := a.applyInitKey(&field, initKeyNewline, 0); err != nil || done {
		t.Fatalf("newline err=%v done=%v", err, done)
	}
	if _, done, err := a.applyInitKey(&field, initKeyChar, 'B'); err != nil || done {
		t.Fatalf("char err=%v done=%v", err, done)
	}
	got, done, err := a.applyInitKey(&field, initKeySubmit, 0)
	if err != nil || !done {
		t.Fatalf("submit err=%v done=%v", err, done)
	}
	if got != "A\nB" {
		t.Fatalf("got %q, want %q", got, "A\nB")
	}
}
