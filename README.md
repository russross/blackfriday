[![Build Status](https://travis-ci.org/miekg/mmark.svg?branch=master)](https://travis-ci.org/miekg/mmark)

Everything that was true of [blackfriday][5], might not be true for mmark anymore.

# Mmark

Mmark is a fork of blackfriday which is a [Markdown][1] processor implemented in
[Go][2]. It supports a number of extensions, inspired by Leanpub, kramdown and
Asciidoc, that allows for large documents to be written. It is specifically
designed to write internet drafts and RFCs for the IETF. With mmark you can create
a single file that serves as input into the XML2RFC processor.

It can currently output HTML5, XML2RFC v2 and XML2RFC v3 XML. Other output
engines could be easily added.

It adds the following syntax elements to [black friday](https://github.com/russross/blackfriday/blob/master/README.md):

* Definition lists;
* Table and codeblock captions;
* Table footer;
* Quote attribution;
* Including other files;
* [TOML][3] titleblock;
* Inline Attribute Lists;
* Indices;
* Citations;
* Abstract;
* Asides;
* Notes;
* Main-, middle- and backmatter divisions;
* Example lists;
* HTML Comment parsing;
* BCP14 (RFC2119) keyword detection;
* Include raw XML references;
* Abbreviations;
* Super- and subscript.

Mmark is forked from blackfriday which started out as a translation from C of [upskirt][4].

A simular effort is [kramdown-rfc2629](https://github.com/cabo/kramdown-rfc2629) from Carsten Bormann.

There is no pretty printed output if you need that pipe the output through `xmllint --format -`.

## Usage

For basic usage, it is as simple as getting your input into a byte
slice and calling:

    output := mmark.MarkdownBasic(input)

This renders it with no extensions enabled. To get a more useful
feature set, use this instead:

    output := mmark.MarkdownCommon(input)


# Extensions

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

    Tables can also have a footer, use equal signs instead of dashes.
    If there are multiple footer line, the first one is used as a
    starting point for the table footer.

    ```
    Name    | Age
    --------|-----:
    Bob     | 27
    Alice   | 23
    ======= | ====
    Charlie | 4
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

*   **Indices**, using `(((item, subitem)))` syntax. To make `item` primary, use
    an `!`: `(((!item, subitem)))`.

*   **Citations**, using the citation syntax from pandoc `[@RFC2535 p. 23]`, the
    citation can either be informative (default) or normative, this can be indicated
    by using the `?` or `!` modifer: `[@!RFC2535]`. Use `[-@RFC1000]` to add the
    cication to the references, but suppress the output in the document.

    If you reference an RFC or I-D the reference will be contructed
    automatically. For I-Ds you may need to add a draft sequence number, which
    can be done as such: `[@?I-D.draft-blah,#06]`. If you have other references
    you can include the raw XML in the document (before the `{backmatter}`).
    Also see **XML references**.

*  **Captions**, table and figure/code block captions. For tables add the string
    `Table: caption text` after the table, this will be rendered as an caption. For
    code blocks you'll need to use `Figure: `

    ```
    Name    | Age
    --------|-----:
    Bob     | 27
    Alice   | 23
    Table: This is a table.
    ```

*  **Quote attribution**, after a blockquote you can optionally use
    `Quote: John Doe -- http://example.org`, where
    the quote will be attributed to John Doe, pointing to the URL.

*  **Notes**, any parapgraph with `N>`

*  **Abstracts**, any paragraph with `A>`

*  **Asides**, any paragraph with `AS>` at the beginning of all lines is an aside.

*  **{frontmatter}/{mainmatter}/{backmatter}** Create useful divisions in your document.

*  **IAL**, kramdown's Inline Attribute List syntax, but took the CommonMark
    proposal, thus without the colon after the brace `{#id .class key=value key="value"}`.

*  **Definitition lists**, the markdown extra syntax.

*  **TOML TitleBlock**, add an extended title block prefixed with `%` in TOML.

*  **Unique anchors**, make anchors unique by adding sequence numbers (-1, -2, etc.) to them.
    All numeric section get an anchor prefixed with `section-`.

*  **Example lists**, a list that is started with `(@good)` is subsequently numbered throughout
    the document. First use is rendered `(1)`, the second one `(2)` and so on.

*  **HTML comments** An HTML comment in the form of `<!-- Miek Gieben: really
    -->` is detected and will be converted to a `cref` with the `source` attribute
    set to "Miek Gieben" and the comment text set to "really".

*  **XML references** Any XML reference fragment included *before* the back matter, can be used
    as a citation reference.

*  **BCP 14** If a RFC 2119 word is found enclosed in `**` it will be rendered as an `<bcp14>`
    element: `**MUST**` becomes `<bcp14>MUST</bcp14>`.

*  **Abbreviations**: See <https://michelf.ca/projects/php-markdown/extra/#abbr>, any text
    defined by:

        *[HTML]: Hyper Text Markup Language

    Allows you to use HTML in the document and it will be expanded to
    `<abbr title="Hyper Text Markup Language">HTML</abbr>`.

* **Super and subscripts**, for superscripts use '^' and for subscripts use '~'. For example:

        H~2~O is a liquid. 2^10^ is 1024.

    Inside a sub/superscript you must escape spaces.
    Thus, if you want the letter P with 'a cat' in subscripts, use P~a\ cat~, not P~a cat~.

# Todo

*   fenced code blocks -> source code with language etc.
*   indentend code blocks -> artwork
*   images -> artwork, use title for caption
    if caption is given, wrap in figure -> otherwise not.
*   Extension to recognize pandoc2rfc indices?

*   make webservers that converts for you
*   cleanups - and loose a bunch of extensions, turn them on per default
    reduce API footprint (hide constants mainly)
*   save original IAL for example lists?
*   add ULink?

# License

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
