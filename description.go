package ics

import (
	"regexp"
	"strings"
)

// htmlAltDescDoctype is the doctype Outlook itself emits for X-ALT-DESC.
const htmlAltDescDoctype = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">`

var (
	// htmlTagRegex matches only known HTML tags, so plain text that happens to
	// contain angle brackets (e.g. "John <john@acme.com>") is not treated as
	// HTML. The tag name must be followed by whitespace, "/" or ">".
	htmlTagRegex      = regexp.MustCompile(`(?i)</?(?:a|abbr|b|blockquote|body|br|center|code|div|em|font|h[1-6]|head|hr|html|i|img|li|ol|p|pre|s|small|span|strike|strong|sub|sup|table|tbody|td|tfoot|th|thead|tr|u|ul)(?:\s[^>]*)?/?>`)
	htmlDocumentRegex = regexp.MustCompile(`(?i)<html[\s>]`)

	newlineNormalizer = strings.NewReplacer("\r\n", "\n", "\r", "\n")
	textEscaper       = strings.NewReplacer(`\`, `\\`, ";", `\;`, ",", `\,`, "\n", `\n`)
)

// normalizeNewlines converts CRLF and bare CR line endings to LF.
func normalizeNewlines(s string) string {
	return newlineNormalizer.Replace(s)
}

// escapeText escapes a value for an RFC 5545 TEXT property: line endings are
// normalized, then backslash, semicolon, comma and newline are escaped.
func escapeText(s string) string {
	return textEscaper.Replace(normalizeNewlines(s))
}

// htmlAltDescription returns the escaped value for
// X-ALT-DESC;FMTTYPE=text/html when the description contains HTML, and ""
// otherwise.
//
// Google Calendar renders HTML found in DESCRIPTION, but Outlook treats
// DESCRIPTION as plain text and only renders HTML from X-ALT-DESC, which must
// be a complete HTML document. A description that is already a complete
// document is used as is. A fragment is wrapped in one, with its newlines
// turned into <br> so line breaks match what Google shows from DESCRIPTION.
func htmlAltDescription(description string) string {
	if !htmlTagRegex.MatchString(description) {
		return ""
	}
	description = normalizeNewlines(description)
	if htmlDocumentRegex.MatchString(description) {
		return escapeText(description)
	}
	body := strings.ReplaceAll(description, "\n", "<br>\n")
	return escapeText(htmlAltDescDoctype + "\n<HTML><BODY>" + body + "</BODY></HTML>")
}
