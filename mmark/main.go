//
// Mmark Markdown Processor, based up Blackfriday Markdown Processor
// Available at http://github.com/russross/mmark
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
// Copyright © 2014 Miek Gieben <miek@miek.nl>.

// Example front-end for command-line use

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/miekg/mmark"
)

var githash string

const DEFAULT_TITLE = ""

func main() {
	// parse command-line options
	var page, xml, xml2, toml, rfc7328, version bool
	var css, head string

	flag.BoolVar(&page, "page", false, "generate a standalone HTML page")
	flag.BoolVar(&xml, "xml", false, "generate xml2rfc v3 output")
	flag.BoolVar(&xml2, "xml2", false, "generate xml2rfc v2 output")
	flag.BoolVar(&version, "version", false, "show mmark version")
	flag.StringVar(&css, "css", "", "link to a CSS stylesheet (implies -page)")
	flag.StringVar(&head, "head", "", "link to HTML to be included in head (implies -page)")

	flag.StringVar(&mmark.CitationsID, "bib-id", mmark.CitationsID, "ID bibliography URL")
	flag.StringVar(&mmark.CitationsRFC, "bib-rfc", mmark.CitationsRFC, "RFC bibliography URL")

	flag.BoolVar(&toml, "toml", false, "input file is xml2rfc XML which is convert to TOML titleblock")
	flag.BoolVar(&rfc7328, "rfc7328", false, "parse RFC 7328 style input")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Mmark Markdown Processor"+
			"\nAvailable at http://github.com/miekg/mmark\n\n"+
			"Copyright © 2014 Miek Gieben <miek@miek.nl>\n"+
			"Copyright © 2011 Russ Ross <russ@russross.com>\n"+
			"Distributed under the Simplified BSD License\n\n"+
			"Usage:\n"+
			"  %s [options] [inputfile [outputfile]]\n\n"+
			"Options:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if version{
		if githash != "" {
			githash = "+" + githash
		}
		fmt.Printf("%s%s\n", mmark.Version, githash)
		return
	}

	// enforce implied options
	if css != "" {
		page = true
	}
	if head != "" {
		page = true
	}

	// read the input
	var input []byte
	var err error
	args := flag.Args()
	switch len(args) {
	case 0:
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			log.Fatalf("error reading from standard input: %v", err)
		}
	case 1, 2:
		if input, err = ioutil.ReadFile(args[0]); err != nil {
			log.Fatalf("error reading from %s: %s", args[0], err)
		}
	default:
		flag.Usage()
		return
	}

	// separate mode for parsing XML to TOML
	if toml {
		parseXMLtoTOML(input)
		return
	}

	// set up options
	extensions := 0
	extensions |= mmark.EXTENSION_TABLES
	extensions |= mmark.EXTENSION_FENCED_CODE
	extensions |= mmark.EXTENSION_AUTOLINK
	extensions |= mmark.EXTENSION_SPACE_HEADERS
	extensions |= mmark.EXTENSION_CITATION
	extensions |= mmark.EXTENSION_TITLEBLOCK_TOML
	extensions |= mmark.EXTENSION_HEADER_IDS
	extensions |= mmark.EXTENSION_AUTO_HEADER_IDS
	extensions |= mmark.EXTENSION_UNIQUE_HEADER_IDS
	extensions |= mmark.EXTENSION_FOOTNOTES
	extensions |= mmark.EXTENSION_SHORT_REF
	extensions |= mmark.EXTENSION_INCLUDE
	extensions |= mmark.EXTENSION_PARTS
	extensions |= mmark.EXTENSION_ABBREVIATIONS
	extensions |= mmark.EXTENSION_DEFINITION_LISTS

	if rfc7328 {
		extensions |= mmark.EXTENSION_RFC7328
	}

	var renderer mmark.Renderer
	xmlFlags := 0
	switch {
	case xml:
		if page {
			xmlFlags = mmark.XML_STANDALONE
		}
		renderer = mmark.XmlRenderer(xmlFlags)
	case xml2:
		if page {
			xmlFlags = mmark.XML2_STANDALONE
		}
		renderer = mmark.Xml2Renderer(xmlFlags)
	default:
		// render the data into HTML
		htmlFlags := 0
		if page {
			htmlFlags |= mmark.HTML_COMPLETE_PAGE
		}
		renderer = mmark.HtmlRenderer(htmlFlags, css, head)
	}

	// parse and render
	output := mmark.Parse(input, renderer, extensions).Bytes()

	// output the result
	out := os.Stdout
	if len(args) == 2 {
		if out, err = os.Create(args[1]); err != nil {
			log.Fatalf("error creating %s: %v", args[1], err)
		}
		defer out.Close()
	}

	if _, err = out.Write(output); err != nil {
		log.Fatalf("error writing output:", err)
	}
}
