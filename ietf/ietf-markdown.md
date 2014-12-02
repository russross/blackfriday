% title = "A Markdown Syntax XML2RFV V3"
% abbrev = "Markdown for XML2RFV v3"
% category = "info"
% docname = "draft-markdown-xml2rfc-v3-00"
% ipr = "trust200902"
%
% date = "2014-02-02T00:00:00-00:00"
% area = "General"
% workgroup = "RFC Beautification Working Group"
% keyword = ["Internet-Draft", "Markdown", "XML2RFV v3]
%
% [[author]]
% initial = "R."
% surname = "Gieben"
% fullname = "R. (Miek) Gieben"

AB> This document shows a markdown syntax for use in creating I-Ds. The document
AB> is converted to the canonical RFC format: XML2RFC v3 [@I-D.hoffman-xml2rfc].

{mainmatter}

# Introduction

We use the "normal" markdown and try to stay as close as possible to the CommonMark work that
is been carried out at the moment. The goal of this specifiction is to have one document (modulo
possible reference file) that can be translated to XML.

## Document Meta Data

The document meta data is speficied in a titleblock (prefix with `%`) which is typeset in [TOML](https://github.com/toml-lang/toml).

## Citations

Citation use the synax from Pandoc: [@smith04], but is more strict, the `@`-sign must be the first
character after the open brace. 
More general the syntax is `[@REF(MOD:PATH) TEXT]`, where

REF: 
:   is the reference anchor;
MOD:
:   is either n for normative or i for informative references. The default is i when left out;
PATH: 
:   the name of the reference file;
TEXT:
:   text used to reference.

So a fully specified reference would look like this:

    [@RFC1035(n:bib/reference.RFC.1035.xml) Section 1]

... knows about RFC and I-D references and will use default filenames for the citations. When referencing
an RFC, the default filename will be `reference.RFC.<NUMBER>.xml` and when referencing I-Ds,
the filename will become `reference.I-D.draft-<name>.xml`. (This does not have a sequence number).

The reference section(s) are outputted automatically when there are citations in the document.

## Abstract, Notes and Asides

An aBstract is created by prefixing a paragraph with `AB>`, Notes with `N>` and an aside with `A>`.

## Extra attributes

Adopted and slightly adapted the kramdown IAL syntax. We've dropped the colon inside the
braces, which seems to be the syntax CommonMark is going to use. So it just is `{.class key=value}`
to add attributes to block elements.

## Document divisions

The front-, middle- and back-sections are signalled with:

`{frontmatter}`, `{middlematter}` and `{backmatter}`. When there is a TOML title block the frontmatter
is automatically openened.

## Index

An index can be created by tripe parentheses, like so: `(((Cats, Tiger)))`.

## Including Files

Use the syntax `{{file}}` which will get `file` included.

# BCP14 keywords

Uppercase in asterisks `*MUST*`.

# Figures

# Artwork

# Sourecode
