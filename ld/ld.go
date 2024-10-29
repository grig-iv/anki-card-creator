package ld

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/proxy"
)

type ldPage interface{}

var (
	client *http.Client
)

func init() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	transport := &http.Transport{DialContext: dialContext, DisableKeepAlives: true}

	client = &http.Client{Transport: transport}
}

func loadPage(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	return client.Do(req)
}

func ParseUrl(url string) (ldPage, error) {
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
		return ldPage(parseWordPage(doc)), nil
	}

	if htmlquery.FindOne(doc, `//h1[@class="search_title"]`) != nil {
		return ldPage(parseSearchPage(doc)), nil
	}

	return nil, errors.New("Unknown page: " + url)
}
