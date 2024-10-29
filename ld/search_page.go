package ld

import (
	"fmt"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type SearchPage struct {
	Results []string
}

const (
	searchUrlFormat = "https://www.ldoceonline.com/search/english/direct/?q=%s"
)

func Search(text string) (Page, error) {
	return ParseUrl(fmt.Sprintf(searchUrlFormat, text))
}

func parseSearchPage(doc *html.Node) SearchPage {
	page := SearchPage{}

	suggestions := htmlquery.Find(doc, `//ul[@class="didyoumean"]/li`)

	page.Results = make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		page.Results = append(page.Results, innerTextTrim(s))
	}

	return page
}
