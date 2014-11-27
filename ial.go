// IAL implements

package blackfriday

// One or more of these can be attached to block elements

type IAL struct {
	id    string            // #id
	class []string          // 0 or more .class
	attr  map[string]string // key=value pairs
}

// Parsing and thus detecting an IAL. Return a valid *IAL or nil.
func NewIAL(data []byte) *IAL {
	return &IAL{id:string(data)}
}
