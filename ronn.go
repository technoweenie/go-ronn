package ronn

import (
  "github.com/russross/blackfriday"
  "regexp"
)

func HTML(input []byte) string {
  output := string(Render(HtmlRenderer(), input))
  return `<div class="mp">` + "\n" + output + "\n</div>\n"
}

func Render(renderer blackfriday.Renderer, input []byte) []byte {
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_HEADER_IDS

	return blackfriday.Markdown(input, renderer, extensions)
}

func SniffHeader(input string) (string, string, string) {
  if match := nameWithSectionRE.FindStringSubmatch(input); match != nil {
    return match[1], match[2], match[3]
  }

  if match := nameRE.FindStringSubmatch(input); match != nil {
    return match[1], "", match[2]
  }

  return "", "", input
}

var (
  // name(section) -- description
  nameWithSectionRE = regexp.MustCompile(`([\w_.\[\]~+=@:-]+)\s*\((\d\w*)\)\s*-+\s*(.*)`)

  // name -- description
  nameRE = regexp.MustCompile(`([\w_.\[\]~+=@:-]+)\s+-+\s+(.*)`)
)
