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
		"<t>\nstuff\n</t>\n<t>\n\n</t>\n",

		"{frontmatter}\ntext\n",
		"<t>\n\ntext\n</t>\n",
	}
	doTestsBlockXML(t, tests, EXTENSION_MATTER)
}
