package ronn

import (
  "bytes"
  "github.com/russross/blackfriday"
  "fmt"
)

func HtmlRenderer() blackfriday.Renderer {
  htmlFlags := 0
  htmlFlags |= blackfriday.HTML_USE_XHTML
  return &customHtmlRenderer{
    Renderer: blackfriday.HtmlRenderer(htmlFlags, "", ""),
  }
}

type customHtmlRenderer struct {
  blackfriday.Renderer
  seenHeader bool
}

/*
<h2 id="NAME">NAME</h2>
<p class="man-name">
  <code>simple</code> - <span class="man-whatis">a simple ron example</span>
</p>
*/

func (r *customHtmlRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
  if level != 1 || r.seenHeader {
    r.Renderer.Header(out, text, level, id)
    return
  }

  current := out.Bytes()

  r.Renderer.Header(&bytes.Buffer{}, text, level, id)

  extra := out.Bytes()[len(current):]
  name, section, desc := SniffHeader(string(extra))

  out.Reset()
  out.Write(current)

  if len(name) + len(section) > 0 {
    out.WriteString(`<h2 id="NAME">NAME</h2>`)
    out.WriteString("\n")
    out.WriteString(`<p class="man-name">`)
    out.WriteString("\n")
    out.WriteString(fmt.Sprintf(`  <code>%s</code> - <span class="man-whatis">%s</span>`, name, desc))
    out.WriteString("\n")
    out.WriteString("</p>\n")
  } else {
    r.Renderer.Header(out, text, level, id)
  }

  r.seenHeader = true
}
