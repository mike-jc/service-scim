package errorsSdks

import (
	"encoding/json"
	"fmt"
	"service-scim/models/sdks/keeper"
)

type KeeperInterface interface {
	SetError(text string)
	Response() *modelsSdkKeeper.Response

	error
}

type KeeperError struct {
	KeeperInterface

	text     string
	response *modelsSdkKeeper.Response
}

func (e *KeeperError) SetError(text string) {
	e.text = text
}

func (e *KeeperError) Error() string {
	return e.text
}

func (e *KeeperError) Response() *modelsSdkKeeper.Response {
	return e.response
}

func NewKeeperError(text string, response *modelsSdkKeeper.Response) *KeeperError {
	if response == nil {
		response = &modelsSdkKeeper.Response{}
	}
	return &KeeperError{
		text:     text,
		response: response,
	}
}

func NewKeeperErrorFromAbstract(err Interface) *KeeperError {
	var jBody modelsSdkKeeper.Error
	var newResponse *modelsSdkKeeper.Response

	text := err.Error()
	restResponse := err.Response()

	if jErr := json.Unmarshal([]byte(restResponse.Body), &jBody); jErr == nil {
		newResponse = &modelsSdkKeeper.Response{
			Code: restResponse.Code,
			Body: jBody,
		}
		text += fmt.Sprintf(" (%s)", jBody.Description)
	} else {
		newResponse = &modelsSdkKeeper.Response{
			Code:         restResponse.Code,
			ParsingError: fmt.Errorf("Can not parse the response: %s. Original body: %s", jErr.Error(), restResponse.Body),
		}
	}

	return &KeeperError{
		text:     text,
		response: newResponse,
	}
}
