// IAL implements

package blackfriday

// One or more of these can be attached to block elements

type IAL struct {
	id    string            // #id
	class []string          // 0 or more .class
	attr  map[string]string // key=value pairs
}

// parsing and thus detecting one

func NewIAL(data []byte) *IAL {
	return nil
}
