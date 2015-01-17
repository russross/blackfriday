# Figures in mmark

A `<figure>` encapsulates source code (v3) or artwork (v2/v3).
We allow multiple artworks/sources in a figure, but there is
currently no syntax to allow for this. Although Scholary Markdown
uses the following, start a axt section, with the special name
`Figure:` and in that section include the figure you want.

A source code or artwork is wrapped in a figure if a caption
is given.

An image is wrapped in a figure is the optional title is used.

## Proposal

*   A Fenced Code Block will becomes a source code in v3 and artwork in v2.
    We can use the language to signal the type.

        ``` c
        printf("%s\n", "hello");
        ```

*   An Indented Code Block becomes artwork in v3 and artwork in v2. The only way
    to indicate the type is by using an IAL. So one has to use:

        {type="ascii-art"}
            +-----+
            | ART |
            +-----+

    v3 allows the usage of a src= attribute to link to external files with images.
    In this proposal we use the image syntax for that.

*   An image `![Alt text](/path/to/img.jpg "Optional title")`, will be converted
    to an artwork with a `src` attribute in v3. The extension of the included
    image file will be used to generate the `type` attribute (see the table below).

    If the "Optional title" is specified the generated artwork will be wrapped in a
    figure with name set to "Optional title"

    extension | type
    ----------|-----
      txt     | ascii-art
      hex     | hex-dump
      ...     | ....

    Creating an artwork with an anchor will become:

        {#fig:id}
        ![](/path/to/art.txt "Optional title")

    For v2 this presents difficulties as there is no way to display any of this. The v2
    renderer will output a warning and not output anything in this case.

*   Grouping artworks and code blocks into figures. Scholary markdown has a neat syntax
    for this. It uses a special section syntax and all images in that section become
    subfigures of a larger figure. As mmark already uses `.#` to start special section
    we can do the same in a cleaner way. Start a special section use `.# Figure` and
    place multiple figures and artworks in the section.

    What do you do with non sourcecode/artwork elements in such a paragraph, just output
    them (there are no non-compliant markdown documents).

    Basic usage:

        .# Figure {#fig:id}

        {type="ascii-art"}
            +-----+
            | ART |
            +-----+
        Figure: The last caption specified will be used.

        ``` c
        printf("%s\n", "hello");
        ```

        Figure: Caption you will see, for both figures.

    In v2 this is not supported so the above will result in two figures.
