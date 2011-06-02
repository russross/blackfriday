//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
//

//
// Markdown 1.0.3 reference tests
//

package blackfriday

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func runReferenceMarkdown(input string) string {
	renderer := HtmlRenderer(0)
	return string(Markdown([]byte(input), renderer, 0))
}

// disregard dos vs. unix line endings differences
func normalizeEol(s string) string {
	return strings.Replace(s, "\r\n", "\n", -1)
}

func doFileTests(t *testing.T, files []string) {
	for _, basename := range files {
		fn := filepath.Join("upskirtref", basename+".text")
		actualdata, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Errorf("Couldn't open '%s', error: %v\n", fn, err)
			continue
		}
		fn = filepath.Join("upskirtref", basename+"_upskirt_ref.html")
		expecteddata, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Errorf("Couldn't open '%s', error: %v\n", fn, err)
			continue
		}

		actual := string(actualdata)
		actual = normalizeEol(string(runReferenceMarkdown(actual)))
		expected := normalizeEol(string(expecteddata))
		if actual != expected {
			t.Errorf("\n    [%#v]\nExpected[%#v]\nActual  [%#v]",
				basename+".text", expected, actual)
		}
	}
}

func TestReference(t *testing.T) {
	files := []string{
		"Amps and angle encoding",
		"Auto links",
		"Backslash escapes",
		"Blockquotes with code blocks",
		"Code Blocks",
		"Code Spans",
		"Hard-wrapped paragraphs with list-like lines",
		"Horizontal rules",
		"Inline HTML (Advanced)",
		"Inline HTML (Simple)",
		"Inline HTML comments",
		"Links, inline style",
		"Links, reference style",
		"Links, shortcut references",
		"Literal quotes in titles",
		"Markdown Documentation - Basics",
		"Markdown Documentation - Syntax",
		"Nested blockquotes",
		"Ordered and unordered lists",
		"Strong and em together",
		"Tabs",
		"Tidyness",
	}
	doFileTests(t, files)
}
