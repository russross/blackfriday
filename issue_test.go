package mmark

import "testing"

func TestIssue55(t *testing.T) {
	tests := []string{
		"абвгдеёжзийклмнопрстуфх",
		"<p>абвгдеёжзийклмнопрстуфх</p>\n",
	}

	doTestsBlock(t, tests, 0)
}

func TestIssue59(t *testing.T) {
	// Need renderer option flag as well, which isn't carried through to the actual
	// tests. For now, just check the parsing.
	tests := []string{
		"stuff\n\n{frontmatter} stuff",
		"<t>\nstuff\n</t>\n<t>\n{frontmatter} stuff\n</t>\n",

		"stuff\n\n{frontmatter}\n",
		"<t>\nstuff\n</t>\n",

		"{frontmatter}\ntext\n",
		"<t>\ntext\n</t>\n",
	}
	doTestsBlockXML(t, tests, EXTENSION_MATTER)
}

func TestIssue73(t *testing.T) {
	tests := []string{
		`* [foo](http://bar)

(@good)  Example

As (@good) says
`,
		"<ul>\n<li><eref target=\"http://bar\">foo</eref></li>\n</ul>\n<ol group=\"good\">\n<li>Example</li>\n</ol>\n<t>\nAs (1) says\n</t>\n",

		`Mmark
: A Markdown-superset converter

* [foo](http://bar)`,
		"<dl>\n<dt>Mmark</dt>\n<dd>A Markdown-superset converter</dd>\n</dl>\n<ul>\n<li><eref target=\"http://bar\">foo</eref></li>\n</ul>\n",

		`Mmark
: A Markdown-superset converter

 * [foo](http://bar)`,
		"<dl>\n<dt>Mmark</dt>\n<dd><t>\nA Markdown-superset converter\n</t>\n<ul>\n<li><eref target=\"http://bar\">foo</eref></li>\n</ul></dd>\n</dl>\n",

		`Mmark
: A Markdown-superset converter

1. [foo](http://bar)`,
		"<dl>\n<dt>Mmark</dt>\n<dd>A Markdown-superset converter</dd>\n</dl>\n<ol>\n<li><eref target=\"http://bar\">foo</eref></li>\n</ol>\n",
	}
	doTestsBlockXML(t, tests, commonXmlExtensions)
}

func TestIssue82(t *testing.T) {
	tests := []string{
		`Option Type

: 8-bit identifier of the type of option.

  ~~~~~
  HEX         act  chg  rest
  ~~~~~
  Figure: Scenic Routing Option Type`,
		"<dl>\n<dt>Option Type</dt>\n<dd><t>\n8-bit identifier of the type of option.\n</t></dd>\n</dl>\n<figure>\n<name>Scenic Routing Option Type</name>\n<artwork>\nHEX         act  chg  rest\n</artwork>\n</figure>\n",
	}

	doTestsBlockXML(t, tests, commonXmlExtensions)
}
