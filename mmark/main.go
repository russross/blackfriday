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
	"os"
	"runtime/pprof"
	"strings"

	"github.com/miekg/mmark"
)

const DEFAULT_TITLE = ""

func main() {
	// parse command-line options
	var page, toc, toconly, xhtml, xml, xml2, commonmark bool
	var css, cpuprofile string
	var repeat int
	flag.BoolVar(&page, "page", false, "generate a standalone HTML page")
	flag.BoolVar(&toc, "toc", false, "generate a table of contents (implies -xml=false)")
	flag.BoolVar(&toconly, "toconly", false, "generate a table of contents only (implies -toc)")
	flag.BoolVar(&xhtml, "xhtml", true, "use XHTML-style tags in HTML output")
	flag.BoolVar(&xml, "xml", false, "generate XML2RFC v3 output")
	flag.BoolVar(&xml2, "xml2", false, "generate XML2RFC v2 output")
	flag.BoolVar(&commonmark, "commonmark", false, "input is commonmark")
	flag.StringVar(&css, "css", "", "link to a CSS stylesheet (implies -page)")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to a file")
	flag.IntVar(&repeat, "repeat", 1, "process the input multiple times (for benchmarking)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Mmark Markdown Processor v"+mmark.VERSION+
			"\nAvailable at http://github.com/miekg/mmark\n\n"+
			"Copyright © 2014 Miek Gieben <miek@miek.nl>\n"+
			"Copyright © 2011 Russ Ross <russ@russross.com>\n"+
			"Distributed under the Simplified BSD License\n"+
			"See website for details\n\n"+
			"Usage:\n"+
			"  %s [options] [inputfile [outputfile]]\n\n"+
			"Options:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// enforce implied options
	if css != "" {
		page = true
	}
	if toconly {
		toc = true
	}
	if toc {
		xml = false
	}

	// turn on profiling?
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// read the input
	var input []byte
	var err error
	args := flag.Args()
	switch len(args) {
	case 0:
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, "error reading from Stdin:", err)
			os.Exit(-1)
		}
	case 1, 2:
		if input, err = ioutil.ReadFile(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "error reading from", args[0], ":", err)
			os.Exit(-1)
		}
	default:
		flag.Usage()
		os.Exit(-1)
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
	extensions |= mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	extensions |= mmark.EXTENSION_INCLUDE

	if commonmark {
		extensions &= ^mmark.EXTENSION_AUTO_HEADER_IDS
		extensions &= ^mmark.EXTENSION_AUTOLINK
		extensions |= mmark.EXTENSION_LAX_HTML_BLOCKS
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
		if xhtml {
			htmlFlags |= mmark.HTML_USE_XHTML
		}
		title := ""
		if page {
			htmlFlags |= mmark.HTML_COMPLETE_PAGE
			title = getTitle(input)
		}
		if toconly {
			htmlFlags |= mmark.HTML_OMIT_CONTENTS
		}
		if toc {
			htmlFlags |= mmark.HTML_TOC
		}
		renderer = mmark.HtmlRenderer(htmlFlags, title, css)
	}

	// parse and render
	var output []byte
	for i := 0; i < repeat; i++ {
		// TODO(miek): io.Copy
		output = mmark.Parse(input, renderer, extensions).Bytes()
	}

	// output the result
	var out *os.File
	if len(args) == 2 {
		if out, err = os.Create(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "error creating %s: %v", args[1], err)
			os.Exit(-1)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if _, err = out.Write(output); err != nil {
		fmt.Fprintln(os.Stderr, "error writing output:", err)
		os.Exit(-1)
	}
}

// try to guess the title from the input buffer
// just check if it starts with an <h1> element and use that
func getTitle(input []byte) string {
	i := 0

	// skip blank lines
	for i < len(input) && (input[i] == '\n' || input[i] == '\r') {
		i++
	}
	if i >= len(input) {
		return DEFAULT_TITLE
	}
	if input[i] == '\r' && i+1 < len(input) && input[i+1] == '\n' {
		i++
	}

	// find the first line
	start := i
	for i < len(input) && input[i] != '\n' && input[i] != '\r' {
		i++
	}
	line1 := input[start:i]
	if input[i] == '\r' && i+1 < len(input) && input[i+1] == '\n' {
		i++
	}
	i++

	// check for a prefix header
	if len(line1) >= 3 && line1[0] == '#' && (line1[1] == ' ' || line1[1] == '\t') {
		return strings.TrimSpace(string(line1[2:]))
	}

	// check for an underlined header
	if i >= len(input) || input[i] != '=' {
		return DEFAULT_TITLE
	}
	for i < len(input) && input[i] == '=' {
		i++
	}
	for i < len(input) && (input[i] == ' ' || input[i] == '\t') {
		i++
	}
	if i >= len(input) || (input[i] != '\n' && input[i] != '\r') {
		return DEFAULT_TITLE
	}

	return strings.TrimSpace(string(line1))
}
