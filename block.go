//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
// Functions to parse block-level elements.
//

package blackfriday

import (
	"bytes"
)

// Parse block-level data.
// Note: this function and many that it calls assume that
// the input buffer ends with a newline.
func (parser *Parser) parseBlock(out *bytes.Buffer, data []byte) {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		panic("parseBlock input is missing terminating newline")
	}

	// this is called recursively: enforce a maximum depth
	if parser.nesting >= parser.maxNesting {
		return
	}
	parser.nesting++

	// parse out one block-level construct at a time
	for len(data) > 0 {
		// prefixed header:
		//
		// # Header 1
		// ## Header 2
		// ...
		// ###### Header 6
		if parser.isPrefixHeader(data) {
			data = data[parser.blockPrefixHeader(out, data):]
			continue
		}

		// block of preformatted HTML:
		//
		// <div>
		//     ...
		// </div>
		if data[0] == '<' {
			if i := parser.blockHtml(out, data, true); i > 0 {
				data = data[i:]
				continue
			}
		}

		// blank lines.  note: returns the # of bytes to skip
		if i := parser.isEmpty(data); i > 0 {
			data = data[i:]
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
		if parser.blockCodePrefix(data) > 0 {
			data = data[parser.blockCode(out, data):]
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
		if parser.flags&EXTENSION_FENCED_CODE != 0 {
			if i := parser.blockFencedCode(out, data); i > 0 {
				data = data[i:]
				continue
			}
		}

		// horizontal rule:
		//
		// ------
		// or
		// ******
		// or
		// ______
		if parser.isHRule(data) {
			parser.r.HRule(out)
			var i int
			for i = 0; data[i] != '\n'; i++ {
			}
			data = data[i:]
			continue
		}

		// block quote:
		//
		// > A big quote I found somewhere
		// > on the web
		if parser.blockQuotePrefix(data) > 0 {
			data = data[parser.blockQuote(out, data):]
			continue
		}

		// table:
		//
		// Name  | Age | Phone
		// ------|-----|---------
		// Bob   | 31  | 555-1234
		// Alice | 27  | 555-4321
		if parser.flags&EXTENSION_TABLES != 0 {
			if i := parser.blockTable(out, data); i > 0 {
				data = data[i:]
				continue
			}
		}

		// an itemized/unordered list:
		//
		// * Item 1
		// * Item 2
		//
		// also works with + or -
		if parser.blockUliPrefix(data) > 0 {
			data = data[parser.blockList(out, data, 0):]
			continue
		}

		// a numbered/ordered list:
		//
		// 1. Item 1
		// 2. Item 2
		if parser.blockOliPrefix(data) > 0 {
			data = data[parser.blockList(out, data, LIST_TYPE_ORDERED):]
			continue
		}

		// anything else must look like a normal paragraph
		// note: this finds underlined headers, too
		data = data[parser.blockParagraph(out, data):]
	}

	parser.nesting--
}

func (parser *Parser) isPrefixHeader(data []byte) bool {
	if data[0] != '#' {
		return false
	}

	if parser.flags&EXTENSION_SPACE_HEADERS != 0 {
		level := 0
		for level < 6 && data[level] == '#' {
			level++
		}
		if data[level] != ' ' {
			return false
		}
	}
	return true
}

func (parser *Parser) blockPrefixHeader(out *bytes.Buffer, data []byte) int {
	level := 0
	for level < 6 && data[level] == '#' {
		level++
	}
	i, end := 0, 0
	for i = level; data[i] == ' '; i++ {
	}
	for end = i; data[end] != '\n'; end++ {
	}
	skip := end
	for end > 0 && data[end-1] == '#' {
		end--
	}
	for end > 0 && data[end-1] == ' ' {
		end--
	}
	if end > i {
		work := func() bool {
			parser.parseInline(out, data[i:end])
			return true
		}
		parser.r.Header(out, work, level)
	}
	return skip
}

func (parser *Parser) isUnderlinedHeader(data []byte) int {
	// test of level 1 header
	if data[0] == '=' {
		i := 1
		for data[i] == '=' {
			i++
		}
		for data[i] == ' ' {
			i++
		}
		if data[i] == '\n' {
			return 1
		} else {
			return 0
		}
	}

	// test of level 2 header
	if data[0] == '-' {
		i := 1
		for data[i] == '-' {
			i++
		}
		for data[i] == ' ' {
			i++
		}
		if data[i] == '\n' {
			return 2
		} else {
			return 0
		}
	}

	return 0
}

func (parser *Parser) blockHtml(out *bytes.Buffer, data []byte, doRender bool) int {
	var i, j int

	// identify the opening tag
	if data[0] != '<' {
		return 0
	}
	curtag, tagfound := parser.blockHtmlFindTag(data[1:])

	// handle special cases
	if !tagfound {
		// check for an HTML comment
		if size := parser.blockHtmlComment(out, data, doRender); size > 0 {
			return size
		}

		// check for an <hr> tag
		if size := parser.blockHtmlHr(out, data, doRender); size > 0 {
			return size
		}

		// no special case recognized
		return 0
	}

	// look for an unindented matching closing tag
	// followed by a blank line
	found := false
	/*
		closetag := []byte("\n</" + curtag + ">")
		j = len(curtag) + 1
		for !found {
			// scan for a closing tag at the beginning of a line
			if skip := bytes.Index(data[j:], closetag); skip >= 0 {
				j += skip + len(closetag)
			} else {
				break
			}

			// see if it is the only thing on the line
			if skip := parser.isEmpty(data[j:]); skip > 0 {
				// see if it is followed by a blank line/eof
				j += skip
				if j >= len(data) {
					found = true
					i = j
				} else {
					if skip := parser.isEmpty(data[j:]); skip > 0 {
						j += skip
						found = true
						i = j
					}
				}
			}
		}
	*/

	// if not found, try a second pass looking for indented match
	// but not if tag is "ins" or "del" (following original Markdown.pl)
	if !found && curtag != "ins" && curtag != "del" {
		i = 1
		for i < len(data) {
			i++
			for i < len(data) && !(data[i-1] == '<' && data[i] == '/') {
				i++
			}

			if i+2+len(curtag) >= len(data) {
				break
			}

			j = parser.blockHtmlFindEnd(curtag, data[i-1:])

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
	if doRender {
		// trim newlines
		end := i
		for end > 0 && data[end-1] == '\n' {
			end--
		}
		parser.r.BlockHtml(out, data[:end])
	}

	return i
}

// HTML comment, lax form
func (parser *Parser) blockHtmlComment(out *bytes.Buffer, data []byte, doRender bool) int {
	if data[0] != '<' || data[1] != '!' || data[2] != '-' || data[3] != '-' {
		return 0
	}

	i := 5

	// scan for an end-of-comment marker, across lines if necessary
	for i < len(data) && !(data[i-2] == '-' && data[i-1] == '-' && data[i] == '>') {
		i++
	}
	i++

	// no end-of-comment marker
	if i >= len(data) {
		return 0
	}

	// needs to end with a blank line
	if j := parser.isEmpty(data[i:]); j > 0 {
		size := i + j
		if doRender {
			// trim trailing newlines
			end := size
			for end > 0 && data[end-1] == '\n' {
				end--
			}
			parser.r.BlockHtml(out, data[:end])
		}
		return size
	}

	return 0
}

// HR, which is the only self-closing block tag considered
func (parser *Parser) blockHtmlHr(out *bytes.Buffer, data []byte, doRender bool) int {
	if data[0] != '<' || (data[1] != 'h' && data[1] != 'H') || (data[2] != 'r' && data[2] != 'R') {
		return 0
	}
	if data[3] != ' ' && data[3] != '/' && data[3] != '>' {
		// not an <hr> tag after all; at least not a valid one
		return 0
	}

	i := 3
	for data[i] != '>' && data[i] != '\n' {
		i++
	}

	if data[i] == '>' {
		i++
		if j := parser.isEmpty(data[i:]); j > 0 {
			size := i + j
			if doRender {
				// trim newlines
				end := size
				for end > 0 && data[end-1] == '\n' {
					end--
				}
				parser.r.BlockHtml(out, data[:end])
			}
			return size
		}
	}

	return 0
}

func (parser *Parser) blockHtmlFindTag(data []byte) (string, bool) {
	i := 0
	for isalnum(data[i]) {
		i++
	}
	key := string(data[:i])
	if blockTags[key] {
		return key, true
	}
	return "", false
}

func (parser *Parser) blockHtmlFindEnd(tag string, data []byte) int {
	// assume data[0] == '<' && data[1] == '/' already tested

	// check if tag is a match
	closetag := []byte("</" + tag + ">")
	if !bytes.HasPrefix(data, closetag) {
		return 0
	}
	i := len(closetag)

	// check that the rest of the line is blank
	skip := 0
	if skip = parser.isEmpty(data[i:]); skip == 0 {
		return 0
	}
	i += skip
	skip = 0

	if i >= len(data) {
		return i
	}

	if parser.flags&EXTENSION_LAX_HTML_BLOCKS != 0 {
		return i
	}
	if skip = parser.isEmpty(data[i:]); skip == 0 {
		// following line must be blank
		return 0
	}

	return i + skip
}

func (parser *Parser) isEmpty(data []byte) int {
	// it is okay to call isEmpty on an empty buffer
	if len(data) == 0 {
		return 0
	}

	var i int
	for i = 0; data[i] != '\n'; i++ {
		if data[i] != ' ' {
			return 0
		}
	}
	return i + 1
}

func (parser *Parser) isHRule(data []byte) bool {
	i := 0

	// skip up to three spaces
	for i < 3 && data[i] == ' ' {
		i++
	}

	// look at the hrule char
	if data[i] != '*' && data[i] != '-' && data[i] != '_' {
		return false
	}
	c := data[i]

	// the whole line must be the char or whitespace
	n := 0
	for data[i] != '\n' {
		switch {
		case data[i] == c:
			n++
		case data[i] != ' ':
			return false
		}
		i++
	}

	return n >= 3
}

func (parser *Parser) isFencedCode(data []byte, syntax **string, oldmarker string) (skip int, marker string) {
	i, size := 0, 0
	skip = 0

	// skip up to three spaces
	for i < 3 && data[i] == ' ' {
		i++
	}

	// check for the marker characters: ~ or `
	if data[i] != '~' && data[i] != '`' {
		return
	}

	c := data[i]

	// the whole line must be the same char or whitespace
	for data[i] == c {
		size++
		i++
	}

	// the marker char must occur at least 3 times
	if size < 3 {
		return
	}
	marker = string(data[i-size : i])

	// if this is the end marker, it must match the beginning marker
	if oldmarker != "" && marker != oldmarker {
		return
	}

	if syntax != nil {
		syn := 0

		for data[i] == ' ' {
			i++
		}

		syntaxStart := i

		if data[i] == '{' {
			i++
			syntaxStart++

			for data[i] != '}' && data[i] != '\n' {
				syn++
				i++
			}

			if data[i] != '}' {
				return
			}

			// strip all whitespace at the beginning and the end
			// of the {} block
			for syn > 0 && isspace(data[syntaxStart]) {
				syntaxStart++
				syn--
			}

			for syn > 0 && isspace(data[syntaxStart+syn-1]) {
				syn--
			}

			i++
		} else {
			for !isspace(data[i]) {
				syn++
				i++
			}
		}

		language := string(data[syntaxStart : syntaxStart+syn])
		*syntax = &language
	}

	for data[i] == ' ' {
		i++
	}
	if data[i] != '\n' {
		return
	}

	skip = i + 1
	return
}

func (parser *Parser) blockFencedCode(out *bytes.Buffer, data []byte) int {
	var lang *string
	beg, marker := parser.isFencedCode(data, &lang, "")
	if beg == 0 || beg >= len(data) {
		return 0
	}

	var work bytes.Buffer

	for {
		// safe to assume beg < len(data)

		// check for the end of the code block
		fenceEnd, _ := parser.isFencedCode(data[beg:], nil, marker)
		if fenceEnd != 0 {
			beg += fenceEnd
			break
		}

		// copy the current line
		end := beg
		for data[end] != '\n' {
			end++
		}
		end++

		// did we reach the end of the buffer without a closing marker?
		if end >= len(data) {
			return 0
		}

		// verbatim copy to the working buffer
		work.Write(data[beg:end])
		beg = end
	}

	syntax := ""
	if lang != nil {
		syntax = *lang
	}

	parser.r.BlockCode(out, work.Bytes(), syntax)

	return beg
}

func (parser *Parser) blockTable(out *bytes.Buffer, data []byte) int {
	var header bytes.Buffer
	i, columns := parser.blockTableHeader(&header, data)
	if i == 0 {
		return 0
	}

	var body bytes.Buffer

	for i < len(data) {
		pipes, rowStart := 0, i
		for ; data[i] != '\n'; i++ {
			if data[i] == '|' {
				pipes++
			}
		}

		if pipes == 0 {
			i = rowStart
			break
		}

		// include the newline in data sent to blockTableRow
		i++
		parser.blockTableRow(&body, data[rowStart:i], columns)
	}

	parser.r.Table(out, header.Bytes(), body.Bytes(), columns)

	return i
}

func (parser *Parser) blockTableHeader(out *bytes.Buffer, data []byte) (size int, columns []int) {
	i := 0
	colCount := 1
	for i = 0; data[i] != '\n'; i++ {
		if data[i] == '|' {
			colCount++
		}
	}

	// doesn't look like a table header
	if colCount == 1 {
		return
	}

	// include the newline in the data sent to blockTableRow
	header := data[:i+1]

	// column count ignores pipes at beginning or end of line
	if data[0] == '|' {
		colCount--
	}
	if i > 2 && data[i-1] == '|' {
		colCount--
	}

	columns = make([]int, colCount)

	// move on to the header underline
	i++
	if i >= len(data) {
		return
	}

	if data[i] == '|' {
		i++
	}
	for data[i] == ' ' {
		i++
	}

	// each column header is of form: / *:?-+:? *|/ with # dashes + # colons >= 3
	// and trailing | optional on last column
	col := 0
	for data[i] != '\n' {
		dashes := 0

		if data[i] == ':' {
			i++
			columns[col] |= TABLE_ALIGNMENT_LEFT
			dashes++
		}
		for data[i] == '-' {
			i++
			dashes++
		}
		if data[i] == ':' {
			i++
			columns[col] |= TABLE_ALIGNMENT_RIGHT
			dashes++
		}
		for data[i] == ' ' {
			i++
		}

		// end of column test is messy
		switch {
		case dashes < 3:
			// not a valid column
			return

		case data[i] == '|':
			// marker found, now skip past trailing whitespace
			col++
			i++
			for data[i] == ' ' {
				i++
			}

		case data[i] != '|' && col+1 < colCount:
			// something else found where marker was required
			return

		case data[i] == '\n':
			// marker is optional for the last column
			col++

		default:
			// trailing junk found after last column
			return
		}
	}
	if col != colCount {
		return
	}

	parser.blockTableRow(out, header, columns)
	size = i + 1
	return
}

func (parser *Parser) blockTableRow(out *bytes.Buffer, data []byte, columns []int) {
	i, col := 0, 0
	var rowWork bytes.Buffer

	if data[i] == '|' {
		i++
	}

	for col = 0; col < len(columns) && data[i] != '\n'; col++ {
		for data[i] == ' ' {
			i++
		}

		cellStart := i

		for data[i] != '|' && data[i] != '\n' {
			i++
		}

		cellEnd := i
		i++

		for cellEnd > cellStart && data[cellEnd-1] == ' ' {
			cellEnd--
		}

		var cellWork bytes.Buffer
		parser.parseInline(&cellWork, data[cellStart:cellEnd])
		parser.r.TableCell(&rowWork, cellWork.Bytes(), columns[col])
	}

	// pad it out with empty columns to get the right number
	for ; col < len(columns); col++ {
		parser.r.TableCell(&rowWork, nil, columns[col])
	}

	// silently ignore rows with too many cells

	parser.r.TableRow(out, rowWork.Bytes())
}

// returns blockquote prefix length
func (parser *Parser) blockQuotePrefix(data []byte) int {
	i := 0
	for i < 3 && data[i] == ' ' {
		i++
	}
	if data[i] == '>' {
		if data[i+1] == ' ' {
			return i + 2
		}
		return i + 1
	}
	return 0
}

// parse a blockquote fragment
func (parser *Parser) blockQuote(out *bytes.Buffer, data []byte) int {
	var raw bytes.Buffer
	beg, end := 0, 0
	for beg < len(data) {
		end = beg
		for data[end] != '\n' {
			end++
		}
		end++

		if pre := parser.blockQuotePrefix(data[beg:]); pre > 0 {
			// string the prefix
			beg += pre
		} else {
			// blockquote ends with at least one blank line
			// followed by something without a blockquote prefix
			if parser.isEmpty(data[beg:]) > 0 &&
				(end >= len(data) ||
					(parser.blockQuotePrefix(data[end:]) == 0 && parser.isEmpty(data[end:]) == 0)) {
				break
			}
		}

		// this line is part of the blockquote
		raw.Write(data[beg:end])
		beg = end
	}

	var cooked bytes.Buffer
	parser.parseBlock(&cooked, raw.Bytes())
	parser.r.BlockQuote(out, cooked.Bytes())
	return end
}

// returns prefix length for block code
func (parser *Parser) blockCodePrefix(data []byte) int {
	if data[0] == ' ' && data[1] == ' ' && data[2] == ' ' && data[3] == ' ' {
		return 4
	}
	return 0
}

func (parser *Parser) blockCode(out *bytes.Buffer, data []byte) int {
	var work bytes.Buffer

	i := 0
	for i < len(data) {
		beg := i
		for data[i] != '\n' {
			i++
		}
		i++

		blankline := parser.isEmpty(data[beg:i]) > 0
		if pre := parser.blockCodePrefix(data[beg:i]); pre > 0 {
			beg += pre
		} else {
			if !blankline {
				// non-empty, non-prefixed line breaks the pre
				i = beg
				break
			}
		}

		// verbatim copy to the working buffeu
		if blankline {
			work.WriteByte('\n')
		} else {
			work.Write(data[beg:i])
		}
	}

	// trim all the \n off the end of work
	workbytes := work.Bytes()
	eol := len(workbytes)
	for eol > 0 && workbytes[eol-1] == '\n' {
		eol--
	}
	if eol != len(workbytes) {
		work.Truncate(eol)
	}

	work.WriteByte('\n')

	parser.r.BlockCode(out, work.Bytes(), "")

	return i
}

// returns unordered list item prefix
func (parser *Parser) blockUliPrefix(data []byte) int {
	i := 0

	// start with up to 3 spaces
	for i < 3 && data[i] == ' ' {
		i++
	}

	// need a *, +, or - followed by a space/tab
	if (data[i] != '*' && data[i] != '+' && data[i] != '-') ||
		data[i+1] != ' ' {
		return 0
	}
	return i + 2
}

// returns ordered list item prefix
func (parser *Parser) blockOliPrefix(data []byte) int {
	i := 0

	// start with up to 3 spaces
	for i < 3 && data[i] == ' ' {
		i++
	}

	// count the digits
	start := i
	for data[i] >= '0' && data[i] <= '9' {
		i++
	}

	// we need >= 1 digits followed by a dot and a space/tab
	if start == i || data[i] != '.' || data[i+1] != ' ' {
		return 0
	}
	return i + 2
}

// parse ordered or unordered list block
func (parser *Parser) blockList(out *bytes.Buffer, data []byte, flags int) int {
	i := 0
	flags |= LIST_ITEM_BEGINNING_OF_LIST
	work := func() bool {
		for i < len(data) {
			skip := parser.blockListItem(out, data[i:], &flags)
			i += skip

			if skip == 0 || flags&LIST_ITEM_END_OF_LIST != 0 {
				break
			}
			flags &= ^LIST_ITEM_BEGINNING_OF_LIST
		}
		return true
	}

	parser.r.List(out, work, flags)
	return i
}

// Parse a single list item.
// Assumes initial prefix is already removed if this is a sublist.
func (parser *Parser) blockListItem(out *bytes.Buffer, data []byte, flags *int) int {
	// keep track of the indentation of the first line
	itemIndent := 0
	for itemIndent < 3 && data[itemIndent] == ' ' {
		itemIndent++
	}

	i := parser.blockUliPrefix(data)
	if i == 0 {
		i = parser.blockOliPrefix(data)
	}
	if i == 0 {
		return 0
	}

	// skip leading whitespace on first line
	for data[i] == ' ' {
		i++
	}

	// find the end of the line
	line := i
	for data[i-1] != '\n' {
		i++
	}

	// get working buffer
	var raw bytes.Buffer

	// put the first line into the working buffer
	raw.Write(data[line:i])
	line = i

	// process the following lines
	containsBlankLine := false
	sublist := 0

loop:
	for line < len(data) {
		i++

		// find the end of this line
		for data[i-1] != '\n' {
			i++
		}

		// if it is an empty line, guess that it is part of this item
		// and move on to the next line
		if parser.isEmpty(data[line:i]) > 0 {
			containsBlankLine = true
			line = i
			continue
		}

		// calculate the indentation
		indent := 0
		for indent < 4 && line+indent < i && data[line+indent] == ' ' {
			indent++
		}

		chunk := data[line+indent : i]

		// evaluate how this line fits in
		switch {
		// is this a nested list item?
		case (parser.blockUliPrefix(chunk) > 0 && !parser.isHRule(chunk)) ||
			parser.blockOliPrefix(chunk) > 0:

			if containsBlankLine {
				*flags |= LIST_ITEM_CONTAINS_BLOCK
			}

			// to be a nested list, it must be indented more
			// if not, it is the next item in the same list
			if indent <= itemIndent {
				break loop
			}

			// is this the first item in the the nested list?
			if sublist == 0 {
				sublist = raw.Len()
			}

		// is this a nested prefix header?
		case parser.isPrefixHeader(chunk):
			// if the header is not indented, it is not nested in the list
			// and thus ends the list
			if containsBlankLine && indent < 4 {
				*flags |= LIST_ITEM_END_OF_LIST
				break loop
			}
			*flags |= LIST_ITEM_CONTAINS_BLOCK

		// anything following an empty line is only part
		// of this item if it is indented 4 spaces
		// (regardless of the indentation of the beginning of the item)
		case containsBlankLine && indent < 4:
			*flags |= LIST_ITEM_END_OF_LIST
			break loop

		// a blank line means this should be parsed as a block
		case containsBlankLine:
			raw.WriteByte('\n')
			*flags |= LIST_ITEM_CONTAINS_BLOCK
		}

		containsBlankLine = false

		// add the line into the working buffer without prefix
		raw.Write(data[line+indent : i])
		line = i
	}

	// render the contents of the list item
	rawBytes := raw.Bytes()
	var cooked bytes.Buffer
	if *flags&LIST_ITEM_CONTAINS_BLOCK != 0 {
		// intermediate render of block li
		if sublist > 0 {
			parser.parseBlock(&cooked, rawBytes[:sublist])
			parser.parseBlock(&cooked, rawBytes[sublist:])
		} else {
			parser.parseBlock(&cooked, rawBytes)
		}
	} else {
		// intermediate render of inline li
		if sublist > 0 {
			parser.parseInline(&cooked, rawBytes[:sublist])
			parser.parseBlock(&cooked, rawBytes[sublist:])
		} else {
			parser.parseInline(&cooked, rawBytes)
		}
	}

	// render the actual list item
	cookedBytes := cooked.Bytes()
	parsedEnd := len(cookedBytes)

	// strip trailing newlines
	for parsedEnd > 0 && cookedBytes[parsedEnd-1] == '\n' {
		parsedEnd--
	}
	parser.r.ListItem(out, cookedBytes[:parsedEnd], *flags)

	return line
}

// render a single paragraph that has already been parsed out
func (parser *Parser) renderParagraph(out *bytes.Buffer, data []byte) {
	if len(data) == 0 {
		return
	}

	// trim leading spaces
	beg := 0
	for data[beg] == ' ' {
		beg++
	}

	// trim trailing newline
	end := len(data) - 1

	// trim trailing spaces
	for end > beg && data[end-1] == ' ' {
		end--
	}

	work := func() bool {
		parser.parseInline(out, data[beg:end])
		return true
	}
	parser.r.Paragraph(out, work)
}

func (parser *Parser) blockParagraph(out *bytes.Buffer, data []byte) int {
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
		if n := parser.isEmpty(current); n > 0 {
			parser.renderParagraph(out, data[:i])
			return i + n
		}

		// an underline under some text marks a header, so our paragraph ended on prev line
		if i > 0 {
			if level := parser.isUnderlinedHeader(current); level > 0 {
				// render the paragraph
				parser.renderParagraph(out, data[:prev])

				// ignore leading and trailing whitespace
				eol := i - 1
				for prev < eol && data[prev] == ' ' {
					prev++
				}
				for eol > prev && data[eol-1] == ' ' {
					eol--
				}

				// render the header
				// this ugly double closure avoids forcing variables onto the heap
				work := func(o *bytes.Buffer, p *Parser, d []byte) func() bool {
					return func() bool {
						p.parseInline(o, d)
						return true
					}
				}(out, parser, data[prev:eol])
				parser.r.Header(out, work, level)

				// find the end of the underline
				for data[i] != '\n' {
					i++
				}
				return i
			}
		}

		// if the next line starts a block of HTML, then the paragraph ends here
		if parser.flags&EXTENSION_LAX_HTML_BLOCKS != 0 {
			if data[i] == '<' && parser.blockHtml(out, current, false) > 0 {
				// rewind to before the HTML block
				parser.renderParagraph(out, data[:i])
				return i
			}
		}

		// if there's a prefixed header or a horizontal rule after this, paragraph is over
		if parser.isPrefixHeader(current) || parser.isHRule(current) {
			parser.renderParagraph(out, data[:i])
			return i
		}

		// otherwise, scan to the beginning of the next line
		for data[i] != '\n' {
			i++
		}
		i++
	}

	parser.renderParagraph(out, data[:i])
	return i
}
