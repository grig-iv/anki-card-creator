package anki

import "encoding/json"

type modelNamesResult []string

func ModelNames() (modelNamesResult, error) {
	resp, err := SendRequest("modelNames", nil)
	if err != nil {
		return nil, err
	}

	if resp.HasError() {
		return nil, resp.ToError()
	}

	result := make([]string, 0)
	err = json.Unmarshal([]byte(resp.Result), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
