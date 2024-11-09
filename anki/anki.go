package anki

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type Request struct {
	Action  string      `json:"action"`
	Version int         `json:"version"`
	Params  interface{} `json:"params,omitempty"`
}

type Response struct {
	Error string `json:"error"`
}

func (r Response) HasError() bool {
	return r.Error != ""
}

func (r Response) ToError() error {
	if r.HasError() {
		return errors.New(r.Error)
	}
	return nil
}

func IsRunning() (bool, error) {
	err := SendRequest("version", nil)
	return err == nil, err
}

func SendRequest(name string, params interface{}) error {
	request := Request{
		Action:  name,
		Version: 6,
		Params:  params,
	}

	requestByte, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("SendRequest: %s\n", requestByte)

	resp, err := http.Post(
		"http://localhost:8765",
		"application/json",
		bytes.NewReader(requestByte),
	)

	if err != nil {
		return err
	}

	response := &Response{}
	responseStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(responseStr), response)
	if err != nil {
		return err
	}

	if response.HasError() {
		return response.ToError()
	}

	return nil
}
