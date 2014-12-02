% title = "Example"
% abbrev = "ex"
% category = "info"
% docname = "draft-example-rfc-markdown-00"
% ipr = "trust200902"
%
% date = "2014-02-02T00:00:00-00:00"
% area = "General"
% workgroup = "Making RFC easier"
% keyword = ["Internet-Draft", "Markdown"]
%
% [[author]]
% initial = "R."
% surname = "Gieben"
% fullname = "R. (Miek) Gieben"

AB> abstract
AB> more

{mainmatter}

# Introduction

# Markdown for XML2RFC v3

Or... how to write Markdown to generate XML that is valid XML2RFC v3. These
are some assorted notes and ideas.

Goal:

1. Self contained file with all information to generate a complete I-D;
1. CommonMarkdown, with some extensions.

## Document meta data

Use a titleblock with TOML data. The titleblock needs to be prefixed with %

## Citations

Citation use the synax from Pandoc: [@smith04], but ... is more strict, the `@`-sign must be the first
character after the open brace. [we allow multiple citations: [@smith04; @doe99] TODO].
More general the syntax is `[@REF(MOD:PATH) TEXT]`, where

* REF: is the reference anchor;
* MOD: is either n for normative or i for informative references. The default is i when left
    out;
* PATH: the name of the reference file;
* TEXT: text used to reference.

So a full reference would look like this:

    [@RFC1035(n:bib/reference.RFC.1035.xml) Section 1]

... knows about RFC and I-D references and will use default filenames for the citations. When referencing
an RFC, like `@RFC2535` the filename will be `reference.RFC.2535.xml` and when referencing I-Ds, like
`I-D.gieben-bla` the filename will become `reference.I-D.gieben-blap.xml`

The reference section(s) are outputted automatically, when there are citations in the document.

## Abstract, notes and asides

Abstract is done with prefixing a paragraph with `AB>`, notes: `N>` and aside with `A>`.

## Extra attributes

Adopted and slightly adapted the kramdown IAL syntax. We've dropped the colon inside the
braces. So it just is `{.class key=value}` to add attributes to block elements.

## Document divisions

The front-, middle- and back-sections are signalled with:

`{frontmatter}`, `{middlematter}` and `{backmatter}`, where a document automatically starts
with a frontmatter.

## Index

Using tripe parentheses, like so: `(((Cats, Tiger)))`.

## Including Markdown files

Use the syntax {{file}} will get `file` included.

All of the above use the extension syntax, i.e. starting with a { .

# BCP14 keywords

Uppercase in asterisks `*MUST*`.

# Figures

# Artwork

Figure is used to group sourcecode or artwork.

            sourcecode
    figure /
           \
            artwork

figure has caption and title (`<name>`) and such.

artwork is svg type:

        ``` svg

        ```

# Comments

C>
C>
C>

becomes cref's (if needed)
