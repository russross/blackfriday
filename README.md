[![Build Status](https://travis-ci.org/miekg/mmark.svg?branch=master)](https://travis-ci.org/miekg/mmark)

Everything that was true of [blackfriday][5], might not be true for mmark anymore.

# Mmark

Mmark is a fork of blackfriday which is a [Markdown][1] processor implemented in
[Go][2]. It supports a number of extensions, inspired by Leanpub, kramdown and
Asciidoc, that allows for large documents to be written. It is specifically
designed to write internet drafts and RFCs for the IETF.

It can currently output HTML5 or XML2RFC v3 XML. Other output engines could be
easily added.

It is fast and it supports the following extensions:

* tables
* definition lists
    lists markers: all from pandoc (roman, letters, and default list markers)
* smart punctuation
* substitutions
* [TOML][3] titleblock
* including other markdown files
* indices
* main-, middle- and backmatter divisions
* Table and Codeblock caption.
    After a codeblock use: Figure: fjfjkfdjkd
* escaping chars (as in pandoc?)
* citations
* abstract
* asides
* IAL, inline attribute list
* pay attention to the first number of ordered list, `4.` start a list with 4. instead of 1.
* EXAMPLE_LISTS (@good) (@good)

Mmark is forked from blackfriday which started out as a translation from C of [upskirt][4].

## Enviroment variables?

`MMARK_REFDIR`: is set prefix
`MMARK_REFONLINE`, use online references

## Usage

For basic usage, it is as simple as getting your input into a byte
slice and calling:

    output := mmark.MarkdownBasic(input)

This renders it with no extensions enabled. To get a more useful
feature set, use this instead:

    output := mmark.MarkdownCommon(input)

Features
--------

All features of upskirt are supported, including:

*   **Compatibility**. The Markdown v1.0.3 test suite passes with
    the `--tidy` option.  Without `--tidy`, the differences are
    mostly in whitespace and entity escaping, where blackfriday is
    more consistent and cleaner.

*   **Common extensions**, including table support, fenced code
    blocks, autolinks, strikethroughs, non-strict emphasis, etc.

*   **Safety**. Blackfriday is paranoid when parsing, making it safe
    to feed untrusted user input without fear of bad things
    happening. The test suite stress tests this and there are no
    known inputs that make it crash.  If you find one, please let me
    know and send me the input that does it.

    NOTE: "safety" in this context means *runtime safety only*. In order to
    protect yourself agains JavaScript injection in untrusted content, see
    [this example](https://github.com/russross/blackfriday#sanitize-untrusted-content).

*   **Fast processing**. It is fast enough to render on-demand in
    most web applications without having to cache the output.

*   **Thread safety**. You can run multiple parsers in different
    goroutines without ill effect. There is no dependence on global
    shared state.

*   **Minimal dependencies**. Blackfriday only depends on standard
    library packages in Go. The source code is pretty
    self-contained, so it is easy to add to any project, including
    Google App Engine projects. And TOML!

*   **Standards compliant**. Output successfully validates using the
    W3C validation tool for HTML 4.01 and XHTML 1.0 Transitional.


Extensions
----------

In addition to the standard markdown syntax, this package
implements the following extensions:

*   **Intra-word emphasis supression**. The `_` character is
    commonly used inside words when discussing code, so having
    markdown interpret it as an emphasis command is usually the
    wrong thing. Blackfriday lets you treat all emphasis markers as
    normal characters when they occur inside a word.

*   **Tables**. Tables can be created by drawing them in the input
    using a simple syntax:

    ```
    Name    | Age
    --------|-----:
    Bob     | 27
    Alice   | 23
    ```

*   **Fenced code blocks**. In addition to the normal 4-space
    indentation to mark code blocks, you can explicitly mark them
    and supply a language (to make syntax highlighting simple). Just
    mark it like this:

        ``` go
        func getTrue() bool {
            return true
        }
        ```

    You can use 3 or more backticks to mark the beginning of the
    block, and the same number to mark the end of the block.

*   **Autolinking**. Blackfriday can find URLs that have not been
    explicitly marked as links and turn them into links.

*   **Strikethrough**. Use two tildes (`~~`) to mark text that
    should be crossed out.

*   **Hard line breaks**. With this extension enabled (it is off by
    default in the `MarkdownBasic` and `MarkdownCommon` convenience
    functions), newlines in the input translate into line breaks in
    the output.

*   **Smart quotes**. Smartypants-style punctuation substitution is
    supported, turning normal double- and single-quote marks into
    curly quotes, etc.

*   **LaTeX-style dash parsing** is an additional option, where `--`
    is translated into `&ndash;`, and `---` is translated into
    `&mdash;`. This differs from most smartypants processors, which
    turn a single hyphen into an ndash and a double hyphen into an
    mdash.

*   **Smart fractions**, where anything that looks like a fraction
    is translated into suitable HTML (instead of just a few special
    cases like most smartypant processors). For example, `4/5`
    becomes `<sup>4</sup>&frasl;<sub>5</sub>`, which renders as
    <sup>4</sup>&frasl;<sub>5</sub>.

*   **Includes**, support including files with `{{filename}}` syntax.

*   **Indices**, using `(((item, subitem)))` syntax.

*   **Citations**, using the citation syntax from pandoc `[@RFC2535 p. 23]`, the citation
    can either be informative (default) or normative, this can be indicated by using
    the `i` or `n` modifer: `[@RFC2535,n]`. Use [`-@RFC1000]` to add the cication to the references, but
    suppress the output in the document.

    To make the references work you can optionally include a filename:
    `[@RFC233,n,bib/reference.RFC.2525.xml]`. If you reference an RFC or ID
    the filename (or URL?) will be contructed automatically.

    A citation is informative by default.

*  **Notes**, any parapgraph with `N>`

*  **Abstracts**, any paragraph with `A>`

*  **Asides**, any paragraph with `AS>` at the beginning of all lines is an aside.

*  **{frontmatter}/{mainmatter}/{backmatter}** Create useful divisions in your document.

*  **IAL**, kramdown's Inline Attribute List syntax, but took the commonMark
    proposal, thus without the colon `{#id .class key=value key="value"}`.

*  **Definitition lists**, the markdown extra syntax, short syntax is not supported (yet).

*  **TOML TitleBlock**, add an extended title block prefixed with % in TOML.

*  **Unique anchors**, make anchors unique by adding sequence numbers to them.

Todo
----

*   Check outputted XML for wellformedness
*   More unit testing
*   pretty print XML
*   fenced code blocks -> source code with language etc.
*   indentend code blocks -> artwork
*   images -> artwork, use title for caption
    if caption is given, wrap in figure -> otherwise not.
*   cleanups - and loose a bunch of extensions, turn them on per default
*   reduce API footprint (hide constants mainly)
*   if some attribute is specified check if the IAL also specifies it (then overule it)
*   use a bunch of env vars to tweak aspects of the mmark binary
*   make webservers that converts for you

License
-------

Mmark is a fork of blackfriday, hence is shares it's license.

Mmark is distributed under the Simplified BSD License:

> Copyright © 2011 Russ Ross
> Copyright © 2014 Miek Gieben
> All rights reserved.
>
> Redistribution and use in source and binary forms, with or without
> modification, are permitted provided that the following conditions
> are met:
>
> 1.  Redistributions of source code must retain the above copyright
>     notice, this list of conditions and the following disclaimer.
>
> 2.  Redistributions in binary form must reproduce the above
>     copyright notice, this list of conditions and the following
>     disclaimer in the documentation and/or other materials provided with
>     the distribution.
>
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
> "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
> LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
> FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
> COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
> INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
> BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
> LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
> CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
> LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
> ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
> POSSIBILITY OF SUCH DAMAGE.


   [1]: http://daringfireball.net/projects/markdown/ "Markdown"
   [2]: http://golang.org/ "Go Language"
   [3]: https://github.com/toml-lang/toml "TOML"
   [4]: http://github.com/tanoku/upskirt "Upskirt"
   [5]: http://github.com/russross/blackfriday "Blackfriday"
