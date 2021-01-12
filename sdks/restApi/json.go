package restApi

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"runtime/debug"
	"service-scim/errors/sdks"
	"service-scim/models/sdks/jsonApi"
	"service-scim/sdks"
)

type RestJson struct {
	AbstractRestApi
}

func (j *RestJson) tuneRequest(request *httplib.BeegoHTTPRequest, headers interface{}) {
	j.AbstractRestApi.tuneRequest(request, headers)

	request.Header("Accept", "application/json;charset=UTF-8")
	request.Header("Content-Type", "application/json;charset=UTF-8")
}

func (j *RestJson) sendRequest(request *httplib.BeegoHTTPRequest, result interface{}) errorsSdks.Interface {
	output, err := j.AbstractRestApi.sendRequest(request)
	sdks.LogMain.Log(logger.CreateInfo("Request was sent").
		SetCode("sdk.http.request").
		SetInstance(j.instanceDomain).
		AddData("url", request.GetRequest().URL).
		AddData("requestDump", j.requestDump(request)).
		AddData("output", output).
		AddData("stackTrace", string(debug.Stack())))
	if err != nil {
		return err
	}

	if result != nil {
		if jErr := json.Unmarshal([]byte(output), result); jErr != nil {
			err := errorsSdks.NewError("Can not parse the response: "+jErr.Error(), 0, output)
			sdks.LogMain.Log(logger.CreateError(err.Error()).SetErrorCode(15010).
				SetInstance(j.instanceDomain).
				AddData("requestDump", j.requestDump(request)).
				AddData("url", request.GetRequest().URL).
				AddData("response", err.Response()).
				AddData("stackTrace", string(debug.Stack())))
			return err
		}
	}
	return nil
}

func (j *RestJson) Get(url string, params interface{}, headers interface{}, result interface{}) errorsSdks.Interface {
	request := httplib.Get(url)
	j.tuneRequest(request, headers)
	j.addQueryParams(request, params)

	if err := j.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}

func (j *RestJson) Post(url string, params interface{}, headers interface{}, result interface{}) errorsSdks.Interface {
	request := httplib.Post(url)
	j.tuneRequest(request, headers)

	if params != nil {
		body, jErr := json.Marshal(params)
		if jErr != nil {
			return errorsSdks.NewError(jErr.Error(), 0, "")
		}
		request.Body(body)
	}

	if err := j.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}

func (j *RestJson) Patch(url string, params interface{}, headers interface{}, result interface{}) errorsSdks.Interface {
	request := httplib.NewBeegoRequest(url, "PATCH")
	j.tuneRequest(request, headers)

	if params != nil {
		body, jErr := json.Marshal(params)
		if jErr != nil {
			return errorsSdks.NewError(jErr.Error(), 0, "")
		}
		request.Body(body)
	}

	if err := j.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}
func (j *RestJson) Put(url string, params interface{}, headers interface{}, result interface{}) errorsSdks.Interface {
	request := httplib.Put(url)
	j.tuneRequest(request, headers)

	if params != nil {
		body, jErr := json.Marshal(params)
		if jErr != nil {
			return errorsSdks.NewError(jErr.Error(), 0, "")
		}
		request.Body(body)
	}

	if err := j.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}

func (j *RestJson) Delete(url string, headers interface{}) errorsSdks.Interface {
	request := httplib.Delete(url)
	j.tuneRequest(request, headers)

	if err := j.sendRequest(request, nil); err != nil {
		return err
	}
	return nil
}

// bilt JSON API request with the given data
func (j *RestJson) builtJsonApiRequest(attributes interface{}, typeName string) *modelsSdkJsonApi.Request {
	return &modelsSdkJsonApi.Request{
		Data: modelsSdkJsonApi.RequestData{
			Type:       typeName,
			Attributes: attributes,
		},
	}
}

// parse JSON API response and return an object basing on data from that
func (j *RestJson) parseJsonApiResponse(response *modelsSdkJsonApi.Response, object interface{}, typeName string) error {
	if response.Data.Type != typeName {
		return fmt.Errorf("Incorrect response type %s, should be %s", response.Data.Type, typeName)
	}
	if jErr := json.Unmarshal(response.Data.Attributes, object); jErr != nil {
		return fmt.Errorf("Can not unmarshall JSON in response: %s", jErr)
	}
	return nil
}
