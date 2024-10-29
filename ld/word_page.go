package ld

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

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
	grammar    string
	geo        string
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

func parseWordPage(doc *html.Node) wordPage {
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

	return page
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
		crossRefNode := htmlquery.FindOne(node, `//a[@class="crossRef"]`)
		if crossRefNode != nil {
			entry.senses = append(entry.senses, parseCrossRef(crossRefNode))
		} else {
			entry.senses = append(entry.senses, parseSense(node))
		}
	}

	return entry
}

func parseCrossRef(node *html.Node) crossRefSense {
	crossRef := crossRefSense{}

	crossRef.text = innerTextTrim(node)
	crossRef.ref = innerTextTrim(htmlquery.FindOne(node, "//@href"))

	return crossRef
}

func parseSense(node *html.Node) sense {
	sense := sense{}

	signpost := htmlquery.FindOne(node, `//span[@class="SIGNPOST"]`)
	if signpost != nil {
		sense.signpost = innerTextTrim(signpost)
	}

	grammar := htmlquery.FindOne(node, `//span[@class="GRAM"]`)
	if grammar != nil {
		sense.grammar = strings.Trim(htmlquery.InnerText(grammar), "[ ]")
	}

	geo := htmlquery.FindOne(node, `//span[@class="GEO"]`)
	if geo != nil {
		sense.geo = innerTextTrim(geo)
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

	if hasClass(node.Parent, "ColloExa") {
		example.colloquial = innerTextTrim(node.Parent.FirstChild)
	}

	return example
}
