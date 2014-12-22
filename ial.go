// IAL implements inline attributes.

package mmark

import "bytes"
import (
	"sort"
	"strings"
)

// One or more of these can be attached to block elements

type InlineAttr struct {
	id    string            // #id
	class map[string]bool   // 0 or more .class
	attr  map[string]string // key=value pairs
}

func newInlineAttr() *InlineAttr {
	return &InlineAttr{class: make(map[string]bool), attr: make(map[string]string)}
}

// Parsing and thus detecting an IAL. Return a valid *IAL or nil.
// IAL can have #id, .class or key=value element seperated by spaces, that may be escaped
func (p *parser) isInlineAttr(data []byte) int {
	esc := false
	quote := false
	ialB := 0
	ial := newInlineAttr()
	for i := 0; i < len(data); i++ {
		switch data[i] {
		case ' ':
			if quote {
				continue
			}
			chunk := data[ialB+1 : i]
			if len(chunk) == 0 {
				ialB = i
				continue
			}
			switch {
			case chunk[0] == '.':
				ial.class[string(chunk[1:])] = true
			case chunk[0] == '#':
				ial.id = string(chunk[1:])
			default:
				k, v := parseKeyValue(chunk)
				if k != "" {
					ial.attr[k] = v
				}
			}
			ialB = i
		case '"':
			if esc {
				esc = !esc
				continue
			}
			quote = !quote
		case '\\':
			esc = !esc
		case '}':
			if esc {
				esc = !esc
				continue
			}
			// if this is mainmatter, frontmatter, or backmatter it isn't an IAL.
			s := string(data[1:i])
			switch s {
			case "frontmatter":
				fallthrough
			case "mainmatter":
				fallthrough
			case "backmatter":
				return 0
			}
			chunk := data[ialB+1 : i]
			if len(chunk) == 0 {
				return i + 1
			}
			switch {
			case chunk[0] == '.':
				ial.class[string(chunk[1:])] = true
			case chunk[0] == '#':
				ial.id = string(chunk[1:])
			default:
				k, v := parseKeyValue(chunk)
				if k != "" {
					ial.attr[k] = v
				}
			}
			p.ial = p.ial.add(ial)
			return i + 1
		default:
			esc = false
		}
	}
	return 0
}

func parseKeyValue(chunk []byte) (string, string) {
	chunks := bytes.SplitN(chunk, []byte{'='}, 2)
	if len(chunks) != 2 {
		return "", ""
	}
	chunks[1] = bytes.Replace(chunks[1], []byte{'"'}, nil, -1)
	return string(chunks[0]), string(chunks[1])
}

// Add IAL to another, overwriting the #id, collapsing classes and attributes
func (i *InlineAttr) add(j *InlineAttr) *InlineAttr {
	if i == nil {
		return j
	}
	if j.id != "" {
		i.id = j.id
	}
	for k, c := range j.class {
		i.class[k] = c
	}
	for k, a := range j.attr {
		i.attr[k] = a
	}
	return i
}

// String renders an IAL and returns a string that can be included in the tag:
// class="class" anchor="id" key="value". The string s has a space as the first character.k
func (i *InlineAttr) String() (s string) {
	if i == nil {
		return ""
	}

	// some fluff needed to make this all sorted.
	if i.id != "" {
		s = " anchor=\"" + i.id + "\""
	}

	keys := make([]string, 0, len(i.class))
	for k, _ := range i.class {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) > 0 {
		s += " class=\"" + strings.Join(keys, " ") + "\""
	}

	keys = keys[:0]
	for k, _ := range i.attr {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	attr := make([]string, len(keys))
	for j, k := range keys {
		v := i.attr[k]
		attr[j] = k + "=\"" + v + "\""
	}
	if len(keys) > 0 {
		s += " " + strings.Join(attr, " ")
	}
	return s
}

// GetOrDefaultAttr set the value under key if is is not set or
// use the value already in there. The boolean returns indicates
// if the value has been overwritten.
func (i *InlineAttr) GetOrDefaultAttr(key, def string) bool {
	v := i.attr[key]
	if v != "" {
		return false
	}
	if def == "" {
		return false
	}
	i.attr[key] = def
	return true
}

//
func (i *InlineAttr) GetOrDefaultId(id string) bool {
	if i.id != "" {
		return false
	}
	if id == "" {
		return false
	}
	i.id = id
	return true
}
