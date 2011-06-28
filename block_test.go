//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
//

//
// Unit tests for block parsing
//

package blackfriday

import (
	"testing"
)

func runMarkdownBlock(input string, extensions uint32) string {
	html_flags := 0
	html_flags |= HTML_USE_XHTML

	renderer := HtmlRenderer(html_flags)

	return string(Markdown([]byte(input), renderer, extensions))
}

func doTestsBlock(t *testing.T, tests []string, extensions uint32) {
	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		expected := tests[i+1]
		actual := runMarkdownBlock(input, extensions)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				input, expected, actual)
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
		"<ul>\n<li><p>List</p>\n\n<h1>Header</h1></li>\n<li><p>List</p></li>\n</ul>\n",

		"* List\n#Header\n* List\n",
		"<ul>\n<li><p>List</p>\n\n<h1>Header</h1></li>\n<li><p>List</p></li>\n</ul>\n",

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
		"<ul>\n<li><p>List</p>\n\n<h1>Header</h1></li>\n<li><p>List</p></li>\n</ul>\n",

		"* List\n#Header\n* List\n",
		"<ul>\n<li>List\n#Header</li>\n<li>List</li>\n</ul>\n",

		"*   List\n    * Nested list\n    # Nested header\n",
		"<ul>\n<li><p>List</p>\n\n<ul>\n<li><p>Nested list</p>\n\n" +
			"<h1>Nested header</h1></li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, EXTENSION_SPACE_HEADERS)
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
		"<ul>\n<li><p>Yin</p></li>\n<li><p>Yang</p></li>\n</ul>\n",

		"* Ting\n\n* Bong\n* Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n<li><p>Bong</p></li>\n<li><p>Goo</p></li>\n</ul>\n",

		"+ Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"+ Yin\n+ Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"+ Ting\n+ Bong\n+ Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"+ Yin\n\n+ Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n<li><p>Yang</p></li>\n</ul>\n",

		"+ Ting\n\n+ Bong\n+ Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n<li><p>Bong</p></li>\n<li><p>Goo</p></li>\n</ul>\n",

		"- Hello\n",
		"<ul>\n<li>Hello</li>\n</ul>\n",

		"- Yin\n- Yang\n",
		"<ul>\n<li>Yin</li>\n<li>Yang</li>\n</ul>\n",

		"- Ting\n- Bong\n- Goo\n",
		"<ul>\n<li>Ting</li>\n<li>Bong</li>\n<li>Goo</li>\n</ul>\n",

		"- Yin\n\n- Yang\n",
		"<ul>\n<li><p>Yin</p></li>\n<li><p>Yang</p></li>\n</ul>\n",

		"- Ting\n\n- Bong\n- Goo\n",
		"<ul>\n<li><p>Ting</p></li>\n<li><p>Bong</p></li>\n<li><p>Goo</p></li>\n</ul>\n",

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
		"<ol>\n<li><p>Yin</p></li>\n<li><p>Yang</p></li>\n</ol>\n",

		"1. Ting\n\n2. Bong\n3. Goo\n",
		"<ol>\n<li><p>Ting</p></li>\n<li><p>Bong</p></li>\n<li><p>Goo</p></li>\n</ol>\n",

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
