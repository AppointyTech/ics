package ics

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// unfold reverses RFC 5545 folding and CRLF line endings, returning the
// output in the LF form the templates render, without a trailing newline.
func unfold(s string) string {
	s = strings.ReplaceAll(s, "\r\n ", "")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.TrimSuffix(s, "\n")
}

func Test_foldLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "short", in: "SUMMARY:Hi", want: "SUMMARY:Hi"},
		{name: "exactly 75", in: strings.Repeat("a", 75), want: strings.Repeat("a", 75)},
		{name: "76", in: strings.Repeat("a", 76), want: strings.Repeat("a", 75) + "\r\n a"},
		{
			name: "continuation lines hold 74 octets after the space",
			in:   strings.Repeat("a", 75+74+1),
			want: strings.Repeat("a", 75) + "\r\n " + strings.Repeat("a", 74) + "\r\n a",
		},
		{
			name: "multi-byte character is not split",
			in:   strings.Repeat("a", 74) + "é" + "b",
			want: strings.Repeat("a", 74) + "\r\n éb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := foldLine(tt.in); got != tt.want {
				t.Errorf("foldLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerateOutputIsFoldedCRLF(t *testing.T) {
	event := &Event{
		Summary:     "Event Name",
		Description: strings.Repeat("<p>Ünïcödé paragraph, with text</p>\n", 10),
		DtStart:     time.Date(2019, time.January, 1, 9, 0, 0, 0, time.UTC),
		DtEnd:       time.Date(2019, time.January, 1, 9, 30, 0, 0, time.UTC),
	}
	gotIcs, err := Generate("com.calendar.my", event)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.HasSuffix(gotIcs, "END:VCALENDAR\r\n") {
		t.Errorf("output does not end with CRLF: %q", gotIcs[len(gotIcs)-20:])
	}
	lines := strings.Split(strings.TrimSuffix(gotIcs, "\r\n"), "\r\n")
	for i, line := range lines {
		if strings.ContainsAny(line, "\r\n") {
			t.Errorf("line %d contains a bare CR or LF: %q", i, line)
		}
		if len(line) > maxLineOctets {
			t.Errorf("line %d is %d octets: %q", i, len(line), line)
		}
		if !utf8.ValidString(line) {
			t.Errorf("line %d splits a UTF-8 character: %q", i, line)
		}
	}
	if !strings.Contains(unfold(gotIcs), "\nDESCRIPTION:"+escapeText(event.Description)+"\n") {
		t.Errorf("DESCRIPTION does not unfold to the original value")
	}
}
