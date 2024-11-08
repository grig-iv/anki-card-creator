package ankiConnect

type CreateModelParams struct {
	ModelName     string         `json:"modelName"`
	Css           string         `json:"css,omitempty"`
	InOrderFields []string       `json:"inOrderFields,omitempty"`
	IsCloze       bool           `json:"isCloze,omitempty"`
	CardTemplates []CardTemplate `json:"cardTemplates,omitempty"`
}

type CardTemplate struct {
	Back  string `json:"back,omitempty"`
	Front string `json:"front,omitempty"`
	Name  string `json:"name,omitempty"`
}

func CreateModel(params CreateModelParams) error {
	_, err := SendRequest("createModel", params)
	return err
}
