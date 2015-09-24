Mmark is a markdown dialect for writing RFCs or books. The main reason it exists
is because the IETF is developing a new XML format (version 3) for writing RFCs.
This new format is *way* more powerful, so sadly
[Pandoc2rfc](https://tools.ietf.org/html/rfc7328.html) does not cut it anymore.

Right now [mmark](https://github.com/miekg/mmark) fully supports writing
RFCs in XML2RFC version 2, but anything written can also be converted to
the (work-in-progress) XML2RFC version 3 format. In a series of blog posts
I will detail some of the features, such as list styling while remaining upwards
compatible.

In XML2RFC you can "style" a list by using:

    <list style="format (%I)">

Here we say use Roman numerals inclosed in braces: `(I) (II)`, etc. In XML2RFC
v3 this becomes:

    <ul type="(%I)">

Mmark knows about these differences and allows you so say the following:

    {style="format (%I)" type="(%I)"}
    * Item 1
    * Item 2
    * Item 3

Where for the v2 output `type` is discarded and vice versa. In the v2 text output
this is rendered as:

       (I)     Item 1

       (II)    Item 2

       (III)   Item 3
