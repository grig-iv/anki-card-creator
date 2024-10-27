package ld

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/proxy"
)

type ldPage interface{}

type wordPage struct {
	title   string
	entries []dictEntry
}

type dictEntry struct {
	hyphenation   string
	pronunciation string
	partOfSpeach  string
	grammar       string

	frequency frequency

	britAudioUrl     string
	americanAudioUrl string

	senses []ldSense
}

type frequency uint8

const (
	none frequency = iota
	low
	mid
	high
)

type ldSense interface{}

type sense struct {
	signpost   string
	definition string
	synonyms   string
	examples   []example
}

type example struct {
	text       string
	audioUrl   string
	colloquial string
}

type crossRefSense struct {
	ref  string
	text string
}

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

	return parse(page.Body)
}

func parse(body io.Reader) (ldPage, error) {
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return nil, err
	}

	page := wordPage{}

	titleNode := htmlquery.FindOne(doc, `//h1[@class="pagetitle"]`)
	if titleNode != nil {
		page.title = innerTextTrim(titleNode)
	}

	page.entries = make([]dictEntry, 0)
	for _, node := range htmlquery.Find(doc, `//span[@class="dictentry"]`) {
		intro := htmlquery.FindOne(node, `//span[@class="dictionary_intro span"]`)
		if intro != nil && innerTextTrim(intro) == "From Longman Business Dictionary" {
			break
		}

		page.entries = append(page.entries, parseEntry(node))
	}

	return ldPage(page), nil
}

func parseEntry(node *html.Node) dictEntry {
	entry := dictEntry{}

	hyphenation := htmlquery.FindOne(node, `//span[@class="HYPHENATION"]`)
	if hyphenation != nil {
		entry.hyphenation = innerTextTrim(hyphenation)
	}

	pronunciation := htmlquery.FindOne(node, `//span[@class="PronCodes"]`)
	if pronunciation != nil {
		entry.pronunciation = innerTextTrim(pronunciation)
	}

	partOfSpeach := htmlquery.FindOne(node, `//span[@class="POS"]`)
	if partOfSpeach != nil {
		entry.partOfSpeach = innerTextTrim(partOfSpeach)
	}

	frequency := htmlquery.FindOne(node, `//span[@class="tooltip LEVEL"]`)
	if frequency != nil {
		switch innerTextTrim(frequency) {
		case "●○○":
			entry.frequency = low
		case "●●○":
			entry.frequency = mid
		case "●●●":
			entry.frequency = high
		}
	}

	grammar := htmlquery.FindOne(node, `//span[@class="GRAM"]`)
	if grammar != nil {
		entry.grammar = strings.Trim(htmlquery.InnerText(grammar), "[ ]")
	}

	for _, node := range htmlquery.Find(node, `//span[@data-src-mp3]`) {
		attr := htmlquery.FindOne(node, `//@data-src-mp3`)
		src := innerTextTrim(attr)
		src = src[:strings.Index(src, "?")]

		if hasClass(node, "brefile") {
			entry.britAudioUrl = src
			continue
		}

		if hasClass(node, "amefile") {
			entry.americanAudioUrl = src
			continue
		}
	}

	entry.senses = make([]ldSense, 0)
	for _, node := range htmlquery.Find(node, `//span[@class="Sense"]`) {
		entry.senses = append(entry.senses, parseSense(node))
	}

	return entry
}

func parseSense(node *html.Node) sense {
	sense := sense{}

	signpost := htmlquery.FindOne(node, `//span[@class="SIGNPOST"]`)
	if signpost != nil {
		sense.signpost = innerTextTrim(signpost)
	}

	definition := htmlquery.FindOne(node, `//span[@class="DEF"]`)
	if definition != nil {
		sense.definition = innerTextTrim(definition)
	}

	synonyms := htmlquery.FindOne(node, `//span[@class="SYN"]`)
	if synonyms != nil {
		sense.synonyms = innerTextTrim(synonyms)[4:]
	}

	sense.examples = make([]example, 0)
	for _, node := range htmlquery.Find(node, `//span[@class="EXAMPLE"]`) {
		sense.examples = append(sense.examples, parseExample(node))
	}

	return sense
}

func parseExample(node *html.Node) example {
	example := example{}

	example.text = innerTextTrim(node)

	scrAttr := htmlquery.FindOne(node, `//span/@data-src-mp3`)
	if scrAttr != nil {
		src := innerTextTrim(scrAttr)
		example.audioUrl = src[:strings.Index(src, "?")]
	}

	return example
}