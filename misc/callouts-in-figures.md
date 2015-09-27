Title: Callouts in Figures in Mmark
Date: 2015-09-27 14:02
Summary: Callout in Figures in Mmark markdown
Tags: Markdown,XML2RFC,callout,callouts,styling,mmark

The use of callouts in Mmark came from
[Asciidoc](http://www.methods.co.nz/asciidoc/chunked/ch20.html), they are "...
a mechanism for annotating verbatim text". In mmark they work like this.

In a codeblocks you can use `<1>`, or `<2>`, etc. to create a callout. After the
codeblock/figure you can reference it. An example:

        Code  <1>
        More  <1>
        Not a callout \<3>

    As you can see in <1> but not in \<1>. There is no <3>.

This will be rendered as:

                             Code  <1>
                             More  <2>
                             Not a callout <3>

    As you can see in (1, 2) but not in <1>.  There is no <3>.

A callout can be escaped with a backslash. The backslash will be removed in the
output (both in sourcecode and text). The callout identifiers will be remembered
until the next code block.

Note that callouts are only detected with the IAL `{callout="yes"}` before the
codeblock.

Now, you don't usually want to globber your sourcecode with callouts as this will
lead to code that does not compile. To fix this the callout needs to be placed
in a comment, but then your source show useless empty comments. To fix this
mmark can optionally detect (and remove!) the comment and the callout, leaving
your example pristine. This can be enabled by setting `{callout="//"}` for
instance. The allowed comment patterns are `//`, `#` and `;`.
