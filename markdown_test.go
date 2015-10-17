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
	expect := `<p>XYZYX</p><p><pre><code class="language-go">12345678</code></pre></p>`

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

func TestIncludeCodeblockInList(t *testing.T) {
	fs := virtualFS{
		"main.md": `
1. Alpha
	1. Beta <{{test.go}}
2. Gamma <{{test.go}}
	* Delta
		* Iota
		<{{test.go}}
3. Kappa
`,
		"test.go": "123\n\t456\n789",
	}

	expect := `<ol><li><p>Alpha</p><ol><li>Beta <code>go123	456789</code></li></ol></li><li><p>Gamma <code>go123	456789</code></p><ul><li>Delta<ul><li>Iota<code>go123	456789</code></li></ul></li></ul></li><li><p>Kappa</p></li></ol>`
	r := HtmlRenderer(0, "", "")
	p := newParser(fs, r, EXTENSION_INCLUDE)
	input, err := p.fs.readFile("main.md")
	if err != nil {
		t.Error(err)
	}

	first := firstPass(p, input, 0)
	second := secondPass(p, first.Bytes(), 0)
	result := strings.Replace(second.String(), "\n", "", -1)
	if result != expect {
		t.Errorf("got\n%s\nexpected\n%s\n", result, expect)
	}
}

func TestCodeblockInList(t *testing.T) {
	qqq := "```"
	fs := virtualFS{
		"main.md": `
1. Alpha
  1. Beta ` + qqq + ` go
123456789
` + qqq + `
2. Gamma ` + qqq + ` go
123456789
` + qqq + `
  * Delta
    * Iota
      ` + qqq + ` go
      123456789
      ` + qqq + `
3. Kappa
`,
		"test.go": "123456789",
	}
	expect := `<ol><li>Alpha<ol><li>Beta <code>go123456789</code></li></ol></li><li>Gamma <code>go123456789</code><ul><li>Delta</li><li>Iota<code>go123456789</code></li></ul></li><li>Kappa</li></ol>`

	r := HtmlRenderer(0, "", "")
	p := newParser(fs, r, EXTENSION_INCLUDE)
	input, err := p.fs.readFile("main.md")
	if err != nil {
		t.Error(err)
	}

	first := firstPass(p, input, 0)
	second := secondPass(p, first.Bytes(), 0)
	result := strings.Replace(second.String(), "\n", "", -1)
	if result != expect {
		t.Errorf("got\n%s\nexpected\n%s\n", result, expect)
	}
}
