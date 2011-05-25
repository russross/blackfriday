package main

import (
	"bytes"
	"fmt"
	"html"
	"unicode"
)

const (
	MKDA_NOT_AUTOLINK = iota
	MKDA_NORMAL
	MKDA_EMAIL
)

const (
	MKDEXT_NO_INTRA_EMPHASIS = 1 << iota
	MKDEXT_TABLES
	MKDEXT_FENCED_CODE
	MKDEXT_AUTOLINK
	MKDEXT_STRIKETHROUGH
	MKDEXT_LAX_HTML_BLOCKS
	MKDEXT_SPACE_HEADERS
)

const (
	_ = iota
	MKD_LIST_ORDERED
	MKD_LI_BLOCK // <li> containing block data
	MKD_LI_END   = 8
)

const (
	MKD_TABLE_ALIGN_L = 1 << iota
	MKD_TABLE_ALIGN_R
	MKD_TABLE_ALIGN_CENTER = (MKD_TABLE_ALIGN_L | MKD_TABLE_ALIGN_R)
)

var block_tags = map[string]int{
	"p":          1, // 0
	"dl":         2,
	"h1":         2,
	"h2":         2,
	"h3":         2,
	"h4":         2,
	"h5":         2,
	"h6":         2,
	"ol":         2,
	"ul":         2,
	"del":        3, // 10
	"div":        3,
	"ins":        3, // 12
	"pre":        3,
	"form":       4,
	"math":       4,
	"table":      5,
	"iframe":     6,
	"script":     6,
	"fieldset":   8,
	"noscript":   8,
	"blockquote": 10,
}

// functions for rendering parsed data
type mkd_renderer struct {
	// block-level callbacks---nil skips the block
	blockcode  func(ob *bytes.Buffer, text []byte, lang string, opaque interface{})
	blockquote func(ob *bytes.Buffer, text []byte, opaque interface{})
	blockhtml  func(ob *bytes.Buffer, text []byte, opaque interface{})
	header     func(ob *bytes.Buffer, text []byte, level int, opaque interface{})
	hrule      func(ob *bytes.Buffer, opaque interface{})
	list       func(ob *bytes.Buffer, text []byte, flags int, opaque interface{})
	listitem   func(ob *bytes.Buffer, text []byte, flags int, opaque interface{})
	paragraph  func(ob *bytes.Buffer, text []byte, opaque interface{})
	table      func(ob *bytes.Buffer, header []byte, body []byte, opaque interface{})
	table_row  func(ob *bytes.Buffer, text []byte, opaque interface{})
	table_cell func(ob *bytes.Buffer, text []byte, flags int, opaque interface{})

	// user data---passed back to every callback
	opaque interface{}
}

type render struct {
	mk        mkd_renderer
	ext_flags uint32
	// ...
}

func parse_inline(work *bytes.Buffer, rndr *render, data []byte) {
	// TODO: inline rendering
	work.Write(data)
}

// parse block-level data
func parse_block(ob *bytes.Buffer, rndr *render, data []byte) {
	// TODO: quit if max_nesting exceeded

	for len(data) > 0 {
		if is_atxheader(rndr, data) {
			data = data[parse_atxheader(ob, rndr, data):]
			continue
		}
		if data[0] == '<' && rndr.mk.blockhtml != nil {
			if i := parse_htmlblock(ob, rndr, data, true); i > 0 {
				data = data[i:]
				continue
			}
		}
		if i := is_empty(data); i > 0 {
			data = data[i:]
			continue
		}
		if is_hrule(data) {
			if rndr.mk.hrule != nil {
				rndr.mk.hrule(ob, rndr.mk.opaque)
			}
			var i int
			for i = 0; i < len(data) && data[i] != '\n'; i++ {
			}
			data = data[i:]
			continue
		}
		if rndr.ext_flags&MKDEXT_FENCED_CODE != 0 {
			if i := parse_fencedcode(ob, rndr, data); i > 0 {
				data = data[i:]
				continue
			}
		}
		if rndr.ext_flags&MKDEXT_TABLES != 0 {
			if i := parse_table(ob, rndr, data); i > 0 {
				data = data[i:]
				continue
			}
		}
		if prefix_quote(data) > 0 {
			data = data[parse_blockquote(ob, rndr, data):]
			continue
		}
		if prefix_code(data) > 0 {
			data = data[parse_blockcode(ob, rndr, data):]
			continue
		}
		if prefix_uli(data) > 0 {
			data = data[parse_list(ob, rndr, data, 0):]
			continue
		}

		data = data[parse_paragraph(ob, rndr, data):]
	}
}

func is_atxheader(rndr *render, data []byte) bool {
	if data[0] != '#' {
		return false
	}

	if rndr.ext_flags&MKDEXT_SPACE_HEADERS != 0 {
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

func parse_atxheader(ob *bytes.Buffer, rndr *render, data []byte) int {
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
		work := bytes.NewBuffer(nil)
		parse_inline(work, rndr, data[i:end])
		if rndr.mk.header != nil {
			rndr.mk.header(ob, work.Bytes(), level, rndr.mk.opaque)
		}
	}
	return skip
}

func is_headerline(data []byte) int {
    i := 0

    // test of level 1 header
    if data[i] == '=' {
        for i = 1; i < len(data) && data[i] == '='; i++ {}
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
        for i = 1; i < len(data) && data[i] == '-'; i++ {}
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

func parse_htmlblock(ob *bytes.Buffer, rndr *render, data []byte, do_render bool) int {
	var i, j int

	// identification of the opening tag
	if len(data) < 2 || data[0] != '<' {
		return 0
	}
	curtag, tagfound := find_block_tag(data[1:])

	// handling of special cases
	if !tagfound {

		// HTML comment, laxist form
		if len(data) > 5 && data[1] == '!' && data[2] == '-' && data[3] == '-' {
			i = 5

			for i < len(data) && !(data[i-2] == '-' && data[i-1] == '-' && data[i] == '>') {
				i++
			}
			i++

			if i < len(data) {
				j = is_empty(data[i:])
			}

			if j > 0 {
				size := i + j
				if do_render && rndr.mk.blockhtml != nil {
					rndr.mk.blockhtml(ob, data[:size], rndr.mk.opaque)
				}
				return size
			}
		}

		// HR, which is the only self-closing block tag considered
		if len(data) > 4 && (data[i] == 'h' || data[1] == 'H') && (data[2] == 'r' || data[2] == 'R') {
			i = 3
			for i < len(data) && data[i] != '>' {
				i++
			}

			if i+1 < len(data) {
				i++
				j = is_empty(data[i:])
				if j > 0 {
					size := i + j
					if do_render && rndr.mk.blockhtml != nil {
						rndr.mk.blockhtml(ob, data[:size], rndr.mk.opaque)
					}
					return size
				}
			}
		}

		// no special case recognized
		return 0
	}

	// looking for an unindented matching closing tag
	//      followed by a blank line
	i = 1
	found := false

	// if not found, trying a second pass looking for indented match
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

			j = htmlblock_end(curtag, rndr, data[i-1:])

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
	if do_render && rndr.mk.blockhtml != nil {
		rndr.mk.blockhtml(ob, data[:i], rndr.mk.opaque)
	}

	return i
}

func find_block_tag(data []byte) (string, bool) {
	i := 0
	for i < len(data) && ((data[i] >= '0' && data[i] <= '9') || (data[i] >= 'A' && data[i] <= 'Z') || (data[i] >= 'a' && data[i] <= 'z')) {
		i++
	}
	if i >= len(data) {
		return "", false
	}
	key := string(data[:i])
	if _, ok := block_tags[key]; ok {
		return key, true
	}
	return "", false
}

func htmlblock_end(tag string, rndr *render, data []byte) int {
	// assuming data[0] == '<' && data[1] == '/' already tested

	// checking tag is a match
	if len(tag)+3 >= len(data) || bytes.Compare(data[2:2+len(tag)], []byte(tag)) != 0 || data[len(tag)+2] != '>' {
		return 0
	}

	// checking white lines
	i := len(tag) + 3
	w := 0
	if i < len(data) {
		if w = is_empty(data[i:]); w == 0 {
			return 0 // non-blank after tag
		}
	}
	i += w
	w = 0

	if rndr.ext_flags&MKDEXT_LAX_HTML_BLOCKS != 0 {
		if i < len(data) {
			w = is_empty(data[i:])
		}
	} else {
		if i < len(data) {
			if w = is_empty(data[i:]); w == 0 {
				return 0 // non-blank line after tag line
			}
		}
	}

	return i + w
}

func is_empty(data []byte) int {
	var i int
	for i = 0; i < len(data) && data[i] != '\n'; i++ {
		if data[i] != ' ' && data[i] != '\t' {
			return 0
		}
	}
	return i + 1
}

func is_hrule(data []byte) bool {
	// skipping initial spaces
	if len(data) < 3 {
		return false
	}
	i := 0
	if data[0] == ' ' {
		i++
		if data[1] == ' ' {
			i++
			if data[2] == ' ' {
				i++
			}
		}
	}

	// looking at the hrule char
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

func is_codefence(data []byte, syntax **string) int {
	i, n := 0, 0

	// skipping initial spaces
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

	// looking at the hrule char
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
			for syn > 0 && unicode.IsSpace(int(data[syntax_start])) {
				syntax_start++
				syn--
			}

			for syn > 0 && unicode.IsSpace(int(data[syntax_start+syn-1])) {
				syn--
			}

			i++
		} else {
			for i < len(data) && !unicode.IsSpace(int(data[i])) {
				syn++
				i++
			}
		}

		language := string(data[syntax_start : syntax_start+syn])
		*syntax = &language
	}

	for i < len(data) && data[i] != '\n' {
		if !unicode.IsSpace(int(data[i])) {
			return 0
		}
		i++
	}

	return i + 1
}

func parse_fencedcode(ob *bytes.Buffer, rndr *render, data []byte) int {
	var lang *string
	beg := is_codefence(data, &lang)
	if beg == 0 {
		return 0
	}

	work := bytes.NewBuffer(nil)

	for beg < len(data) {
		fence_end := is_codefence(data[beg:], nil)
		if fence_end != 0 {
			beg += fence_end
			break
		}

		var end int
		for end = beg + 1; end < len(data) && data[end-1] != '\n'; end++ {
		}

		if beg < end {
			// verbatim copy to the working buffer, escaping entities
			if is_empty(data[beg:]) > 0 {
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

	if rndr.mk.blockcode != nil {
		syntax := ""
		if lang != nil {
			syntax = *lang
		}

		rndr.mk.blockcode(ob, work.Bytes(), syntax, rndr.mk.opaque)
	}

	return beg
}

func parse_table(ob *bytes.Buffer, rndr *render, data []byte) int {
	header_work := bytes.NewBuffer(nil)
	i, columns, col_data := parse_table_header(header_work, rndr, data)
	if i > 0 {
		body_work := bytes.NewBuffer(nil)

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

			parse_table_row(body_work, rndr, data[row_start:i], columns, col_data)
			i++
		}

		if rndr.mk.table != nil {
			rndr.mk.table(ob, header_work.Bytes(), body_work.Bytes(), rndr.mk.opaque)
		}
	}

	return i
}

func parse_table_header(ob *bytes.Buffer, rndr *render, data []byte) (size int, columns int, column_data []int) {
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
			column_data[col] |= MKD_TABLE_ALIGN_L
			dashes++
		}

		for i < under_end && data[i] == '-' {
			i++
			dashes++
		}

		if i < under_end && data[i] == ':' {
			i++
			column_data[col] |= MKD_TABLE_ALIGN_R
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

	parse_table_row(ob, rndr, data[:header_end], columns, column_data)
	size = under_end + 1
	return
}

func parse_table_row(ob *bytes.Buffer, rndr *render, data []byte, columns int, col_data []int) {
	i, col := 0, 0
	row_work := bytes.NewBuffer(nil)

	if i < len(data) && data[i] == '|' {
		i++
	}

	for col = 0; col < columns && i < len(data); col++ {
		for i < len(data) && unicode.IsSpace(int(data[i])) {
			i++
		}

		cell_start := i

		for i < len(data) && data[i] != '|' {
			i++
		}

		cell_end := i - 1

		for cell_end > cell_start && unicode.IsSpace(int(data[cell_end])) {
			cell_end--
		}

		cell_work := bytes.NewBuffer(nil)
		parse_inline(cell_work, rndr, data[cell_start:cell_end+1])

		if rndr.mk.table_cell != nil {
			cdata := 0
			if col < len(col_data) {
				cdata = col_data[col]
			}
			rndr.mk.table_cell(row_work, cell_work.Bytes(), cdata, rndr.mk.opaque)
		}

		i++
	}

	for ; col < columns; col++ {
		empty_cell := []byte{}
		if rndr.mk.table_cell != nil {
			cdata := 0
			if col < len(col_data) {
				cdata = col_data[col]
			}
			rndr.mk.table_cell(row_work, empty_cell, cdata, rndr.mk.opaque)
		}
	}

	if rndr.mk.table_row != nil {
		rndr.mk.table_row(ob, row_work.Bytes(), rndr.mk.opaque)
	}
}

// returns blockquote prefix length
func prefix_quote(data []byte) int {
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

// handles parsing of a blockquote fragment
func parse_blockquote(ob *bytes.Buffer, rndr *render, data []byte) int {
	out := bytes.NewBuffer(nil)
	work := bytes.NewBuffer(nil)
	beg, end := 0, 0
	for beg < len(data) {
		for end = beg + 1; end < len(data) && data[end-1] != '\n'; end++ {
		}

		if pre := prefix_quote(data[beg:]); pre > 0 {
			beg += pre // skipping prefix
		} else {
			// empty line followed by non-quote line
			if is_empty(data[beg:]) > 0 && (end >= len(data) || (prefix_quote(data[end:]) == 0 && is_empty(data[end:]) == 0)) {
				break
			}
		}

		if beg < end { // copy into the in-place working buffer
			work.Write(data[beg:end])
		}
		beg = end
	}

	parse_block(out, rndr, work.Bytes())
	if rndr.mk.blockquote != nil {
		rndr.mk.blockquote(ob, out.Bytes(), rndr.mk.opaque)
	}
	return end
}

// returns prefix length for block code
func prefix_code(data []byte) int {
	if len(data) > 0 && data[0] == '\t' {
		return 1
	}
	if len(data) > 3 && data[0] == ' ' && data[1] == ' ' && data[2] == ' ' && data[3] == ' ' {
		return 4
	}
	return 0
}

func parse_blockcode(ob *bytes.Buffer, rndr *render, data []byte) int {
	work := bytes.NewBuffer(nil)

	beg, end := 0, 0
	for beg < len(data) {
		for end = beg + 1; end < len(data) && data[end-1] != '\n'; end++ {
		}

		chunk := data[beg:end]
		if pre := prefix_code(chunk); pre > 0 {
			beg += pre
		} else {
			if is_empty(chunk) == 0 {
				// non-empty non-prefixed line breaks the pre
				break
			}
		}

		if beg < end {
			// verbatim copy to the working buffer, escaping entities
			if is_empty(chunk) > 0 {
				work.WriteByte('\n')
			} else {
				work.Write(chunk)
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
		work = bytes.NewBuffer(workbytes[:len(workbytes)-n])
	}

	work.WriteByte('\n')

	if rndr.mk.blockcode != nil {
		rndr.mk.blockcode(ob, work.Bytes(), "", rndr.mk.opaque)
	}

	return beg
}

// returns unordered list item prefix
func prefix_uli(data []byte) int {
	i := 0
	for i < len(data) && i < 3 && data[i] == ' ' {
		i++
	}
	if i+1 >= len(data) || (data[i] != '*' && data[i] != '+' && data[i] != '-') || (data[i+1] != ' ' && data[i+1] != '\t') {
		return 0
	}
	return i + 2
}

// returns ordered list item prefix
func prefix_oli(data []byte) int {
	i := 0
	for i < len(data) && i < 3 && data[i] == ' ' {
		i++
	}
	if i >= len(data) || data[i] < '0' || data[i] > '9' {
		return 0
	}
	for i < len(data) && data[i] >= '0' && data[i] <= '9' {
		i++
	}
	if i+1 >= len(data) || data[i] != '.' || (data[i+1] != ' ' && data[i+1] != '\t') {
		return 0
	}
	return i + 2
}

// parsing ordered or unordered list block
func parse_list(ob *bytes.Buffer, rndr *render, data []byte, flags int) int {
	work := bytes.NewBuffer(nil)

	i, j, flags := 0, 0, 0
	for i < len(data) {
		j, flags = parse_listitem(work, rndr, data[i:], flags)
		i += j

		if j == 0 || flags&MKD_LI_END != 0 {
			break
		}
	}

	if rndr.mk.list != nil {
		rndr.mk.list(ob, work.Bytes(), flags, rndr.mk.opaque)
	}
	return i
}

// parsing a single list item
// assuming initial prefix is already removed
func parse_listitem(ob *bytes.Buffer, rndr *render, data []byte, flags_in int) (size int, flags int) {
	size, flags = 0, flags_in

	// keeping book of the first indentation prefix
	beg, end, pre, sublist, orgpre, i := 0, 0, 0, 0, 0, 0

	for orgpre < 3 && orgpre < len(data) && data[orgpre] == ' ' {
		orgpre++
	}

	beg = prefix_uli(data)
	if beg == 0 {
		beg = prefix_oli(data)
	}
	if beg == 0 {
		return
	}

	// skipping to the beginning of the following line
	end = beg
	for end < len(data) && data[end-1] != '\n' {
		end++
	}

	// getting working buffers
	work := bytes.NewBuffer(nil)
	inter := bytes.NewBuffer(nil)

	// putting the first line into the working buffer
	work.Write(data[beg:end])
	beg = end

	// process the following lines
	in_empty, has_inside_empty := false, false
	for beg < len(data) {
		end++

		for end < len(data) && data[end-1] != '\n' {
			end++
		}

		// process an empty line
		if is_empty(data[beg:end]) > 0 {
			in_empty = true
			beg = end
			continue
		}

		// calculating the indentation
		i = 0
		for i < 4 && beg+i < end && data[beg+i] == ' ' {
			i++
		}

		pre = i
		if data[beg] == '\t' {
			i = 1
			pre = 8
		}

		// checking for a new item
		chunk := data[beg+i : end]
		if (prefix_uli(chunk) > 0 && !is_hrule(chunk)) || prefix_oli(chunk) > 0 {
			if in_empty {
				has_inside_empty = true
			}

			if pre == orgpre { // the following item must have
				break // the same indentation
			}

			if sublist == 0 {
				sublist = work.Len()
			}
		} else {
			// joining only indented stuff after empty lines
			if in_empty && i < 4 && data[beg] != '\t' {
				flags |= MKD_LI_END
				break
			} else {
				if in_empty {
					work.WriteByte('\n')
					has_inside_empty = true
				}
			}
		}

		in_empty = false

		// adding the line without prefix into the working buffer
		work.Write(data[beg+i : end])
		beg = end
	}

	// render of li contents
	if has_inside_empty {
		flags |= MKD_LI_BLOCK
	}

	workbytes := work.Bytes()
	if flags&MKD_LI_BLOCK != 0 {
		// intermediate render of block li
		if sublist > 0 && sublist < len(workbytes) {
			parse_block(inter, rndr, workbytes[:sublist])
			parse_block(inter, rndr, workbytes[sublist:])
		} else {
			parse_block(inter, rndr, workbytes)
		}
	} else {
		// intermediate render of inline li
		if sublist > 0 && sublist < len(workbytes) {
			parse_inline(inter, rndr, workbytes[:sublist])
			parse_inline(inter, rndr, workbytes[sublist:])
		} else {
			parse_inline(inter, rndr, workbytes)
		}
	}

	// render of li itself
	if rndr.mk.listitem != nil {
		rndr.mk.listitem(ob, inter.Bytes(), flags, rndr.mk.opaque)
	}

	size = beg
	return
}

func parse_paragraph(ob *bytes.Buffer, rndr *render, data []byte) int {
    i, end, level := 0, 0, 0

    for i < len(data) {
        for end = i + 1; end < len(data) && data[end-1] != '\n'; end++ {}

        if is_empty(data[i:]) > 0 {
            break
        }
        if level = is_headerline(data[i:]); level > 0 {
            break
        }

        if rndr.ext_flags & MKDEXT_LAX_HTML_BLOCKS  != 0 {
            if data[i] == '<' && rndr.mk.blockhtml != nil && parse_htmlblock(ob, rndr, data[i:], false) > 0 {
                end = i
                break
            }
        }

        if is_atxheader(rndr, data[i:]) || is_hrule(data[i:]) {
            end = i
            break
        }

        i = end
    }

    work := data
    size := i
    for size > 0 && work[size-1] == '\n' {
        size--
    }
    
    if level == 0 {
        tmp := bytes.NewBuffer(nil)
        parse_inline(tmp, rndr, work[:size])
        if rndr.mk.paragraph != nil {
            rndr.mk.paragraph(ob, tmp.Bytes(), rndr.mk.opaque)
        }
    } else {
        if size > 0 {
            beg := 0
            i = size
            size--

            for size > 0 && work[size] != '\n' {
                size--
            }

            beg = size + 1
            for size > 0 && work[size-1] == '\n' {
                size--
            }

            if size > 0 {
                tmp := bytes.NewBuffer(nil)
                parse_inline(tmp, rndr, work[:size])
                if rndr.mk.paragraph != nil {
                    rndr.mk.paragraph(ob, tmp.Bytes(), rndr.mk.opaque)
                }

                work = work[beg:]
                size = i - beg
            } else {
                size = i
            }
        }

        header_work := bytes.NewBuffer(nil)
        parse_inline(header_work, rndr, work[:size])

        if rndr.mk.header != nil {
            rndr.mk.header(ob, header_work.Bytes(), level, rndr.mk.opaque)
        }
    }

    return end
}


//
//
// HTML rendering
//
//

const (
	HTML_SKIP_HTML = 1 << iota
	HTML_SKIP_STYLE
	HTML_SKIP_IMAGES
	HTML_SKIP_LINKS
	HTML_EXPAND_TABS
	HTML_SAFELINK
	HTML_TOC
	HTML_HARD_WRAP
	HTML_GITHUB_BLOCKCODE
	HTML_USE_XHTML
)

type html_renderopts struct {
	toc_data struct {
		header_count  int
		current_level int
	}
	flags     uint32
	close_tag string
}

func attr_escape(ob *bytes.Buffer, src []byte) {
	ob.WriteString(html.EscapeString(string(src)))
}

func rndr_header(ob *bytes.Buffer, text []byte, level int, opaque interface{}) {
	options := opaque.(*html_renderopts)

	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}

	if options.flags&HTML_TOC != 0 {
		ob.WriteString(fmt.Sprintf("<h%d id=\"toc_%d\">", level, options.toc_data.header_count))
		options.toc_data.header_count++
	} else {
		ob.WriteString(fmt.Sprintf("<h%d>", level))
	}

	ob.Write(text)
	ob.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func rndr_raw_block(ob *bytes.Buffer, text []byte, opaque interface{}) {
	sz := len(text)
	for sz > 0 && text[sz-1] == '\n' {
		sz--
	}
	org := 0
	for org < sz && text[org] == '\n' {
		org++
	}
	if org >= sz {
		return
	}
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.Write(text[org:sz])
	ob.WriteByte('\n')
}

func rndr_hrule(ob *bytes.Buffer, opaque interface{}) {
	options := opaque.(*html_renderopts)

	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.WriteString("<hr")
	ob.WriteString(options.close_tag)
}

func rndr_blockcode(ob *bytes.Buffer, text []byte, lang string, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}

	if lang != "" {
		ob.WriteString("<pre><code class=\"")

		for i, cls := 0, 0; i < len(lang); i, cls = i+1, cls+1 {
			for i < len(lang) && unicode.IsSpace(int(lang[i])) {
				i++
			}

			if i < len(lang) {
				org := i
				for i < len(lang) && !unicode.IsSpace(int(lang[i])) {
					i++
				}

				if lang[org] == '.' {
					org++
				}

				if cls > 0 {
					ob.WriteByte(' ')
				}
				attr_escape(ob, []byte(lang[org:]))
			}
		}

		ob.WriteString("\">")
	} else {
		ob.WriteString("<pre><code>")
	}

	if len(text) > 0 {
		attr_escape(ob, text)
	}

	ob.WriteString("</code></pre>\n")
}

func rndr_table(ob *bytes.Buffer, header []byte, body []byte, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.WriteString("<table><thead>\n")
	ob.Write(header)
	ob.WriteString("\n</thead><tbody>\n")
	ob.Write(body)
	ob.WriteString("\n</tbody></table>")
}

func rndr_tablerow(ob *bytes.Buffer, text []byte, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.WriteString("<tr>\n")
	ob.Write(text)
	ob.WriteString("\n</tr>")
}

func rndr_tablecell(ob *bytes.Buffer, text []byte, align int, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	switch align {
	case MKD_TABLE_ALIGN_L:
		ob.WriteString("<td align=\"left\">")
	case MKD_TABLE_ALIGN_R:
		ob.WriteString("<td align=\"right\">")
	case MKD_TABLE_ALIGN_CENTER:
		ob.WriteString("<td align=\"center\">")
	default:
		ob.WriteString("<td>")
	}

	ob.Write(text)
	ob.WriteString("</td>")
}

func rndr_list(ob *bytes.Buffer, text []byte, flags int, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	if flags&MKD_LIST_ORDERED != 0 {
		ob.WriteString("<ol>\n")
	} else {
		ob.WriteString("<ul>\n")
	}
	ob.Write(text)
	if flags&MKD_LIST_ORDERED != 0 {
		ob.WriteString("</ol>\n")
	} else {
		ob.WriteString("</ul>\n")
	}
}

func rndr_listitem(ob *bytes.Buffer, text []byte, flags int, opaque interface{}) {
	ob.WriteString("<li>")
	size := len(text)
	for size > 0 && text[size-1] == '\n' {
		size--
	}
	ob.Write(text[:size])
	ob.WriteString("</li>\n")
}

func rndr_paragraph(ob *bytes.Buffer, text []byte, opaque interface{}) {
	options := opaque.(*html_renderopts)
    i := 0

    if ob.Len() > 0 {
        ob.WriteByte('\n')
    }

    if len(text) == 0 {
        return
    }

    for i < len(text) && unicode.IsSpace(int(text[i])) {
        i++
    }

    if i == len(text) {
        return
    }

    ob.WriteString("<p>")
    if options.flags & HTML_HARD_WRAP != 0 {
        for i < len(text) {
            org := i
            for i < len(text) && text[i] != '\n' {
                i++
            }

            if i > org {
                ob.Write(text[org:i])
            }

            if i >= len(text) {
                break
            }

            ob.WriteString("<br>")
            ob.WriteString(options.close_tag)
            i++
        }
    } else {
        ob.Write(text[i:])
    }
    ob.WriteString("</p>\n")
}


func main() {
	ob := bytes.NewBuffer(nil)
	input := "##Header##\n"
	input += "\n"
	input += "----------\n"
	input += "\n"
	input += "Underlined header\n"
	input += "-----------------\n"
	input += "\n"
	input += "<p>Some block html\n"
	input += "</p>\n"
	input += "\n"
	input += "Score | Grade\n"
	input += "------|------\n"
	input += "94    | A\n"
	input += "85    | B\n"
	input += "74    | C\n"
	input += "65    | D\n"
	input += "\n"
	input += "``` go\n"
	input += "func fib(n int) int {\n"
	input += "    if n <= 1 {\n"
	input += "        return n\n"
	input += "    }\n"
	input += "    return n * fib(n-1)\n"
	input += "}\n"
	input += "```\n"
	input += "\n"
	input += "> A blockquote\n"
	input += "> or something like that\n"
	input += "> With a table | of two columns\n"
	input += "> -------------|---------------\n"
	input += "> key          | value \n"
	input += "\n"
	input += "\n"
	input += "Some **bold** Some *italic* and [a link][1] \n"
	input += "\n"
	input += "A little code sample\n"
	input += "\n"
	input += "    </head>\n"
	input += "    <title>Web Page Title</title>\n"
	input += "    </head>\n"
	input += "\n"
	input += "A picture\n"
	input += "\n"
	input += "![alt text][2]\n"
	input += "\n"
	input += "A list\n"
	input += "\n"
	input += "- apples\n"
	input += "- oranges\n"
	input += "- eggs\n"
	input += "\n"
	input += "A numbered list\n"
	input += "\n"
	input += "1. a\n"
	input += "2. b\n"
	input += "3. c\n"
	input += "\n"
	input += "A little quote\n"
	input += "\n"
	input += "> It is now time for all good men to come to the aid of their country. \n"
	input += "\n"
	input += "A final paragraph.\n"
	input += "\n"
	input += "  [1]: http://www.google.com\n"
	input += "  [2]: http://www.google.com/intl/en_ALL/images/logo.gif\n"

	ib := []byte(input)
	rndrer := new(mkd_renderer)
	rndrer.blockcode = rndr_blockcode
	rndrer.blockhtml = rndr_raw_block
	rndrer.header = rndr_header
	rndrer.hrule = rndr_hrule
	rndrer.list = rndr_list
	rndrer.listitem = rndr_listitem
    rndrer.paragraph = rndr_paragraph
	rndrer.table = rndr_table
	rndrer.table_row = rndr_tablerow
	rndrer.table_cell = rndr_tablecell
	rndrer.opaque = &html_renderopts{close_tag: " />"}
	var extensions uint32 = MKDEXT_FENCED_CODE | MKDEXT_TABLES
	Ups_markdown(ob, ib, rndrer, extensions)
	fmt.Print(ob.String())
}

func Ups_markdown(ob *bytes.Buffer, ib []byte, rndrer *mkd_renderer, extensions uint32) {

	/* filling the render structure */
	if rndrer == nil {
		return
	}

	rndr := &render{*rndrer, extensions}

	parse_block(ob, rndr, ib)
}
