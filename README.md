# Mmark

Mmark is a powerful markdown processor Go geared for writing IETF document. It is, however also
suited for writing books and other technical documentation.

Further documentation can be [found at my site](https://miek.nl/tags/mmark/). A complete syntax
document [is being created](https://github.com/miekg/mmark/wiki/Syntax). That syntax doc will also
be mirrored on my website.

With Mmark your can write RFCs using markdown. Mmark (written in Go) provides an advanced markdown
dialect that processes a single file to produce internet-drafts in XML format. Internet-drafts
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

* More enumerated lists.
* Table and codeblock captions.
* Table footers.
* Subfigures.
* Quote attribution.
* Including other files.
* TOML] titleblock.
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

There is no pretty printed output if you need that pipe the output through `xmllint --format -`.

## Usage

In the mmark subdirectory you can build the mmark tool:

    % cd mmark
    % go build
    % ./mmark -h
    Mmark Markdown Processor v1.0
    ...

To output v2 xml just give it a markdown file and:

    % ./mmark/mmark -xml2 -page mmark2rfc.md

Making a draft in text form:

    % ./mmark/mmark -xml2 -page mmark2rfc.md > x.xml \
    && xml2rfc --text x.xml \
    && rm x.xml && mv x.txt mmark2rfc.txt

Outputting v3 xml is done with the `-xml` switch. There is not yet
a processor for this XML, but you should be able to validate the
resulting XML against the schema from the xml2rfc v3 draft. I'm trying
to stay current with the latest draft for the V3 spec:
<https://tools.ietf.org/html/draft-hoffman-xml2rfc-23>

> BEIGN MOVED TO THE SYNTAX DOCUMENT

*   **Tables**. Tables can be created by drawing them in the input
    using a simple syntax:

    ```
    Name    | Age
    --------|-----:
    Bob     | 27
    Alice   | 23
    ```

    Tables can also have a footer, use equal signs instead of dashes for
    the separator.
    If there are multiple footer lines, the first one is used as a
    starting point for the table footer.

    ```
    Name    | Age
    --------|-----:
    Bob     | 27
    Alice   | 23
    ======= | ====
    Charlie | 4
    ```

    If a table is started with a *block table header*, which starts
    with a pipe or plus sign and a minimum of three dashes,
    it is a **Block Table**. A block table may include block level elements in each
    (body) cell. If we want to start a new cell reuse the block table header
    syntax. In the example below we include a list in one of the cells.

    ```
    |+-----------------------------------------------|
    | Default aligned |Left aligned| Center aligned  |
    |-----------------|:-----------|:---------------:|
    | First body part |Second cell | Third cell      |
    | Second line     |foo         | **strong**      |
    | Third line      |quux        | baz             |
    |------------------------------------------------|
    | Second body     |            | 1. Item2        |
    | 2 line          |            | 2. Item2        |
    |================================================|
    | Footer row      |            |                 |
    |-----------------+------------+-----------------|
    ```

    Note that the header and footer can't contain block level elements.

    Row spanning is supported as well, by using the
    [multiple pipe syntax](http://bywordapp.com/markdown/guide.html#section-mmd-tables-colspanning).

*   **Subfigure**. Fenced code blocks and indented code block can be
    grouped into a single figure containing both (or more) elements.
    Use the special quote prefix `F>` for this.
*   **Subfigures**, any paragraph prefix with `F>` will wrap all images and
    code in a single figure.



*   **Code Block Includes**, use the syntax `<{{code/hello.c}}[address]`, where
    address is the syntax described in <https://godoc.org/golang.org/x/tools/present/>, the
    OMIT keyword in the code also works.

    So including a code snippet will work like so:

        <{{test.go}}[/START OMIT/,/END OMIT/]

    where `test.go` looks like this:

    ``` go
    tedious_code = boring_function()
    // START OMIT
    interesting_code = fascinating_function()
    // END OMIT
    ```
    To aid in including HTML or XML fragments, where the `OMIT` key words is
    probably embedded in comments, lines which end in `OMIT -->` are also excluded.

    Of course the captioning works here as well:

        <{{test.go}}[/START OMIT/,/END OMIT/]
        Figure: A sample program.

    The address may be omitted: `<{{test.go}}` is legal as well.

    Note that the special `prefix` attribute can be set in an IAL and it
    will be used to prefix each line with the value of `prefix`.

        {prefix="S"}
        <{{test.go}}

    Will cause `test.go` to be included with each line being prefixed with `S`.

*   **Definitition lists**, the markdown extra syntax.

         Apple
         :   Pomaceous fruit of plants of the genus Malus in
             the family Rosaceae.

         Orange
         :   The fruit of an evergreen tree of the genus Citrus.

*   **Enumerated lists**, roman, uppercase roman and normal letters can be used
     to start lists. Note that you'll need two space after the list counter:

         a.  Item2
         b.  Item2

*   **Callouts**, in codeblocks you can use `<number>` to create a callout, later you can
     reference it:

             Code  <1>
             More  <1>
             Not a callout \<3>

         As you can see in <1> but not in \<1>. There is no <3>.

     You can escape a callout with a backslash. The backslash will be removed
     in the output (both in source code and text). The callout identifiers will be remembered until
     the next code block. The above would render as:

                 Code <1>
                 Code <2>
                 Not a callout <3>

             As you can see in (1, 2) but not in <1>. There is no <3>.

     Note that callouts are only detected with the IAL `{callout="yes"}` or any other
     non-empty value is defined before the code block.
     Now, you don't usually want to clobber your source code with callouts as this will
     lead to code that does not compile. To fix this the callout needs to be placed
     in a comment, but then your source show useless empty comments. To fix this mmark
     can optionally detect (and remove!) the comment and the callout, leaving your
     example pristine. This can be enabled by setting `{callout="//"}` for instance.
     The allowed comment patterns are `//`, `#` and `;`.
