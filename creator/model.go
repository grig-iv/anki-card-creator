package creator

import "github.com/grig-iv/anki-card-creator/anki"

const modelName = "AnkiCardCreator Basic"

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
    <p>{{example_text}}</p>
    
    {{exapple_audio_url}}
</div>
`

const sentenceTemplateBack = `
<div>
    <h2>{{word}}</h2>
    
    {{word_audio_url}}
    
    <p><strong>Hyphenation:</strong> {{hyphenation}}</p>
    <p><strong>Pronunciation:</strong> {{pronunciation}}</p>
    <p><strong>Part of Speech:</strong> {{part_of_speach}}</p>
    <p><strong>Grammar:</strong> {{grammar}}</p>
    
    {{#signpost}}
        <p><strong>Signpost:</strong> {{signpost}}</p>
    {{/signpost}}
    
    {{#sense_grammar}}
        <p><strong>Sense Grammar:</strong> {{sense_grammar}}</p>
    {{/sense_grammar}}
    
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
		ModelName:     modelName,
		Css:           modelCss,
		InOrderFields: cardFields(),
		IsCloze:       false,
		CardTemplates: []anki.CardTemplate{
			{
				Name:  "Sentence",
				Front: sentenceTemplateFront,
				Back:  sentenceTemplateBack,
			},
		},
	})
}
