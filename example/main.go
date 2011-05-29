//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
//

//
//
// Example front-end for command-line use
//
//

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
    "github.com/russross/blackfriday"
	"os"
)

func main() {
	// read the input
	var input []byte
	var err os.Error
	switch len(os.Args) {
	case 1:
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
			os.Exit(-1)
		}
	case 2, 3:
		if input, err = ioutil.ReadFile(os.Args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from", os.Args[1], ":", err)
			os.Exit(-1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[inputfile [outputfile]]")
		os.Exit(-1)
	}

	// set up options
	output := bytes.NewBuffer(nil)
	var extensions uint32
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	html_flags := 0
	html_flags |= blackfriday.HTML_USE_XHTML
	html_flags |= blackfriday.HTML_USE_SMARTYPANTS
	html_flags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	html_flags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

    // render the data
	blackfriday.Markdown(output, input, blackfriday.HtmlRenderer(html_flags), extensions)

	// output the result
	if len(os.Args) == 3 {
		if err = ioutil.WriteFile(os.Args[2], output.Bytes(), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to", os.Args[2], ":", err)
			os.Exit(-1)
		}
	} else {
		if _, err = os.Stdout.Write(output.Bytes()); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to Stdout:", err)
			os.Exit(-1)
		}
	}
}
