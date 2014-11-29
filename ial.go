// IAL implements

package blackfriday

// One or more of these can be attached to block elements

type IAL struct {
	id    string            // #id
	class []string          // 0 or more .class
	attr  map[string]string // key=value pairs
}

// Parsing and thus detecting an IAL. Return a valid *IAL or nil.
// we are on the openening brace
func (p *parser) isIAL(data []byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == '}' {
			// if this is mainmatter, frontmatter, or backmatter it
			// isn't an IAL.
			s := string(data[1:i])
			switch s {
			case "frontmatter":
				fallthrough
			case "mainmatter":
				fallthrough
			case "backmatter":
				return 0
			}
			p.ial = append(p.ial, &IAL{id: s})
			return i + 1
		}
	}
	return 0
}
