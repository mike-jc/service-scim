package auth

import (
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego/context"
	"strings"
)

type HttpBasic struct {
	Abstract

	user     string
	password string
}

func (a *HttpBasic) SetCredentials(user, password string) {
	a.user = user
	a.password = password
}

func (a *HttpBasic) Auth(ctx *context.Context) error {
	if a.user == "" {
		return errors.New("Credentials for HTTP Basic authorization are not set")
	}

	if authHeader := ctx.Input.Header("Authorization"); authHeader == "" {
		return errors.New("No 'Authorization' header in request")
	} else if strings.Index(authHeader, "Basic ") != 0 {
		return errors.New("'Authorization' header is not for HTTP Basic authorization: " + authHeader)
	} else {
		encodedCredentials := strings.Replace(authHeader, "Basic ", "", 1)
		if credentials, err := base64.StdEncoding.DecodeString(encodedCredentials); err != nil {
			return errors.New("Can not decode credentials for HTTP Basic authorization. Encoded credentials: " + encodedCredentials)
		} else {
			if parts := strings.Split(string(credentials), ":"); len(parts) != 2 {
				return errors.New("Wrong format of credentials for HTTP Basic authorization: " + string(credentials))
			} else if parts[0] != a.user || parts[1] != a.password {
				return errors.New("Invalid credentials for HTTP Basic authorization")
			}
		}
	}

	return nil
}
