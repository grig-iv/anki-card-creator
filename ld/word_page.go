package ld

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type WordPage struct {
	Title   string
	Entries []DictEntry
}

type DictEntry struct {
	Hyphenation   string
	Pronunciation string
	PartOfSpeach  string
	Grammar       string

	Frequency Frequency

	BritAudioUrl     string
	AmericanAudioUrl string

	Senses []LdSense
}

type Frequency uint8

const (
	None Frequency = iota
	Low
	Mid
	High
)

type LdSense interface{}

type Sense struct {
	Signpost   string
	Grammar    string
	Geo        string
	Definition string
	Synonyms   string
	Examples   []Example
}

type Example struct {
	Text       string
	AudioUrl   string
	Colloquial string
}

type CrossRefSense struct {
	Ref  string
	Text string
}

func parseWordPage(doc *html.Node) WordPage {
	page := WordPage{}

	titleNode := htmlquery.FindOne(doc, `//h1[@class="pagetitle"]`)
	if titleNode != nil {
		page.Title = innerTextTrim(titleNode)
	}

	page.Entries = make([]DictEntry, 0)
	for _, node := range htmlquery.Find(doc, `//span[@class="dictentry"]`) {
		intro := htmlquery.FindOne(node, `//span[@class="dictionary_intro span"]`)
		if intro != nil && innerTextTrim(intro) == "From Longman Business Dictionary" {
			break
		}

		page.Entries = append(page.Entries, parseEntry(node))
	}

	return page
}

func parseEntry(node *html.Node) DictEntry {
	entry := DictEntry{}

	hyphenation := htmlquery.FindOne(node, `//span[@class="HYPHENATION"]`)
	if hyphenation != nil {
		entry.Hyphenation = innerTextTrim(hyphenation)
	}

	pronunciation := htmlquery.FindOne(node, `//span[@class="PronCodes"]`)
	if pronunciation != nil {
		entry.Pronunciation = innerTextTrim(pronunciation)
	}

	partOfSpeach := htmlquery.FindOne(node, `//span[@class="POS"]`)
	if partOfSpeach != nil {
		entry.PartOfSpeach = innerTextTrim(partOfSpeach)
	}

	frequency := htmlquery.FindOne(node, `//span[@class="tooltip LEVEL"]`)
	if frequency != nil {
		switch innerTextTrim(frequency) {
		case "●○○":
			entry.Frequency = Low
		case "●●○":
			entry.Frequency = Mid
		case "●●●":
			entry.Frequency = High
		}
	}

	grammar := htmlquery.FindOne(node, `//span[@class="GRAM"]`)
	if grammar != nil {
		entry.Grammar = strings.Trim(htmlquery.InnerText(grammar), "[ ]")
	}

	for _, node := range htmlquery.Find(node, `//span[@data-src-mp3]`) {
		attr := htmlquery.FindOne(node, `//@data-src-mp3`)
		src := innerTextTrim(attr)
		src = src[:strings.Index(src, "?")]

		if hasClass(node, "brefile") {
			entry.BritAudioUrl = src
			continue
		}

		if hasClass(node, "amefile") {
			entry.AmericanAudioUrl = src
			continue
		}
	}

	entry.Senses = make([]LdSense, 0)
	for _, node := range htmlquery.Find(node, `//span[@class="Sense"]`) {
		crossRefNode := htmlquery.FindOne(node, `//a[@class="crossRef"]`)
		if crossRefNode != nil {
			entry.Senses = append(entry.Senses, parseCrossRef(crossRefNode))
		} else {
			entry.Senses = append(entry.Senses, parseSense(node))
		}
	}

	return entry
}

func parseCrossRef(node *html.Node) CrossRefSense {
	crossRef := CrossRefSense{}

	crossRef.Text = innerTextTrim(node)
	crossRef.Ref = innerTextTrim(htmlquery.FindOne(node, "//@href"))

	return crossRef
}

func parseSense(node *html.Node) Sense {
	sense := Sense{}

	signpost := htmlquery.FindOne(node, `//span[@class="SIGNPOST"]`)
	if signpost != nil {
		sense.Signpost = innerTextTrim(signpost)
	}

	grammar := htmlquery.FindOne(node, `//span[@class="GRAM"]`)
	if grammar != nil {
		sense.Grammar = strings.Trim(htmlquery.InnerText(grammar), "[ ]")
	}

	geo := htmlquery.FindOne(node, `//span[@class="GEO"]`)
	if geo != nil {
		sense.Geo = innerTextTrim(geo)
	}

	definition := htmlquery.FindOne(node, `//span[@class="DEF"]`)
	if definition != nil {
		sense.Definition = innerTextTrim(definition)
	}

	synonyms := htmlquery.FindOne(node, `//span[@class="SYN"]`)
	if synonyms != nil {
		sense.Synonyms = innerTextTrim(synonyms)[4:]
	}

	sense.Examples = make([]Example, 0)
	for _, node := range htmlquery.Find(node, `//span[@class="EXAMPLE"]`) {
		sense.Examples = append(sense.Examples, parseExample(node))
	}

	return sense
}

func parseExample(node *html.Node) Example {
	example := Example{}

	example.Text = innerTextTrim(node)

	scrAttr := htmlquery.FindOne(node, `//span/@data-src-mp3`)
	if scrAttr != nil {
		src := innerTextTrim(scrAttr)
		example.AudioUrl = src[:strings.Index(src, "?")]
	}

	if hasClass(node.Parent, "ColloExa") {
		example.Colloquial = innerTextTrim(node.Parent.FirstChild)
	}

	return example
}
