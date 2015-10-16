package mmark

import (
	"strings"
	"testing"
)

func TestNestedInclude(t *testing.T) {
	fs := virtualFS{
		"A.md":    "{{B.md}}",
		"B.md":    "{{C.md}}",
		"C.md":    "XYZYX\n\n<{{test.go}}[/START OMIT/,/END OMIT/]\n",
		"test.go": "abcdef\n// START OMIT\n12345678\n// END OMIT\n",
	}
	expect := "<p>XYZYX</p><p><pre><code class=\"language-go\">12345678</code></pre></p>"

	r := HtmlRenderer(0, "", "")
	p := newParser(fs, r, EXTENSION_INCLUDE)
	input, err := p.fs.readFile("A.md")
	if err != nil {
		t.Error(err)
	}

	first := firstPass(p, input, 0)
	second := secondPass(p, first.Bytes(), 0)
	result := strings.Replace(second.String(), "\n", "", -1)
	if result != expect {
		t.Errorf("got `%s`\nexpected `%s`", result, expect)
	}
}
