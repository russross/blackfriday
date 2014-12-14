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

Mmark [@mmark] is a markdown processor, it supports the basic markdown syntax and has been
extended to support more (advanced?) features needed to write larger, structured documents
such as RFC and I-Ds. The goals of mmark can be stated as:

{style="format (%I)"}
1. Self contained: a single file can be converted to XML2RFC v2 or (v3);
2. Make the markdown source code look as natural as possible.

...Fork of blackfriday...

# TOML header

Mmark uses TOML [@!toml] document header to specify the document's meta data. Each line of this
header must start with an `%`.

# Citations

A citation can be entered by using the syntax from pandoc [@!pandoc]: `[@reference]`,
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

Any paragraph prefix with `A> ` is an abstract.

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

## Quotes

After a quote (a paragraph prefixed with `> `) you can add the caption:

    Quote: Name -- URI for attribution

In v3 this is used in the block quote attributes, for v2 it is discarded.

# Inline Attribute Lists

# Miscellaneous

## Example Lists

## Citations

## HTML Comment

If a HTML comment contains `--`, it will be rendered as a `cref` comment in the resulting
XML file. Typically `<!-- Miek Gieben -- you want to include the next paragraph? -->`.

## RFC 2119 Keywords

## Including Files

Files can be included using ``{{filename}}``, `filename` is relative to the current working
directory.

# XML2RFC V3 features

The v3 syntax adds some new features, those can already be used in mmark (even for documents targetting
v2 -- but there they will be faked with the limited constructs of v2 syntax).

## Asides

## Notes

## Images


# Converting from RFC 7328 Syntax

<!-- reference we need to include -->

<reference anchor='mmark' target="http://github.com/miekg/mmark">
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

<reference anchor='toml' target="https://github.com/toml-lang/toml">
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

<reference anchor='pandoc' target="http://johnmacfarlane.net/pandoc/">
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
