package ics

import (
	"strings"
	"testing"
	"time"
)

func Test_escapeText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "plain", in: "Event Description", want: "Event Description"},
		{name: "special characters", in: `a\b;c,d`, want: `a\\b\;c\,d`},
		{name: "LF", in: "a\nb", want: `a\nb`},
		{name: "CRLF", in: "a\r\nb", want: `a\nb`},
		{name: "bare CR", in: "a\rb", want: `a\nb`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeText(tt.in); got != tt.want {
				t.Errorf("escapeText(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func Test_htmlAltDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{name: "empty", description: "", want: ""},
		{name: "plain text", description: "Just text\nmore", want: ""},
		{name: "comparison is not a tag", description: "a < b and c > d", want: ""},
		{name: "email address is not a tag", description: "Contact John <john@acme.com>\nThanks", want: ""},
		{name: "single-letter email is not a tag", description: "Mail <a@acme.com>", want: ""},
		{name: "bracketed word is not a tag", description: "Call me <tomorrow>", want: ""},
		{name: "unknown tag-like word is not a tag", description: "Hello <name>", want: ""},
		{
			name:        "fragment is wrapped and newlines become <br>",
			description: "<p>Hi <b>there</b></p>\r\n<a href=\"https://example.com\">Join</a>",
			want:        `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">\n<HTML><BODY><p>Hi <b>there</b></p><br>\n<a href="https://example.com">Join</a></BODY></HTML>`,
		},
		{
			name:        "self-closing tag is detected",
			description: "Line one<br/>Line two",
			want:        `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">\n<HTML><BODY>Line one<br/>Line two</BODY></HTML>`,
		},
		{
			name:        "special characters are escaped",
			description: `<p style="color:red;font-weight:bold">A&nbsp;B, C\D</p>`,
			want:        `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">\n<HTML><BODY><p style="color:red\;font-weight:bold">A&nbsp\;B\, C\\D</p></BODY></HTML>`,
		},
		{
			name:        "complete document is kept as is",
			description: "<html>\n<body>Hi</body></html>",
			want:        `<html>\n<body>Hi</body></html>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := htmlAltDescription(tt.description); got != tt.want {
				t.Errorf("htmlAltDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateWithHTMLDescription(t *testing.T) {
	newEvent := func() *Event {
		return &Event{
			Class:       Classification_PUBLIC,
			Summary:     "Event Name",
			Description: "<p>Hi</p>\n<b>there</b>",
			Status:      EventStatus_CONFIRMED,
			DtStart:     time.Date(2019, time.January, 1, 9, 0, 0, 0, time.UTC),
			DtEnd:       time.Date(2019, time.January, 1, 9, 30, 0, 0, time.UTC),
			UID:         "123123123",
		}
	}
	want := "\nDESCRIPTION:<p>Hi</p>\\n<b>there</b>\n" +
		`X-ALT-DESC;FMTTYPE=text/html:<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">\n<HTML><BODY><p>Hi</p><br>\n<b>there</b></BODY></HTML>` +
		"\nCLASS:PUBLIC\n"

	gotIcs, err := Generate("com.calendar.my", newEvent())
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(gotIcs, want) {
		t.Errorf("Generate() = %v, want it to contain %v", gotIcs, want)
	}

	gotIcs, err = newEvent().Generate("com.calendar.my")
	if err != nil {
		t.Fatalf("Event.Generate() error = %v", err)
	}
	if !strings.Contains(gotIcs, want) {
		t.Errorf("Event.Generate() = %v, want it to contain %v", gotIcs, want)
	}
}

func TestGeneratePlainTextWithAngleBracketsHasNoAltDescription(t *testing.T) {
	event := &Event{
		Summary:     "Event Name",
		Description: "Contact John <john@acme.com>, thanks",
		DtStart:     time.Date(2019, time.January, 1, 9, 0, 0, 0, time.UTC),
		DtEnd:       time.Date(2019, time.January, 1, 9, 30, 0, 0, time.UTC),
	}
	gotIcs, err := Generate("com.calendar.my", event)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if strings.Contains(gotIcs, "X-ALT-DESC") {
		t.Errorf("Generate() added X-ALT-DESC for plain text: %v", gotIcs)
	}
	if !strings.Contains(gotIcs, "\nDESCRIPTION:Contact John <john@acme.com>\\, thanks\n") {
		t.Errorf("Generate() DESCRIPTION not escaped as expected: %v", gotIcs)
	}
}
