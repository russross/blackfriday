//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
//
// Markdown parsing and processing
//
//

package blackfriday

import (
	"bytes"
	"utf8"
)

const VERSION = "0.6"

// These are the supported markdown parsing extensions.
// OR these values together to select multiple extensions.
const (
	EXTENSION_NO_INTRA_EMPHASIS = 1 << iota
	EXTENSION_TABLES
	EXTENSION_FENCED_CODE
	EXTENSION_AUTOLINK
	EXTENSION_STRIKETHROUGH
	EXTENSION_LAX_HTML_BLOCKS
	EXTENSION_SPACE_HEADERS
	EXTENSION_HARD_LINE_BREAK
	EXTENSION_NO_EXPAND_TABS
	EXTENSION_TAB_SIZE_EIGHT
)

// These are the possible flag values for the link renderer.
// Only a single one of these values will be used; they are not ORed together.
// These are mostly of interest if you are writing a new output format.
const (
	LINK_TYPE_NOT_AUTOLINK = iota
	LINK_TYPE_NORMAL
	LINK_TYPE_EMAIL
)

// These are the possible flag values for the listitem renderer.
// Multiple flag values may be ORed together.
// These are mostly of interest if you are writing a new output format.
const (
	LIST_TYPE_ORDERED = 1 << iota
	LIST_ITEM_CONTAINS_BLOCK
	LIST_ITEM_END_OF_LIST
)

// These are the possible flag values for the table cell renderer.
// Only a single one of these values will be used; they are not ORed together.
// These are mostly of interest if you are writing a new output format.
const (
	TABLE_ALIGNMENT_LEFT = 1 << iota
	TABLE_ALIGNMENT_RIGHT
	TABLE_ALIGNMENT_CENTER = (TABLE_ALIGNMENT_LEFT | TABLE_ALIGNMENT_RIGHT)
)

// The size of a tab stop.
const (
	TAB_SIZE_DEFAULT = 4
	TAB_SIZE_EIGHT   = 8
)

// These are the tags that are recognized as HTML block tags.
// Any of these can be included in markdown text without special escaping.
var blockTags = map[string]bool{
	"p":          true,
	"dl":         true,
	"h1":         true,
	"h2":         true,
	"h3":         true,
	"h4":         true,
	"h5":         true,
	"h6":         true,
	"ol":         true,
	"ul":         true,
	"del":        true,
	"div":        true,
	"ins":        true,
	"pre":        true,
	"form":       true,
	"math":       true,
	"table":      true,
	"iframe":     true,
	"script":     true,
	"fieldset":   true,
	"noscript":   true,
	"blockquote": true,
}

// This interface defines the rendering interface.
// This is mostly of interest if you are implementing a new rendering format.
// Currently Html and Latex implementations are provided
type Renderer interface {
	// block-level callbacks
	BlockCode(out *bytes.Buffer, text []byte, lang string)
	BlockQuote(out *bytes.Buffer, text []byte)
	BlockHtml(out *bytes.Buffer, text []byte)
	Header(out *bytes.Buffer, text func() bool, level int)
	HRule(out *bytes.Buffer)
	List(out *bytes.Buffer, text func() bool, flags int)
	ListItem(out *bytes.Buffer, text []byte, flags int)
	Paragraph(out *bytes.Buffer, text func() bool)
	Table(out *bytes.Buffer, header []byte, body []byte, columnData []int)
	TableRow(out *bytes.Buffer, text []byte)
	TableCell(out *bytes.Buffer, text []byte, flags int)

	// Span-level callbacks
	AutoLink(out *bytes.Buffer, link []byte, kind int)
	CodeSpan(out *bytes.Buffer, text []byte)
	DoubleEmphasis(out *bytes.Buffer, text []byte)
	Emphasis(out *bytes.Buffer, text []byte)
	Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)
	LineBreak(out *bytes.Buffer)
	Link(out *bytes.Buffer, link []byte, title []byte, content []byte)
	RawHtmlTag(out *bytes.Buffer, tag []byte)
	TripleEmphasis(out *bytes.Buffer, text []byte)
	StrikeThrough(out *bytes.Buffer, text []byte)

	// Low-level callbacks
	Entity(out *bytes.Buffer, entity []byte)
	NormalText(out *bytes.Buffer, text []byte)

	// Header and footer
	DocumentHeader(out *bytes.Buffer)
	DocumentFooter(out *bytes.Buffer)
}

// Callback functions for inline parsing. One such function is defined
// for each character that triggers a response when parsing inline data.
type inlineParser func(parser *Parser, out *bytes.Buffer, data []byte, offset int) int

// The main parser object.
// This is constructed by the Markdown function and
// contains state used during the parsing process.
type Parser struct {
	r          Renderer
	refs       map[string]*reference
	inline     [256]inlineParser
	flags      int
	nesting    int
	maxNesting int
	insideLink bool
}


//
//
// Public interface
//
//

// Call Markdown with no extensions
func MarkdownBasic(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := HTML_USE_XHTML
	renderer := HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0

	return Markdown(input, renderer, extensions)
}

// Call Markdown with most useful extensions enabled
func MarkdownCommon(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := 0
	htmlFlags |= HTML_USE_XHTML
	htmlFlags |= HTML_USE_SMARTYPANTS
	htmlFlags |= HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= HTML_SMARTYPANTS_LATEX_DASHES
	renderer := HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= EXTENSION_NO_INTRA_EMPHASIS
	extensions |= EXTENSION_TABLES
	extensions |= EXTENSION_FENCED_CODE
	extensions |= EXTENSION_AUTOLINK
	extensions |= EXTENSION_STRIKETHROUGH
	extensions |= EXTENSION_SPACE_HEADERS

	return Markdown(input, renderer, extensions)
}

// Parse and render a block of markdown-encoded text.
// The renderer is used to format the output, and extensions dictates which
// non-standard extensions are enabled.
func Markdown(input []byte, renderer Renderer, extensions int) []byte {
	// no point in parsing if we can't render
	if renderer == nil {
		return nil
	}

	// fill in the render structure
	parser := new(Parser)
	parser.r = renderer
	parser.flags = extensions
	parser.refs = make(map[string]*reference)
	parser.maxNesting = 16
	parser.insideLink = false

	// register inline parsers
	parser.inline['*'] = inlineEmphasis
	parser.inline['_'] = inlineEmphasis
	if extensions&EXTENSION_STRIKETHROUGH != 0 {
		parser.inline['~'] = inlineEmphasis
	}
	parser.inline['`'] = inlineCodeSpan
	parser.inline['\n'] = inlineLineBreak
	parser.inline['['] = inlineLink
	parser.inline['<'] = inlineLAngle
	parser.inline['\\'] = inlineEscape
	parser.inline['&'] = inlineEntity

	if extensions&EXTENSION_AUTOLINK != 0 {
		parser.inline[':'] = inlineAutoLink
	}

	first := firstPass(parser, input)
	second := secondPass(parser, first)

	return second
}

// first pass:
// - extract references
// - expand tabs
// - normalize newlines
// - copy everything else
func firstPass(parser *Parser, input []byte) []byte {
	var out bytes.Buffer
	tabSize := TAB_SIZE_DEFAULT
	if parser.flags&EXTENSION_TAB_SIZE_EIGHT != 0 {
		tabSize = TAB_SIZE_EIGHT
	}
	beg, end := 0, 0
	for beg < len(input) { // iterate over lines
		if end = isReference(parser, input[beg:]); end > 0 {
			beg += end
		} else { // skip to the next line
			end = beg
			for end < len(input) && input[end] != '\n' && input[end] != '\r' {
				end++
			}

			// add the line body if present
			if end > beg {
				if parser.flags&EXTENSION_NO_EXPAND_TABS == 0 {
					expandTabs(&out, input[beg:end], tabSize)
				} else {
					out.Write(input[beg:end])
				}
			}
			out.WriteByte('\n')

			if end < len(input) && input[end] == '\r' {
				end++
			}
			if end < len(input) && input[end] == '\n' {
				end++
			}

			beg = end
		}
	}
	return out.Bytes()
}

// second pass: actual rendering
func secondPass(parser *Parser, input []byte) []byte {
	var output bytes.Buffer

	parser.r.DocumentHeader(&output)
	parser.parseBlock(&output, input)
	parser.r.DocumentFooter(&output)

	if parser.nesting != 0 {
		panic("Nesting level did not end at zero")
	}

	return output.Bytes()
}


//
// Link references
//
// This section implements support for references that (usually) appear
// as footnotes in a document, and can be referenced anywhere in the document.
// The basic format is:
//
//    [1]: http://www.google.com/ "Google"
//    [2]: http://www.github.com/ "Github"
//
// Anywhere in the document, the reference can be linked by referring to its
// label, i.e., 1 and 2 in this example, as in:
//
//    This library is hosted on [Github][2], a git hosting site.

// References are parsed and stored in this struct.
type reference struct {
	link  []byte
	title []byte
}

// Check whether or not data starts with a reference link.
// If so, it is parsed and stored in the list of references
// (in the render struct).
// Returns the number of bytes to skip to move past it,
// or zero if the first line is not a reference.
func isReference(parser *Parser, data []byte) int {
	// up to 3 optional leading spaces
	if len(data) < 4 {
		return 0
	}
	i := 0
	for i < 3 && data[i] == ' ' {
		i++
	}

	// id part: anything but a newline between brackets
	if data[i] != '[' {
		return 0
	}
	i++
	idOffset := i
	for i < len(data) && data[i] != '\n' && data[i] != '\r' && data[i] != ']' {
		i++
	}
	if i >= len(data) || data[i] != ']' {
		return 0
	}
	idEnd := i

	// spacer: colon (space | tab)* newline? (space | tab)*
	i++
	if i >= len(data) || data[i] != ':' {
		return 0
	}
	i++
	for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
		i++
	}
	if i < len(data) && (data[i] == '\n' || data[i] == '\r') {
		i++
		if i < len(data) && data[i] == '\n' && data[i-1] == '\r' {
			i++
		}
	}
	for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
		i++
	}
	if i >= len(data) {
		return 0
	}

	// link: whitespace-free sequence, optionally between angle brackets
	if data[i] == '<' {
		i++
	}
	linkOffset := i
	for i < len(data) && data[i] != ' ' && data[i] != '\t' && data[i] != '\n' && data[i] != '\r' {
		i++
	}
	linkEnd := i
	if data[linkOffset] == '<' && data[linkEnd-1] == '>' {
		linkOffset++
		linkEnd--
	}

	// optional spacer: (space | tab)* (newline | '\'' | '"' | '(' )
	for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
		i++
	}
	if i < len(data) && data[i] != '\n' && data[i] != '\r' && data[i] != '\'' && data[i] != '"' && data[i] != '(' {
		return 0
	}

	// compute end-of-line
	lineEnd := 0
	if i >= len(data) || data[i] == '\r' || data[i] == '\n' {
		lineEnd = i
	}
	if i+1 < len(data) && data[i] == '\r' && data[i+1] == '\n' {
		lineEnd++
	}

	// optional (space|tab)* spacer after a newline
	if lineEnd > 0 {
		i = lineEnd + 1
		for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
			i++
		}
	}

	// optional title: any non-newline sequence enclosed in '"() alone on its line
	titleOffset, titleEnd := 0, 0
	if i+1 < len(data) && (data[i] == '\'' || data[i] == '"' || data[i] == '(') {
		i++
		titleOffset = i

		// look for EOL
		for i < len(data) && data[i] != '\n' && data[i] != '\r' {
			i++
		}
		if i+1 < len(data) && data[i] == '\n' && data[i+1] == '\r' {
			titleEnd = i + 1
		} else {
			titleEnd = i
		}

		// step back
		i--
		for i > titleOffset && (data[i] == ' ' || data[i] == '\t') {
			i--
		}
		if i > titleOffset && (data[i] == '\'' || data[i] == '"' || data[i] == ')') {
			lineEnd = titleEnd
			titleEnd = i
		}
	}
	if lineEnd == 0 { // garbage after the link
		return 0
	}

	// a valid ref has been found
	if parser == nil {
		return lineEnd
	}

	// id matches are case-insensitive
	id := string(bytes.ToLower(data[idOffset:idEnd]))
	parser.refs[id] = &reference{
		link:  data[linkOffset:linkEnd],
		title: data[titleOffset:titleEnd],
	}

	return lineEnd
}


//
//
// Miscellaneous helper functions
//
//


// Test if a character is a punctuation symbol.
// Taken from a private function in regexp in the stdlib.
func ispunct(c byte) bool {
	for _, r := range []byte("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		if c == r {
			return true
		}
	}
	return false
}

// Test if a character is a whitespace character.
func isspace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// Test if a character is a letter or a digit.
// TODO: check when this is looking for ASCII alnum and when it should use unicode
func isalnum(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// Replace tab characters with spaces, aligning to the next TAB_SIZE column.
// always ends output with a newline
func expandTabs(out *bytes.Buffer, line []byte, tabSize int) {
	// first, check for common cases: no tabs, or only tabs at beginning of line
	i, prefix := 0, 0
	slowcase := false
	for i = 0; i < len(line); i++ {
		if line[i] == '\t' {
			if prefix == i {
				prefix++
			} else {
				slowcase = true
				break
			}
		}
	}

	// no need to decode runes if all tabs are at the beginning of the line
	if !slowcase {
		for i = 0; i < prefix*tabSize; i++ {
			out.WriteByte(' ')
		}
		out.Write(line[prefix:])
		return
	}

	// the slow case: we need to count runes to figure out how
	// many spaces to insert for each tab
	column := 0
	i = 0
	for i < len(line) {
		start := i
		for i < len(line) && line[i] != '\t' {
			_, size := utf8.DecodeRune(line[i:])
			i += size
			column++
		}

		if i > start {
			out.Write(line[start:i])
		}

		if i >= len(line) {
			break
		}

		for {
			out.WriteByte(' ')
			column++
			if column%tabSize == 0 {
				break
			}
		}

		i++
	}
}
