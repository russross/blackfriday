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

type smartypants_data struct {
	in_squote bool
	in_dquote bool
}

func word_boundary(c byte) bool {
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

func smartypants_quotes(ob *bytes.Buffer, previous_char byte, next_char byte, quote byte, is_open *bool) bool {
	switch {
	// edge of the buffer is likely to be a tag that we don't get to see,
	// so we assume there is text there
	case word_boundary(previous_char) && previous_char != 0 && next_char == 0:
		*is_open = true
	case previous_char == 0 && word_boundary(next_char) && next_char != 0:
		*is_open = false
	case word_boundary(previous_char) && !word_boundary(next_char):
		*is_open = true
	case !word_boundary(previous_char) && word_boundary(next_char):
		*is_open = false
	case !word_boundary(previous_char) && !word_boundary(next_char):
		*is_open = true
	default:
		*is_open = !*is_open
	}

	ob.WriteByte('&')
	if *is_open {
		ob.WriteByte('l')
	} else {
		ob.WriteByte('r')
	}
	ob.WriteByte(quote)
	ob.WriteString("quo;")
	return true
}

func smartypants_cb__squote(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	if len(text) >= 2 {
		t1 := tolower(text[1])

		if t1 == '\'' {
			next_char := byte(0)
			if len(text) >= 3 {
				next_char = text[2]
			}
			if smartypants_quotes(ob, previous_char, next_char, 'd', &smrt.in_dquote) {
				return 1
			}
		}

		if (t1 == 's' || t1 == 't' || t1 == 'm' || t1 == 'd') && (len(text) < 3 || word_boundary(text[2])) {
			ob.WriteString("&rsquo;")
			return 0
		}

		if len(text) >= 3 {
			t2 := tolower(text[2])

			if ((t1 == 'r' && t2 == 'e') || (t1 == 'l' && t2 == 'l') || (t1 == 'v' && t2 == 'e')) && (len(text) < 4 || word_boundary(text[3])) {
				ob.WriteString("&rsquo;")
				return 0
			}
		}
	}

	next_char := byte(0)
	if len(text) > 1 {
		next_char = text[1]
	}
	if smartypants_quotes(ob, previous_char, next_char, 's', &smrt.in_squote) {
		return 0
	}

	ob.WriteByte(text[0])
	return 0
}

func smartypants_cb__parens(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
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

func smartypants_cb__dash(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	if len(text) >= 2 {
		if text[1] == '-' {
			ob.WriteString("&mdash;")
			return 1
		}

		if word_boundary(previous_char) && word_boundary(text[1]) {
			ob.WriteString("&ndash;")
			return 0
		}
	}

	ob.WriteByte(text[0])
	return 0
}

func smartypants_cb__dash_latex(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
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

func smartypants_cb__amp(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	if bytes.HasPrefix(text, []byte("&quot;")) {
		next_char := byte(0)
		if len(text) >= 7 {
			next_char = text[6]
		}
		if smartypants_quotes(ob, previous_char, next_char, 'd', &smrt.in_dquote) {
			return 5
		}
	}

	if bytes.HasPrefix(text, []byte("&#0;")) {
		return 3
	}

	ob.WriteByte('&')
	return 0
}

func smartypants_cb__period(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
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

func smartypants_cb__backtick(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	if len(text) >= 2 && text[1] == '`' {
		next_char := byte(0)
		if len(text) >= 3 {
			next_char = text[2]
		}
		if smartypants_quotes(ob, previous_char, next_char, 'd', &smrt.in_dquote) {
			return 1
		}
	}

	return 0
}

func smartypants_cb__number_generic(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	if word_boundary(previous_char) && len(text) >= 3 {
		// is it of the form digits/digits(word boundary)?, i.e., \d+/\d+\b
		num_end := 0
		for len(text) > num_end && isdigit(text[num_end]) {
			num_end++
		}
		if num_end == 0 {
			return 0
		}
		if len(text) < num_end+2 || text[num_end] != '/' {
			return 0
		}
		den_end := num_end + 1
		for len(text) > den_end && isdigit(text[den_end]) {
			den_end++
		}
		if den_end == num_end+1 {
			return 0
		}
		if len(text) == den_end || word_boundary(text[den_end]) {
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

func smartypants_cb__number(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	if word_boundary(previous_char) && len(text) >= 3 {
		if text[0] == '1' && text[1] == '/' && text[2] == '2' {
			if len(text) < 4 || word_boundary(text[3]) {
				ob.WriteString("&frac12;")
				return 2
			}
		}

		if text[0] == '1' && text[1] == '/' && text[2] == '4' {
			if len(text) < 4 || word_boundary(text[3]) || (len(text) >= 5 && tolower(text[3]) == 't' && tolower(text[4]) == 'h') {
				ob.WriteString("&frac14;")
				return 2
			}
		}

		if text[0] == '3' && text[1] == '/' && text[2] == '4' {
			if len(text) < 4 || word_boundary(text[3]) || (len(text) >= 6 && tolower(text[3]) == 't' && tolower(text[4]) == 'h' && tolower(text[5]) == 's') {
				ob.WriteString("&frac34;")
				return 2
			}
		}
	}

	ob.WriteByte(text[0])
	return 0
}

func smartypants_cb__dquote(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	next_char := byte(0)
	if len(text) > 1 {
		next_char = text[1]
	}
	if !smartypants_quotes(ob, previous_char, next_char, 'd', &smrt.in_dquote) {
		ob.WriteString("&quot;")
	}

	return 0
}

func smartypants_cb__ltag(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int {
	i := 0

	for i < len(text) && text[i] != '>' {
		i++
	}

	ob.Write(text[:i+1])
	return i
}

type smartypants_cb func(ob *bytes.Buffer, smrt *smartypants_data, previous_char byte, text []byte) int

type SmartypantsRenderer [256]smartypants_cb

func Smartypants(flags int) *SmartypantsRenderer {
	r := new(SmartypantsRenderer)
	r['"'] = smartypants_cb__dquote
	r['&'] = smartypants_cb__amp
	r['\''] = smartypants_cb__squote
	r['('] = smartypants_cb__parens
	if flags&HTML_SMARTYPANTS_LATEX_DASHES == 0 {
		r['-'] = smartypants_cb__dash
	} else {
		r['-'] = smartypants_cb__dash_latex
	}
	r['.'] = smartypants_cb__period
	if flags&HTML_SMARTYPANTS_FRACTIONS == 0 {
		r['1'] = smartypants_cb__number
		r['3'] = smartypants_cb__number
	} else {
		for ch := '1'; ch <= '9'; ch++ {
			r[ch] = smartypants_cb__number_generic
		}
	}
	r['<'] = smartypants_cb__ltag
	r['`'] = smartypants_cb__backtick
	return r
}

func rndr_smartypants(ob *bytes.Buffer, text []byte, opaque interface{}) {
	options := opaque.(*htmlOptions)
	smrt := smartypants_data{false, false}

	mark := 0
	for i := 0; i < len(text); i++ {
		if action := options.smartypants[text[i]]; action != nil {
			if i > mark {
				ob.Write(text[mark:i])
			}

			previous_char := byte(0)
			if i > 0 {
				previous_char = text[i-1]
			}
			i += action(ob, &smrt, previous_char, text[i:])
			mark = i + 1
		}
	}

	if mark < len(text) {
		ob.Write(text[mark:])
	}
}
