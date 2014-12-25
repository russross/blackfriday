// Unit tests for block parsing

package mmark

import (
	"testing"
)

func runMarkdownBlock(input string, extensions int) string {
	htmlFlags := 0
	htmlFlags |= HTML_USE_XHTML

	renderer := HtmlRenderer(htmlFlags, "", "")

	return string(Markdown([]byte(input), renderer, extensions))
}

func runMarkdownBlockXML(input string, extensions int) string {
	xmlFlags := 0

	extensions |= commonXmlExtensions
	extensions |= EXTENSION_UNIQUE_HEADER_IDS
	renderer := XmlRenderer(xmlFlags)

	return string(Markdown([]byte(input), renderer, extensions))
}

func doTestsBlock(t *testing.T, tests []string, extensions int) {
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
		actual := runMarkdownBlock(candidate, extensions)
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

func doTestsBlockXML(t *testing.T, tests []string, extensions int) {
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
		actual := runMarkdownBlockXML(candidate, extensions)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownBlockXML(candidate, extensions)
				}
			}
		}
	}
}

func TestPrefixHeaderNoExtensions(t *testing.T) {
	var tests = []string{
		"# Header 1\n",
		"<h1>Header 1</h1>\n",

		"## Header 2\n",
		"<h2>Header 2</h2>\n",

		"### Header 3\n",
		"<h3>Header 3</h3>\n",

		"#### Header 4\n",
		"<h4>Header 4</h4>\n",

		"##### Header 5\n",
		"<h5>Header 5</h5>\n",

		"###### Header 6\n",
		"<h6>Header 6</h6>\n",

		"####### Header 7\n",
		"<h6># Header 7</h6>\n",

		"#Header 1\n",
		"<h1>Header 1</h1>\n",

		"##Header 2\n",
		"<h2>Header 2</h2>\n",

		"###Header 3\n",
		"<h3>Header 3</h3>\n",

		"####Header 4\n",
		"<h4>Header 4</h4>\n",

		"#####Header 5\n",
		"<h5>Header 5</h5>\n",

		"######Header 6\n",
		"<h6>Header 6</h6>\n",

		"#######Header 7\n",
		"<h6>#Header 7</h6>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1>Header 1</h1>\n\n<p>Goodbye</p>\n",

		"* List\n# Header\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1>Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"* List\n#Header\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1>Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"*   List\n    * Nested list\n    # Nested header\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li><p>Nested list</p>\n\n" +
			"<h1>Nested header</h1></li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, 0)
}

func TestPrefixHeaderSpaceExtension(t *testing.T) {
	var tests = []string{
		"# Header 1\n",
		"<h1>Header 1</h1>\n",

		"## Header 2\n",
		"<h2>Header 2</h2>\n",

		"### Header 3\n",
		"<h3>Header 3</h3>\n",

		"#### Header 4\n",
		"<h4>Header 4</h4>\n",

		"##### Header 5\n",
		"<h5>Header 5</h5>\n",

		"###### Header 6\n",
		"<h6>Header 6</h6>\n",

		"####### Header 7\n",
		"<p>####### Header 7</p>\n",

		"#Header 1\n",
		"<p>#Header 1</p>\n",

		"##Header 2\n",
		"<p>##Header 2</p>\n",

		"###Header 3\n",
		"<p>###Header 3</p>\n",

		"####Header 4\n",
		"<p>####Header 4</p>\n",

		"#####Header 5\n",
		"<p>#####Header 5</p>\n",

		"######Header 6\n",
		"<p>######Header 6</p>\n",

		"#######Header 7\n",
		"<p>#######Header 7</p>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1>Header 1</h1>\n\n<p>Goodbye</p>\n",

		"* List\n# Header\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1>Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"* List\n#Header\n* List\n",
		"<ul>\n<li>List\n#Header</li>\n<li>List</li>\n</ul>\n",

		"*   List\n    * Nested list\n    # Nested header\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li><p>Nested list</p>\n\n" +
			"<h1>Nested header</h1></li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, EXTENSION_SPACE_HEADERS)
}

func TestPrefixHeaderIdExtension(t *testing.T) {
	var tests = []string{
		"# Header 1 {#someid}\n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"# Header 1 {#someid}   \n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"# Header 1         {#someid}\n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"# Header 1 {#someid\n",
		"<h1>Header 1 {#someid</h1>\n",

		"# Header 1 {#someid\n",
		"<h1>Header 1 {#someid</h1>\n",

		"# Header 1 {#someid}}\n",
		"<h1 id=\"someid\">Header 1</h1>\n\n<p>}</p>\n",

		"## Header 2 {#someid}\n",
		"<h2 id=\"someid\">Header 2</h2>\n",

		"### Header 3 {#someid}\n",
		"<h3 id=\"someid\">Header 3</h3>\n",

		"#### Header 4 {#someid}\n",
		"<h4 id=\"someid\">Header 4</h4>\n",

		"##### Header 5 {#someid}\n",
		"<h5 id=\"someid\">Header 5</h5>\n",

		"###### Header 6 {#someid}\n",
		"<h6 id=\"someid\">Header 6</h6>\n",

		"####### Header 7 {#someid}\n",
		"<h6 id=\"someid\"># Header 7</h6>\n",

		"# Header 1 # {#someid}\n",
		"<h1 id=\"someid\">Header 1</h1>\n",

		"## Header 2 ## {#someid}\n",
		"<h2 id=\"someid\">Header 2</h2>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1>Header 1</h1>\n\n<p>Goodbye</p>\n",

		"* List\n# Header {#someid}\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1 id=\"someid\">Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"* List\n#Header {#someid}\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1 id=\"someid\">Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"*   List\n    * Nested list\n    # Nested header {#someid}\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li><p>Nested list</p>\n\n" +
			"<h1 id=\"someid\">Nested header</h1></li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, EXTENSION_HEADER_IDS)
}

func TestPrefixAutoHeaderIdExtension(t *testing.T) {
	var tests = []string{
		"# Header 1\n",
		"<h1 id=\"header-1\">Header 1</h1>\n",

		"# Header 1   \n",
		"<h1 id=\"header-1\">Header 1</h1>\n",

		"## Header 2\n",
		"<h2 id=\"header-2\">Header 2</h2>\n",

		"### Header 3\n",
		"<h3 id=\"header-3\">Header 3</h3>\n",

		"#### Header 4\n",
		"<h4 id=\"header-4\">Header 4</h4>\n",

		"##### Header 5\n",
		"<h5 id=\"header-5\">Header 5</h5>\n",

		"###### Header 6\n",
		"<h6 id=\"header-6\">Header 6</h6>\n",

		"####### Header 7\n",
		"<h6 id=\"-header-7\"># Header 7</h6>\n",

		"Hello\n# Header 1\nGoodbye\n",
		"<p>Hello</p>\n\n<h1 id=\"header-1\">Header 1</h1>\n\n<p>Goodbye</p>\n",

		"* List\n# Header\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1 id=\"header\">Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"* List\n#Header\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1 id=\"header\">Header</h1></li>\n\n<li><p>List</p></li>\n</ul>\n",

		"*   List\n    * Nested list\n    # Nested header\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li><p>Nested list</p>\n\n" +
			"<h1 id=\"nested-header\">Nested header</h1></li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, EXTENSION_AUTO_HEADER_IDS)
}

func TestUnderlineHeaders(t *testing.T) {
	var tests = []string{
		"Header 1\n========\n",
		"<h1>Header 1</h1>\n",

		"Header 2\n--------\n",
		"<h2>Header 2</h2>\n",

		"A\n=\n",
		"<h1>A</h1>\n",

		"B\n-\n",
		"<h2>B</h2>\n",

		"Paragraph\nHeader\n=\n",
		"<p>Paragraph</p>\n\n<h1>Header</h1>\n",

		"Header\n===\nParagraph\n",
		"<h1>Header</h1>\n\n<p>Paragraph</p>\n",

		"Header\n===\nAnother header\n---\n",
		"<h1>Header</h1>\n\n<h2>Another header</h2>\n",

		"   Header\n======\n",
		"<h1>Header</h1>\n",

		"    Code\n========\n",
		"<pre><code>Code\n</code></pre>\n\n<p>========</p>\n",

		"Header with *inline*\n=====\n",
		"<h1>Header with <em>inline</em></h1>\n",

		"*   List\n    * Sublist\n    Not a header\n    ------\n",
		"<ul>\n<li>List\n\n<ul>\n<li>Sublist\nNot a header\n------</li>\n</ul></li>\n</ul>\n",

		"Paragraph\n\n\n\n\nHeader\n===\n",
		"<p>Paragraph</p>\n\n<h1>Header</h1>\n",

		"Trailing space \n====        \n\n",
		"<h1>Trailing space</h1>\n",

		"Trailing spaces\n====        \n\n",
		"<h1>Trailing spaces</h1>\n",

		"Double underline\n=====\n=====\n",
		"<h1>Double underline</h1>\n\n<p>=====</p>\n",
	}
	doTestsBlock(t, tests, 0)
}

func TestUnderlineHeadersAutoIDs(t *testing.T) {
	var tests = []string{
		"Header 1\n========\n",
		"<h1 id=\"header-1\">Header 1</h1>\n",

		"Header 2\n--------\n",
		"<h2 id=\"header-2\">Header 2</h2>\n",

		"A\n=\n",
		"<h1 id=\"a\">A</h1>\n",

		"B\n-\n",
		"<h2 id=\"b\">B</h2>\n",

		"Paragraph\nHeader\n=\n",
		"<p>Paragraph</p>\n\n<h1 id=\"header\">Header</h1>\n",

		"Header\n===\nParagraph\n",
		"<h1 id=\"header\">Header</h1>\n\n<p>Paragraph</p>\n",

		"Header\n===\nAnother header\n---\n",
		"<h1 id=\"header\">Header</h1>\n\n<h2 id=\"another-header\">Another header</h2>\n",

		"   Header\n======\n",
		"<h1 id=\"header\">Header</h1>\n",

		"Header with *inline*\n=====\n",
		"<h1 id=\"header-with-inline\">Header with <em>inline</em></h1>\n",

		"Paragraph\n\n\n\n\nHeader\n===\n",
		"<p>Paragraph</p>\n\n<h1 id=\"header\">Header</h1>\n",

		"Trailing space \n====        \n\n",
		"<h1 id=\"trailing-space\">Trailing space</h1>\n",

		"Trailing spaces\n====        \n\n",
		"<h1 id=\"trailing-spaces\">Trailing spaces</h1>\n",

		"Double underline\n=====\n=====\n",
		"<h1 id=\"double-underline\">Double underline</h1>\n\n<p>=====</p>\n",
	}
	doTestsBlock(t, tests, EXTENSION_AUTO_HEADER_IDS)
}

func TestHorizontalRule(t *testing.T) {
	var tests = []string{
		"-\n",
		"<p>-</p>\n",

		"--\n",
		"<p>--</p>\n",

		"---\n",
		"<hr />\n",

		"----\n",
		"<hr />\n",

		"*\n",
		"<p>*</p>\n",

		"**\n",
		"<p>**</p>\n",

		"***\n",
		"<hr />\n",

		"****\n",
		"<hr />\n",

		"_\n",
		"<p>_</p>\n",

		"__\n",
		"<p>__</p>\n",

		"___\n",
		"<hr />\n",

		"____\n",
		"<hr />\n",

		"-*-\n",
		"<p>-*-</p>\n",

		"- - -\n",
		"<hr />\n",

		"* * *\n",
		"<hr />\n",

		"_ _ _\n",
		"<hr />\n",

		"-----*\n",
		"<p>-----*</p>\n",

		"   ------   \n",
		"<hr />\n",

		"Hello\n***\n",
		"<p>Hello</p>\n\n<hr />\n",

		"---\n***\n___\n",
		"<hr />\n\n<hr />\n\n<hr />\n",
	}
	doTestsBlock(t, tests, 0)
}

func TestUnorderedList(t *testing.T) {
	var tests = []string{
		"* Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"* Yin\n* Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"* Ting\n* Bong\n* Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"* Yin\n\n* Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ul>\n",

		"* Ting\n\n* Bong\n* Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ul>\n",

		"+ Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"+ Yin\n+ Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"+ Ting\n+ Bong\n+ Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"+ Yin\n\n+ Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ul>\n",

		"+ Ting\n\n+ Bong\n+ Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ul>\n",

		"- Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"- Yin\n- Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"- Ting\n- Bong\n- Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"- Yin\n\n- Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ul>\n",

		"- Ting\n\n- Bong\n- Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ul>\n",

		"*Hello\n",
		"<p>*Hello</p>\n",

		"*   Hello \n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"*   Hello \n    Next line \n",
		"<ul>\n<li>Hello\nNext line</li>\n</ul>\n",

		"Paragraph\n* No linebreak\n",
		"<p>Paragraph\n* No linebreak</p>\n",

		"Paragraph\n\n* Linebreak\n",
		"<p>Paragraph</p>\n\n<ul>\n<li>Linebreak</li>\n</ul>\n",

		"*   List\n    * Nested list\n",
		"<ul>\n<li>List\n\n<ul>\n<li>Nested list</li>\n</ul></li>\n</ul>\n",

		"*   List\n\n    * Nested list\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li>Nested list</li>\n</ul></li>\n</ul>\n",

		"*   List\n    Second line\n\n    + Nested\n",
		"<ul>\n<li><p>List\nSecond line</p>\n\n<ul>\n<li>Nested</li>\n</ul></li>\n</ul>\n",

		"*   List\n    + Nested\n\n    Continued\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li>Nested</li>\n</ul>\n\n<p>Continued</p></li>\n</ul>\n",

		"*   List\n   * shallow indent\n",
		"<ul>\n<li>List\n\n<ul>\n<li>shallow indent</li>\n</ul></li>\n</ul>\n",

		"* List\n" +
			" * shallow indent\n" +
			"  * part of second list\n" +
			"   * still second\n" +
			"    * almost there\n" +
			"     * third level\n",
		"<ul>\n" +
			"<li>List\n\n" +
			"<ul>\n" +
			"<li>shallow indent</li>\n" +
			"<li>part of second list</li>\n" +
			"<li>still second</li>\n" +
			"<li>almost there\n\n" +
			"<ul>\n" +
			"<li>third level</li>\n" +
			"</ul></li>\n" +
			"</ul></li>\n" +
			"</ul>\n",

		"* List\n        extra indent, same paragraph\n",
		"<ul>\n<li>List\n    extra indent, same paragraph</li>\n</ul>\n",

		"* List\n\n        code block\n",
		"<ul>\n<li><p>List</p>\n\n<pre><code>code block\n</code></pre></li>\n</ul>\n",

		"* List\n\n          code block with spaces\n",
		"<ul>\n<li><p>List</p>\n\n<pre><code>  code block with spaces\n</code></pre></li>\n</ul>\n",

		"* List\n\n    * sublist\n\n    normal text\n\n    * another sublist\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li>sublist</li>\n</ul>\n\n<p>normal text</p>\n\n<ul>\n<li>another sublist</li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, 0)
}

func TestOrderedList(t *testing.T) {
	var tests = []string{
		"1. Hello\n",
		"<ol>\n<li>Hello</li>\n</ol>\n",

		"1. Yin\n2. Yang\n",
		"<ol>\n<li>Yin</li>\n<li>Yang</li>\n</ol>\n",

		"1. Ting\n2. Bong\n3. Goo\n",
		"<ol>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ol>\n",

		"1. Yin\n\n2. Yang\n",
		"<ol>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ol>\n",

		"1. Ting\n\n2. Bong\n3. Goo\n",
		"<ol>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ol>\n",

		"1 Hello\n",
		"<p>1 Hello</p>\n",

		"1.Hello\n",
		"<p>1.Hello</p>\n",

		"1.  Hello \n",
		"<ol>\n<li>Hello</li>\n</ol>\n",

		"1.  Hello \n    Next line \n",
		"<ol>\n<li>Hello\nNext line</li>\n</ol>\n",

		"Paragraph\n1. No linebreak\n",
		"<p>Paragraph\n1. No linebreak</p>\n",

		"Paragraph\n\n1. Linebreak\n",
		"<p>Paragraph</p>\n\n<ol>\n<li>Linebreak</li>\n</ol>\n",

		"1.  List\n    1. Nested list\n",
		"<ol>\n<li>List\n\n<ol>\n<li>Nested list</li>\n</ol></li>\n</ol>\n",

		"1.  List\n\n    1. Nested list\n",
		"<ol>\n<li><p>List</p>\n\n<ol>\n<li>Nested list</li>\n</ol></li>\n</ol>\n",

		"1.  List\n    Second line\n\n    1. Nested\n",
		"<ol>\n<li><p>List\nSecond line</p>\n\n<ol>\n<li>Nested</li>\n</ol></li>\n</ol>\n",

		"1.  List\n    1. Nested\n\n    Continued\n",
		"<ol>\n<li><p>List</p>\n\n<ol>\n<li>Nested</li>\n</ol>\n\n<p>Continued</p></li>\n</ol>\n",

		"1.  List\n   1. shallow indent\n",
		"<ol>\n<li>List\n\n<ol>\n<li>shallow indent</li>\n</ol></li>\n</ol>\n",

		"1. List\n" +
			" 1. shallow indent\n" +
			"  2. part of second list\n" +
			"   3. still second\n" +
			"    4. almost there\n" +
			"     1. third level\n",
		"<ol>\n" +
			"<li>List\n\n" +
			"<ol>\n" +
			"<li>shallow indent</li>\n" +
			"<li>part of second list</li>\n" +
			"<li>still second</li>\n" +
			"<li>almost there\n\n" +
			"<ol>\n" +
			"<li>third level</li>\n" +
			"</ol></li>\n" +
			"</ol></li>\n" +
			"</ol>\n",

		"1. List\n        extra indent, same paragraph\n",
		"<ol>\n<li>List\n    extra indent, same paragraph</li>\n</ol>\n",

		"1. List\n\n        code block\n",
		"<ol>\n<li><p>List</p>\n\n<pre><code>code block\n</code></pre></li>\n</ol>\n",

		"1. List\n\n          code block with spaces\n",
		"<ol>\n<li><p>List</p>\n\n<pre><code>  code block with spaces\n</code></pre></li>\n</ol>\n",

		"1. List\n    * Mixted list\n",
		"<ol>\n<li>List\n\n<ul>\n<li>Mixted list</li>\n</ul></li>\n</ol>\n",

		"1. List\n * Mixed list\n",
		"<ol>\n<li>List\n\n<ul>\n<li>Mixed list</li>\n</ul></li>\n</ol>\n",

		"* Start with unordered\n 1. Ordered\n",
		"<ul>\n<li>Start with unordered\n\n<ol>\n<li>Ordered</li>\n</ol></li>\n</ul>\n",

		"* Start with unordered\n    1. Ordered\n",
		"<ul>\n<li>Start with unordered\n\n<ol>\n<li>Ordered</li>\n</ol></li>\n</ul>\n",

		"1. numbers\n1. are ignored\n",
		"<ol>\n<li>numbers</li>\n<li>are ignored</li>\n</ol>\n",
	}
	doTestsBlock(t, tests, 0)
}

func TestAbbreviation(t *testing.T) {
	var tests = []string{
		"*[HTML]: Hyper Text Markup Language\nHTML is cool",
		"<p><abbr title=\"Hyper Text Markup Language\">HTML</abbr> is cool</p>\n",

		"*[HTML]:\nHTML is cool",
		"<p><abbr>HTML</abbr> is cool</p>\n",

		"*[H]:\n cool    H is",
		"<p>cool    <abbr>H</abbr> is</p>\n",

		"*[H]:\n cool    H is   ",
		"<p>cool    <abbr>H</abbr> is</p>\n",

		"*[H]:   \n cool    H is   ",
		"<p>cool    <abbr>H</abbr> is</p>\n",

		"*[H]: aa  \n cool    H is better yet some more words  ",
		"<p>cool    <abbr title=\"aa\">H</abbr> is better yet some more words</p>\n",

		"*[H] aa  \n cool    H is   ",
		"<p>*[H] aa<br />\n cool    H is</p>\n",

		"*[HTML]: \"Hyper Text Markup Language\"\nHTML is cool",
		"<p><abbr title=\"\"Hyper Text Markup Language\"\">HTML</abbr> is cool</p>\n",
	}
	doTestsBlock(t, tests, EXTENSION_ABBREVIATIONS)
}

func TestPreformattedHtml(t *testing.T) {
	var tests = []string{
		"<div></div>\n",
		"<div></div>\n",

		"<div>\n</div>\n",
		"<div>\n</div>\n",

		"<div>\n</div>\nParagraph\n",
		"<p><div>\n</div>\nParagraph</p>\n",

		"<div class=\"foo\">\n</div>\n",
		"<div class=\"foo\">\n</div>\n",

		"<div>\nAnything here\n</div>\n",
		"<div>\nAnything here\n</div>\n",

		"<div>\n  Anything here\n</div>\n",
		"<div>\n  Anything here\n</div>\n",

		"<div>\nAnything here\n  </div>\n",
		"<div>\nAnything here\n  </div>\n",

		"<div>\nThis is *not* &proceessed\n</div>\n",
		"<div>\nThis is *not* &proceessed\n</div>\n",

		"<faketag>\n  Something\n</faketag>\n",
		"<p><faketag>\n  Something\n</faketag></p>\n",

		"<div>\n  Something here\n</divv>\n",
		"<p><div>\n  Something here\n</divv></p>\n",

		"Paragraph\n<div>\nHere? >&<\n</div>\n",
		"<p>Paragraph\n<div>\nHere? &gt;&amp;&lt;\n</div></p>\n",

		"Paragraph\n\n<div>\nHow about here? >&<\n</div>\n",
		"<p>Paragraph</p>\n\n<div>\nHow about here? >&<\n</div>\n",

		"Paragraph\n<div>\nHere? >&<\n</div>\nAnd here?\n",
		"<p>Paragraph\n<div>\nHere? &gt;&amp;&lt;\n</div>\nAnd here?</p>\n",

		"Paragraph\n\n<div>\nHow about here? >&<\n</div>\nAnd here?\n",
		"<p>Paragraph</p>\n\n<p><div>\nHow about here? &gt;&amp;&lt;\n</div>\nAnd here?</p>\n",

		"Paragraph\n<div>\nHere? >&<\n</div>\n\nAnd here?\n",
		"<p>Paragraph\n<div>\nHere? &gt;&amp;&lt;\n</div></p>\n\n<p>And here?</p>\n",

		"Paragraph\n\n<div>\nHow about here? >&<\n</div>\n\nAnd here?\n",
		"<p>Paragraph</p>\n\n<div>\nHow about here? >&<\n</div>\n\n<p>And here?</p>\n",
	}
	doTestsBlock(t, tests, 0)
}

func TestPreformattedHtmlLax(t *testing.T) {
	var tests = []string{
		"Paragraph\n<div>\nHere? >&<\n</div>\n",
		"<p>Paragraph</p>\n\n<div>\nHere? >&<\n</div>\n",

		"Paragraph\n\n<div>\nHow about here? >&<\n</div>\n",
		"<p>Paragraph</p>\n\n<div>\nHow about here? >&<\n</div>\n",

		"Paragraph\n<div>\nHere? >&<\n</div>\nAnd here?\n",
		"<p>Paragraph</p>\n\n<div>\nHere? >&<\n</div>\n\n<p>And here?</p>\n",

		"Paragraph\n\n<div>\nHow about here? >&<\n</div>\nAnd here?\n",
		"<p>Paragraph</p>\n\n<div>\nHow about here? >&<\n</div>\n\n<p>And here?</p>\n",

		"Paragraph\n<div>\nHere? >&<\n</div>\n\nAnd here?\n",
		"<p>Paragraph</p>\n\n<div>\nHere? >&<\n</div>\n\n<p>And here?</p>\n",

		"Paragraph\n\n<div>\nHow about here? >&<\n</div>\n\nAnd here?\n",
		"<p>Paragraph</p>\n\n<div>\nHow about here? >&<\n</div>\n\n<p>And here?</p>\n",
	}
	doTestsBlock(t, tests, EXTENSION_LAX_HTML_BLOCKS)
}

func TestFencedCodeBlock(t *testing.T) {
	var tests = []string{
		"``` go\nfunc foo() bool {\n\treturn true;\n}\n```\n",
		"<pre><code class=\"language-go\">func foo() bool {\n\treturn true;\n}\n</code></pre>\n",

		"``` c\n/* special & char < > \" escaping */\n```\n",
		"<pre><code class=\"language-c\">/* special &amp; char &lt; &gt; &quot; escaping */\n</code></pre>\n",

		"``` c\nno *inline* processing ~~of text~~\n```\n",
		"<pre><code class=\"language-c\">no *inline* processing ~~of text~~\n</code></pre>\n",

		"```\nNo language\n```\n",
		"<pre><code>No language\n</code></pre>\n",

		"``` {ocaml}\nlanguage in braces\n```\n",
		"<pre><code class=\"language-ocaml\">language in braces\n</code></pre>\n",

		"```    {ocaml}      \nwith extra whitespace\n```\n",
		"<pre><code class=\"language-ocaml\">with extra whitespace\n</code></pre>\n",

		"```{   ocaml   }\nwith extra whitespace\n```\n",
		"<pre><code class=\"language-ocaml\">with extra whitespace\n</code></pre>\n",

		"~ ~~ java\nWith whitespace\n~~~\n",
		"<p>~ ~~ java\nWith whitespace\n~~~</p>\n",

		"~~\nonly two\n~~\n",
		"<p>~~\nonly two\n~~</p>\n",

		"```` python\nextra\n````\n",
		"<pre><code class=\"language-python\">extra\n</code></pre>\n",

		"~~~ perl\nthree to start, four to end\n~~~~\n",
		"<pre><code class=\"language-perl\">three to start, four to end\n</code></pre>\n",

		"~~~~ perl\nfour to start, three to end\n~~~\n",
		"<p>~~~~ perl\nfour to start, three to end\n~~~</p>\n",

		"~~~ bash\ntildes\n~~~\n",
		"<pre><code class=\"language-bash\">tildes\n</code></pre>\n",

		"``` lisp\nno ending\n",
		"<p>``` lisp\nno ending</p>\n",

		"~~~ lisp\nend with language\n~~~ lisp\n",
		"<p>~~~ lisp\nend with language\n~~~ lisp</p>\n",

		"```\nmismatched begin and end\n~~~\n",
		"<p>```\nmismatched begin and end\n~~~</p>\n",

		"~~~\nmismatched begin and end\n```\n",
		"<p>~~~\nmismatched begin and end\n```</p>\n",

		"   ``` oz\nleading spaces\n```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		"  ``` oz\nleading spaces\n ```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		" ``` oz\nleading spaces\n  ```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		"``` oz\nleading spaces\n   ```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		"    ``` oz\nleading spaces\n    ```\n",
		"<pre><code>``` oz\n</code></pre>\n\n<p>leading spaces\n    ```</p>\n",

		"Bla bla\n\n``` oz\ncode blocks breakup paragraphs\n```\n\nBla Bla\n",
		"<p>Bla bla</p>\n\n<pre><code class=\"language-oz\">code blocks breakup paragraphs\n</code></pre>\n\n<p>Bla Bla</p>\n",

		"Some text before a fenced code block\n``` oz\ncode blocks breakup paragraphs\n```\nAnd some text after a fenced code block",
		"<p>Some text before a fenced code block</p>\n\n<pre><code class=\"language-oz\">code blocks breakup paragraphs\n</code></pre>\n\n<p>And some text after a fenced code block</p>\n",

		"`",
		"<p>`</p>\n",

		"Bla bla\n\n``` oz\ncode blocks breakup paragraphs\n```\n\nBla Bla\n\n``` oz\nmultiple code blocks work okay\n```\n\nBla Bla\n",
		"<p>Bla bla</p>\n\n<pre><code class=\"language-oz\">code blocks breakup paragraphs\n</code></pre>\n\n<p>Bla Bla</p>\n\n<pre><code class=\"language-oz\">multiple code blocks work okay\n</code></pre>\n\n<p>Bla Bla</p>\n",

		"Some text before a fenced code block\n``` oz\ncode blocks breakup paragraphs\n```\nSome text in between\n``` oz\nmultiple code blocks work okay\n```\nAnd some text after a fenced code block",
		"<p>Some text before a fenced code block</p>\n\n<pre><code class=\"language-oz\">code blocks breakup paragraphs\n</code></pre>\n\n<p>Some text in between</p>\n\n<pre><code class=\"language-oz\">multiple code blocks work okay\n</code></pre>\n\n<p>And some text after a fenced code block</p>\n",
	}
	doTestsBlock(t, tests, EXTENSION_FENCED_CODE)
}

func TestTable(t *testing.T) {
	var tests = []string{
		"a | b\n---|---\nc | d\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td>c</td>\n<td>d</td>\n</tr>\n</tbody>\n</table>\n",

		"a | b\n---|--\nc | d\n",
		"<p>a | b\n---|--\nc | d</p>\n",

		"|a|b|c|d|\n|----|----|----|---|\n|e|f|g|h|\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n<th>c</th>\n<th>d</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td>e</td>\n<td>f</td>\n<td>g</td>\n<td>h</td>\n</tr>\n</tbody>\n</table>\n",

		"*a*|__b__|[c](C)|d\n---|---|---|---\ne|f|g|h\n",
		"<table>\n<thead>\n<tr>\n<th><em>a</em></th>\n<th><strong>b</strong></th>\n<th><a href=\"C\">c</a></th>\n<th>d</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td>e</td>\n<td>f</td>\n<td>g</td>\n<td>h</td>\n</tr>\n</tbody>\n</table>\n",

		"a|b|c\n---|---|---\nd|e|f\ng|h\ni|j|k|l|m\nn|o|p\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n<th>c</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td>d</td>\n<td>e</td>\n<td>f</td>\n</tr>\n\n" +
			"<tr>\n<td>g</td>\n<td>h</td>\n<td></td>\n</tr>\n\n" +
			"<tr>\n<td>i</td>\n<td>j</td>\n<td>k</td>\n</tr>\n\n" +
			"<tr>\n<td>n</td>\n<td>o</td>\n<td>p</td>\n</tr>\n</tbody>\n</table>\n",

		"a|b|c\n---|---|---\n*d*|__e__|f\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n<th>c</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td><em>d</em></td>\n<td><strong>e</strong></td>\n<td>f</td>\n</tr>\n</tbody>\n</table>\n",

		"a|b|c|d\n:--|--:|:-:|---\ne|f|g|h\n",
		"<table>\n<thead>\n<tr>\n<th align=\"left\">a</th>\n<th align=\"right\">b</th>\n" +
			"<th align=\"center\">c</th>\n<th>d</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td align=\"left\">e</td>\n<td align=\"right\">f</td>\n" +
			"<td align=\"center\">g</td>\n<td>h</td>\n</tr>\n</tbody>\n</table>\n",

		"a|b|c\n---|---|---\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n<th>c</th>\n</tr>\n</thead>\n\n<tbody>\n</tbody>\n</table>\n",

		"a| b|c | d | e\n---|---|---|---|---\nf| g|h | i |j\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n<th>c</th>\n<th>d</th>\n<th>e</th>\n</tr>\n</thead>\n\n" +
			"<tbody>\n<tr>\n<td>f</td>\n<td>g</td>\n<td>h</td>\n<td>i</td>\n<td>j</td>\n</tr>\n</tbody>\n</table>\n",

		"a|b\\|c|d\n---|---|---\nf|g\\|h|i\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b|c</th>\n<th>d</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>f</td>\n<td>g|h</td>\n<td>i</td>\n</tr>\n</tbody>\n</table>\n",

		"a   | c\n--- | ---:\nd   | e\nf   | g\n==  | =\nh   | j\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th align=\"right\">c</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>d</td>\n<td align=\"right\">e</td>\n</tr>\n\n<tr>\n<td>f</td>\n<td align=\"right\">g</td>\n</tr>\n</tbody>\n<tfoot>\n<tr>\n<td>h</td>\n<td align=\"right\">j</td>\n</tr>\n</tfoot>\n</table>\n",

		"a   | c\n--- | --:\nd   | e\nf   | g\n==  | =\n==  | =\nh   | j\n",
		"<table>\n<thead>\n<tr>\n<th>a</th>\n<th align=\"right\">c</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>d</td>\n<td align=\"right\">e</td>\n</tr>\n\n<tr>\n<td>f</td>\n<td align=\"right\">g</td>\n</tr>\n</tbody>\n<tfoot>\n<tr>\n<td>==</td>\n<td align=\"right\">=</td>\n</tr>\n\n<tr>\n<td>h</td>\n<td align=\"right\">j</td>\n</tr>\n</tfoot>\n</table>\n",
	}
	doTestsBlock(t, tests, EXTENSION_TABLES)
}

func TestBlockTable(t *testing.T) {
	var tests = []string{
		"|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n| Second |foo     |\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n|--------+--------+\n| Footer | Footer |\n|--------+--------+\n",
		"<table>\n<thead>\n<tr>\n<th>Defaul</th>\n<th>Left ald</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td><p>Second</p>\n</td>\n<td><p>foo</p>\n</td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n<td><ol>\n<li>Ite</li>\n<li>Ite</li>\n</ol>\n</td>\n</tr>\n\n<tr>\n<td><p>Footer</p>\n</td>\n<td><p>Footer</p>\n</td>\n</tr>\n</tbody>\n</table>\n",
		"|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n| Second |foo     |\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n|--------+--------+\n| Footer | Footer |\n|--------+--------+\n",
		"<table>\n<thead>\n<tr>\n<th>Defaul</th>\n<th>Left ald</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td><p>Second</p>\n</td>\n<td><p>foo</p>\n</td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n<td><ol>\n<li>Ite</li>\n<li>Ite</li>\n</ol>\n</td>\n</tr>\n\n<tr>\n<td><p>Footer</p>\n</td>\n<td><p>Footer</p>\n</td>\n</tr>\n</tbody>\n</table>\n",
		"|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n| Second |foo     |\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n|--------+--------+\n| Footer | Footer |\n|--------+--------+\nTable: this is a table\n",
		"<table>\n<caption>\nthis is a table\n</caption>\n<thead>\n<tr>\n<th>Defaul</th>\n<th>Left ald</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td><p>Second</p>\n</td>\n<td><p>foo</p>\n</td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n<td><ol>\n<li>Ite</li>\n<li>Ite</li>\n</ol>\n</td>\n</tr>\n\n<tr>\n<td><p>Footer</p>\n</td>\n<td><p>Footer</p>\n</td>\n</tr>\n</tbody>\n</table>\n",
		"|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n| Second |foo     |\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n+========+========+\n| Footer | Footer |\n|--------+--------+\nTable: this is a table\n",
		"<table>\n<caption>\nthis is a table\n</caption>\n<thead>\n<tr>\n<th>Defaul</th>\n<th>Left ald</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td><p>Second</p>\n</td>\n<td><p>foo</p>\n</td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n<td><ol>\n<li>Ite</li>\n<li>Ite</li>\n</ol>\n</td>\n</tr>\n</tbody>\n<tfoot>\n<tr>\n<td>Footer</td>\n<td>Footer</td>\n</tr>\n</tfoot>\n</table>\n",
		"|--------+--------+\n|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n+========+========+\n|--------+--------+\nTable: this is a table\n",
		"<table>\n<caption>\nthis is a table\n</caption>\n<thead>\n</thead>\n\n<tbody>\n<tr>\n<td></td>\n</tr>\n\n<tr>\n<td><p>Defaul</p>\n</td>\n</tr>\n\n<tr>\n<td></td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n</tr>\n</tbody>\n</table>\n",
		"|--------+--------+\n|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n+========+========+\n+========+========+\n|--------+--------+\nTable: this is a table\n",
		"<table>\n<thead>\n</thead>\n\n<tbody>\n<tr>\n<td></td>\n</tr>\n\n<tr>\n<td><p>Defaul</p>\n</td>\n</tr>\n\n<tr>\n<td></td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n</tr>\n</tbody>\n</table>\n\n<p>+========+========+\n|--------+--------+\nTable: this is a table</p>\n",
		"|--------+--------+\n| Defaul |Left ald|\n|--------|--------|\n|--------+--------+\n| Second | 2. Ite |\n| 2 line | 3. Ite |\n+========+========+\n|--------+--------+\nTable: this is a table\n",
		"<table>\n<caption>\nthis is a table\n</caption>\n<thead>\n<tr>\n<th>Defaul</th>\n<th>Left ald</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td></td>\n<td></td>\n</tr>\n\n<tr>\n<td><p>Second\n2 line</p>\n</td>\n<td><ol>\n<li>Ite</li>\n<li>Ite</li>\n</ol>\n</td>\n</tr>\n</tbody>\n</table>\n",
	}

	doTestsBlock(t, tests, EXTENSION_TABLES)
}

func TestUnorderedListWith_EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK(t *testing.T) {
	var tests = []string{
		"* Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"* Yin\n* Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"* Ting\n* Bong\n* Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"* Yin\n\n* Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ul>\n",

		"* Ting\n\n* Bong\n* Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ul>\n",

		"+ Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"+ Yin\n+ Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"+ Ting\n+ Bong\n+ Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"+ Yin\n\n+ Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ul>\n",

		"+ Ting\n\n+ Bong\n+ Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ul>\n",

		"- Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"- Yin\n- Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"- Ting\n- Bong\n- Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"- Yin\n\n- Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ul>\n",

		"- Ting\n\n- Bong\n- Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ul>\n",

		"*Hello\n",
		"<p>*Hello</p>\n",

		"*   Hello \n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"*   Hello \n    Next line \n",
		"<ul>\n<li>Hello\nNext line</li>\n</ul>\n",

		"Paragraph\n* No linebreak\n",
		"<p>Paragraph</p>\n\n<ul>\n<li>No linebreak</li>\n</ul>\n",

		"Paragraph\n\n* Linebreak\n",
		"<p>Paragraph</p>\n\n<ul>\n<li>Linebreak</li>\n</ul>\n",

		"*   List\n    * Nested list\n",
		"<ul>\n<li>List\n\n<ul>\n<li>Nested list</li>\n</ul></li>\n</ul>\n",

		"*   List\n\n    * Nested list\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li>Nested list</li>\n</ul></li>\n</ul>\n",

		"*   List\n    Second line\n\n    + Nested\n",
		"<ul>\n<li><p>List\nSecond line</p>\n\n<ul>\n<li>Nested</li>\n</ul></li>\n</ul>\n",

		"*   List\n    + Nested\n\n    Continued\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li>Nested</li>\n</ul>\n\n<p>Continued</p></li>\n</ul>\n",

		"*   List\n   * shallow indent\n",
		"<ul>\n<li>List\n\n<ul>\n<li>shallow indent</li>\n</ul></li>\n</ul>\n",

		"* List\n" +
			" * shallow indent\n" +
			"  * part of second list\n" +
			"   * still second\n" +
			"    * almost there\n" +
			"     * third level\n",
		"<ul>\n" +
			"<li>List\n\n" +
			"<ul>\n" +
			"<li>shallow indent</li>\n" +
			"<li>part of second list</li>\n" +
			"<li>still second</li>\n" +
			"<li>almost there\n\n" +
			"<ul>\n" +
			"<li>third level</li>\n" +
			"</ul></li>\n" +
			"</ul></li>\n" +
			"</ul>\n",

		"* List\n        extra indent, same paragraph\n",
		"<ul>\n<li>List\n    extra indent, same paragraph</li>\n</ul>\n",

		"* List\n\n        code block\n",
		"<ul>\n<li><p>List</p>\n\n<pre><code>code block\n</code></pre></li>\n</ul>\n",

		"* List\n\n          code block with spaces\n",
		"<ul>\n<li><p>List</p>\n\n<pre><code>  code block with spaces\n</code></pre></li>\n</ul>\n",

		"* List\n\n    * sublist\n\n    normal text\n\n    * another sublist\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li>sublist</li>\n</ul>\n\n<p>normal text</p>\n\n<ul>\n<li>another sublist</li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK)
}

func TestOrderedList_EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK(t *testing.T) {
	var tests = []string{
		"1. Hello\n",
		"<ol>\n<li>Hello</li>\n</ol>\n",

		"1. Yin\n2. Yang\n",
		"<ol>\n<li>Yin</li>\n<li>Yang</li>\n</ol>\n",

		"1. Ting\n2. Bong\n3. Goo\n",
		"<ol>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ol>\n",

		"1. Yin\n\n2. Yang\n",
		"<ol>\n<li><p>Yin</p></li>\n\n<li><p>Yang</p></li>\n</ol>\n",

		"1. Ting\n\n2. Bong\n3. Goo\n",
		"<ol>\n<li><p>Ting</p></li>\n\n<li><p>Bong</p></li>\n\n<li><p>Goo</p></li>\n</ol>\n",

		"1 Hello\n",
		"<p>1 Hello</p>\n",

		"1.Hello\n",
		"<p>1.Hello</p>\n",

		"1.  Hello \n",
		"<ol>\n<li>Hello</li>\n</ol>\n",

		"1.  Hello \n    Next line \n",
		"<ol>\n<li>Hello\nNext line</li>\n</ol>\n",

		"Paragraph\n1. No linebreak\n",
		"<p>Paragraph</p>\n\n<ol>\n<li>No linebreak</li>\n</ol>\n",

		"Paragraph\n\n1. Linebreak\n",
		"<p>Paragraph</p>\n\n<ol>\n<li>Linebreak</li>\n</ol>\n",

		"1.  List\n    1. Nested list\n",
		"<ol>\n<li>List\n\n<ol>\n<li>Nested list</li>\n</ol></li>\n</ol>\n",

		"1.  List\n\n    1. Nested list\n",
		"<ol>\n<li><p>List</p>\n\n<ol>\n<li>Nested list</li>\n</ol></li>\n</ol>\n",

		"1.  List\n    Second line\n\n    1. Nested\n",
		"<ol>\n<li><p>List\nSecond line</p>\n\n<ol>\n<li>Nested</li>\n</ol></li>\n</ol>\n",

		"1.  List\n    1. Nested\n\n    Continued\n",
		"<ol>\n<li><p>List</p>\n\n<ol>\n<li>Nested</li>\n</ol>\n\n<p>Continued</p></li>\n</ol>\n",

		"1.  List\n   1. shallow indent\n",
		"<ol>\n<li>List\n\n<ol>\n<li>shallow indent</li>\n</ol></li>\n</ol>\n",

		"1. List\n" +
			" 1. shallow indent\n" +
			"  2. part of second list\n" +
			"   3. still second\n" +
			"    4. almost there\n" +
			"     1. third level\n",
		"<ol>\n" +
			"<li>List\n\n" +
			"<ol>\n" +
			"<li>shallow indent</li>\n" +
			"<li>part of second list</li>\n" +
			"<li>still second</li>\n" +
			"<li>almost there\n\n" +
			"<ol>\n" +
			"<li>third level</li>\n" +
			"</ol></li>\n" +
			"</ol></li>\n" +
			"</ol>\n",

		"1. List\n        extra indent, same paragraph\n",
		"<ol>\n<li>List\n    extra indent, same paragraph</li>\n</ol>\n",

		"1. List\n\n        code block\n",
		"<ol>\n<li><p>List</p>\n\n<pre><code>code block\n</code></pre></li>\n</ol>\n",

		"1. List\n\n          code block with spaces\n",
		"<ol>\n<li><p>List</p>\n\n<pre><code>  code block with spaces\n</code></pre></li>\n</ol>\n",

		"1. List\n    * Mixted list\n",
		"<ol>\n<li>List\n\n<ul>\n<li>Mixted list</li>\n</ul></li>\n</ol>\n",

		"1. List\n * Mixed list\n",
		"<ol>\n<li>List\n\n<ul>\n<li>Mixed list</li>\n</ul></li>\n</ol>\n",

		"* Start with unordered\n 1. Ordered\n",
		"<ul>\n<li>Start with unordered\n\n<ol>\n<li>Ordered</li>\n</ol></li>\n</ul>\n",

		"* Start with unordered\n    1. Ordered\n",
		"<ul>\n<li>Start with unordered\n\n<ol>\n<li>Ordered</li>\n</ol></li>\n</ul>\n",

		"1. numbers\n1. are ignored\n",
		"<ol>\n<li>numbers</li>\n<li>are ignored</li>\n</ol>\n",
	}
	doTestsBlock(t, tests, EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK)
}

func TestFencedCodeBlock_EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK(t *testing.T) {
	var tests = []string{
		"``` go\nfunc foo() bool {\n\treturn true;\n}\n```\n",
		"<pre><code class=\"language-go\">func foo() bool {\n\treturn true;\n}\n</code></pre>\n",

		"``` c\n/* special & char < > \" escaping */\n```\n",
		"<pre><code class=\"language-c\">/* special &amp; char &lt; &gt; &quot; escaping */\n</code></pre>\n",

		"``` c\nno *inline* processing ~~of text~~\n```\n",
		"<pre><code class=\"language-c\">no *inline* processing ~~of text~~\n</code></pre>\n",

		"```\nNo language\n```\n",
		"<pre><code>No language\n</code></pre>\n",

		"``` {ocaml}\nlanguage in braces\n```\n",
		"<pre><code class=\"language-ocaml\">language in braces\n</code></pre>\n",

		"```    {ocaml}      \nwith extra whitespace\n```\n",
		"<pre><code class=\"language-ocaml\">with extra whitespace\n</code></pre>\n",

		"```{   ocaml   }\nwith extra whitespace\n```\n",
		"<pre><code class=\"language-ocaml\">with extra whitespace\n</code></pre>\n",

		"~ ~~ java\nWith whitespace\n~~~\n",
		"<p>~ ~~ java\nWith whitespace\n~~~</p>\n",

		"~~\nonly two\n~~\n",
		"<p>~~\nonly two\n~~</p>\n",

		"```` python\nextra\n````\n",
		"<pre><code class=\"language-python\">extra\n</code></pre>\n",

		"~~~ perl\nthree to start, four to end\n~~~~\n",
		"<pre><code class=\"language-perl\">three to start, four to end\n</code></pre>\n",

		"~~~~ perl\nfour to start, three to end\n~~~\n",
		"<p>~~~~ perl\nfour to start, three to end\n~~~</p>\n",

		"~~~ bash\ntildes\n~~~\n",
		"<pre><code class=\"language-bash\">tildes\n</code></pre>\n",

		"``` lisp\nno ending\n",
		"<p>``` lisp\nno ending</p>\n",

		"~~~ lisp\nend with language\n~~~ lisp\n",
		"<p>~~~ lisp\nend with language\n~~~ lisp</p>\n",

		"```\nmismatched begin and end\n~~~\n",
		"<p>```\nmismatched begin and end\n~~~</p>\n",

		"~~~\nmismatched begin and end\n```\n",
		"<p>~~~\nmismatched begin and end\n```</p>\n",

		"   ``` oz\nleading spaces\n```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		"  ``` oz\nleading spaces\n ```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		" ``` oz\nleading spaces\n  ```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		"``` oz\nleading spaces\n   ```\n",
		"<pre><code class=\"language-oz\">leading spaces\n</code></pre>\n",

		"    ``` oz\nleading spaces\n    ```\n",
		"<pre><code>``` oz\n</code></pre>\n\n<p>leading spaces</p>\n\n<pre><code>```\n</code></pre>\n",
	}
	doTestsBlock(t, tests, EXTENSION_FENCED_CODE|EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK)
}

func TestDefinitionListXML(t *testing.T) {
	var tests = []string{
		"Term1\n:   Hi There",
		"<dl>\n<dt>Term1</dt>\n<dd>Hi There</dd>\n</dl>\n",

		"Term1\n:   Yin\nTerm2\n:   Yang\n",
		"<dl>\n<dt>Term1</dt>\n<dd>Yin</dd>\n<dt>Term2</dt>\n<dd>Yang</dd>\n</dl>\n",

		// fix sourcecode/artwork here.
		//		`Term 1
		//:   This is a definition with two paragraphs. Lorem ipsum
		//
		//    Vestibulum enim wisi, viverra nec, fringilla in, laoreet
		//    vitae, risus.
		//
		//Term 2
		//:   This definition has a code block, a blockquote and a list.
		//
		//        code block.
		//
		//    > block quote
		//    > on two lines.
		//
		//    1.  first list item
		//    2.  second list item`,
		//
		//		"<dl>\n<dt>Term 1</dt>\n<dd><t>This is a definition with two paragraphs. Lorem ipsum</t>\n<t>Vestibulum enim wisi, viverra nec, fringilla in, laoreet\nvitae, risus.</t></dd>\n<dt>Term 2</dt>\n<dd><t>This definition has a code block, a blockquote and a list.</t>\n<sourcecode>\ncode block.\n</sourcecode>\n<blockquote>\n<t>block quote\non two lines.</t>\n</blockquote>\n<ol>\n<li>first list item</li>\n<li>second list item</li>\n</ol></dd>\n</dl>\n",
		//
		`Apple
:   Pomaceous fruit of plants of the genus Malus in
    the family Rosaceae.

Orange and *Apples*
:   The thing of an evergreen tree of the genus Citrus.`,
		"<dl>\n<dt>Apple</dt>\n<dd><t>Pomaceous fruit of plants of the genus Malus in\nthe family Rosaceae.</t></dd>\n<dt>Orange and <em>Apples</em></dt>\n<dd><t>The thing of an evergreen tree of the genus Citrus.</t></dd>\n</dl>\n",
	}
	doTestsBlockXML(t, tests, 0)
}

func TestAbstractNoteAsideXML(t *testing.T) {
	var tests = []string{
		"A> begin of abstract\nA>\nA> this is an abstract\n",
		"<abstract>\n<t>begin of abstract</t>\n<t>this is an abstract</t>\n</abstract>\n",

		"N> begin of note\nN> this is a note\n",
		"<note>\n<t>begin of note\nthis is a note</t>\n</note>\n",

		"AS> begin of aside\nAS> this is an aside\n",
		"<aside>\n<t>begin of aside\nthis is an aside</t>\n</aside>\n",
	}
	doTestsBlockXML(t, tests, 0)
}

func TestOrderedListStartXML(t *testing.T) {
	var tests = []string{
		"1. hello\n1. hello\n\ndivide\n\n4. hello\n5. hello\n\ndivide\n\n 7. hello\n5. hello\n",
		"<ol>\n<li>hello</li>\n<li>hello</li>\n</ol>\n<t>divide</t>\n<ol start=\"4\">\n<li>hello</li>\n<li>hello</li>\n</ol>\n<t>divide</t>\n<ol>\n<li>hello</li>\n<li>hello</li>\n</ol>\n",
	}
	doTestsBlockXML(t, tests, 0)
}

func TestInsideQuoteXML(t *testing.T) {
	var tests = []string{
	// need to fix header in quote
	//	"N> # hello\nN>\n N> text\n",
	//	"<note>\n<name>hello</name>\n<t>text</t>\n</note>\n",
	}
	doTestsBlockXML(t, tests, 0)
}

func TestInlineAttr(t *testing.T) {
	var tests = []string{
		"{attribution=\"BLA BLA\" .green}\n{bla=BLA}\n{more=\"ALB ALB\" #ref:quote .yellow}\n> Hallo2\n> Hallo3\n\nThis is no one `{source='BLEIP'}` on of them\n\n{evenmore=\"BLE BLE\"}\n> Hallo6\n> Hallo7",
		"<blockquote anchor=\"ref:quote\" class=\"green yellow\" attribution=\"BLA BLA\" bla=\"BLA\" more=\"ALB ALB\">\n<t>Hallo2\nHallo3</t>\n</blockquote>\n<t>This is no one <tt>{source='BLEIP'}</tt> on of them</t>\n<blockquote evenmore=\"BLE BLE\">\n<t>Hallo6\nHallo7</t>\n</blockquote>\n",

		"{style=\"format REQ(%c)\" start=\"4\"}\n1. Term1\n2. Term2",
		"<ol start=\"4\" style=\"format REQ(%c)\">\n<li>Term1</li>\n<li>Term2</li>\n</ol>\n",

		"    {style=\"format REQ(%c)\" start=\"4\"}\n1. Term1\n2. Term2",
		"<artwork>\n{style=\"format REQ(%c)\" start=\"4\"}\n</artwork>\n<ol>\n<li>Term1</li>\n<li>Term2</li>\n</ol>\n",

		"{.green #ref1}\n# hallo\n\n{.yellow}\n# hallo {#ref2}\n\n{.blue #ref3}\n# hallo {#ref4}\n",
		"\n<section anchor=\"ref1\" class=\"green\"><name>hallo</name>\n</section>\n\n<section anchor=\"ref2\" class=\"yellow\"><name>hallo</name>\n</section>\n\n<section anchor=\"ref3\" class=\"blue\"><name>hallo</name>\n",
	}
	doTestsBlockXML(t, tests, 0)
}

func TestOrderedExampleListXML(t *testing.T) {
	var tests = []string{
		`(@good)  Example1

(@good)  Example2

As this illustrates

(@good)  Example3

As we can see from (@good) some things never change. Some we do not when we reference (@good1).
`,
		"<ol group=\"good\">\n<li><t>Example1</t></li>\n<li><t>Example2</t></li>\n</ol>\n<t>As this illustrates</t>\n<ol group=\"good\">\n<li>Example3</li>\n</ol>\n<t>As we can see from (2) some things never change. Some we do not when we reference (@good1).</t>\n",
	}
	doTestsBlockXML(t, tests, 0)
}

// TODO:
// Roman lists
// Letter prefix
// figure caption
// table caption
// frontmatter
// broken IALs
