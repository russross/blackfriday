package mmark

import (
	"bytes"
	"time"

	"github.com/BurntSushi/toml"
)

// Sentinal for not set.
const piNotSet = "__mmark_toml_pi_not_set"

type author struct {
	Initials     string
	Surname      string
	Fullname     string
	Organization string
	Role         string
	Ascii        string
	Address      address
}

type address struct {
	Phone  string
	Email  string
	Uri    string
	Postal addressPostal
}

type addressPostal struct {
	Street     string
	City       string
	Code       string
	Country    string
	PostalLine []string
}

// PIs the processing instructions... TODO(miek): more.
var PIs = []string{"toc", "symrefs", "sortrefs", "compact", "subcompact", "private", "topblock", "header", "footer", "comments"}

type pi struct {
	Toc        string
	Symrefs    string
	Sortrefs   string
	Compact    string
	Subcompact string
	Private    string
	Topblock   string
	Comments	string // Typeset cref in the text
	Header	   string // Top-Left header, usually Internet-Draft.
	Footer	   string // Bottom-Center footer, usually Expires ...
}

type title struct {
	Title  string
	Abbrev string

	DocName        string
	Ipr            string
	Category       string
	Obsoletes      []string
	Updates        []string
	PI             pi // Processing Instructions
	SubmissionType string

	Date      time.Time
	Area      string
	Workgroup string
	Keyword   []string
	Author    []author
}

func (p *parser) titleBlockTOML(out *bytes.Buffer, data []byte) title {
	data = bytes.TrimPrefix(data, []byte("%"))
	data = bytes.Replace(data, []byte("\n%"), []byte("\n"), -1)
	var block title
	block.PI.Header = piNotSet
	block.PI.Footer = piNotSet
	if _, err := toml.Decode(string(data), &block); err != nil {
		printf(p, "error in TOML titleblock: %s", err.Error())
		return block // never an error when encoding markdown
	}
	return block
}
