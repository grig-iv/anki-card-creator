package ld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	page, err := ParseUrl("https://www.ldoceonline.com/dictionary/programming")
	if err != nil {
		t.Fatal(err)
	}

	expected := WordPage{
		Title: "programming",
		Entries: []DictEntry{
			{
				Hyphenation:      "pro‧gram‧ming",
				Pronunciation:    "/ˈprəʊɡræmɪŋ $ ˈproʊ-/",
				Frequency:        Low,
				PartOfSpeach:     "noun",
				Grammar:          "uncountable",
				BritAudioUrl:     "https://www.ldoceonline.com/media/english/breProns/programming0205.mp3",
				AmericanAudioUrl: "https://www.ldoceonline.com/media/english/ameProns/laadprogramming.mp3",
				Senses: []LdSense{
					LdSense(Sense{
						Definition: "the activity of writing programs for computers, or something written by a programmer",
						Synonyms:   "",
						Examples: []Example{{
							Text:       "a course in computer programming",
							AudioUrl:   "https://www.ldoceonline.com/media/english/exaProns/p008-000299324.mp3",
							Colloquial: "",
						}},
					}),
					LdSense(Sense{
						Definition: "television or radio programmes, or the planning of these broadcasts",
						Synonyms:   "",
						Examples: []Example{{
							Text:       "The Winter Olympics received over 160 hours of television programming.",
							AudioUrl:   "https://www.ldoceonline.com/media/english/exaProns/p008-000299329.mp3",
							Colloquial: "",
						}},
					}),
				},
			},
		},
	}

	assert.Equal(t, expected, page)
}

func TestParseUrl_Mental(t *testing.T) {
	page, err := ParseUrl("https://www.ldoceonline.com/dictionary/mental")
	if err != nil {
		t.Fatal(err)
	}

	wordPage, _ := page.(WordPage)

	sense1, _ := wordPage.Entries[0].Senses[1].(Sense)
	assert.Equal(t, "only before noun", sense1.Grammar)
	assert.Equal(t, "mental picture/image", sense1.Examples[2].Colloquial)

	sense2, _ := wordPage.Entries[0].Senses[2].(CrossRefSense)
	assert.Equal(t, "make a mental note", sense2.Text)
	assert.Equal(t, "/dictionary/make-a-mental-note", sense2.Ref)

	sense3, _ := wordPage.Entries[0].Senses[5].(Sense)
	assert.Equal(t, "British English", sense3.Geo)
}
