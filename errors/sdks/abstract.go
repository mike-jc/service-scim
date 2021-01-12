package errorsSdks

import (
	"service-scim/models/sdks"
)

type Interface interface {
	SetError(text string)
	Response() modelsSdks.RestResponse

	error
}

type Error struct {
	Interface

	text     string
	response modelsSdks.RestResponse
}

func (e *Error) SetError(text string) {
	e.text = text
}

func (e *Error) Error() string {
	return e.text
}

func (e *Error) Response() modelsSdks.RestResponse {
	return e.response
}

func NewError(text string, code int, body string) *Error {
	return &Error{
		text: text,
		response: modelsSdks.RestResponse{
			Code: code,
			Body: body,
		},
	}
}
