package ics

import (
	"regexp"
	"strings"
)

// htmlAltDescDoctype is the doctype Outlook itself emits for X-ALT-DESC.
const htmlAltDescDoctype = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">`

var (
	htmlTagRegex      = regexp.MustCompile(`<[a-zA-Z!/][^>]*>`)
	htmlDocumentRegex = regexp.MustCompile(`(?i)<html[\s>]`)
)

func escapeNewlines(s string) string {
	return strings.Join(strings.Split(s, "\n"), `\n`)
}

// htmlAltDescription returns the value for X-ALT-DESC;FMTTYPE=text/html when
// the description contains HTML, and "" otherwise.
//
// Google Calendar renders HTML found in DESCRIPTION, but Outlook treats
// DESCRIPTION as plain text and only renders HTML from X-ALT-DESC, which must
// be a complete HTML document. A description that is already a complete
// document is used as is; a fragment is wrapped in one.
func htmlAltDescription(description string) string {
	if !htmlTagRegex.MatchString(description) {
		return ""
	}
	if htmlDocumentRegex.MatchString(description) {
		return escapeNewlines(description)
	}
	return escapeNewlines(htmlAltDescDoctype + "\n<HTML><BODY>" + description + "</BODY></HTML>")
}
