// Resolve internal references. If the class chapter has been
// set we take the number of the chapter which is put in the
// <span class="chapter-number"> of the header's text.
// Otherwise we take the section's title and substitute that, just
// like \namref in LaTeX.
function resolveXrefs() {
    var els = document.getElementsByTagName("a");
    for (var i = 0, l = els.length; i < l; i++) {
        var el = els[i];
        var href = el.getAttribute("href")
        if (href.charAt(0) == '#' && el.text == '') {
            id = document.getElementById(href.slice(1));
            if (id == null) {
                continue;
            }
            var nodes = id.childNodes;
            switch (nodes[0].nodeName) {
                case '#text':
                    text = nodes[0].nodeValue;
                    el.text = text;
                    break;
                case 'SPAN':
                    el.text = nodes[0].innerHTML;
                    break;
            }
        }
    }
}

// numberChapters numbers all chapter (h1 with class="chapter")
// and adds a number in the text. If the class "nonumber" is also
// present we skip the numbering.
// If the class "appendix" is *also* present we resort to letters.
function numberChapters() {
    var els = document.getElementsByTagName("h1");
    var chapCount = 1;
    var appxCount = 0;
    var appendix='ABCDEFGHIJKLMNOPQRSTUVWXYZ';
Elements:
    for (var i = 0, l = els.length; i < l; i++) {
        var app = false;
        var el = els[i];
        var classList = el.className.split(/\s+/);
        for (var k = 0; k < classList.length; k++) {
            if (classList[k] == 'nonumber' || classList[k] != 'chapter' ) {
                 continue Elements;
            }
            if (classList[k] == 'appendix') {
                app = true;
            }
        }
        var span = document.createElement('span');
        if (app) {
            var node = document.createTextNode(appendix.charAt[appxCount % 26]);
            appxCount++;
        } else {
            var node = document.createTextNode(chapCount);
            chapCount++;
        }
        span.appendChild(node);
        el.insertBefore(span, el.childNodes[0]);
    }
}
