//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
//

//
//
// SmartyPants rendering
//
//

package blackfriday

import (
	"bytes"
)

type smartypantsData struct {
	inSingleQuote bool
	inDoubleQuote bool
}

func wordBoundary(c byte) bool {
	return c == 0 || isspace(c) || ispunct(c)
}

func tolower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c - 'A' + 'a'
	}
	return c
}

func isdigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func smartQuotesHelper(ob *bytes.Buffer, previousChar byte, nextChar byte, quote byte, isOpen *bool) bool {
	switch {
	// edge of the buffer is likely to be a tag that we don't get to see,
	// so we assume there is text there
	case wordBoundary(previousChar) && previousChar != 0 && nextChar == 0:
		*isOpen = true
	case previousChar == 0 && wordBoundary(nextChar) && nextChar != 0:
		*isOpen = false
	case wordBoundary(previousChar) && !wordBoundary(nextChar):
		*isOpen = true
	case !wordBoundary(previousChar) && wordBoundary(nextChar):
		*isOpen = false
	case !wordBoundary(previousChar) && !wordBoundary(nextChar):
		*isOpen = true
	default:
		*isOpen = !*isOpen
	}

	ob.WriteByte('&')
	if *isOpen {
		ob.WriteByte('l')
	} else {
		ob.WriteByte('r')
	}
	ob.WriteByte(quote)
	ob.WriteString("quo;")
	return true
}

func smartSquote(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if len(text) >= 2 {
		t1 := tolower(text[1])

		if t1 == '\'' {
			nextChar := byte(0)
			if len(text) >= 3 {
				nextChar = text[2]
			}
			if smartQuotesHelper(ob, previousChar, nextChar, 'd', &smrt.inDoubleQuote) {
				return 1
			}
		}

		if (t1 == 's' || t1 == 't' || t1 == 'm' || t1 == 'd') && (len(text) < 3 || wordBoundary(text[2])) {
			ob.WriteString("&rsquo;")
			return 0
		}

		if len(text) >= 3 {
			t2 := tolower(text[2])

			if ((t1 == 'r' && t2 == 'e') || (t1 == 'l' && t2 == 'l') || (t1 == 'v' && t2 == 'e')) && (len(text) < 4 || wordBoundary(text[3])) {
				ob.WriteString("&rsquo;")
				return 0
			}
		}
	}

	nextChar := byte(0)
	if len(text) > 1 {
		nextChar = text[1]
	}
	if smartQuotesHelper(ob, previousChar, nextChar, 's', &smrt.inSingleQuote) {
		return 0
	}

	ob.WriteByte(text[0])
	return 0
}

func smartParens(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if len(text) >= 3 {
		t1 := tolower(text[1])
		t2 := tolower(text[2])

		if t1 == 'c' && t2 == ')' {
			ob.WriteString("&copy;")
			return 2
		}

		if t1 == 'r' && t2 == ')' {
			ob.WriteString("&reg;")
			return 2
		}

		if len(text) >= 4 && t1 == 't' && t2 == 'm' && text[3] == ')' {
			ob.WriteString("&trade;")
			return 3
		}
	}

	ob.WriteByte(text[0])
	return 0
}

func smartDash(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if len(text) >= 2 {
		if text[1] == '-' {
			ob.WriteString("&mdash;")
			return 1
		}

		if wordBoundary(previousChar) && wordBoundary(text[1]) {
			ob.WriteString("&ndash;")
			return 0
		}
	}

	ob.WriteByte(text[0])
	return 0
}

func smartDashLatex(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if len(text) >= 3 && text[1] == '-' && text[2] == '-' {
		ob.WriteString("&mdash;")
		return 2
	}
	if len(text) >= 2 && text[1] == '-' {
		ob.WriteString("&ndash;")
		return 1
	}

	ob.WriteByte(text[0])
	return 0
}

func smartAmp(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if bytes.HasPrefix(text, []byte("&quot;")) {
		nextChar := byte(0)
		if len(text) >= 7 {
			nextChar = text[6]
		}
		if smartQuotesHelper(ob, previousChar, nextChar, 'd', &smrt.inDoubleQuote) {
			return 5
		}
	}

	if bytes.HasPrefix(text, []byte("&#0;")) {
		return 3
	}

	ob.WriteByte('&')
	return 0
}

func smartPeriod(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if len(text) >= 3 && text[1] == '.' && text[2] == '.' {
		ob.WriteString("&hellip;")
		return 2
	}

	if len(text) >= 5 && text[1] == ' ' && text[2] == '.' && text[3] == ' ' && text[4] == '.' {
		ob.WriteString("&hellip;")
		return 4
	}

	ob.WriteByte(text[0])
	return 0
}

func smartBacktick(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if len(text) >= 2 && text[1] == '`' {
		nextChar := byte(0)
		if len(text) >= 3 {
			nextChar = text[2]
		}
		if smartQuotesHelper(ob, previousChar, nextChar, 'd', &smrt.inDoubleQuote) {
			return 1
		}
	}

	return 0
}

func smartNumberGeneric(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if wordBoundary(previousChar) && len(text) >= 3 {
		// is it of the form digits/digits(word boundary)?, i.e., \d+/\d+\b
		num_end := 0
		for len(text) > num_end && isdigit(text[num_end]) {
			num_end++
		}
		if num_end == 0 {
			ob.WriteByte(text[0])
			return 0
		}
		if len(text) < num_end+2 || text[num_end] != '/' {
			ob.WriteByte(text[0])
			return 0
		}
		den_end := num_end + 1
		for len(text) > den_end && isdigit(text[den_end]) {
			den_end++
		}
		if den_end == num_end+1 {
			ob.WriteByte(text[0])
			return 0
		}
		if len(text) == den_end || wordBoundary(text[den_end]) {
			ob.WriteString("<sup>")
			ob.Write(text[:num_end])
			ob.WriteString("</sup>&frasl;<sub>")
			ob.Write(text[num_end+1 : den_end])
			ob.WriteString("</sub>")
			return den_end - 1
		}
	}

	ob.WriteByte(text[0])
	return 0
}

func smartNumber(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	if wordBoundary(previousChar) && len(text) >= 3 {
		if text[0] == '1' && text[1] == '/' && text[2] == '2' {
			if len(text) < 4 || wordBoundary(text[3]) {
				ob.WriteString("&frac12;")
				return 2
			}
		}

		if text[0] == '1' && text[1] == '/' && text[2] == '4' {
			if len(text) < 4 || wordBoundary(text[3]) || (len(text) >= 5 && tolower(text[3]) == 't' && tolower(text[4]) == 'h') {
				ob.WriteString("&frac14;")
				return 2
			}
		}

		if text[0] == '3' && text[1] == '/' && text[2] == '4' {
			if len(text) < 4 || wordBoundary(text[3]) || (len(text) >= 6 && tolower(text[3]) == 't' && tolower(text[4]) == 'h' && tolower(text[5]) == 's') {
				ob.WriteString("&frac34;")
				return 2
			}
		}
	}

	ob.WriteByte(text[0])
	return 0
}

func smartDquote(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	nextChar := byte(0)
	if len(text) > 1 {
		nextChar = text[1]
	}
	if !smartQuotesHelper(ob, previousChar, nextChar, 'd', &smrt.inDoubleQuote) {
		ob.WriteString("&quot;")
	}

	return 0
}

func smartLtag(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int {
	i := 0

	for i < len(text) && text[i] != '>' {
		i++
	}

	ob.Write(text[:i+1])
	return i
}

type smartCallback func(ob *bytes.Buffer, smrt *smartypantsData, previousChar byte, text []byte) int

type SmartypantsRenderer [256]smartCallback

func Smartypants(flags int) *SmartypantsRenderer {
	r := new(SmartypantsRenderer)
	r['"'] = smartDquote
	r['&'] = smartAmp
	r['\''] = smartSquote
	r['('] = smartParens
	if flags&HTML_SMARTYPANTS_LATEX_DASHES == 0 {
		r['-'] = smartDash
	} else {
		r['-'] = smartDashLatex
	}
	r['.'] = smartPeriod
	if flags&HTML_SMARTYPANTS_FRACTIONS == 0 {
		r['1'] = smartNumber
		r['3'] = smartNumber
	} else {
		for ch := '1'; ch <= '9'; ch++ {
			r[ch] = smartNumberGeneric
		}
	}
	r['<'] = smartLtag
	r['`'] = smartBacktick
	return r
}

func htmlSmartypants(ob *bytes.Buffer, text []byte, opaque interface{}) {
	options := opaque.(*htmlOptions)
	smrt := smartypantsData{false, false}

	// first do normal entity escaping
	var escaped bytes.Buffer
	attrEscape(&escaped, text)
	text = escaped.Bytes()

	mark := 0
	for i := 0; i < len(text); i++ {
		if action := options.smartypants[text[i]]; action != nil {
			if i > mark {
				ob.Write(text[mark:i])
			}

			previousChar := byte(0)
			if i > 0 {
				previousChar = text[i-1]
			}
			i += action(ob, &smrt, previousChar, text[i:])
			mark = i + 1
		}
	}

	if mark < len(text) {
		ob.Write(text[mark:])
	}
}
