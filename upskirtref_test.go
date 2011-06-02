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

func runReferenceMarkdown(input []byte, extensions uint32) string {
	renderer := HtmlRenderer(0)
	return string(Markdown(input, renderer, extensions))
}

// disregard dos vs. unix line endings differences
func normalizeEol(s string) string {
	return strings.Replace(s, "\r\n", "\n", -1)
}

// when re-generating reference files, keep the newlines in dos
// format to avoid unnecessary diffs
func unnormalizeEol(s string) string {
	return strings.Replace(s, "\n", "\r\n", -1)
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

		actual := runReferenceMarkdown(actualdata, 0)
		//Note: uncommenting lines below re-generates reference files. Those
		//must be inspected manually to verify no mistakes have been introduced
		//actual = unnormalizeEol(actual)
		//ioutil.WriteFile(filepath.Join("upskirtref", basename+"_upskirt_ref.html"), []byte(actual), 0666)
		actual = normalizeEol(actual)
		expected := normalizeEol(string(expecteddata))
		if actual != expected {
			t.Errorf("\n    [%#v]\nExpected[%#v]\nActual  [%#v]",
				basename+".text", expected, actual)
		}
	}
}

// benchmark with all extensions enabled
func BenchmarkFile(b *testing.B) {
	b.StopTimer()
	// using the largest file we have
	fn := filepath.Join("upskirtref", "Markdown Documentation - Syntax.text")
	data, err := ioutil.ReadFile(fn)
	b.StartTimer()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		runReferenceMarkdown(data, 0xff)
	}
}

// benchmark with no extensions enabled
func BenchmarkFileNoExt(b *testing.B) {
	b.StopTimer()
	// using the largest file we have
	fn := filepath.Join("upskirtref", "Markdown Documentation - Syntax.text")
	data, err := ioutil.ReadFile(fn)
	b.StartTimer()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		runReferenceMarkdown(data, 0)
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
