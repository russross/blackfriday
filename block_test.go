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

		"*   List\n    * Sublist\n    Not a header\n    ------\n",
		"<ul>\n<li>List\n\n<ul>\n<li>Sublist\nNot a header\n------</li>\n</ul></li>\n</ul>\n",
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

		"*   List\n    * Sublist\n    Not a header\n    ------\n",
		"<ul>\n<li>List\n\n<ul>\n<li>Sublist\nNot a header\n------</li>\n</ul></li>\n</ul>\n",
	}
	doTestsBlock(t, tests, EXTENSION_SPACE_HEADERS)
}
