package ics

import (
	"strings"
	"unicode/utf8"
)

// maxLineOctets is the RFC 5545 limit for a content line, excluding the
// line break.
const maxLineOctets = 75

// toContentLines converts rendered output, which uses LF line endings, into
// RFC 5545 content lines: every line ends with CRLF and lines longer than
// 75 octets are folded (CRLF followed by a single space).
func toContentLines(s string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	var b strings.Builder
	b.Grow(len(s) + len(s)/maxLineOctets*3 + len(lines))
	for _, line := range lines {
		b.WriteString(foldLine(line))
		b.WriteString("\r\n")
	}
	return b.String()
}

// foldLine splits a content line into chunks of at most 75 octets, joined by
// CRLF and a space. The leading space of a continuation line counts towards
// its 75 octets. A multi-byte UTF-8 character is never split.
func foldLine(line string) string {
	if len(line) <= maxLineOctets {
		return line
	}
	var b strings.Builder
	limit := maxLineOctets
	for len(line) > limit {
		cut := limit
		for cut > 0 && !utf8.RuneStart(line[cut]) {
			cut--
		}
		b.WriteString(line[:cut])
		b.WriteString("\r\n ")
		line = line[cut:]
		limit = maxLineOctets - 1
	}
	b.WriteString(line)
	return b.String()
}
