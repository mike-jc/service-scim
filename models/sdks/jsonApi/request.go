package modelsSdkJsonApi

type Request struct {
	Data RequestData `json:"data"`
}

type RequestData struct {
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes"`
}
