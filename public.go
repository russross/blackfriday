// Public interface

package mmark

// MarkdownBasic is a convenience function for simple rendering.
// It processes markdown input with no extensions enabled.
func MarkdownBasic(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := HTML_USE_XHTML
	renderer := HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0

	return Markdown(input, renderer, extensions)
}

// Call Markdown with most useful extensions enabled
// MarkdownCommon is a convenience function for simple rendering.
// It processes markdown input with common extensions enabled, including:
//
// * Smartypants processing with smart fractions and LaTeX dashes
//
// * Intra-word emphasis suppression
//
// * Tables
//
// * Fenced code blocks
//
// * Autolinking
//
// * Strikethrough support
//
// * Strict header parsing
//
// * Custom Header IDs
func MarkdownCommon(input []byte) []byte {
	renderer := HtmlRenderer(commonHtmlFlags, "", "")
	return Markdown(input, renderer, commonExtensions)
}

// XML2RFC v3 output.
func MarkdownXml2rfc(input []byte) []byte {
	renderer := HtmlRenderer(XML_STANDALONE, "", "")
	return Markdown(input, renderer, commonXmlExtensions)
}

// XML2RFC v2 output.
func MarkdownXml22rfc(input []byte) []byte {
	renderer := HtmlRenderer(XML2_STANDALONE, "", "")
	return Markdown(input, renderer, commonXmlExtensions)
}
