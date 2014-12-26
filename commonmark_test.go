// Unit tests for commonmark parsing

package mmark

import (
	"testing"
)

func runMarkdownCommonMark(input string, extensions int) string {
	htmlFlags := 0
	htmlFlags |= HTML_USE_XHTML

	renderer := HtmlRenderer(htmlFlags, "", "")

	return string(Markdown([]byte(input), renderer, extensions))
}

func doTestsCommonMark(t *testing.T, tests []string, extensions int) {
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
		actual := runMarkdownCommonMark(candidate, extensions)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownBlock(candidate, extensions)
				}
			}
		}
	}
}

func TestPrefixHeaderCommonMark_29(t *testing.T) {
	var tests = []string{
"# hallo\n\n # hallo\n\n  # hallo\n\n   # hallo\n\n    # hallo\n",
"<h1>hallo</h1>\n\n<h1>hallo</h1>\n\n<h1>hallo</h1>\n\n<h1>hallo</h1>\n\n<pre><code># hallo\n</code></pre>\n",
	}
	doTestsCommonMark(t, tests, 0)
}
