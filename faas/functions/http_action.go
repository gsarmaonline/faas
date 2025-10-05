package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gsarmaonline/faas/faas/intf"
)

const (
	// Http Methods
	GetHttpMethod  = HttpMethodT("GET")
	PostHttpMethod = HttpMethodT("POST")
)

type (
	HttpMethodT string

	HttpInput struct {
		Url         string      `json:"url"`
		Method      HttpMethodT `json:"method"`
		RequestBody interface{} `json:"request_body"`
	}

	HttpAction struct {
		Input HttpInput
	}
)

func NewHttpAction() (httpAction *HttpAction) {
	return &HttpAction{}
}

func (httpAction HttpAction) GetConfig() intf.FunctionConfig {
	return intf.FunctionConfig{Name: "http"}
}

func (httpAction *HttpAction) ParsePayload(payload intf.Payload) error {
	processedInput := HttpInput{
		Url:    payload["url"].(string),
		Method: HttpMethodT(payload["method"].(string)),
	}

	// Optional request body
	if requestBody, exists := payload["request_body"]; exists {
		processedInput.RequestBody = requestBody
	}

	httpAction.Input = processedInput
	return nil
}

func (httpAction HttpAction) Validate() (err error) {
	if httpAction.Input.Url == "" {
		return fmt.Errorf("missing required field: url")
	}
	if httpAction.Input.Method == "" {
		return fmt.Errorf("missing required field: method")
	}
	return nil
}

func (httpAction HttpAction) Execute() (output intf.FunctionOutput, err error) {
	var (
		client   *http.Client
		req      *http.Request
		reqBody  io.Reader
		payloadB []byte
	)

	client = &http.Client{}
	if payloadB, err = json.Marshal(httpAction.Input.RequestBody); err != nil {
		return
	}
	reqBody = bytes.NewBuffer(payloadB)
	if req, err = http.NewRequest(string(httpAction.Input.Method), httpAction.Input.Url, reqBody); err != nil {
		return
	}
	if _, err = client.Do(req); err != nil {
		return
	}

	return
}
