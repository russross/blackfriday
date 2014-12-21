// Unit tests for complete doc parsing

package mmark

import (
	"testing"
)

func runMarkdownStandAloneXML(input string, extensions int) string {
	xmlFlags := XML_STANDALONE

	extensions |= commonXmlExtensions
	extensions |= EXTENSION_UNIQUE_HEADER_IDS
	renderer := XmlRenderer(xmlFlags)

	return string(Markdown([]byte(input), renderer, extensions))
}

func doTestsStandAlone(t *testing.T, tests []string, extensions int) {
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
		actual := runMarkdownStandAloneXML(candidate, extensions)
		if actual != expected {
			t.Errorf("\nInput   [%#v]\nExpected[%#v]\nActual  [%#v]",
				candidate, expected, actual)
		}
	}
}

func TestExampleDocSmall(t *testing.T) {
	var tests = []string{
		`% title = "Signaling Type of Denial via Delegation Signer Records"
% abbrev = "DS Denial Type Signalling"
% docName = "example-00"
% ipr = "trust200902"
% category = "info"
%
% date = 2014-12-01T00:00:00Z
% area = "Network"
% keyword = ["DNSSEC"]
%
% [[author]]
% initials = "R."
% surname = "Gieben"
% fullname = "R. (Miek) Gieben"
% organization = "Google"
%   [author.address]
%   email = "miek@google.com"

A> This will become a paragraph in the abstract.

{mainmatter}

# Introduction

The DS Resource Record [@RFC3658]...
particular a digest of the **DNSKEY**, *the* algorithm used for signature of the

> This transition method is best described as a hack.
Quote: Miek Gieben -- http://miek.nl/

In this document, the key words "MUST", "MUST NOT", "REQUIRED",
"SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY",
and "OPTIONAL" are to be interpreted as described in [@!RFC2119].

<!-- Miek -- are you sure you want to include this stuff? -->

{align="left"}
VALUE   | Digest Type        | Status
------: | :----------------- | -------------
 0      | Reserved           | -
 1      | SHA-1              | MANDATORY
=====   | ======             | ====
 TBD    | DT-SHA-256         | OPTIONAL
Table: As shown here.

*[HTML]: Hyper Text Markup Language

What HTML says could not be denied.

{backmatter}

# Other Options

H~2~O is a liquid.  2^10^ is 1024.

A.  Item1
B.  Item2
`,
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<t>% title = \"Signaling Type of Denial via Delegation Signer Records\"\n% abbrev = \"DS Denial Type Signalling\"\n% docName = \"example-00\"\n% ipr = \"trust200902\"\n% category = \"info\"\n%\n% date = 2014-12-01T00:00:00Z\n% area = \"Network\"\n% keyword = [\"DNSSEC\"]\n%\n% [[author]]\n% initials = \"R.\"\n% surname = \"Gieben\"\n% fullname = \"R. (Miek) Gieben\"\n% organization = \"Google\"\n%   [author.address]\n%   email = \"miek@google.com\"</t>\n<abstract>\n<t>This will become a paragraph in the abstract.</t>\n</abstract>\n</front>\n\n<middle>\n\n<section anchor=\"introduction\"><name>Introduction</name>\n<t>The DS Resource Record <xref target=\"RFC3658\"/>...\nparticular a digest of the <strong>DNSKEY</strong>, <em>the</em> algorithm used for signature of the</t>\n<blockquote cite=\"Miek Gieben\" quotedFrom=\"http://miek.nl/\">\n<t>This transition method is best described as a hack.</t>\n</blockquote>\n<t>In this document, the key words \"MUST\", \"MUST NOT\", \"REQUIRED\",\n\"SHALL\", \"SHALL NOT\", \"SHOULD\", \"SHOULD NOT\", \"RECOMMENDED\", \"MAY\",\nand \"OPTIONAL\" are to be interpreted as described in <xref target=\"RFC2119\"/>.</t>\n<t><cref source=\"Miek\">are you sure you want to include this stuff?</cref></t>\n<table align=\"left\">\n<name>As shown here.\n</name>\n<thead>\n<tr><th align=\"right\">VALUE</th><th align=\"left\">Digest Type</th><th align=\"center\">Status</th></tr>\n</thead>\n<tr><td>0</td><td>Reserved</td><td>-</td></tr>\n<tr><td>1</td><td>SHA-1</td><td>MANDATORY</td></tr>\n<tfoot>\n<tr><th align=\"right\">VALUE</th><th align=\"left\">Digest Type</th><th align=\"center\">Status</th></tr>\n</tfoot>\n</table>\n<t>What HTML says could not be denied.</t>\n</section>\n\n</middle>\n<back>\n<references title=\"Informative References\">\n<xi:include href=\"reference.RFC.3658.xml\"/>\n<references title=\"Normative References\">\n<xi:include href=\"reference.RFC.2119.xml\"/>\n</references>\n\n<section anchor=\"other-options\"><name>Other Options</name>\n<t>H<sub>2</sub>O is a liquid.  2<sup>10</sup> is 1024.</t>\n<ol>\n<li>Item1</li>\n<li>Item2</li>\n</ol>\n</section>\n\n</back>\n</rfc>\n",
	}

	doTestsStandAlone(t, tests, 0)
}
