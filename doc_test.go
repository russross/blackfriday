// Parse document fragments.

package mmark

import "testing"

var extensions = EXTENSION_TABLES | EXTENSION_FENCED_CODE | EXTENSION_AUTOLINK | EXTENSION_SPACE_HEADERS | EXTENSION_CITATION | EXTENSION_TITLEBLOCK_TOML | EXTENSION_HEADER_IDS | EXTENSION_AUTO_HEADER_IDS | EXTENSION_UNIQUE_HEADER_IDS | EXTENSION_FOOTNOTES | EXTENSION_SHORT_REF | EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK | EXTENSION_INCLUDE | EXTENSION_PARTS | EXTENSION_ABBREVIATIONS

func TestCarriageReturn(t *testing.T) {
	var tests = []string{".# Abstract\n\rThis document\n\r# More\n\rand more\n\r{#fig-a}\n\r```\n\rfigure\n\r```\n\rFigure: Traditional Media Server\n\r\n\r{#fig-b}\n\r```\n\rfigure\n\r```\n\rFigure: Endpoint\n\r",
		"\n<abstract>\n<t>\nThis document\n</t>\n</abstract>\n\n\n<section anchor=\"more\">\n<name>More</name>\n<t>\nand more\n</t>\n<artwork anchor=\"fig-a\">\n\nfigure\n\n</artwork>\n<t>\nFigure: Traditional Media Server\n</t>\n<artwork anchor=\"fig-b\">\n\nfigure\n\n</artwork>\n<t>\nFigure: Endpoint\n</t>\n</section>\n",
	}
	doTestsBlockXML(t, tests, extensions)
}

func TestFigureCaption(t *testing.T) {
	var tests = []string{
		// This checks the *single* newline after the Figure: caption.
		".# Abstract\nThis document\n# More\nand more\n{#fig-a}\n```\nfigure\n```\nFigure: Traditional Media Server\n\n{#fig-b}\n```\nfigure\n```\n",
		"\n<abstract>\n<t>\nThis document\n</t>\n</abstract>\n\n\n<section anchor=\"more\">\n<name>More</name>\n<t>\nand more\n\n</t>\n<figure anchor=\"fig-a\">\n<name>Traditional Media Server\n</name>\n<artwork>\nfigure\n</artwork>\n</figure>\n<artwork anchor=\"fig-b\">\nfigure\n</artwork>\n</section>\n",
	}
	doTestsBlockXML(t, tests, extensions)
}
