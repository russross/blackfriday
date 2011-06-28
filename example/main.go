//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
//
// Example front-end for command-line use
//
//

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"github.com/russross/blackfriday"
	"os"
	"runtime/pprof"
)

func main() {
	// parse command-line options
	var page, xhtml, latex, smartypants, latexdashes, fractions bool
	var css, cpuprofile string
	var repeat int
	flag.BoolVar(&page, "page", false,
		"Generate a standalone HTML page (implies -latex=false)")
	flag.BoolVar(&xhtml, "xhtml", true,
		"Use XHTML-style tags in HTML output")
	flag.BoolVar(&latex, "latex", false,
		"Generate LaTeX output instead of HTML")
	flag.BoolVar(&smartypants, "smartypants", true,
		"Apply smartypants-style substitutions")
	flag.BoolVar(&latexdashes, "latexdashes", true,
		"Use LaTeX-style dash rules for smartypants")
	flag.BoolVar(&fractions, "fractions", true,
		"Use improved fraction rules for smartypants")
	flag.StringVar(&css, "css", "",
		"Link to a CSS stylesheet (implies -page)")
	flag.StringVar(&cpuprofile, "cpuprofile", "",
		"Write cpu profile to a file")
	flag.IntVar(&repeat, "repeat", 1,
		"Process the input multiple times (for benchmarking)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Blackfriday Markdown Processor v"+blackfriday.VERSION+
			"\nAvailable at http://github.com/russross/blackfriday\n\n"+
			"Copyright © 2011 Russ Ross <russ@russross.com>\n"+
			"Distributed under the Simplified BSD License\n"+
			"See website for details\n\n"+
			"Usage:\n"+
			"  %s [options] [inputfile [outputfile]]\n\n"+
			"Options:\n",os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// enforce implied options
	if css != "" {
		page = true
	}
	if page {
		latex = false
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
	var err os.Error
	args := flag.Args()
	switch len(args) {
	case 0:
		if input, err = ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
			os.Exit(-1)
		}
	case 1, 2:
		if input, err = ioutil.ReadFile(args[0]); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from", args[0], ":", err)
			os.Exit(-1)
		}
	default:
		flag.Usage()
		os.Exit(-1)
	}

	// set up options
	var extensions uint32
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	var renderer *blackfriday.Renderer
	if latex {
		// render the data into LaTeX
		renderer = blackfriday.LatexRenderer(0)
	} else {
		// render the data into HTML
		html_flags := 0
		if xhtml {
			html_flags |= blackfriday.HTML_USE_XHTML
		}
		if smartypants {
			html_flags |= blackfriday.HTML_USE_SMARTYPANTS
		}
		if fractions {
			html_flags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
		}
		if latexdashes {
			html_flags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
		}
		renderer = blackfriday.HtmlRenderer(html_flags)
	}

	// parse and render
	var output []byte
	for i := 0; i < repeat; i++ {
		output = blackfriday.Markdown(input, renderer, extensions)
	}

	// output the result
	var out *os.File
	if len(args) == 2 {
		if out, err = os.Create(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", args[1], err)
			os.Exit(-1)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if page {
		// if it starts with an <h1>, make that the title
		title := ""
		if bytes.HasPrefix(output, []byte("<h1>")) {
			end := 0
			// we know the buffer ends with a newline, so no need to check bounds
			for output[end] != '\n' {
				end++
			}
			if bytes.HasSuffix(output[:end], []byte("</h1>")) {
				title = string(output[len("<h1>") : end-len("</h1>")])
			}
		}

		ending := ""
		if xhtml {
			fmt.Fprint(out, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" ")
			fmt.Fprintln(out, "\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">")
			fmt.Fprintln(out, "<html xmlns=\"http://www.w3.org/1999/xhtml\">")
			ending = " /"
		} else {
			fmt.Fprint(out, "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" ")
			fmt.Fprintln(out, "\"http://www.w3.org/TR/html4/strict.dtd\">")
			fmt.Fprintln(out, "<html>")
		}
		fmt.Fprintln(out, "<head>")
		fmt.Fprintf(out, "  <title>%s</title>\n", title)
		fmt.Fprintf(out, "  <meta name=\"GENERATOR\" content=\"Blackfriday Markdown Processor v%s\"%s>\n",
			blackfriday.VERSION, ending)
		fmt.Fprintf(out, "  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"%s>\n",
			ending)
		if css != "" {
			fmt.Fprintf(out, "  <link rel=\"stylesheet\" type=\"text/css\" href=\"%s\" />\n", css)
		}
		fmt.Fprintln(out, "</head>")
		fmt.Fprintln(out, "<body>")
	}
	if _, err = out.Write(output); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing output:", err)
		os.Exit(-1)
	}
	if page {
		fmt.Fprintln(out, "</body>")
		fmt.Fprintln(out, "</html>")
	}
}
