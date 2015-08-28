package mmark

import (
	"bytes"
	"fmt"
)

// xml2rfc.go contains common code and variables that is shared
// between xml2rfcv[23].go.

var (
	// URLs where we can find the references for IDs and RFCs.
	// These have been know to change. These are the current ones.
	// (2015-08-27).
	CitationsID  = "http://xml2rfc.ietf.org/public/rfc/bibxml3/"
	CitationsRFC = "http://xml2rfc.ietf.org/public/rfc/bibxml/"
)

const (
	referenceRFC      = "reference.RFC."
	referenceID       = "reference.I-D.draft-"
	referenceIDLatest = "reference.I-D."
	ext               = ".xml"
)

// referenceFile creates a .xml filename for the citation c.
// For I-D references like '[@?I-D.ietf-dane-openpgpkey#02]' it will
// create http://<CitationsID>/reference.I-D.draft-ietf-dane-openpgpkey-02.xml
// without an sequence number it becomes:
// http://<CitationsID>/reference.I-D.ietf-dane-openpgpkey.xml
func referenceFile(c *citation) string {
	if len(c.link) < 4 {
		return ""
	}
	switch string(c.link[:3]) {
	case "RFC":
		return CitationsRFC + referenceRFC + string(c.link[3:]) + ext
	case "I-D":
		seq := ""
		if c.seq != -1 {
			seq = "-" + fmt.Sprintf("%02d", c.seq)
			return CitationsID + referenceID + string(c.link[4:]) + seq + ext
		}
		return CitationsID + referenceIDLatest + string(c.link[4:]) + ext
	}
	return ""
}

var entityConvert = map[byte][]byte{
	'<': []byte("&lt;"),
	'>': []byte("&gt;"),
	'&': []byte("&amp;"),
	//	'\'': []byte("&apos;"),
	//	'"': []byte("&quot;"),
}

func writeEntity(out *bytes.Buffer, text []byte) {
	for i := 0; i < len(text); i++ {
		if s, ok := entityConvert[text[i]]; ok {
			out.Write(s)
			continue
		}
		out.WriteByte(text[i])
	}
}

// sanitizeXML strips XML from a string.
func sanitizeXML(s []byte) []byte {
	inTag := false
	j := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			inTag = true
			continue
		}
		if s[i] == '>' {
			inTag = false
			continue
		}
		if !inTag {
			s[j] = s[i]
			j++
		}
	}
	return s[:j]
}

// writeSanitizeXML strips XML from a string and writes
// to out.
func writeSanitizeXML(out *bytes.Buffer, s []byte) {
	inTag := false
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			inTag = true
			continue
		}
		if s[i] == '>' {
			inTag = false
			continue
		}
		if !inTag {
			out.WriteByte(s[i])
		}
	}
}
