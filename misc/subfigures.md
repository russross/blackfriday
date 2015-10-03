Title: Subfigures in mmark
Date: 2015-10-03 13:24
Summary: Typesetting subfigures in Mmark markdown
Tags: Markdown,XML2RFC,figures,subfigures,styling

For XML2RFC version 3 there is a new [extended figure element](https://tools.ietf.org/html/draft-hoffman-xml2rfc-23#section-2.25),
namely: it can multiple `<artwork>` or `<sourcecode>` (new in v3) elements.
This in turn meant that mmark needed some way to signal "these artworks belong together".
Now [Scholarly markdown](http://scholarlymarkdown.com/Scholarly-Markdown-Guide.html)
has a syntax for this. It uses a special section syntax and all images in that section become
subfigures of a larger figure. Disadvantage of this syntax is that it can not be
used in lists. Hence mmark uses a quote like solution, just like asides (`A>`) and notes (`N>`)
but for figures: we prefix the entire paragraph with `F>` .

Basic usage:

    F> {type="ascii-art"}
    F>      +-----+
    F>      | ART |
    F>      +-----+
    F>  Figure: This caption is ignored in v3, but used in v2.
    F>
    F>  ``` c
    F>  printf("%s\n", "hello");
    F>  ```
    F>
    Figure: Caption for both figures in v3 (in v2 this is ignored).

In v2 this is not supported so the above will result in one figure. Yes one,
because the fenced (with the `printf`) code block does not have a caption, so it
will not be wrapped in a figure.

To summerize in v2 the inner captions *are* used and the outer one is discarded, for v3 it
is the other way around.

The figure from above will be rendered in v2 as:


                                  +-----+
                                  | ART |
                                  +-----+

              This caption is ignored in v3, but used in v2.

                         printf("%s\n", "hello");

There is no v3 renderer (yet). Adding a second caption results in:

    F> {type="ascii-art"}
    F>      +-----+
    F>      | ART |
    F>      +-----+
    F>  Figure: This caption is ignored in v3, but used in v2.
    F>
    F>  ``` c
    F>  printf("%s\n", "hello");
    F>  Figure: Another caption that is ignored in v3, but used in v2.
    F>  ```
    F>
    Figure: Caption for both figures in v3 (in v2 this is ignored).

Becomes:

                                  +-----+
                                  | ART |
                                  +-----+

              This caption is ignored in v3, but used in v2.

                         printf("%s\n", "hello");

          Another caption that is ignored in v3, but used in v2.
