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

	expected := wordPage{
		title: "programming",
		entries: []dictEntry{
			{
				hyphenation:      "pro‧gram‧ming",
				pronunciation:    "/ˈprəʊɡræmɪŋ $ ˈproʊ-/",
				frequency:        low,
				partOfSpeach:     "noun",
				grammar:          "uncountable",
				britAudioUrl:     "https://www.ldoceonline.com/media/english/breProns/programming0205.mp3",
				americanAudioUrl: "https://www.ldoceonline.com/media/english/ameProns/laadprogramming.mp3",
				senses: []ldSense{
					ldSense(sense{
						definition: "the activity of writing programs for computers, or something written by a programmer",
						synonyms:   "",
						examples: []example{{
							text:       "a course in computer programming",
							audioUrl:   "https://www.ldoceonline.com/media/english/exaProns/p008-000299324.mp3",
							colloquial: "",
						}},
					}),
					ldSense(sense{
						definition: "television or radio programmes, or the planning of these broadcasts",
						synonyms:   "",
						examples: []example{{
							text:       "The Winter Olympics received over 160 hours of television programming.",
							audioUrl:   "https://www.ldoceonline.com/media/english/exaProns/p008-000299329.mp3",
							colloquial: "",
						}},
					}),
				},
			},
		},
	}

	assert.Equal(t, expected, page)
}

func TestParseUrl_Mind(t *testing.T) {
	page, err := ParseUrl("https://www.ldoceonline.com/dictionary/mind")
	if err != nil {
		t.Fatal(err)
	}

	wordPage, _ := page.(wordPage)

	sense1, _ := wordPage.entries[0].senses[0].(sense)
	assert.Equal(t, "countable, uncountable", sense1.grammar)
	assert.Equal(t, "in somebody’s mind", sense1.examples[3].colloquial)

	sense2, _ := wordPage.entries[0].senses[13].(sense)
	assert.Equal(t, "intelligent person", sense2.signpost)
	assert.Equal(t, "brain", sense2.synonyms)
}
