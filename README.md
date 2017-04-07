# Mmark

Mmark is a powerful markdown processor Go geared towards writing IETF documents. It is, however,
also suited for writing books and other technical documentation.

Further documentation can be [found at my site](https://miek.nl/tags/mmark/). A complete syntax
document [can be found here](https://github.com/miekg/mmark/wiki/Syntax).

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

Example documents written in Mmark can be found in the `rfc/` directory.

Mmark adds the following syntax elements to [black friday](https://github.com/russross/blackfriday/blob/master/README.md):

* TOML titleblock.
* Including other files.
* More enumerated lists and task-lists.
* Table and codeblock captions.
* Quote attribution (quote "captions").
* Table footers, header and block tables.
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
    1.3.4

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

## Running from a Docker Container

To use mmark and xml2rfc without installing and configuring the separate software packages and
dependencies, you can also use mmark via a Docker container.  You can use
https://github.com/paulej/rfctools to build a Docker image or you can use the one already
created and available in Docker Hub (https://hub.docker.com/r/paulej/rfctools/).

To use Docker, you invoke commands like this:

    % docker run --rm paulej/rfctools mmark -version
    1.3.4

To output a v2 XML file as demonstrated in the previous section, use this command:

    % docker run --rm -v $(pwd):/rfc paulej/rfctools mmark -xml2 -page mmark2rfc.md

Making a draft in text form: 

    % docker run --rm -v $(pwd):/rfc paulej/rfctools mmark -xml2 -page mmark2rfc.md >x.xml \
    && docker run --rm -v $(pwd):/rfc -v $HOME/.cache/xml2rfc:/var/cache/xml2rfc \
    --user=$(id -u):$(id -g) paulej/rfctools xml2rfc --text x.xml \
    && rm x.xml && mv x.xml mmark2rfc.txt

The docker container expects source files to be stored in /rfc, so the above command maps
the current directory to /rfc.  Likewise, xml2rfc will attempt to write cache files to
/var/cache/xml2rfc, so the above command maps the user's cache directory in the container.

Note also that the xml2rfc program will write an output file that will be owned by "root".
To prevent that (and the cache files) from being owned by root, we instruct docker to run
using the user's default user ID and group ID via the --user switch.
    
There is a script available called "md2rfc" simplifies the above to this:

    % docker run --rm -v $(pwd):/rfc -v $HOME/.cache/xml2rfc:/var/cache/xml2rfc \
    --user=$(id -u):$(id -g) paulej/rfctools mmark2rfc.md

Still appreciating that is a lot to type, there is an example directory in the
https://github.com/paulej/rfctools repository that contains a Makefile.  Place your .md file
in a directory along with that Makefile and just type "make" to produce the .txt file.

## Syntax

See the [syntax](https://github.com/miekg/mmark/wiki/Syntax) document on all syntax elements that
are supported by Mmark.

[1]: https://daringfireball.net/projects/markdown/ "Markdown"
[2]: https://golang.org/ "Go Language"
