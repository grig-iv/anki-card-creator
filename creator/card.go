package creator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/grig-iv/anki-card-creator/anki"
)

type Card struct {
	Word            string `json:"word"`
	Hyphenation     string `json:"hyphenation"`
	Pronunciation   string `json:"pronunciation"`
	PartOfSpeach    string `json:"part_of_speach"`
	Grammar         string `json:"grammar"`
	WordAudioUrl    string `json:"word_audio"`
	Signpost        string `json:"signpost"`
	SenseGrammar    string `json:"sense_grammar"`
	Geo             string `json:"geo"`
	Definition      string `json:"definition"`
	Synoyms         string `json:"synoyms"`
	ExampleText     string `json:"example_text"`
	ExampleAudioUrl string `json:"example_audio"`
}

func CreateCard(card Card, tags []string) error {
	fieldsMap := card.toMap()

	if card.WordAudioUrl != "" {
		fieldsMap["word_audio"] = downloadAudio(card.WordAudioUrl)
	}

	if card.ExampleAudioUrl != "" {
		fieldsMap["example_audio"] = downloadAudio(card.ExampleAudioUrl)
	}

	return anki.AddNote(anki.AddNoteParams{
		DeckName:  deckName,
		ModelName: modelName,
		Fields:    fieldsMap,
		Options: anki.AddNoteOptions{
			AllowDuplicate: false,
			DuplicateScope: deckName,
			DuplicateScopeOptions: anki.DuplicateScopeOptions{
				DeckName:       modelName,
				CheckChildren:  false,
				CheckAllModels: false,
			},
		},
		Tags:    tags,
		Audio:   []anki.Media{},
		Video:   []anki.Media{},
		Picture: []anki.Media{},
	})
}

func downloadAudio(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		return ""
	}

	const mediaFolder = "/home/grig/.local/share/Anki2/User 1/collection.media/"
	audioName := path.Base(url)
	audioPath := path.Join(mediaFolder, audioName)
	audioFile, err := os.OpenFile(audioPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Sprintf("[sound:%s]", audioName)
		}

		return ""
	}

	_, err = io.Copy(audioFile, resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("[sound:%s]", audioName)
}

func createAnkiConnectAudio(card Card) []anki.Media {
	audio := make([]anki.Media, 0)

	if card.WordAudioUrl != "" {
		media := anki.Media{
			Url:      card.WordAudioUrl,
			Filename: path.Base(card.WordAudioUrl),
			SkipHash: "",
			Fields:   []string{"word_audio"},
		}
		audio = append(audio, media)
	}

	if card.ExampleAudioUrl != "" {
		media := anki.Media{
			Url:      card.ExampleAudioUrl,
			Filename: path.Base(card.ExampleAudioUrl),
			SkipHash: "",
			Fields:   []string{"example_audio"},
		}
		audio = append(audio, media)
	}

	return audio
}

func (c Card) toMap() map[string]string {
	cType := reflect.TypeOf(c)
	cValue := reflect.ValueOf(c)

	fieldMap := make(map[string]string, cType.NumField())
	for i := range cType.NumField() {
		field := cType.Field(i)
		value := cValue.Field(i)
		fieldName := strings.Replace(string(field.Tag), "json:", "", 1)
		fieldName = strings.Trim(fieldName, `"`)
		fieldMap[fieldName] = value.String()
	}

	return fieldMap
}

func cardFields() []string {
	fieldMap := Card{}.toMap()

	fields := make([]string, 0, len(fieldMap))
	for k := range fieldMap {
		fields = append(fields, k)
	}

	return fields
}
