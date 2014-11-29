package blackfriday

import (
	"bytes"
	"time"

	"github.com/BurntSushi/toml"
)

type author struct {
	Initials     string 
	Surname      string
	Fullname     string
	Organization string
	Address      address
}

type address struct {
	Phone  string
	Email  string
	Uri    string
	Postal addressPostal
}

type addressPostal struct {
	Street  string
	City    string
	Code    string
	Country string
}

type title struct {
	Title     string
	Abbrev    string
	Date      time.Time
	Area      string
	Workgroup string
	Keyword   []string
	Author    []author
}

func (p *parser) titleBlockTOML(out *bytes.Buffer, data []byte, doRender bool) int {
	data = bytes.TrimPrefix(data, []byte("% "))
	data = bytes.Replace(data, []byte("\n% "), []byte("\n"), -1)
	var block title
	if _, err := toml.Decode(string(data), &block); err != nil {
		println(err.Error())
		return 0 // never an error when encoding markdown
	}
	p.titleblock = block // Not needed to save this value
	return 0
}
