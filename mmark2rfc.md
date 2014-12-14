% Title = "Using mmark to create I-Ds and RFCs"
% abbrev = "mmark2rfc"
% category = "info"
% docName = "draft-gieben-mmark2rfc-00"
% ipr= "trust200902"
% date = 2014-12-10T00:00:00Z
% area = "Internet"
% workgroup = ""
% keyword = ["markdown", "xml", "mmark"]
%
% [[author]]
% initials="R."
% surname="Gieben"
% fullname="R. (Miek) Gieben"
% organization = "Google"
%   [author.address]
%   email = "miek@google.com"

A> This document describes an markdown variant called mmark [@?mmark] that can
A> be used to create RFC documents. The aim of mmark is to make writing document
A> as natural as possible, while providing a lot of power on how to structure and layout
A> the document.
A>
A> Of course the
A> [source of this document](https://raw.githubusercontent.com/miekg/mmark/master/mmark2rfc.md)
A> provides an excellent example.

{mainmatter}

# Introduction

Mmark [@mmark] is a markdown processor. It supports the basic markdown syntax and has been
extended to support more (syntax) features needed to write larger, structured documents
such as RFC and I-Ds. The extra syntax elements have been copied from kramdown, leanpub,
markdown extra and Asciidoc.

The goals of mmark can be stated as:

{style="format (%I)"}
1. Self contained: a single file can be converted to XML2RFC v2 or (v3);
2. Make the markdown source code look as natural as possible.

Mmark is a fork of blackfriday [@blackfriday] written in Golang.

Using Figure 1 from [@!RFC7328], mmark can be positioned as follows:

{#fig:mmark}

     +-------------------+   pandoc   +---------+
     | ALMOST PLAIN TEXT |   ------>  | DOCBOOK |
     +-------------------+            +---------+
                   |      \                 |
     non-existent  |       \_________       | xsltproc
       faster way  |         mmark   \      |
                   v                  v     v
           +------------+    xml2rfc  +---------+
           | PLAIN TEXT |  <--------  |   XML   |
           +------------+             +---------+
Figure: Mmark skips the conversion to DOCBOOK and directly outputs XML2RFC XML.

Note that [kramdown-2629](https://github.com/cabo/kramdown-rfc2629) fills the same niche as mmark.

# TOML header

Mmark uses TOML [@!toml] document header to specify the document's meta data. Each line of this
header must start with an `%`.

# Citations

A citation can be entered by using the syntax from pandoc [@pandoc]: `[@reference]`,
such a reference is "informative" by default. Making a reference informative or normative
can be done with a `?` and `!` respectively: `[@!reference]` is a normative reference.

For RFC and I-Ds the references are generated automatically, although for I-Ds you might
need to include a draft version in the reference `[@?I-D.draft-blah,#06]`, creates an
informative reference to the sixth version of draft-blah.

If the need arises an XML reference fragment can be included, not that this needs to happen
before the back matter is started, because that is the point when the references are outputted.

# Document divisions

Using `{mainmatter}` on a line by itself starts the main matter (middle) of the document, `{backmatter}`
starts the appendix. There is also a `{frontmatter}` that starts the front matter (front) of the document,
but is normally not need because the TOML header ([](#toml-header)) starts that by default.

# Abstract

Any paragraph prefix with `A> ` is an abstract. This is similar to asides and notes
([](#asides) , [](#notes)) work.

# Captions

Captions for tables are supported in pandoc, but not for code block. In mmark you can put a
caption under either (even after a fenced code block). Referencing these element (and thus
creating an document `id` for them), is done with an IAL ([](#inline-attribute-lists):

    {#identifier}
    Name    | Age
    --------|-----:
    Bob     | 27
    Alice   | 23

An empty line between the IAL and the table of indented code block is allowed.

## Figures

Any text directly after the figure starting with `Figure: ` is used as the caption.

## Table

A table caption is signalled by using `Table: ` directly after the table.
The table syntax used that one of
[Markdown Extra](https://michelf.ca/projects/php-markdown/extra/#table).

## Quotes

After a quote (a paragraph prefixed with `> `) you can add a caption:

    Quote: Name -- URI for attribution

In v3 this is used in the block quote attributes, for v2 it is discarded.

# Inline Attribute Lists

This borrows from [kramdown][http://kramdown.gettalong.org/syntax.html#block-ials], which
the difference that the colon is dropped and each IAL must be typeset *before* the block element.
Added an anchor to blockquote can be done like so:

    {#quote:ref1}
    > A block quote

You can specify classes with `.class` (although these are not used when converting to XML2RFC), and
arbitrary key value pairs where each key becomes an attribute.

# Miscellaneous

## Example Lists

This is the example list syntax
[from pandoc](http://johnmacfarlane.net/pandoc/README.html#extension-example_lists).

The reference syntax `(@list-id)` is *not* supported.

## HTML Comment

If a HTML comment contains `--`, it will be rendered as a `cref` comment in the resulting
XML file. Typically `<!-- Miek Gieben -- you want to include the next paragraph? -->`.

## RFC 2119 Keywords

Any [@?RFC2119] keyword used with strong emphasis *and* in uppercase  will be typeset
within `bcp14` tags, that is `**MUST**` becomes `<bcp14>MUST</bcp14`, but `**must**` will not.

## Including Files

Files can be included using ``{{filename}}``, `filename` is relative to the current working
directory it not absolute.

# XML2RFC V3 features

The v3 syntax adds some new features, those can already be used in mmark (even for documents targetting
v2 -- but there they will be faked with the limited constructs of v2 syntax).

## Asides

Any paragraph prefixed with `AS> `.

## Notes

Any paragraph prefixed with `N> `.

## Images


# Converting from RFC 7328 Syntax

> The author is pondering an automated conversion mechanism.

The markdown syntax in [@!RFC7328] is slightly more liberal than the one from mmark.

## Citations

Citations in [@RFC7238] are done by using references: `[](#RFC5155)` in mmark you use
a proper citation `[@RFC5155]` (with possibly some metadata).

## Definition Lists

Pandoc allows an empty line between the term and the definition:

    Original owner name:

    :   the owner name corresponding to a hashed owner name if hashing is used. Or
        the owner name as-is if no hashing is used.

Mmark does not, use:

    Original owner name:
    :   the owner name corresponding to a hashed owner name if hashing is used. Or
        the owner name as-is if no hashing is used.

## Figure Captions

Instead of using the footnote syntax below a figure, in mmark you can just use 'Figure: '.

                         1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 3 3
     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |   Hash Alg.   |     Flags     |  Iterations   | Salt Length   |
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    |                             Salt                              /
    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    ^[fig:dnsnxtparam-wire::The NEXTPARAM on-the-wire format.]

And use an IAL to define an anchor:

    {#fig:dnsnxtparam-wire}

                            1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 3 3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
        |   Hash Alg.   |     Flags     |  Iterations   | Salt Length   |
        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
        |                             Salt                              /
        +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
    Figure: The NEXTPARAM on-the-wire format.

## Table Caption

Table captions work identically to figure captions, except that you use: `Table: Caption text.`

## Indicess

Use `(((term, subterm)))` instead of using footnotes. Making `term` a primary term, is done
with `(((!term, subterm)))`.

## Ordered Lists with Custom Counters

Use an IAL for that:

    {format="REQ(%c)"}
    1. Term1
    1. Term2



<!-- reference we need to include -->

<reference anchor='mmark' target='http://github.com/miekg/mmark'>
    <front>
        <title abbrev='mmark'>Mmark git repository</title>
        <author initials='R.' surname='Gieben' fullname='R. (Miek) Gieben'>
            <address>
                <email>miek@miek.nl</email>
            </address>
        </author>
        <date year='2014' month='December'/>
    </front>
</reference>

<reference anchor='blackfriday' target='http://github.com/russross/blackfriday'>
    <front>
        <title abbrev='mmark'>Blackfriday git repository</title>
        <author initials='' surname='' fullname=''>
            <address>
                <email>miek@miek.nl</email>
            </address>
        </author>
        <date year='2011' month='November'/>
    </front>
</reference>

<reference anchor='toml' target='https://github.com/toml-lang/toml'>
    <front>
        <title abbrev='mmark'>TOML git repository</title>
        <author initials='T.' surname='Preston-Werner' fullname='Tom Preston-Werner'>
            <address>
                <email></email>
                </address>
            </author>
        <date year='2013' month='March' />
    </front>
</reference>

<reference anchor='pandoc' target='http://johnmacfarlane.net/pandoc/'>
    <front>
        <title>Pandoc, a universal document converter</title>
        <author initials='J.' surname='MacFarlane' fullname='John MacFarlane'>
            <organization>University of California, Berkeley</organization>
            <address>
                <email>jgm@berkeley.edu</email>
                <uri>http://johnmacfarlane.net/</uri>
            </address>
        </author>
        <date year='2006' />
    </front>
</reference>

{backmatter}
