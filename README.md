Black Friday
============

This is an implementation of John Gruber's [markdown][1] in [Go][2].
It is a translation of the [upskirt][3] library written in C with a
few minor changes. It retains the paranoia of the original (it is
careful not to trust its input, and as such it should be safe to
feed it arbitrary user-supplied inputs). It also retains the
emphasis on high performance, and the source is almost as ugly as
the original.

HTML output is currently supported, along with Smartpants
extensions.


Installation
------------

Assuming you have recent version of Go installed, along with git:

    goinstall github.com/russross/blackfriday

will download, compile, and install the package into
`$GOROOT/src/pkg/github.com/russross/blackfriday`.

Check out `example/main.go` for an example of how to use it. Run
`gomake` in that directory to build a simple command-line markdown
tool:

    cd $GOROOT/src/pkg/github.com/russross/blackfriday
    gomake markdown

will build the binary `markdown` in the `example` directory.


Extensions
----------

In addition to the extensions offered by upskirt, this package
implements two additional Smartypants options:

*   LaTeX-style dash parsing, where `--` is translated into
    `&ndash;`
*   Generic fractions, where anything that looks like a fraction
    is translated into suitable HTML (instead of just a few special
    cases). For example, `4/5` becomes `<sup>4</sup>&frasl;<sub>5</sub>`


Todo
----

*   Code cleanup
*   Better code documentation
*   Implement a LaTeX backend


   [1]: http://daringfireball.net/projects/markdown/ "Markdown"
   [2]: http://golang.org/ "Go Language"
   [3]: http://github.com/tanoku/upskirt "Upskirt"
