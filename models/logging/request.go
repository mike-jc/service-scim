package modelsLogging

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method  string      `json:"method"`
	Url     *url.URL    `json:"url"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}
