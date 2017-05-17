package mmark

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// xml2rfc.go contains common code and variables that is shared
// between xml2rfcv[23].go.

var (
	// These have been known to change, these are the current ones (2015-08-27).
	CitationsBase = "https://xml2rfc.tools.ietf.org/public/rfc/"

	// CitationsRFC is the URL where the citations for RFCs are.
	CitationsRFC = CitationsBase + "bibxml/"

	// CitationsANSI is the URL where the citations for ANSI documents are.
	CitationsANSI = CitationsBase + "bibxml2/"

	// CitationsID is the URL where the citations for I-Ds are.
	CitationsID = CitationsBase + "bibxml3/"

	// CitationsW3C is the URL where the citations for W3C documents are.
	CitationsW3C = CitationsBase + "bibxml4/"
)

const (
	referenceRFC      = "reference.RFC."
	referenceID       = "reference.I-D.draft-"
	referenceIDLatest = "reference.I-D."
	referenceW3C      = "reference."
	referenceANSI     = "reference."
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
	switch {
	case bytes.HasPrefix(c.link, []byte("RFC")):
		return CitationsRFC + referenceRFC + string(c.link[3:]) + ext
	case bytes.HasPrefix(c.link, []byte("I-D")):
		seq := ""
		if c.seq != -1 {
			seq = "-" + fmt.Sprintf("%02d", c.seq)
			return CitationsID + referenceID + string(c.link[4:]) + seq + ext
		}
		return CitationsID + referenceIDLatest + string(c.link[4:]) + ext
	case bytes.HasPrefix(c.link, []byte("W3C")):
		return CitationsW3C + referenceW3C + string(c.link) + ext
	case bytes.HasPrefix(c.link, []byte("ANSI")):
		fallthrough
	case bytes.HasPrefix(c.link, []byte("CCITT")):
		fallthrough
	case bytes.HasPrefix(c.link, []byte("FIPS")):
		fallthrough
	case bytes.HasPrefix(c.link, []byte("IEEE")):
		fallthrough
	case bytes.HasPrefix(c.link, []byte("ISO")):
		fallthrough
	case bytes.HasPrefix(c.link, []byte("ITU")):
		fallthrough
	case bytes.HasPrefix(c.link, []byte("PKCS")):
		return CitationsANSI + referenceANSI + string(c.link) + ext

	}
	return ""
}

// countCitationsAndSort returns the number of informative and normative
// references and a string slice with the sorted keys.
func countCitationsAndSort(citations map[string]*citation) (int, int, []string) {
	keys := make([]string, 0, len(citations))
	refi, refn := 0, 0
	for k, c := range citations {
		if c.typ == 'i' {
			refi++
		}
		if c.typ == 'n' {
			refn++
		}

		keys = append(keys, k)
	}
	sort.Strings(keys)
	return refi, refn, keys
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

// titleBlockTOMLAuthor outputs the author from the TOML title block.
func titleBlockTOMLAuthor(out *bytes.Buffer, a author) {
	out.WriteString("<author")

	if a.Role != "" {
		out.WriteString(" role=\"")
		writeEntity(out, []byte(a.Role))
		out.WriteString("\"")
	}

	out.WriteString(" initials=\"")
	writeEntity(out, []byte(a.Initials))
	out.WriteString("\"")

	out.WriteString(" surname=\"")
	writeEntity(out, []byte(a.Surname))
	out.WriteString("\"")

	out.WriteString(" fullname=\"")
	writeEntity(out, []byte(a.Fullname))
	out.WriteString("\">\n")

	abbrev := ""
	if a.OrganizationAbbrev != "" {
		abbrev = " abbrev=\"" + a.OrganizationAbbrev + "\""
	}
	out.WriteString("<organization")
	writeEntity(out, []byte(abbrev))
	out.WriteString(">")
	writeEntity(out, []byte(a.Organization))
	out.WriteString("</organization>\n")

	out.WriteString("<address>\n")
	out.WriteString("<postal>\n")

	out.WriteString("<street>")
	writeEntity(out, []byte(a.Address.Postal.Street))
	out.WriteString("</street>\n")
	for _, street := range a.Address.Postal.Streets {
		out.WriteString("<street>")
		writeEntity(out, []byte(street))
		out.WriteString("</street>\n")
	}

	out.WriteString("<city>")
	writeEntity(out, []byte(a.Address.Postal.City))
	out.WriteString("</city>\n")
	for _, city := range a.Address.Postal.Cities {
		out.WriteString("<city>")
		writeEntity(out, []byte(city))
		out.WriteString("</city>\n")
	}

	out.WriteString("<code>")
	writeEntity(out, []byte(a.Address.Postal.Code))
	out.WriteString("</code>\n")
	for _, code := range a.Address.Postal.Codes {
		out.WriteString("<code>")
		writeEntity(out, []byte(code))
		out.WriteString("</code>\n")
	}

	out.WriteString("<country>")
	writeEntity(out, []byte(a.Address.Postal.Country))
	out.WriteString("</country>\n")
	for _, country := range a.Address.Postal.Countries {
		out.WriteString("<country>")
		writeEntity(out, []byte(country))
		out.WriteString("</country>\n")
	}

	out.WriteString("<region>")
	writeEntity(out, []byte(a.Address.Postal.Region))
	out.WriteString("</region>\n")
	for _, region := range a.Address.Postal.Regions {
		out.WriteString("<region>")
		writeEntity(out, []byte(region))
		out.WriteString("</region>\n")
	}

	out.WriteString("</postal>\n")

	out.WriteString("<phone>" + a.Address.Phone + "</phone>\n")
	out.WriteString("<email>" + a.Address.Email + "</email>\n")
	out.WriteString("<uri>" + a.Address.Uri + "</uri>\n")

	out.WriteString("</address>\n")
	out.WriteString("</author>\n")
}

// titleBlockTOMLDate outputs the date from the TOML title block.
func titleBlockTOMLDate(out *bytes.Buffer, d time.Time) {
	year := ""
	if d.Year() > 0 {
		year = " year=\"" + strconv.Itoa(d.Year()) + "\""
	}
	month := ""
	if d.Month() > 0 {
		month = " month=\"" + time.Month(d.Month()).String() + "\""
	}
	day := ""
	if d.Day() > 0 {
		day = " day=\"" + strconv.Itoa(d.Day()) + "\""
	}
	out.WriteString("<date" + year + month + day + "/>\n\n")
}

// titleBlockTOMLKeyword outputs the keywords from the TOML title block.
func titleBlockTOMLKeyword(out *bytes.Buffer, keywords []string) {
	for _, k := range keywords {
		out.WriteString("<keyword>" + k + "</keyword>\n")
	}
}

// titleBlockTOMLPI returns "yes" or "no" or a stringified number
// for use as process instruction. If version is 3 they are returned
// as attributes for use *inside* the <rfc> tag.
func titleBlockTOMLPI(pi pi, name string, version int) string {
	if version == 2 {
		switch name {
		case "toc":
			return "<?rfc toc=\"" + yesno(pi.Toc, "yes") + "\"?>\n"
		case "symrefs":
			return "<?rfc symrefs=\"" + yesno(pi.Symrefs, "yes") + "\"?>\n"
		case "sortrefs":
			return "<?rfc sortrefs=\"" + yesno(pi.Sortrefs, "yes") + "\"?>\n"
		case "compact":
			return "<?rfc compact=\"" + yesno(pi.Compact, "yes") + "\"?>\n"
		case "topblock":
			return "<?rfc topblock=\"" + yesno(pi.Topblock, "yes") + "\"?>\n"
		case "comments":
			return "<?rfc comments=\"" + yesno(pi.Comments, "no") + "\"?>\n"
		case "subcompact":
			return "<?rfc subcompact=\"" + yesno(pi.Subcompact, "no") + "\"?>\n"
		case "private":
			return "<?rfc private=\"" + yesno(pi.Private, "") + "\"?>\n"
		case "header":
			if pi.Header == piNotSet {
				return ""
			}
			return "<?rfc header=\"" + pi.Header + "\"?>\n"
		case "footer":
			if pi.Footer == piNotSet {
				return ""
			}
			return "<?rfc footer=\"" + pi.Footer + "\"?>\n"
		default:
			printf(nil, "unhandled or unknown PI seen: %s", name)
			return ""
		}
	}
	// version 3
	return ""
}

func yesno(s, def string) string {
	if s == "" {
		return def
	}
	if s == "yes" {
		return "yes"
	}
	return "no"
}
