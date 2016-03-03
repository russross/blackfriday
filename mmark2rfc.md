
> BEING MOVED TO THE SYNTAX document.

# Tables

Tables can be created by drawing them in the input using a simple syntax:

```
Name    | Age
--------|-----:
Bob     | 27
Alice   | 23
```

Tables can also have a footer: use equal signs instead of dashes for the separator,
to start a table footer. If there are multiple footer lines, the first one is used as a
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
(body) cell. If we want to start a new cell use the block table header
syntax. In the example below we include a list in one of the cells.

```
|-----------------+------------+-----------------|
| Default aligned |Left aligned| Center aligned  |
|-----------------|:-----------|:---------------:|
| Bob             |Second cell | Third cell      |
| Alice           |foo         | **strong**      |
| Charlie         |quux        | baz             |
|-----------------+------------+-----------------|
| Bob             | foot       | 1. Item2        |
| Alice           | quuz       | 2. Item2        |
|=================+============+=================|
| Footer row      | more footer| and more        |
|-----------------+------------+-----------------|
```

Note that the header and footer can't contain block level elements.
The table syntax used that one of
[Markdown Extra](https://michelf.ca/projects/php-markdown/extra/#table).

# Lists

## Ordered Lists

The are several ways to start an ordered lists. You can use numbers, roman numbers, letters and uppercase
letters. When using roman numbers and letter you **MUST** use two spaces after the dot or the brace (the
underscore signals a space here):

    a)__
    A)__

Note that mmark (just as @pandoc) pays attention to the starting number of a list (when using decimal numbers), thus
a list started with:

    4) Item4
    5) Item5

Will use for `4` as the starting number.

## Unordered Lists

Unordered lists can be started with `*`, `+` or `-` and follow the normal markdown syntax rules.

# Figures and Images

When an figure has a caption it will be wrapped in `<figure`> tags. A figure can
wrap source code (v3) or artwork (v2/v3).

An image is wrapped in a figure when the optional title syntax is used. But images
are only useful when outputting v3. For v2 the actual image can not be shown, see
(#images-in-v2) for this.

Multiple artworks/sources can be put in one figure. This done by prefixing the
section containing the figures with a figure quote: `F> `.

## Details

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

    v3 allows the usage of a `src` attribute to link to external files with images.
    We use the image syntax for that.

*   An image `![Alt text](/path/to/img.jpg "Optional title")`, will be converted
    to an artwork with a `src` attribute in v3. Again the type needs to be specified
    as an IAL.

    If the "Optional title" is specified the generated artwork will be wrapped in a
    figure with name set to "Optional title"

    Creating an artwork with an anchor and type will become:

        {#fig-id type="ascii-art"}
        ![](/path/to/art.txt "Optional title")

    For v2 this presents difficulties as there is no way to display any of this, see
    (#images-in-v2) for a treatment on how to deal with that.

*   To group artworks and code blocks into figures, we need an extra syntax element.
    [Scholarly markdown] has a neat syntax
    for this. It uses a special section syntax and all images in that section become
    subfigures of a larger figure. Disadvantage of this syntax is that it can not be
    used in lists. Hence we use a quote like solution, just like asides and notes,
    but for figures: we prefix the entire paragraph with `F>` .

    Basic usage:

        F>  {type="ascii-art"}
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

    In v2 this is not supported so the above will result in one figure. Yes one, because
    the fenced code block does not have a caption, so it will not be wrapped in a figure.

    To summerize in v2 the inner captions *are* used and the outer one is discarded, for v3 it
    is the other way around.

    The figure from above will be rendered as:

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
    Figure: Caption for both figures in v3 (in v2 it's ignored).


## Images in v2

Images (real images, not ascii-art) are non-existent in v2, but are allowed in v3. To allow
writers to use images *and* output v2 and v3 formats, the following hack is used in v2 output.
Any image will be converted to a figure with an title attribute set to the "Optional title".
And the url in the image will be type set as a link in the postamble.
So `![](misc/image.xml "Optional title")` will be converted to:

    <figure title="Optional title">
     <artwork>
     </artwork>
      <postamble>
       <eref target="misc/image.xml"/>
      </postamble>
    </figure>

If a image does not have a title, the `figure` is dropped and only the link remains. The default
is to center the entire element. Note that is you don't give the image an anchor, `xml2rfc` won't
typeset it with a `Figure X`, so for an optional "image" rendering, you should use the folowing:

    {#fig-id}
    ![](misc/image.xml "Optional title")

Which when rendered becomes:

{#fig-id}
![](misc/image.xml "Optional title")

Note that ideas to improve/change on this are welcome.
