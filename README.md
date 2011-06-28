Blackfriday
===========

This is an implementation of John Gruber's [markdown][1] in [Go][2].
It is a translation of the [upskirt][3] library written in C with a
few minor changes. It retains the paranoia of the original (it is
careful not to trust its input, and as such it should be safe to
feed it arbitrary user-supplied inputs). It also retains the
emphasis on high performance, and the source is almost as ugly as
the original.

HTML output is currently supported, along with Smartypants
extensions. An experimental LaTeX output engine is also included.


Installation
------------

Assuming you have recent version of Go installed, along with git:

    goinstall github.com/russross/blackfriday

will download, compile, and install the package into
`$GOROOT/src/pkg/github.com/russross/blackfriday`.

Check out `example/main.go` for an example of how to use it. Run
`gomake` in that directory to build a simple command-line markdown
tool:

    cd $GOROOT/src/pkg/github.com/russross/blackfriday/example
    gomake

will build the binary `markdown` in the `example` directory.


Features
--------

All features of upskirt are supported, including:

*   The Markdown v1.0.3 test suite passes with the `--tidy` option.
    Without `--tidy`, the differences appear to be bugs/dubious
    features in the original, mostly related to whitespace.

*   Common extensions, including table support, fenced code blocks,
    autolinks, strikethroughs, non-strict emphasis, etc.

*   Paranoid parsing, making it safe to feed untrusted used input
    without fear of bad things happening. There are still some
    corner cases that are untested, but it is already more strict
    than upskirt (Go's bounds-checking uncovered a few off-by-one
    errors that were present in the C code).

*   Good performance. I have not done rigorous benchmarking, but
    informal testing suggests it is around 3--4x slower than upskirt.

*   Minimal dependencies. Blackfriday only depends on standard
    library packages in Go. The source code is pretty
    self-contained, so it is easy to add to any project.

*   Output successfully validates using the W3C validation tool for
    HTML 4.01 and XHTML 1.0 Transitional.


Extensions
----------

In addition to the extensions offered by upskirt, this package
implements two additional Smartypants options:

*   LaTeX-style dash parsing, where `--` is translated into
    `&ndash;`, and `---` is translated into `&mdash;`
*   Generic fractions, where anything that looks like a fraction is
    translated into suitable HTML (instead of just a few special
    cases).  For example, `4/5` becomes
    `<sup>4</sup>&frasl;<sub>5</sub>`, which renders as
    <sup>4</sup>&frasl;<sub>5</sub>.


LaTeX Output
------------

A rudimentary LaTeX rendering backend is also included. To see an
example of its usage, see `main.go`:

It renders some basic documents, but is only experimental at this
point. In particular, it does not do any inline escaping, so input
that happens to look like LaTeX code will be passed through without
modification.


Todo
----

*   More unit testing
*   Code cleanup
*   Better code documentation
*   Markdown pretty-printer output engine


License
-------

Blackfriday is distributed under the Simplified BSD License:

> Copyright © 2011 Russ Ross \<russ@russross.com\>. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
> 
> The views and conclusions contained in the software and documentation are those of the
> authors and should not be interpreted as representing official policies, either expressed
> or implied, of the copyright holder.


   [1]: http://daringfireball.net/projects/markdown/ "Markdown"
   [2]: http://golang.org/ "Go Language"
   [3]: http://github.com/tanoku/upskirt "Upskirt"
