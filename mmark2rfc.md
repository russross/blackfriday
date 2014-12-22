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
%   [author.address.postal]
%   street = "Buckingham Palace Road"

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
such as RFC and I-Ds.
The mmark syntax is based on the Markdown syntax and has been enhanced with features that are
found in other Markdown implementations like [kramdown], [PHP markdown extra], [pandoc], [leanpub] and
[asciidoc]. 

The goals of mmark can be stated as:

{style="format (%I)"}
1. Self contained: a single file can be converted to XML2RFC v2 or (v3);
2. Make the markdown source code look as natural as possible;
3. Provide seemless upgrade path to XML2RFC v3.

Mmark is a fork of blackfriday [@blackfriday] written in Golang and it is very fast.
Input to mmark must be UTF-8, the output is also UTF-8. Mmark converts tabs to 4 spaces.

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

# Mmark Syntax





# TOML header

Mmark uses TOML [@!toml] document header to specify the document's meta data. Each line of this
header must start with an `% `.

# Citations

A citation can be entered by using the syntax from pandoc [@pandoc]: `[@reference]`,
such a reference is "informative" by default. Making a reference informative or normative
can be done with a `?` and `!` respectively: `[@!reference]` is a normative reference.

For RFC and I-Ds the references are generated automatically, although for I-Ds you might
need to include a draft version in the reference `[@?I-D.draft-blah,#06]`, creates an
informative reference to the sixth version of draft-blah.

If the need arises an XML reference fragment can be included, note that this needs to happen
*before* the back matter is started, because that is the point when the references are outputted
(right now the implementation does not scan the entire file for citations).

# Document divisions

Using `{mainmatter}` on a line by itself starts the main matter (middle) of the document, `{backmatter}`
starts the appendix. There is also a `{frontmatter}` that starts the front matter (front) of the document,
but is normally not need because the TOML header ([](#toml-header)) starts that by default.

# Abstract

Any paragraph prefix with `A> ` is an abstract. This is similar to asides and notes
([](#asides) , [](#notes)) work.

# Captions

Whenever an blockquote, fenced codeblock or image has caption text, the entire block is wrapped
in a `<figure>` and the caption text is put in a `<name>` tag for v3.

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

The reference syntax `(@list-id)` is *not* (yet?) supported.

## HTML Comment

If a HTML comment contains `--`, it will be rendered as a `cref` comment in the resulting
XML file. Typically `<!-- Miek Gieben -- you want to include the next paragraph? -->`.

## Including Files

Files can be included using ``{{filename}}``, `filename` is relative to the current working
directory if it is not absolute.

# XML2RFC V3 features

The v3 syntax adds some new features, those can already be used in mmark (even for documents targetting
v2 -- but there they will be faked with the limited constructs of v2 syntax).

## Asides

Any paragraph prefixed with `AS> `. For v2 this becomes a indentend paragraph.

## Notes

Any paragraph prefixed with `N> `. For v2 this becomes a indentend paragraph.

## RFC 2119 Keywords

Any [@?RFC2119] keyword used with strong emphasis *and* in uppercase  will be typeset
within `bcp14` tags, that is `**MUST**` becomes `<bcp14>MUST</bcp14`, but `**must**` will not.

## Super- and Subscripts

Use H~2~O and 2^10^ is 1024.

## Images


# Converting from RFC 7328 Syntax

Converting from an RFC 7328 ([@!RFC7328]) document can be done using the quick
and dirty [Perl script](https://raw.githubusercontent.com/miekg/mmark/master/convert/parts.pl),
which uses pandoc to output markdown PHP extra and converts that into proper mmark:
(mmark is more like markdown PHP extra, than like pandoc).

    for i in middle.mkd back.mkd; do \
        pandoc --atx-headers -t markdown_phpextra < $i |
        ./parts.pl
    done

Note this:

* Does not convert the abstract to a prefixed paragraph;
* Makes all RFC references normative;
* Handles all figure and table captions and adds references (if appropriate);
* Probably has bugs, so a manual review should be in order.

There is also [titleblock.pl](https://raw.githubusercontent.com/miekg/mmark/master/convert/titleblock.pl)
which can be given an [@RFC7328] `template.xml` file and will output a TOML titleblock, that can
be used as a starting point.

AS> Yes, this uses pandoc and Perl.. why? Becasue if mmark could parse the file by itself, there wasn't much
AS> of problem. Two things are holding this back: mmark cannot parse definition lists with empty spaces and
AS> there isn't renderer that can output markdown syntax.

For now the mmark parser will not get any features that makes it backwards compatible with pandoc2rfc.

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

[kramdown]: http://http://kramdown.gettalong.org/
[leanpub]: https://leanpub.com/help/manual
[asciidoc]: http://www.methods.co.nz/asciidoc/
[PHP markdown extra]: http://michelf.com/projects/php-markdown/extra/
[pandoc]: http://johnmacfarlane.net/pandoc/
