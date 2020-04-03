package blackfriday

import "strings"

// attr - Abstraction for html attribute
type attr []string

// Add - adds one more attribute value
func (a attr) add(value string) attr {
	for _, item := range a {
		if item == value {
			return a
		}
	}
	return append(a, value)
}

// Remove - removes given value from attribute
func (a attr) remove(value string) attr {
	for i := range a {
		if a[i] == value {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func (a attr) String() string {
	return strings.Join(a, " ")
}

// Attributes - store for many attributes
type Attributes struct {
	attrsMap map[string]attr
	keys     []string
}

// NewAttributes - creates new Attributes instance
func NewAttributes() *Attributes {
	return &Attributes{
		attrsMap: make(map[string]attr),
	}
}

// Add - adds attribute if not exists and sets value for it
func (a *Attributes) Add(name, value string) *Attributes {
	if _, ok := a.attrsMap[name]; !ok {
		a.attrsMap[name] = make(attr, 0)
		a.keys = append(a.keys, name)
	}

	a.attrsMap[name] = a.attrsMap[name].add(value)
	return a
}

// Remove - removes attribute by name
func (a *Attributes) Remove(name string) *Attributes {
	for i := range a.keys {
		if a.keys[i] == name {
			a.keys = append(a.keys[:i], a.keys[i+1:]...)
		}
	}

	delete(a.attrsMap, name)
	return a
}

// RemoveValue - removes given value from attribute by name
// If given attribues become empty it alose removes entire attribute
func (a *Attributes) RemoveValue(name, value string) *Attributes {
	if attr, ok := a.attrsMap[name]; ok {
		a.attrsMap[name] = attr.remove(value)
		if len(a.attrsMap[name]) == 0 {
			a.Remove(name)
		}
	}
	return a
}

// Empty - checks if attributes is empty
func (a *Attributes) Empty() bool {
	return len(a.keys) == 0
}

func (a *Attributes) String() string {
	r := []string{}
	for _, attrName := range a.keys {
		r = append(r, attrName+"=\""+a.attrsMap[attrName].String()+"\"")
	}

	return strings.Join(r, " ")
}
