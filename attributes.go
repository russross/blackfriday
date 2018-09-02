package blackfriday

import "strings"

// Attr - Abstraction for html attribute
type Attr []string

// Add - adds one more attribute value
func (a Attr) Add(value string) Attr {
	return append(a, value)
}

// Remove - removes given value from attribute
func (a Attr) Remove(value string) Attr {
	for i := range a {
		if a[i] == value {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func (a Attr) String() string {
	return strings.Join(a, " ")
}

// Attributes - store for many attributes
type Attributes struct {
	attrsMap map[string]Attr
	keys     []string
}

// NewAttributes - creates new Attributes instance
func NewAttributes() *Attributes {
	return &Attributes{
		attrsMap: make(map[string]Attr),
	}
}

// Add - adds attribute if not exists and sets value for it
func (a *Attributes) Add(name, value string) *Attributes {
	if _, ok := a.attrsMap[name]; !ok {
		a.attrsMap[name] = make(Attr, 0)
		a.keys = append(a.keys, name)
	}

	a.attrsMap[name] = a.attrsMap[name].Add(value)
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
		a.attrsMap[name] = attr.Remove(value)
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
