package mmark

import "testing"

func TestIssue55(t *testing.T) {
	tests := []string{
		"абвгдеёжзийклмнопрстуфх",
		"<p>абвгдеёжзийклмнопрстуфх</p>\n",
	}

	doTestsBlock(t, tests, 0)
}
