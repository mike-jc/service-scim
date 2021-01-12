package controllers

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gitlab.com/24sessions/lib-go-logger/logger"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/models"
	"service-scim/models/logging"
	"strconv"
	"strings"
	"sync"
)

type AbstractController struct {
	beego.Controller

	lgLock sync.Mutex
	lg     *logger.Logger

	domain     string
	scimDomain string
	scimConfig *sdksData.ScimContainer
}

func (c *AbstractController) getClientIp() string {
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	if len(ip) > 0 && ip[0] != "[" {
		return ip[0]
	}
	return "127.0.0.1"
}

func (c *AbstractController) GetLogger() *logger.Logger {
	c.lgLock.Lock()
	defer c.lgLock.Unlock()

	if c.lg == nil {
		c.lg = new(logger.Logger).
			SetSubject("anonymous", "").
			SetTraceId(uuid.NewV4String()).
			SetClientIp(c.getClientIp()).
			SetParent(LogMain)
	}
	return c.lg
}

func (c *AbstractController) addCors() {
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE")
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Authorization,Content-type,Accept-Language")
}

func (c *AbstractController) Format() string {
	var format string
	if c.scimConfig != nil && c.scimConfig.ResponseFormat() != "" {
		format = c.scimConfig.ResponseFormat()
	} else {
		format = sdksData.ScimDefaultResponseFormat
	}
	return format
}

func (c *AbstractController) ServeResponse(data interface{}) {
	format := c.Format()
	if format == "json" {
		c.Data["json"] = data
		c.ServeJSON()
	} else if format == "xml" {
		c.Data["xml"] = data
		c.ServeXML()
	} else {
		c.GetLogger().LogFatal(logger.CreateError("Invalid response format").SetCode("app.invalid_format"))
	}
}

func (c *AbstractController) ServeResponseWithStatus(data interface{}, status int) {
	c.Ctx.Output.SetStatus(status)
	c.ServeResponse(data)
}

func (c *AbstractController) SuccessResponse() {
	result := make(map[string]interface{})
	result["status"] = "success"

	c.addCors()
	c.ServeResponse(result)
}

func (c *AbstractController) SuccessResponseWithStatus(status int) {
	c.Ctx.Output.SetStatus(status)
	c.SuccessResponse()
}

func (c *AbstractController) ShowError(message string, status int, errorCode string, reason string, showReason bool) {
	// log the error
	req, reqJson := c.MarshalRequest(c.Ctx)
	logRow := logger.CreateError(fmt.Sprintf("%d: %s, %s", status, message, reason)).
		SetCode(errorCode).
		AddData("request", reqJson).
		AddData("requestBody", req.Body)

	if status < 500 {
		logRow.SetLevel(logger.LOG_LEVEL_NOTICE)
	}

	c.GetLogger().Log(logRow)

	// output the error
	err := &models.Error{
		Schemas:  []string{"urn:ietf:params:scim:api:messages:2.0:Error"},
		Status:   strconv.Itoa(status),
		ScimType: errorCode,
		Detail:   message,
	}
	if showReason {
		err.Detail += ". " + reason
	}

	c.Ctx.Output.Status = status
	c.addCors()
	c.ServeResponse(err)
}

func (c *AbstractController) MarshalRequest(ctx *context.Context) (modelsLogging.Request, string) {
	res := modelsLogging.Request{
		Method:  ctx.Request.Method,
		Url:     ctx.Request.URL,
		Headers: ctx.Request.Header,
	}

	switch ctx.Request.Method {
	case "POST", "PATCH", "PUT":
		requestBody := ctx.Input.RequestBody
		if decodedBody, err := base64.StdEncoding.DecodeString(string(requestBody)); err == nil {
			requestBody = decodedBody
		}
		res.Body = string(requestBody)
	}

	bytes, _ := json.Marshal(res)
	return res, string(bytes)
}

func (c *AbstractController) BasicURL() string {
	port := c.Ctx.Input.Port()
	if port == 80 {
		return c.Ctx.Input.Site()
	} else {
		return c.Ctx.Input.Site() + ":" + strconv.Itoa(c.Ctx.Input.Port())
	}
}

func (c *AbstractController) UnmarshalRequestBody(data interface{}) (result interface{}, format string, err error) {
	// detect request format
	switch c.Ctx.Request.Header.Get("Content-Type") {
	case "application/scim+json", "application/json":
		format = "json"
	case "application/scim+xml", "application/xml":
		format = "xml"
	default:
		if c.scimConfig != nil && c.scimConfig.ResponseFormat() != "" {
			format = c.scimConfig.ResponseFormat()
		} else {
			format = sdksData.ScimDefaultResponseFormat
		}
	}

	requestBody := c.Ctx.Input.RequestBody
	if decodedBody, err := base64.StdEncoding.DecodeString(string(requestBody)); err == nil {
		requestBody = decodedBody
	}

	// unmarshal
	switch format {
	case "json":
		decoder := json.NewDecoder(strings.NewReader(string(requestBody)))
		decoder.UseNumber()
		if jErr := decoder.Decode(data); jErr != nil {
			err = fmt.Errorf("Cannot unserialize JSON: %s", jErr)
		}
		return data, format, nil
	case "xml":
		decoder := xml.NewDecoder(strings.NewReader(string(requestBody)))
		if xErr := decoder.Decode(data); xErr != nil {
			err = fmt.Errorf("Cannot unserialize XML: %s", xErr)
		}
		return data, format, nil
	default:
		err = fmt.Errorf("Unknown request format: %s", format)
		return nil, "", err
	}
}

func (c *AbstractController) UnmarshalRequestBodyToMap() (data map[string]interface{}, format string, err error) {
	data = make(map[string]interface{})
	if _, format, sErr := c.UnmarshalRequestBody(&data); sErr != nil {
		return nil, "", sErr
	} else {
		return data, format, nil
	}
}
