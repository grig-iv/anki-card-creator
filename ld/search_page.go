package ld

import (
	"fmt"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type searchPage struct {
	results []string
}

const (
	searchUrlFormat = "https://www.ldoceonline.com/search/english/direct/?q=%s"
)

func Search(text string) (ldPage, error) {
	return ParseUrl(fmt.Sprintf(searchUrlFormat, text))
}

func parseSearchPage(doc *html.Node) searchPage {
	page := searchPage{}

	suggestions := htmlquery.Find(doc, `//ul[@class="didyoumean"]/li`)

	page.results = make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		page.results = append(page.results, innerTextTrim(s))
	}

	return page
}
