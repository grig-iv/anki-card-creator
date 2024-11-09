package anki

type modelNamesResult []string

func ModelNames() (modelNamesResult, error) {
	result := make([]string, 0)
	err := SendRequest("modelNames", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
