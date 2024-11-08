package ankiConnect

type addNoteParams struct {
	DeckName  string            `json:"deckName"`
	ModelName string            `json:"modelName"`
	Fields    map[string]string `json:"fields"`
	Options   AddNoteOptions    `json:"options"`
	Tags      []string          `json:"tags"`
	Audio     []Media           `json:"audio"`
	Video     []Media           `json:"video"`
	Picture   []Media           `json:"picture"`
}

type Media struct {
	Url      string   `json:"url"`
	Filename string   `json:"filename"`
	SkipHash string   `json:"skipHash"`
	Fields   []string `json:"fields"`
}

type AddNoteOptions struct {
	AllowDuplicate        bool   `json:"allowDuplicate"`
	DuplicateScope        string `json:"duplicateScope"`
	DuplicateScopeOptions struct {
		DeckName       string `json:"deckName"`
		CheckChildren  bool   `json:"checkChildren"`
		CheckAllModels bool   `json:"checkAllModels"`
	} `json:"duplicateScopeOptions"`
}
