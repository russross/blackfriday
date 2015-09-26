package main

// Parse template.xml and output TOML titleblock.

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

// parseXMLtoTOML parses XML to TOML (in a slightly brain dead way).
func parseXMLtoTOML(input []byte) {
	parser := xml.NewDecoder(bytes.NewReader(input))
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
