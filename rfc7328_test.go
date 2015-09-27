// Unit tests for block parsing

package mmark

import "testing"

func init() {
	test = true
}

func runMarkdownBlockXML_rfc7328(input string, extensions int) string {
	xmlFlags := 0

	extensions |= commonXmlExtensions
	extensions |= EXTENSION_UNIQUE_HEADER_IDS
	extensions |= EXTENSION_FOOTNOTES
	extensions |= EXTENSION_RFC7328
	renderer := XmlRenderer(xmlFlags)

	return Parse([]byte(input), renderer, extensions).String()
}

func doTestsBlockXML_rfc7328(t *testing.T, tests []string, extensions int) {
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
		actual := runMarkdownBlockXML_rfc7328(candidate, extensions)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}

		// now test every substring to stress test bounds checking
		if !testing.Short() {
			for start := 0; start < len(input); start++ {
				for end := start + 1; end <= len(input); end++ {
					candidate = input[start:end]
					_ = runMarkdownBlockXML_rfc7328(candidate, extensions)
				}
			}
		}
	}
}

func TestConversionFromRFC7328(t *testing.T) {
	var tests = []string{
		`What is this?

    example.org.        SOA ( ... )
    example.org.        NS  a.example.org.
    a.example.org.      A 192.0.2.1
                        TXT "a record"
    d.example.org.      A 192.0.2.1
                        TXT "d record"
^[fig:the-unsigned::The unsigned "example.org" zone.]
`,
		"<t>\nWhat is this?\n</t>\n<artwork>\nexample.org.        SOA ( ... )\nexample.org.        NS  a.example.org.\na.example.org.      A 192.0.2.1\n                    TXT \"a record\"\nd.example.org.      A 192.0.2.1\n                    TXT \"d record\"\n</artwork>\n<t>\n\n</t>\n",

		"And another one. An index ^[ ^itemindex^ subitem ]",
		"<t>\nAnd another one. An index <iref item=\"itemindex\" subitem=\"subitem\"/>\n</t>\n",

		"Index ^[ ^indexer^   ]",
		"<t>\nIndex <iref item=\"indexer\" subitem=\"\"/>\n</t>\n",
	}
	doTestsBlockXML_rfc7328(t, tests, 0)
}
