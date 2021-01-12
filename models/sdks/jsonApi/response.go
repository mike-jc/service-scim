package modelsSdkJsonApi

import "encoding/json"

type Response struct {
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Type       string          `json:"type"`
	Attributes json.RawMessage `json:"attributes"`
}
