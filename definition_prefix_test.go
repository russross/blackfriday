package blackfriday

import (
	"strings"
	"testing"
)

func TestDefinitionListAfterParagraphLines(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"html-break", "PREFIX<br>\nTERM\n: DEFINITION\n", "<p>PREFIX<br></p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n"},
		{"hard-break", "PREFIX  \nTERM\n: DEFINITION\n", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n"},
		{"soft-break", "PREFIX\nTERM\n: DEFINITION\n", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n"},
		{"multiple-lines", "PREFIX\nMORE\nTERM\n: DEFINITION\n", "<p>PREFIX\nMORE</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n"},
		{"blank-before-description", "PREFIX\nTERM\n\n: DEFINITION\n", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd><p>DEFINITION</p></dd>\n</dl>\n"},
		{"following-paragraph", "PREFIX\nTERM\n: DEFINITION\n\nTAIL\n", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n\n<p>TAIL</p>\n"},
		{"no-final-newline", "PREFIX\nTERM\n: DEFINITION", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n"},
		{"multiple-definitions", "PREFIX\nTERM\n: DEFINITION\n: ANOTHER\n", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n<dd>ANOTHER</dd>\n</dl>\n"},
		{"multiple-terms", "PREFIX\nTERM\n: DEFINITION\n\nOTHER\n: LAST\n", "<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n<dt>OTHER</dt>\n<dd>LAST</dd>\n</dl>\n"},
		{"blockquote", "> PREFIX\n> TERM\n> : DEFINITION\n", "<blockquote>\n<p>PREFIX</p>\n\n<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n</blockquote>\n"},
		{"single-term", "TERM\n: DEFINITION\n", "<dl>\n<dt>TERM</dt>\n<dd>DEFINITION</dd>\n</dl>\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := string(Run([]byte(tc.input), WithExtensions(DefinitionLists)))
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDefinitionListAfterUnicodePrefix(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		input := "é🙂" + newline + "TERM\n: DEFINITION\n\nTAIL\n"
		got := string(Run([]byte(input), WithExtensions(DefinitionLists)))
		for _, word := range []string{"é🙂", "TERM", "DEFINITION", "TAIL"} {
			if n := strings.Count(got, word); n != 1 {
				t.Errorf("newline %q: %q occurs %d times in %q", newline, word, n, got)
			}
		}
	}
}
