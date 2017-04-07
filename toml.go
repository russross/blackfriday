package mmark

import (
	"bytes"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	// Sentinal for not set.
	piNotSet = "__mmark_toml_pi_not_set"

	DefaultIpr  = "trust200902"
	DefaultArea = "Internet"
)

type author struct {
	Initials           string
	Surname            string
	Fullname           string
	Organization       string
	OrganizationAbbrev string `toml:"abbrev"`
	Role               string
	Ascii              string
	Address            address
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
	Region     string
	PostalLine []string

	// Plurals when these need to be specified multiple times.
	Streets   []string
	Cities    []string
	Codes     []string
	Countries []string
	Regions   []string
}

// PIs the processing instructions.
var PIs = []string{"toc", "symrefs", "sortrefs", "compact", "subcompact", "private", "topblock", "header", "footer", "comments"}

type pi struct {
	Toc        string
	Symrefs    string
	Sortrefs   string
	Compact    string
	Subcompact string
	Private    string
	Topblock   string
	Comments   string // Typeset cref's in the text.
	Header     string // Top-Left header, usually Internet-Draft.
	Footer     string // Bottom-Center footer, usually Expires ...
}

type title struct {
	Title  string
	Abbrev string

	DocName        string
	Ipr            string
	Category       string
	Number         int // RFC number
	Obsoletes      []int
	Updates        []int
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

	// Set some sentinels and defaults.
	var block title
	block.PI.Header = piNotSet
	block.PI.Footer = piNotSet
	block.Area = DefaultArea
	block.Ipr = DefaultIpr
	block.Date = time.Now()

	if _, err := toml.Decode(string(data), &block); err != nil {
		printf(p, "error in TOML titleblock: %s", err.Error())
		return block // never an error when encoding markdown
	}
	return block
}
