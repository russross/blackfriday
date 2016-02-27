package mmark

import "testing"

func TestIssueXXX(t *testing.T) {
	tests := []string{
		"абвгдеёжзийклмнопрстуфх",
		"<p>абвгдеёжзийклмнопрстуфх</p>\n",
	}

	doTestsBlock(t, tests, 0)
}
