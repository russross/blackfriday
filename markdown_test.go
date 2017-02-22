//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
// Unit tests for full document parsing and rendering
//

package blackfriday

import (
	"io/ioutil"
	"os"
	"testing"
)

func runMarkdown(input string) string {
	return string(MarkdownCommon([]byte(input)))
}

func doTests(t *testing.T, tests []string) {
	// catch and report panics
	var candidate string
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("\npanic while processing [%#v]: %s\n", candidate, err)
		}
	}()

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runMarkdown(candidate)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdown(candidate)
				}
			}
		}
	}
}

func TestDocument(t *testing.T) {
	var tests = []string{
		// Empty document.
		"",
		"",

		" ",
		"",

		// This shouldn't panic.
		// https://github.com/russross/blackfriday/issues/172
		"[]:<",
		"<p>[]:&lt;</p>\n",

		// This shouldn't panic.
		// https://github.com/russross/blackfriday/issues/173
		"   [",
		"<p>[</p>\n",
	}
	doTests(t, tests)
}

func TestJoinLines(t *testing.T) {
	result := `<h1>标题</h1>

<p>第一行文字。</p>

<p>第二行文字。</p>
`

	file, err := os.Open("testdata/zhJoinLines.text")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	input, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	opt := Options{Extensions: commonExtensions | EXTENSION_JOIN_LINES}
	renderer := HtmlRenderer(commonHtmlFlags, "", "")
	output := MarkdownOptions(input, renderer, opt)

	if string(output) != result {
		t.Error("output dose not match.")
	}
}
