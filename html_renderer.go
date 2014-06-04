package ronn

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
)

func HtmlRenderer(d *Document) *customHtmlRenderer {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_XHTML
	return &customHtmlRenderer{
		Doc:      d,
		Renderer: blackfriday.HtmlRenderer(htmlFlags, "", ""),
	}
}

type customHtmlRenderer struct {
	Doc *Document
	blackfriday.Renderer
	SeenHeader bool
}

/*
<h2 id="NAME">NAME</h2>
<p class="man-name">
  <code>simple</code> - <span class="man-whatis">a simple ron example</span>
</p>
*/

func (r *customHtmlRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	if level != 1 {
		r.SubHeader(out, text, level, id)
		return
	}

	if r.SeenHeader {
		r.Renderer.Header(out, text, level, id)
		return
	}

	extra := r.getText(out, text, func() {
		r.Renderer.Header(&bytes.Buffer{}, text, level, id)
	})

	name, section, tagline := SniffHeader(extra)
	r.Doc.Name = name
	r.Doc.Section = section
	r.Doc.Tagline = tagline

	if len(name)+len(section) > 0 {
		r.ManHeader(out, name, section, tagline)
	} else {
		r.Renderer.Header(out, text, level, id)
	}

	r.SeenHeader = true
}

func (r *customHtmlRenderer) SubHeader(out *bytes.Buffer, text func() bool, level int, id string) {
	extra := r.getText(out, text, func() {
		r.Renderer.Header(&bytes.Buffer{}, text, level, id)
	})

	out.WriteString("\n")
	out.WriteString(fmt.Sprintf(`<h%d id="%s">%s</h%d>`, level, extra, extra, level))
	out.WriteString("\n")
}

func (r *customHtmlRenderer) ManHeader(out *bytes.Buffer, name, section, tagline string) {
	out.WriteString(`<h2 id="NAME">NAME</h2>`)
	out.WriteString("\n")
	out.WriteString(`<p class="man-name">`)
	out.WriteString(fmt.Sprintf("\n  <code>%s</code>", name))
	if len(tagline) > 0 {
		out.WriteString(fmt.Sprintf(` - <span class="man-whatis">%s</span>`, tagline))
	}
	out.WriteString("\n")
	out.WriteString("</p>\n")
}

func (r *customHtmlRenderer) getText(out *bytes.Buffer, text func() bool, cb func()) string {
	current := out.Bytes()
	cb()
	extra := out.Bytes()[len(current):]
	out.Reset()
	out.Write(current)
	return string(extra)
}
