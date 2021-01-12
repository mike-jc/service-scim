package restApi

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"net/http/httputil"
	"runtime/debug"
	"service-scim/errors/sdks"
	"service-scim/sdks"
)

const RequestDumpMaxLength = 1000

type AbstractRestApi struct {
	instanceDomain string
}

func (a *AbstractRestApi) SetInstanceDomain(instanceDomain string) {
	a.instanceDomain = instanceDomain
}

func (a *AbstractRestApi) tuneRequest(request *httplib.BeegoHTTPRequest, headers interface{}) {
	request.Debug(false)

	if headers != nil {
		if headers, ok := headers.(map[string]string); ok {
			for name, value := range headers {
				request.Header(name, value)
			}
		}
	}
}

func (a *AbstractRestApi) addQueryParams(request *httplib.BeegoHTTPRequest, origParams interface{}) {
	if origParams == nil {
		return
	}

	if paramsMap, ok := origParams.(map[string]interface{}); ok {
		for name, item := range paramsMap {
			if valuesMap, ok := item.(map[string]interface{}); ok {
				for mapKey, subValue := range valuesMap {
					if subValueStr := fmt.Sprintf("%v", subValue); len(subValueStr) > 0 {
						request.Param(fmt.Sprintf("%s[%s]", name, mapKey), subValueStr)
					}
				}
			} else if values, ok := item.([]string); ok {
				for _, subValue := range values {
					request.Param(name, subValue)
				}
			} else if value := fmt.Sprintf("%v", item); len(value) > 0 {
				request.Param(name, value)
			}
		}
	}
}

func (a *AbstractRestApi) requestDump(request *httplib.BeegoHTTPRequest) string {
	cutDump := func(d []byte) string {
		if len(d) > RequestDumpMaxLength {
			return string(d[:RequestDumpMaxLength])
		} else {
			return string(d)
		}
	}

	if dump, err := httputil.DumpRequest(request.GetRequest(), true); err == nil {
		return cutDump(dump)
	} else if dump := request.DumpRequest(); len(dump) > 0 {
		return cutDump(dump)
	}
	return ""
}

func (a *AbstractRestApi) sendRequest(request *httplib.BeegoHTTPRequest) (output string, err errorsSdks.Interface) {
	requestDump := a.requestDump(request)

	response, reqErr := request.Response()
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if reqErr != nil {
		err = errorsSdks.NewError(reqErr.Error(), 0, "")
		sdks.LogMain.Log(logger.CreateError(err.Error()).SetErrorCode(15010).
			SetInstance(a.instanceDomain).
			AddData("requestDump", requestDump).
			AddData("url", request.GetRequest().URL).
			AddData("response", err.Response()).
			AddData("stackTrace", string(debug.Stack())))
		return
	}

	if output, reqErr = request.String(); reqErr != nil {
		err = errorsSdks.NewError(reqErr.Error(), 0, "")
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		err = errorsSdks.NewError("Not OK status", response.StatusCode, output)
		sdks.LogMain.Log(logger.CreateError(err.Error()).SetErrorCode(15010).
			SetInstance(a.instanceDomain).
			AddData("requestDump", requestDump).
			AddData("url", request.GetRequest().URL).
			AddData("response", err.Response()).
			AddData("stackTrace", string(debug.Stack())))
		return
	}
	return
}
