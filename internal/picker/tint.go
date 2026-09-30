package picker

import "strings"

// tintLines colors Markdown-like structure in wrapped preview lines without
// changing their text or width: headings, list markers, inline `code`, and
// fenced code blocks. It is a reading aid, not a Markdown renderer.
func tintLines(lines []string, st styles) []string {
	out := make([]string, len(lines))
	inFence, inCode := false, false
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		switch {
		case strings.HasPrefix(trimmed, "```"):
			inFence = !inFence
			inCode = false
			out[i] = st.faint.Render(line)
		case inFence:
			out[i] = st.mdCode.Render(line)
		case strings.HasPrefix(trimmed, "#"):
			inCode = false
			out[i] = st.mdHeading.Render(line)
		case trimmed == "":
			inCode = false
			out[i] = line
		default:
			marker := listMarker(line)
			out[i] = st.mdMarker.Render(line[:marker]) + tintCode(line[marker:], &inCode, st)
		}
	}
	return out
}

// listMarker returns the byte length of leading indentation plus a list
// marker such as "- " or "12. ", or 0 when the line is not a list item.
func listMarker(line string) int {
	indent := hangingIndent(line)
	if indent == len(line)-len(strings.TrimLeft(line, " ")) {
		return 0
	}
	return indent
}

// tintCode colors text between backticks. inCode carries an open span across
// wrapped lines of the same paragraph.
func tintCode(text string, inCode *bool, st styles) string {
	parts := strings.Split(text, "`")
	var b strings.Builder
	for i, part := range parts {
		if i > 0 {
			*inCode = !*inCode
			b.WriteString(st.faint.Render("`"))
		}
		if *inCode {
			b.WriteString(st.mdCode.Render(part))
		} else {
			b.WriteString(part)
		}
	}
	return b.String()
}
