// Unit tests for block parsing

package mmark

import "testing"

var extensions = EXTENSION_TABLES | EXTENSION_FENCED_CODE | EXTENSION_AUTOLINK | EXTENSION_SPACE_HEADERS | EXTENSION_CITATION | EXTENSION_TITLEBLOCK_TOML | EXTENSION_HEADER_IDS | EXTENSION_AUTO_HEADER_IDS | EXTENSION_UNIQUE_HEADER_IDS | EXTENSION_FOOTNOTES | EXTENSION_SHORT_REF | EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK | EXTENSION_INCLUDE | EXTENSION_PARTS | EXTENSION_ABBREVIATIONS

func TestCarriageReturn(t *testing.T) {
	var tests = []string{".# Abstract\n\rThis document\n\r# More\n\rand more\n\r{#fig-a}\n\r```\n\rfigure\n\r```\n\rFigure: Traditional Media Server\n\r{#fig-b}\n\r```\n\rfigure\n\r```\n\rFigure: Endpoint\n\r",
		"\n<abstract>\n<t>\nThis document\n</t>\n</abstract>\n\n\n<section anchor=\"more\"><name>More</name>\n<t>\nand more\n</t>\n<artwork anchor=\"fig-a\">\n\nfigure\n\n</artwork>\n<t>\nFigure: Traditional Media Server\n</t>\n<artwork anchor=\"fig-b\">\n\nfigure\n\n</artwork>\n<t>\nFigure: Endpoint\n</t>\n</section>\n",
	}
	doTestsBlockXML(t, tests, extensions)
}

// fails.
func testAnchorAbstract(t *testing.T) {
	var tests = []string{".# Abstract\nThis document\n# More\nand more\n{#fig-a}\n```\nfigure\n```\nFigure: Traditional Media Server\n{#fig-b}\n```\nfigure\n```\nFigure: Endpoint\n",
		"\n<abstract>\n<t>\nThis document\n</t>\n</abstract>\n\n\n<section anchor=\"more\"><name>More</name>\n<t>\nand more\n</t>\n<artwork anchor=\"fig-a\">\n\nfigure\n\n</artwork>\n<t>\nFigure: Traditional Media Server\n</t>\n<artwork anchor=\"fig-b\">\n\nfigure\n\n</artwork>\n<t>\nFigure: Endpoint\n</t>\n</section>\n",
	}
	doTestsBlockXML(t, tests, extensions)
}
