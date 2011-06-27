//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
//

//
// Functions to parse block-level elements.
//

package blackfriday

import (
	"bytes"
)

// parse block-level data
func parseBlock(out *bytes.Buffer, rndr *render, data []byte) {
	// this is called recursively: enforce a maximum depth
	if rndr.nesting >= rndr.maxNesting {
		return
	}
	rndr.nesting++

	// parse out one block-level construct at a time
	for len(data) > 0 {
		// prefixed header:
		//
		// # Header 1
		// ## Header 2
		// ...
		// ###### Header 6
		if isPrefixHeader(rndr, data) {
			data = data[blockPrefixHeader(out, rndr, data):]
			continue
		}

		// block of preformatted HTML:
		//
		// <div>
		//     ...
		// </div>
		if data[0] == '<' && rndr.mk.BlockHtml != nil {
			if i := blockHtml(out, rndr, data, true); i > 0 {
				data = data[i:]
				continue
			}
		}

		// blank lines.  note: returns the # of bytes to skip
		if i := isEmpty(data); i > 0 {
			data = data[i:]
			continue
		}

		// horizontal rule:
		//
		// ------
		// or
		// ******
		// or
		// ______
		if isHRule(data) {
			if rndr.mk.HRule != nil {
				rndr.mk.HRule(out, rndr.mk.Opaque)
			}
			var i int
			for i = 0; i < len(data) && data[i] != '\n'; i++ {
			}
			data = data[i:]
			continue
		}

		// fenced code block:
		//
		// ``` go
		// func fact(n int) int {
		//     if n <= 1 {
		//         return n
		//     }
		//     return n * fact(n-1)
		// }
		// ```
		if rndr.flags&EXTENSION_FENCED_CODE != 0 {
			if i := blockFencedCode(out, rndr, data); i > 0 {
				data = data[i:]
				continue
			}
		}

		// table:
		//
		// Name  | Age | Phone
		// ------|-----|---------
		// Bob   | 31  | 555-1234
		// Alice | 27  | 555-4321
		if rndr.flags&EXTENSION_TABLES != 0 {
			if i := blockTable(out, rndr, data); i > 0 {
				data = data[i:]
				continue
			}
		}

		// block quote:
		//
		// > A big quote I found somewhere
		// > on the web
		if blockQuotePrefix(data) > 0 {
			data = data[blockQuote(out, rndr, data):]
			continue
		}

		// indented code block:
		//
		//     func max(a, b int) int {
		//         if a > b {
		//             return a
		//         }
		//         return b
		//      }
		if blockCodePrefix(data) > 0 {
			data = data[blockCode(out, rndr, data):]
			continue
		}

		// an itemized/unordered list:
		//
		// * Item 1
		// * Item 2
		//
		// also works with + or -
		if blockUliPrefix(data) > 0 {
			data = data[blockList(out, rndr, data, 0):]
			continue
		}

		// a numbered/ordered list:
		//
		// 1. Item 1
		// 2. Item 2
		if blockOliPrefix(data) > 0 {
			data = data[blockList(out, rndr, data, LIST_TYPE_ORDERED):]
			continue
		}

		// anything else must look like a normal paragraph
		// note: this finds underlined headers, too
		data = data[blockParagraph(out, rndr, data):]
	}

	rndr.nesting--
}

func isPrefixHeader(rndr *render, data []byte) bool {
	if data[0] != '#' {
		return false
	}

	if rndr.flags&EXTENSION_SPACE_HEADERS != 0 {
		level := 0
		for level < len(data) && level < 6 && data[level] == '#' {
			level++
		}
		if level < len(data) && data[level] != ' ' && data[level] != '\t' {
			return false
		}
	}
	return true
}

func blockPrefixHeader(out *bytes.Buffer, rndr *render, data []byte) int {
	level := 0
	for level < len(data) && level < 6 && data[level] == '#' {
		level++
	}
	i, end := 0, 0
	for i = level; i < len(data) && (data[i] == ' ' || data[i] == '\t'); i++ {
	}
	for end = i; end < len(data) && data[end] != '\n'; end++ {
	}
	skip := end
	for end > 0 && data[end-1] == '#' {
		end--
	}
	for end > 0 && (data[end-1] == ' ' || data[end-1] == '\t') {
		end--
	}
	if end > i {
		if rndr.mk.Header != nil {
			work := func() bool {
				parseInline(out, rndr, data[i:end])
				return true
			}
			rndr.mk.Header(out, work, level, rndr.mk.Opaque)
		}
	}
	return skip
}

func isUnderlinedHeader(data []byte) int {
	i := 0

	// test of level 1 header
	if data[i] == '=' {
		for i = 1; i < len(data) && data[i] == '='; i++ {
		}
		for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
			i++
		}
		if i >= len(data) || data[i] == '\n' {
			return 1
		} else {
			return 0
		}
	}

	// test of level 2 header
	if data[i] == '-' {
		for i = 1; i < len(data) && data[i] == '-'; i++ {
		}
		for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
			i++
		}
		if i >= len(data) || data[i] == '\n' {
			return 2
		} else {
			return 0
		}
	}

	return 0
}

func blockHtml(out *bytes.Buffer, rndr *render, data []byte, do_render bool) int {
	var i, j int

	// identify the opening tag
	if len(data) < 2 || data[0] != '<' {
		return 0
	}
	curtag, tagfound := blockHtmlFindTag(data[1:])

	// handle special cases
	if !tagfound {

		// HTML comment, lax form
		if len(data) > 5 && data[1] == '!' && data[2] == '-' && data[3] == '-' {
			i = 5

			for i < len(data) && !(data[i-2] == '-' && data[i-1] == '-' && data[i] == '>') {
				i++
			}
			i++

			if i < len(data) {
				j = isEmpty(data[i:])
			}

			if j > 0 {
				size := i + j
				if do_render && rndr.mk.BlockHtml != nil {
					rndr.mk.BlockHtml(out, data[:size], rndr.mk.Opaque)
				}
				return size
			}
		}

		// HR, which is the only self-closing block tag considered
		if len(data) > 4 &&
			(data[1] == 'h' || data[1] == 'H') &&
			(data[2] == 'r' || data[2] == 'R') {

			i = 3
			for i < len(data) && data[i] != '>' {
				i++
			}

			if i+1 < len(data) {
				i++
				j = isEmpty(data[i:])
				if j > 0 {
					size := i + j
					if do_render && rndr.mk.BlockHtml != nil {
						rndr.mk.BlockHtml(out, data[:size], rndr.mk.Opaque)
					}
					return size
				}
			}
		}

		// no special case recognized
		return 0
	}

	// look for an unindented matching closing tag
	//      followed by a blank line
	i = 1
	found := false

	// if not found, try a second pass looking for indented match
	// but not if tag is "ins" or "del" (following original Markdown.pl)
	if curtag != "ins" && curtag != "del" {
		i = 1
		for i < len(data) {
			i++
			for i < len(data) && !(data[i-1] == '<' && data[i] == '/') {
				i++
			}

			if i+2+len(curtag) >= len(data) {
				break
			}

			j = blockHtmlFindEnd(curtag, rndr, data[i-1:])

			if j > 0 {
				i += j - 1
				found = true
				break
			}
		}
	}

	if !found {
		return 0
	}

	// the end of the block has been found
	if do_render && rndr.mk.BlockHtml != nil {
		rndr.mk.BlockHtml(out, data[:i], rndr.mk.Opaque)
	}

	return i
}

func blockHtmlFindTag(data []byte) (string, bool) {
	i := 0
	for i < len(data) && ((data[i] >= '0' && data[i] <= '9') || (data[i] >= 'A' && data[i] <= 'Z') || (data[i] >= 'a' && data[i] <= 'z')) {
		i++
	}
	if i >= len(data) {
		return "", false
	}
	key := string(data[:i])
	if block_tags[key] {
		return key, true
	}
	return "", false
}

func blockHtmlFindEnd(tag string, rndr *render, data []byte) int {
	// assume data[0] == '<' && data[1] == '/' already tested

	// check if tag is a match
	if len(tag)+3 >= len(data) || bytes.Compare(data[2:2+len(tag)], []byte(tag)) != 0 || data[len(tag)+2] != '>' {
		return 0
	}

	// check white lines
	i := len(tag) + 3
	w := 0
	if i < len(data) {
		if w = isEmpty(data[i:]); w == 0 {
			return 0 // non-blank after tag
		}
	}
	i += w
	w = 0

	if rndr.flags&EXTENSION_LAX_HTML_BLOCKS != 0 {
		if i < len(data) {
			w = isEmpty(data[i:])
		}
	} else {
		if i < len(data) {
			if w = isEmpty(data[i:]); w == 0 {
				return 0 // non-blank line after tag line
			}
		}
	}

	return i + w
}

func isEmpty(data []byte) int {
	var i int
	for i = 0; i < len(data) && data[i] != '\n'; i++ {
		if data[i] != ' ' && data[i] != '\t' {
			return 0
		}
	}
	return i + 1
}

func isHRule(data []byte) bool {
	// skip initial spaces
	if len(data) < 3 {
		return false
	}
	i := 0

	// skip up to three spaces
	for i < 3 && data[i] == ' ' {
		i++
	}

	// look at the hrule char
	if i+2 >= len(data) || (data[i] != '*' && data[i] != '-' && data[i] != '_') {
		return false
	}
	c := data[i]

	// the whole line must be the char or whitespace
	n := 0
	for i < len(data) && data[i] != '\n' {
		switch {
		case data[i] == c:
			n++
		case data[i] != ' ' && data[i] != '\t':
			return false
		}
		i++
	}

	return n >= 3
}

func isFencedCode(data []byte, syntax **string) int {
	i, n := 0, 0

	// skip initial spaces
	if len(data) < 3 {
		return 0
	}
	if data[0] == ' ' {
		i++
		if data[1] == ' ' {
			i++
			if data[2] == ' ' {
				i++
			}
		}
	}

	// look at the hrule char
	if i+2 >= len(data) || !(data[i] == '~' || data[i] == '`') {
		return 0
	}

	c := data[i]

	// the whole line must be the char or whitespace
	for i < len(data) && data[i] == c {
		n++
		i++
	}

	if n < 3 {
		return 0
	}

	if syntax != nil {
		syn := 0

		for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
			i++
		}

		syntax_start := i

		if i < len(data) && data[i] == '{' {
			i++
			syntax_start++

			for i < len(data) && data[i] != '}' && data[i] != '\n' {
				syn++
				i++
			}

			if i == len(data) || data[i] != '}' {
				return 0
			}

			// string all whitespace at the beginning and the end
			// of the {} block
			for syn > 0 && isspace(data[syntax_start]) {
				syntax_start++
				syn--
			}

			for syn > 0 && isspace(data[syntax_start+syn-1]) {
				syn--
			}

			i++
		} else {
			for i < len(data) && !isspace(data[i]) {
				syn++
				i++
			}
		}

		language := string(data[syntax_start : syntax_start+syn])
		*syntax = &language
	}

	for i < len(data) && data[i] != '\n' {
		if !isspace(data[i]) {
			return 0
		}
		i++
	}

	return i + 1
}

func blockFencedCode(out *bytes.Buffer, rndr *render, data []byte) int {
	var lang *string
	beg := isFencedCode(data, &lang)
	if beg == 0 {
		return 0
	}

	var work bytes.Buffer

	for beg < len(data) {
		fence_end := isFencedCode(data[beg:], nil)
		if fence_end != 0 {
			beg += fence_end
			break
		}

		var end int
		for end = beg + 1; end < len(data) && data[end-1] != '\n'; end++ {
		}

		if beg < end {
			// verbatim copy to the working buffer, escaping entities
			if isEmpty(data[beg:]) > 0 {
				work.WriteByte('\n')
			} else {
				work.Write(data[beg:end])
			}
		}
		beg = end
	}

	if work.Len() > 0 && work.Bytes()[work.Len()-1] != '\n' {
		work.WriteByte('\n')
	}

	if rndr.mk.BlockCode != nil {
		syntax := ""
		if lang != nil {
			syntax = *lang
		}

		rndr.mk.BlockCode(out, work.Bytes(), syntax, rndr.mk.Opaque)
	}

	return beg
}

func blockTable(out *bytes.Buffer, rndr *render, data []byte) int {
	var header_work bytes.Buffer
	i, columns, col_data := blockTableHeader(&header_work, rndr, data)
	if i == 0 {
		return 0
	}

	var body_work bytes.Buffer

	for i < len(data) {
		pipes, row_start := 0, i
		for ; i < len(data) && data[i] != '\n'; i++ {
			if data[i] == '|' {
				pipes++
			}
		}

		if pipes == 0 || i == len(data) {
			i = row_start
			break
		}

		blockTableRow(&body_work, rndr, data[row_start:i], columns, col_data)
		i++
	}

	if rndr.mk.Table != nil {
		rndr.mk.Table(out, header_work.Bytes(), body_work.Bytes(), col_data, rndr.mk.Opaque)
	}

	return i
}

func blockTableHeader(out *bytes.Buffer, rndr *render, data []byte) (size int, columns int, column_data []int) {
	i, pipes := 0, 0
	column_data = []int{}
	for i = 0; i < len(data) && data[i] != '\n'; i++ {
		if data[i] == '|' {
			pipes++
		}
	}

	if i == len(data) || pipes == 0 {
		return 0, 0, column_data
	}

	header_end := i

	if data[0] == '|' {
		pipes--
	}

	if i > 2 && data[i-1] == '|' {
		pipes--
	}

	columns = pipes + 1
	column_data = make([]int, columns)

	// parse the header underline
	i++
	if i < len(data) && data[i] == '|' {
		i++
	}

	under_end := i
	for under_end < len(data) && data[under_end] != '\n' {
		under_end++
	}

	col := 0
	for ; col < columns && i < under_end; col++ {
		dashes := 0

		for i < under_end && (data[i] == ' ' || data[i] == '\t') {
			i++
		}

		if data[i] == ':' {
			i++
			column_data[col] |= TABLE_ALIGNMENT_LEFT
			dashes++
		}

		for i < under_end && data[i] == '-' {
			i++
			dashes++
		}

		if i < under_end && data[i] == ':' {
			i++
			column_data[col] |= TABLE_ALIGNMENT_RIGHT
			dashes++
		}

		for i < under_end && (data[i] == ' ' || data[i] == '\t') {
			i++
		}

		if i < under_end && data[i] != '|' {
			break
		}

		if dashes < 3 {
			break
		}

		i++
	}

	if col < columns {
		return 0, 0, column_data
	}

	blockTableRow(out, rndr, data[:header_end], columns, column_data)
	size = under_end + 1
	return
}

func blockTableRow(out *bytes.Buffer, rndr *render, data []byte, columns int, col_data []int) {
	i, col := 0, 0
	var row_work bytes.Buffer

	if i < len(data) && data[i] == '|' {
		i++
	}

	for col = 0; col < columns && i < len(data); col++ {
		for i < len(data) && isspace(data[i]) {
			i++
		}

		cell_start := i

		for i < len(data) && data[i] != '|' {
			i++
		}

		cell_end := i - 1

		for cell_end > cell_start && isspace(data[cell_end]) {
			cell_end--
		}

		var cell_work bytes.Buffer
		parseInline(&cell_work, rndr, data[cell_start:cell_end+1])

		if rndr.mk.TableCell != nil {
			cdata := 0
			if col < len(col_data) {
				cdata = col_data[col]
			}
			rndr.mk.TableCell(&row_work, cell_work.Bytes(), cdata, rndr.mk.Opaque)
		}

		i++
	}

	for ; col < columns; col++ {
		empty_cell := []byte{}
		if rndr.mk.TableCell != nil {
			cdata := 0
			if col < len(col_data) {
				cdata = col_data[col]
			}
			rndr.mk.TableCell(&row_work, empty_cell, cdata, rndr.mk.Opaque)
		}
	}

	if rndr.mk.TableRow != nil {
		rndr.mk.TableRow(out, row_work.Bytes(), rndr.mk.Opaque)
	}
}

// returns blockquote prefix length
func blockQuotePrefix(data []byte) int {
	i := 0
	for i < len(data) && i < 3 && data[i] == ' ' {
		i++
	}
	if i < len(data) && data[i] == '>' {
		if i+1 < len(data) && (data[i+1] == ' ' || data[i+1] == '\t') {
			return i + 2
		}
		return i + 1
	}
	return 0
}

// parse a blockquote fragment
func blockQuote(out *bytes.Buffer, rndr *render, data []byte) int {
	var block bytes.Buffer
	var work bytes.Buffer
	beg, end := 0, 0
	for beg < len(data) {
		for end = beg + 1; end < len(data) && data[end-1] != '\n'; end++ {
		}

		if pre := blockQuotePrefix(data[beg:]); pre > 0 {
			beg += pre // skip prefix
		} else {
			// empty line followed by non-quote line
			if isEmpty(data[beg:]) > 0 && (end >= len(data) || (blockQuotePrefix(data[end:]) == 0 && isEmpty(data[end:]) == 0)) {
				break
			}
		}

		if beg < end { // copy into the in-place working buffer
			work.Write(data[beg:end])
		}
		beg = end
	}

	parseBlock(&block, rndr, work.Bytes())
	if rndr.mk.BlockQuote != nil {
		rndr.mk.BlockQuote(out, block.Bytes(), rndr.mk.Opaque)
	}
	return end
}

// returns prefix length for block code
func blockCodePrefix(data []byte) int {
	if len(data) > 0 && data[0] == '\t' {
		return 1
	}
	if len(data) > 3 && data[0] == ' ' && data[1] == ' ' && data[2] == ' ' && data[3] == ' ' {
		return 4
	}
	return 0
}

func blockCode(out *bytes.Buffer, rndr *render, data []byte) int {
	var work bytes.Buffer

	beg, end := 0, 0
	for beg < len(data) {
		for end = beg + 1; end < len(data) && data[end-1] != '\n'; end++ {
		}

		if pre := blockCodePrefix(data[beg:end]); pre > 0 {
			beg += pre
		} else {
			if isEmpty(data[beg:end]) == 0 {
				// non-empty non-prefixed line breaks the pre
				break
			}
		}

		if beg < end {
			// verbatim copy to the working buffer, escaping entities
			if isEmpty(data[beg:end]) > 0 {
				work.WriteByte('\n')
			} else {
				work.Write(data[beg:end])
			}
		}
		beg = end
	}

	// trim all the \n off the end of work
	workbytes := work.Bytes()
	n := 0
	for len(workbytes) > n && workbytes[len(workbytes)-n-1] == '\n' {
		n++
	}
	if n > 0 {
		work.Truncate(len(workbytes) - n)
	}

	work.WriteByte('\n')

	if rndr.mk.BlockCode != nil {
		rndr.mk.BlockCode(out, work.Bytes(), "", rndr.mk.Opaque)
	}

	return beg
}

// returns unordered list item prefix
func blockUliPrefix(data []byte) int {
	i := 0

	// start with up to 3 spaces
	for i < len(data) && i < 3 && data[i] == ' ' {
		i++
	}

	// need a *, +, or - followed by a space/tab
	if i+1 >= len(data) ||
		(data[i] != '*' && data[i] != '+' && data[i] != '-') ||
		(data[i+1] != ' ' && data[i+1] != '\t') {
		return 0
	}
	return i + 2
}

// returns ordered list item prefix
func blockOliPrefix(data []byte) int {
	i := 0

	// start with up to 3 spaces
	for i < len(data) && i < 3 && data[i] == ' ' {
		i++
	}

	// count the digits
	start := i
	for i < len(data) && data[i] >= '0' && data[i] <= '9' {
		i++
	}

	// we need >= 1 digits followed by a dot and a space/tab
	if start == i || data[i] != '.' || i+1 >= len(data) ||
		(data[i+1] != ' ' && data[i+1] != '\t') {
		return 0
	}
	return i + 2
}

// parse ordered or unordered list block
func blockList(out *bytes.Buffer, rndr *render, data []byte, flags int) int {
	i := 0
	work := func() bool {
		j := 0
		for i < len(data) {
			j = blockListItem(out, rndr, data[i:], &flags)
			i += j

			if j == 0 || flags&LIST_ITEM_END_OF_LIST != 0 {
				break
			}
		}
		return true
	}

	if rndr.mk.List != nil {
		rndr.mk.List(out, work, flags, rndr.mk.Opaque)
	}
	return i
}

// parse a single list item
// assumes initial prefix is already removed
func blockListItem(out *bytes.Buffer, rndr *render, data []byte, flags *int) int {
	// keep track of the first indentation prefix
	beg, end, pre, sublist, orgpre, i := 0, 0, 0, 0, 0, 0

	for orgpre < 3 && orgpre < len(data) && data[orgpre] == ' ' {
		orgpre++
	}

	beg = blockUliPrefix(data)
	if beg == 0 {
		beg = blockOliPrefix(data)
	}
	if beg == 0 {
		return 0
	}

	// skip leading whitespace on first line
	for beg < len(data) && data[beg] == ' ' {
		beg++
	}

	// skip to the beginning of the following line
	end = beg
	for end < len(data) && data[end-1] != '\n' {
		end++
	}

	// get working buffers
	var work bytes.Buffer
	var inter bytes.Buffer

	// put the first line into the working buffer
	work.Write(data[beg:end])
	beg = end

	// process the following lines
	contains_blank_line, contains_block := false, false
	for beg < len(data) {
		end++

		for end < len(data) && data[end-1] != '\n' {
			end++
		}

		// process an empty line
		if isEmpty(data[beg:end]) > 0 {
			contains_blank_line = true
			beg = end
			continue
		}

		// calculate the indentation
		i = 0
		for i < 4 && beg+i < end && data[beg+i] == ' ' {
			i++
		}

		pre = i
		if data[beg] == '\t' {
			i = 1
			pre = 8
		}

		chunk := data[beg+i : end]

		// check for a nested list item
		if (blockUliPrefix(chunk) > 0 && !isHRule(chunk)) || blockOliPrefix(chunk) > 0 {
			if contains_blank_line {
				contains_block = true
			}

			if pre == orgpre { // the following item must have the same indentation
				break
			}

			if sublist == 0 {
				sublist = work.Len()
			}
		} else {
			// how about a nested prefix header?
			if isPrefixHeader(rndr, chunk) {
				// only nest headers that are indented
				if contains_blank_line && i < 4 && data[beg] != '\t' {
					*flags |= LIST_ITEM_END_OF_LIST
					break
				}
				contains_block = true
			} else {
				// only join stuff after empty lines when indented
				if contains_blank_line && i < 4 && data[beg] != '\t' {
					*flags |= LIST_ITEM_END_OF_LIST
					break
				} else {
					if contains_blank_line {
						work.WriteByte('\n')
						contains_block = true
					}
				}
			}
		}

		contains_blank_line = false

		// add the line into the working buffer without prefix
		work.Write(data[beg+i : end])
		beg = end
	}

	// render li contents
	if contains_block {
		*flags |= LIST_ITEM_CONTAINS_BLOCK
	}

	workbytes := work.Bytes()
	if *flags&LIST_ITEM_CONTAINS_BLOCK != 0 {
		// intermediate render of block li
		if sublist > 0 && sublist < len(workbytes) {
			parseBlock(&inter, rndr, workbytes[:sublist])
			parseBlock(&inter, rndr, workbytes[sublist:])
		} else {
			parseBlock(&inter, rndr, workbytes)
		}
	} else {
		// intermediate render of inline li
		if sublist > 0 && sublist < len(workbytes) {
			parseInline(&inter, rndr, workbytes[:sublist])
			parseBlock(&inter, rndr, workbytes[sublist:])
		} else {
			parseInline(&inter, rndr, workbytes)
		}
	}

	// render li itself
	if rndr.mk.ListItem != nil {
		rndr.mk.ListItem(out, inter.Bytes(), *flags, rndr.mk.Opaque)
	}

	return beg
}

// render a single paragraph that has already been parsed out
func renderParagraph(out *bytes.Buffer, rndr *render, data []byte) {
	// trim leading whitespace
	beg := 0
	for beg < len(data) && isspace(data[beg]) {
		beg++
	}

	// trim trailing whitespace
	end := len(data)
	for end > beg && isspace(data[end-1]) {
		end--
	}
	if end == beg || rndr.mk.Paragraph == nil {
		return
	}

	work := func() bool {
		parseInline(out, rndr, data[beg:end])
		return true
	}
	rndr.mk.Paragraph(out, work, rndr.mk.Opaque)
}

func blockParagraph(out *bytes.Buffer, rndr *render, data []byte) int {
	// prev: index of 1st char of previous line
	// line: index of 1st char of current line
	// i: index of cursor/end of current line
	var prev, line, i int

	// keep going until we find something to mark the end of the paragraph
	for i < len(data) {
		// mark the beginning of the current line
		prev = line
		current := data[i:]
		line = i

		// did we find a blank line marking the end of the paragraph?
		if n := isEmpty(current); n > 0 {
			renderParagraph(out, rndr, data[:i])
			return i + n
		}

		// an underline under some text marks a header, so our paragraph ended on prev line
		if i > 0 && rndr.mk.Header != nil {
			if level := isUnderlinedHeader(current); level > 0 {
				// render the paragraph
				renderParagraph(out, rndr, data[:prev])

				// render the header
				// this ugly, convoluted closure avoids forcing variables onto the heap
				work := func(o *bytes.Buffer, r *render, d []byte) func() bool {
					return func() bool {
						parseInline(o, r, d)
						return true
					}
				}(out, rndr, data[prev:i-1])
				rndr.mk.Header(out, work, level, rndr.mk.Opaque)

				// find the end of the underline
				for ; i < len(data) && data[i] != '\n'; i++ {
				}
				return i
			}
		}

		// if the next line starts a block of HTML, then the paragraph ends here
		if rndr.flags&EXTENSION_LAX_HTML_BLOCKS != 0 {
			if data[i] == '<' && rndr.mk.BlockHtml != nil && blockHtml(out, rndr, current, false) > 0 {
				// rewind to before the HTML block
				renderParagraph(out, rndr, data[:i])
				return i
			}
		}

		// if there's a prefixed header or a horizontal rule after this, paragraph is over
		if isPrefixHeader(rndr, current) || isHRule(current) {
			renderParagraph(out, rndr, data[:i])
			return i
		}

		// otherwise, scan to the beginning of the next line
		i++
		for i < len(data) && data[i-1] != '\n' {
			i++
		}
	}

	renderParagraph(out, rndr, data[:i])
	return i
}
