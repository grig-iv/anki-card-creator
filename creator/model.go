package creator

import "github.com/grig-iv/anki-card-creator/anki"

const modelName = "LogmaDictionaryEng"

const modelCss = `
<style>
    div {
        font-family: Arial, sans-serif;
        line-height: 1.6;
    }
    h2 {
        color: #2A6DB0;
    }
    p {
        margin: 5px 0;
    }
    audio {
        margin-top: 10px;
    }
</style>
`

const sentenceTemplateFront = `
<div>
    <h2>Example Sentence</h2>
    <p>{{exampleText}}</p>
    
    {{exappleAudioUrl}}
</div>
`

const sentenceTemplateBack = `
<div>
    <h2>{{word}}</h2>
    
    {{wordAudioUrl}}
    
    <p><strong>Hyphenation:</strong> {{hyphenation}}</p>
    <p><strong>Pronunciation:</strong> {{pronunciation}}</p>
    <p><strong>Part of Speech:</strong> {{partOfSpeach}}</p>
    <p><strong>Grammar:</strong> {{grammar}}</p>
    
    {{#signpost}}
        <p><strong>Signpost:</strong> {{signpost}}</p>
    {{/signpost}}
    
    {{#senseGrammar}}
        <p><strong>Sense Grammar:</strong> {{senseGrammar}}</p>
    {{/senseGrammar}}
    
    {{#geo}}
        <p><strong>Geographical Information:</strong> {{geo}}</p>
    {{/geo}}
    
    <p><strong>Definition:</strong> {{definition}}</p>
    
    {{#synoyms}}
        <p><strong>Synonyms:</strong> {{synoyms}}</p>
    {{/synoyms}}
</div>
`

func CreateModel() error {
	return anki.CreateModel(anki.CreateModelParams{
		ModelName: "LogmaDictionaryEng",
		Css:       modelCss,
		InOrderFields: []string{
			"word",
			"hyphenation",
			"pronunciation",
			"partOfSpeach",
			"grammar",
			"wordAudioUrl",
			"signpost",
			"senseGrammar",
			"geo",
			"definition",
			"synoyms",
			"exampleText",
			"exappleAudioUrl",
		},
		IsCloze: false,
		CardTemplates: []anki.CardTemplate{
			{
				Name:  "Sentence",
				Front: sentenceTemplateFront,
				Back:  sentenceTemplateBack,
			},
		},
	})
}
