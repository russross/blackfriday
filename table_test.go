package mmark

import "testing"

func TestTableColSpan(t *testing.T) {
	var tests = []string{`
| Column 1 | Column 2 | Column 3 |
| -------- | :------: | -------- |
| No span  | Span across two columns ||
| Span accross two columns || No span |`,
		"<table>\n<thead>\n<tr>\n<th>Column 1</th>\n<th align=\"center\">Column 2</th>\n<th>Column 3</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td>No span</td>\n<td align=\"center\" colspan=\"2\">Span across two columns</td>\n</tr>\n\n<tr>\n<td colspan=\"2\">Span accross two columns</td>\n<td>No span</td>\n</tr>\n</tbody>\n</table>\n",

		`
|+--
| Default aligned |Left aligned| Center aligned
|-----------------|:-----------|:---------------:
| First body part more test   || 1. Third cell
| Second line                 || 2. **strong**
| Third line       hallo      || 3. baz
`,
		"<table>\n<thead>\n<tr>\n<th>Default aligned</th>\n<th align=\"left\">Left aligned</th>\n<th align=\"center\">Center aligned</th>\n</tr>\n</thead>\n\n<tbody>\n<tr>\n<td colspan=\"2\"><p>First body part more test\nSecond line\nThird line       hallo</p>\n</td>\n<td align=\"center\"><ol>\n<li>Third cell</li>\n<li><strong>strong</strong></li>\n<li>baz</li>\n</ol>\n</td>\n</tr>\n</tbody>\n</table>\n",
	}
	doTestsBlock(t, tests, EXTENSION_TABLES)
}
