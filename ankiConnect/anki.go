package ankiConnect

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Request struct {
	Action  string      `json:"action"`
	Version int         `json:"version"`
	Params  interface{} `json:"params,omitempty"`
}

type Response struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

func IsRunning() (bool, error) {
	_, err := SendRequest("version", nil)
	return err == nil, err
}

func SendRequest(name string, params interface{}) (*Response, error) {
	request := Request{
		Action:  name,
		Version: 6,
		Params:  params,
	}

	requestStr, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		"http://localhost:8765",
		"application/json",
		bytes.NewReader(requestStr),
	)

	if err != nil {
		return nil, err
	}

	response := &Response{}
	responseStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(responseStr), response)
	if err != nil {
		return nil, err
	}

	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response, nil
}
