package mmark

import "testing"

func TestIALSyntax(t *testing.T) {
	tests := []string{
		"{.class key=\"value\"}\n> Quote\n",
		"<blockquote class=\"class\" key=\"value\">\n<p>Quote</p>\n</blockquote>\n",

		"{.class key=\"va\\}lue\"}\n> Quote\n",
		"<blockquote class=\"class\" key=\"va\\}lue\">\n<p>Quote</p>\n</blockquote>\n",
	}

	doTestsBlock(t, tests, 0)
}
