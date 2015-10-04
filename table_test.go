package mmark

import "testing"

func TestTableColSpan(t *testing.T) {
	var tests = []string{`
| Column 1 | Column 2 | Column 3 |
| -------- | :------: | -------- |
| No span  | Span across two columns ||
| Span accross two columns || No span |`,
		"<table>\n<thead>\n<tr>\n<th>Column 1</th>\n<th align=\"center\">Column 2</th>\n<th>Column 3</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>No span</td>\n<td align=\"center\" colspan=\"2\">Span across two columns</td>\n</tr>\n\n<tr>\n<td colspan=\"2\">Span accross two columns</td>\n<td>No span</td>\n</tr>\n</tbody>\n</table>\n",
	}
	doTestsBlock(t, tests, EXTENSION_TABLES)
}
