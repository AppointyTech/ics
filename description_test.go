package ics

import (
	"strings"
	"testing"
	"time"
)

func Test_htmlAltDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{name: "empty", description: "", want: ""},
		{name: "plain text", description: "Just text\nmore", want: ""},
		{name: "comparison is not a tag", description: "a < b and c > d", want: ""},
		{
			name:        "fragment is wrapped in a complete document",
			description: "<p>Hi <b>there</b></p>\n<a href=\"https://example.com\">Join</a>",
			want:        `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">\n<HTML><BODY><p>Hi <b>there</b></p>\n<a href="https://example.com">Join</a></BODY></HTML>`,
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
		`X-ALT-DESC;FMTTYPE=text/html:<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2//EN">\n<HTML><BODY><p>Hi</p>\n<b>there</b></BODY></HTML>` +
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
