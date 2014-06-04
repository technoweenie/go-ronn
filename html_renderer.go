package ronn

import (
  "bytes"
  "github.com/russross/blackfriday"
  "fmt"
)

func HtmlRenderer(d *Document) *customHtmlRenderer {
  htmlFlags := 0
  htmlFlags |= blackfriday.HTML_USE_XHTML
  return &customHtmlRenderer{
    Doc: d,
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
  if level != 1 || r.SeenHeader {
    r.Renderer.Header(out, text, level, id)
    return
  }

  current := out.Bytes()

  r.Renderer.Header(&bytes.Buffer{}, text, level, id)

  extra := out.Bytes()[len(current):]
  name, section, tagline := SniffHeader(string(extra))
  r.Doc.Name = name
  r.Doc.Section = section
  r.Doc.Tagline = tagline

  out.Reset()
  out.Write(current)

  if len(name) + len(section) > 0 {
    r.ManHeader(out, name, section, tagline)
  } else {
    r.Renderer.Header(out, text, level, id)
  }

  r.SeenHeader = true
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
