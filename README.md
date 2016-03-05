# Mmark

Mmark is a powerful markdown processor Go geared for writing IETF document. It is, however also
suited for writing books and other technical documentation.

Further documentation can be [found at my site](https://miek.nl/tags/mmark/). A complete syntax
document [is being created](https://github.com/miekg/mmark/wiki/Syntax). That syntax doc will also
be mirrored on my website.

With Mmark your can write RFCs using markdown. Mmark (written in Go) provides an advanced markdown
dialect that processes a (single) file to produce internet-drafts in XML format. Internet-drafts
written in mmark can produce xml2rfc v2, xml2rfc v3 and HTML5 output.

It also allows for writing large documents such as technical books, like my [Learning Go
book](https://github.com/miekg/learninggo). Sample text output of this book (when rendered as an
I-D) can [be found
here](https://gist.githubusercontent.com/miekg/0251f3e28652fa603a51/raw/7e0a7028506f7d2948e4ad3091f533711bf5f2a4/learninggo.txt).
It is not perfect due to limitations in xml2rfc version 2. Fully rendered HTML version [can be found
here](http://miek.nl/go).

Mmark is a fork of blackfriday which is a [Markdown][1] processor implemented in [Go][2]. It
supports a number of extensions, inspired by Leanpub, kramdown and Asciidoc, that allows for large
documents to be written.

It adds the following syntax elements to [black friday](https://github.com/russross/blackfriday/blob/master/README.md):

* TOML titleblock.
* Including other files.
* More enumerated lists.
* Table and codeblock captions.
* Quote attribution (quote "captions").
* Table footers, header and block tables
* Subfigures.
* Inline Attribute Lists.
* Indices.
* Citations.
* Abstract/Preface/Notes sections.
* Parts.
* Asides.
* Main-, middle- and backmatter divisions.
* Math support.
* Example lists.
* HTML Comment parsing.
* BCP14 (RFC2119) keyword detection.
* Include raw XML references.
* Abbreviations.
* Super- and subscript.
* Callouts in code blocks.

## Usage

In the mmark subdirectory you can build the mmark tool:

    % cd mmark
    % go build
    % ./mmark -version
    1.3.1

To output v2 xml just give it a markdown file and:

    % ./mmark/mmark -xml2 -page mmark2rfc.md

Making a draft in text form:

    % ./mmark/mmark -xml2 -page mmark2rfc.md > x.xml \
    && xml2rfc --text x.xml \
    && rm x.xml && mv x.txt mmark2rfc.txt

Outputting v3 xml is done with the `-xml` switch. There is not yet a processor for this XML, but you
should be able to validate the resulting XML against the schema from the xml2rfc v3 draft. I'm
trying to stay current with the latest draft for the V3 spec:
<https://tools.ietf.org/html/draft-iab-xml2rfc-03>

## Syntax

See the [syntax](https://github.com/miekg/mmark/wiki/Syntax) document on all syntax elements that
are supported by Mmark.
