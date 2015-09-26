package main

// Parse template.xml and output TOML titleblock.

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

const XML = `<?xml version="1.0" ?>
<!DOCTYPE rfc SYSTEM 'rfc2629.dtd' [
<!ENTITY pandocMiddle PUBLIC '' 'middle.xml'>
<!ENTITY abstract     PUBLIC '' 'abstract.xml'>
<!ENTITY appendix     PUBLIC '' 'appendix.xml'>

<!-- normative -->
<!ENTITY RFC1034    PUBLIC '' 'bib/reference.RFC.1034'>
<!ENTITY RFC2065    PUBLIC '' 'bib/reference.RFC.2065'>
<!ENTITY RFC2308    PUBLIC '' 'bib/reference.RFC.2308'>
<!ENTITY RFC4033    PUBLIC '' 'bib/reference.RFC.4033'>
<!ENTITY RFC4034    PUBLIC '' 'bib/reference.RFC.4034'>
<!ENTITY RFC4035    PUBLIC '' 'bib/reference.RFC.4035'>
<!ENTITY RFC4592    PUBLIC '' 'bib/reference.RFC.4592'>
<!ENTITY RFC4648    PUBLIC '' 'bib/reference.RFC.4648'>
<!ENTITY RFC5155    PUBLIC '' 'bib/reference.RFC.5155'>
<!ENTITY RFC6672    PUBLIC '' 'bib/reference.RFC.6672'>

<!-- informative -->
<!ENTITY RFC2535    PUBLIC '' 'bib/reference.RFC.2535'>
<!ENTITY RFC3655    PUBLIC '' 'bib/reference.RFC.3655'>
<!ENTITY RFC3755    PUBLIC '' 'bib/reference.RFC.3755'>
<!ENTITY RFC4470    PUBLIC '' 'bib/reference.RFC.4470'>
<!ENTITY RFC4956    PUBLIC '' 'bib/reference.RFC.4956'>
<!ENTITY draftdnsnr    PUBLIC '' 'bib/reference.I-D.arends-dnsnr.xml'>
<!ENTITY draftnsec2v2  PUBLIC '' 'bib/reference.I-D.laurie-dnsext-nsec2v2.xml'>
<!ENTITY draftexist    PUBLIC '' 'bib/reference.I-D.ietf-dnsext-not-exsiting-rr.xml'>
<!ENTITY RFC5155Errata PUBLIC '' 'bib/reference.RFC.5155.errata.xml'>
<!ENTITY unbound    PUBLIC '' 'bib/reference.unbound.xml'>
<!ENTITY phreebird     PUBLIC '' 'bib/reference.phreebird.xml'>
]>

<rfc ipr="trust200902" submissionType="independent" category="info" docName="draft-gieben-auth-denial-of-existence-dns-06">
<?rfc toc="yes"?>         <!-- generate a table of contents -->
<?rfc tocompact="no"?>
<?rfc tocdepth="6"?>
<?rfc symrefs="yes"?>     <!-- use anchors instead of numbers for references -->
<?rfc sortrefs="yes" ?>   <!-- alphabetize the references -->
<?rfc rfcedstyle="yes"?>
<?rfc strict="yes"?>
<?rfc autobreaks="yes"?>
<?rfc compact="yes" ?>    <!-- conserve vertical whitespace -->
<?rfc subcompact="no" ?>  <!-- but keep a blank line between list items -->
<front>

<title abbrev="Authenticated Denial in DNS">Authenticated Denial of Existence in the DNS</title>

        <author initials='R.' surname='Gieben' fullname='R. (Miek) Gieben'>
            <organization>Google</organization>
            <address>
                <phone></phone>
                <email>miek@google.com</email>
                <uri></uri>
            </address>
        </author>

        <author initials='W.' surname='Mekking' fullname='W. (Matthijs) Mekking'>
            <organization>NLnet Labs</organization>

            <address>
                <postal>
                    <street>Science Park 400</street>
                    <street></street>
                    <city>Amsterdam</city> <region></region>
                    <code>1098 XH</code>
                    <country>NL</country>
                </postal>

                <phone></phone>
                <email>matthijs@nlnetlabs.nl</email>
                <uri>http://www.nlnetlabs.nl/</uri>
            </address>
        </author>

        <date day="3" month='February' year='2014' />

        <area>Internet</area>
        <keyword>DNSSEC</keyword>
        <keyword>Denial of Existance</keyword>
        <keyword>NSEC</keyword>
        <keyword>NSEC3</keyword>
</front>
</rfc>
`

func main() {
	parser := xml.NewDecoder(strings.NewReader(XML))
	keywords := []string{}
	name := ""
	fmt.Println("% # Quick 'n dirty translated by mmark")
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name = elmt.Name.Local
			switch name {
			case "author":
				fmt.Println("%\n% [[author]]")
				outAttr(elmt.Attr)
			case "rfc":
				fallthrough
			case "title":
				outAttr(elmt.Attr)
			case "address":
				fmt.Println("% [author.address]")
			case "postal":
				fmt.Println("% [author.address.postal]")
			case "date":
				outDate(elmt.Attr)
			}
		case xml.CharData:
			if name == "" {
				continue
			}
			data := xml.CharData(t)
			data = bytes.TrimSpace(data)
			if len(data) == 0 {
				continue
			}
			if name == "keyword" {
				keywords = append(keywords, "\""+string(data)+"\"")
				continue
			}
			outString(name, string(data))
		case xml.EndElement:
			name = ""
		case xml.Comment:
			// don't care
		case xml.ProcInst:
			// don't care
		case xml.Directive:
			// don't care
		default:
		}
	}
	outArray("keyword", keywords)
}

func outString(k, v string) {
	fmt.Printf("%% %s = \"%s\"\n", k, v)
}

func outDate(attr []xml.Attr) {
	/*
		year, month, day := time.Now().Date()
		for _, a := range attr {
			switch a.Name.Local {
			case "day":
				day, err := strconv.Atoi(a.Value)
			case "month":
				_ = a.Value
			case "year":
				year, err := strconv.Atoi(a.Value)
			}
		}
	*/
	fmt.Printf("%%\n%% # TODO date \n%%\n")
}

func outArray(k string, arr []string) {
	all := strings.Join(arr, ", ")
	fmt.Printf("%%\n%% keyword = [ %s ]\n%%", all)
}

func outAttr(attr []xml.Attr) {
	for _, a := range attr {
		outString(a.Name.Local, a.Value)
	}
}
