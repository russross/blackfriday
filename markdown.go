package main

import (
	"fmt"
	"bytes"
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

type mkd_renderer struct {
	blockhtml func(ob *bytes.Buffer, text []byte, opaque interface{})
	header    func(ob *bytes.Buffer, text []byte, level int, opaque interface{})
    hrule func(ob *bytes.Buffer, opaque interface{})
	opaque    interface{}
}

type render struct {
	maker     mkd_renderer
	ext_flags uint32
	// ...
}

func parse_inline(work *bytes.Buffer, rndr *render, data []byte) {
	// TODO: inline rendering
	work.Write(data)
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
		work := new(bytes.Buffer)
		parse_inline(work, rndr, data[i:end])
		if rndr.maker.header != nil {
			rndr.maker.header(ob, work.Bytes(), level, rndr.maker.opaque)
		}
	}
	return skip
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
				if do_render && rndr.maker.blockhtml != nil {
					rndr.maker.blockhtml(ob, data[:size], rndr.maker.opaque)
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
					if do_render && rndr.maker.blockhtml != nil {
						rndr.maker.blockhtml(ob, data[:size], rndr.maker.opaque)
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
	// but not if tag is "ins" or "del" (folloing original Markdown.pl)
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
	if do_render && rndr.maker.blockhtml != nil {
		rndr.maker.blockhtml(ob, data[:i], rndr.maker.opaque)
	}

	return i
}

func parse_block(ob *bytes.Buffer, rndr *render, data []byte) {
	// TODO: quit if max_nesting exceeded

	for len(data) > 0 {
		if is_atxheader(rndr, data) {
			data = data[parse_atxheader(ob, rndr, data):]
			continue
		}
		if data[0] == '<' && rndr.maker.blockhtml != nil {
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
            if rndr.maker.hrule != nil {
                rndr.maker.hrule(ob, rndr.maker.opaque)
            }
            var i int
            for i = 0; i < len(data) && data[i] != '\n'; i++ {}
            data = data[i:]
        }

		data = data[1:]
	}
}

func Ups_markdown(ob *bytes.Buffer, ib []byte, rndrer *mkd_renderer, extensions uint32) {

	/* filling the render structure */
	if rndrer == nil {
		return
	}

	rndr := &render{*rndrer, 0}

	parse_block(ob, rndr, ib)
}

func main() {
	ob := new(bytes.Buffer)
	input := "### Header 3\n-----\n# Header 1 #\n<ul>A list\n</ul>\n"
	ib := bytes.NewBufferString(input).Bytes()
	rndrer := new(mkd_renderer)
	rndrer.blockhtml = rndr_raw_block
	rndrer.header = rndr_header
    rndrer.hrule = rndr_hrule
	rndrer.opaque = &html_renderopts{close_tag:" />"}
	var extensions uint32
	extensions = 0
	Ups_markdown(ob, ib, rndrer, extensions)
	fmt.Print(ob.String())
}


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

	if len(text) > 0 {
		ob.Write(text)
	}
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
