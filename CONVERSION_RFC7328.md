# Converting From RFC 7328 Syntax

Mmark can not directly parse an RFC 7328 style document, but pandoc can, and
pandoc can output a document that *can* be parsed by mmark.

The following (long-ish) commandline allows pandoc to parse the document and output
something mmark can grok (main use here is to convert tables to the mmark table
format):

    cat YOURFILE.md | \
    pandoc --atx-headers -f markdown_phpextra+simple_tables+multiline_tables+grid_tables+superscript \
        -t markdown_phpextra+superscript | \
        sed 's|\\^\[|\^\[|' > YOURFILE_mmark.md

This should deal with most of the constructs used in your original pandoc2rfc
file. But be aware of the following limitations:

* Indices (RFC 7328 Section 6.4), are detected *and* parsed.
* Captions and anchors (Section 6.3) are detected *but not* parsed, instead
  a warning is given that the text should be reworked.
* Makes all RFC references normative.
* Abstract needs to be moved to a `.# Abstract` section.

This leaves the title block, i.e. the `template.xml` from pandoc2rfc, which
should be converted to a TOML titleblock, use `mmark -toml template.xml` for
this. It should output a(n) (in)correct TOML title block that can be used as
a starting point.

## Example

To take a live example, consider [this draft middle piece](https://raw.githubusercontent.com/miekg/denialid/master/middle.mkd),
which later became [RFC 7129](https://tools.ietf.org/html/rfc7129):

    % curl https://raw.githubusercontent.com/miekg/denialid/master/middle.mkd | \
     pandoc --atx-headers \
        -f markdown_phpextra+simple_tables+multiline_tables+grid_tables+superscript \
        -t markdown_phpextra+superscript | \
     sed 's|\\^\[|\^\[|'  | ./mmark/mmark -rfc7328 -xml2

With warnings about the captions and anchors for figures:

    mmark: rfc 7328 style anchor seen: consider adding '{#fig:the-unsigned}' IAL before the figure/table
    mmark: rfc 7328 style caption seen: consider adding 'Figure: The unsigned "example.org" zone.' or 'Table: The unsigned "example.org" zone.' after the figure/table
    mmark: rfc 7328 style anchor seen: consider adding '{#fig:the-nsec-r}' IAL before the figure/table
    mmark: rfc 7328 style caption seen: consider adding 'Figure: The NSEC records of "example.org". The arrows represent NSEC records, starting from the apex.' or 'Table: The NSEC records of "example.org". The arrows represent NSEC records, starting from the apex.' after the figure/table
    ....

## Why Convert?

If you want to submit your document in XML2RFC v3 (unlikely right now), or just
want to make use of the speed, convenience (everything in one document) and extra
features of mmark.

If you are happy with pandoc2rfc then there is no need to convert.
