package ld

import (
	"errors"
	"net/http"

	"github.com/antchfx/htmlquery"
)

type Page interface{}

var (
	HttpClient *http.Client = http.DefaultClient
)

func loadPage(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	return HttpClient.Do(req)
}

func ParseUrl(url string) (Page, error) {
	page, err := loadPage(url)
	if err != nil {
		return nil, err
	}

	defer page.Body.Close()

	doc, err := htmlquery.Parse(page.Body)
	if err != nil {
		return nil, err
	}

	if htmlquery.FindOne(doc, `//h1[@class="pagetitle"]`) != nil {
		return Page(parseWordPage(doc)), nil
	}

	if htmlquery.FindOne(doc, `//h1[@class="search_title"]`) != nil {
		return Page(parseSearchPage(doc)), nil
	}

	return nil, errors.New("Unknown page: " + url)
}
