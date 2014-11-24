// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
// Extended by Miek Gieben <miek@miek.nl> © 2014.

//
// Unit tests for inline parsing
//

package blackfriday

import (
	"testing"
)

func runMarkdownInlineXML(input string, extensions, xmlFlags int) string {
	extensions |= EXTENSION_AUTOLINK
	extensions |= EXTENSION_STRIKETHROUGH
	extensions |= EXTENSION_INDEX
	extensions |= EXTENSION_CITATION

	renderer := XmlRenderer(xmlFlags)

	return string(Markdown([]byte(input), renderer, extensions))
}

func doTestsInlineXML(t *testing.T, tests []string) {
	doTestsInlineParamXML(t, tests, 0, 0)
}

func doTestsInlineParamXML(t *testing.T, tests []string, extensions, xmlFlags int) {
	var candidate string

	for i := 0; i+1 < len(tests); i += 2 {
		input := tests[i]
		candidate = input
		expected := tests[i+1]
		actual := runMarkdownInlineXML(candidate, extensions, xmlFlags)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownInlineXML(candidate, extensions, xmlFlags)
				}
			}
		}
	}
}

func TestIndexXML(t *testing.T) {
	var tests = []string{
		"(((Tiger, Cats)))\n",
		"<t>\n<iref item=\"Tiger\" subitem=\"Cats\"/>\n</t>\n",

		"(((Tiger, Cats))\n",
		"<t>\n(((Tiger, Cats))\n</t>\n",
	}
	doTestsInlineXML(t, tests)
}

func TestCitationXML(t *testing.T) {
	var tests = []string{
		"(((Tiger, Cats)))\n",
		"<t>\n<iref item=\"Tiger\" subitem=\"Cats\"/>\n</t>\n",

	}
	doTestsInlineXML(t, tests)
}
