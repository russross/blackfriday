Inline Attribute Lists

With Inline Attribute Lists (IALs) you can add extra information to block elements, think of (in
HTML) extra classes, attributes or an ID. These will mostly be backend specific, i.e.
a `{color="blue"}` might do something in HTML5 output, but does nothing in the XML v2 output. Worse
yet; the XML might not validate anymore. Mmark will not try to fix this.

However there are some attributes that have special meaning in Mmmark and *will* be treated special
and will work in all renderers.

## BlockCode

* `type`
* `prefix`
* `callout`

In blockcode (i.e. source code output) you can specifiy the langauge of the block:

~~~ go
// Go stuff here
~~~

Mmark will treat the `type` attribute in the same way as this language specifier, i.e. for XML
output this will become: `type="go"` and in HTML this will be`class="language-go`. So the following
is equavalent:

{type="go"}
~~~
// Go stuff here
~~~

Further more there is `prefix=".."` which allows the source code to be prefix with a string (at the
beginning of each line). This is handled by Mmark and the attribute `prefix` is dropped from the
final output.

Then there is `callout=".."` which enables callout (TODO:link) parsing in Mmark. Currently this This
is handled by Mmark and the attribute `callout` is dropped from the final output.

Things like `align` are supported by each renderer and don't require support from Mmark.
