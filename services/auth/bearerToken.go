package auth

import (
	"errors"
	"github.com/astaxie/beego/context"
	"strings"
)

type BearerToken struct {
	Abstract

	token string
}

func (a *BearerToken) SetToken(token string) {
	a.token = token
}

func (a *BearerToken) Auth(ctx *context.Context) error {
	if a.token == "" {
		return errors.New("Token for Bearer Token authorization is not set")
	}

	if authHeader := ctx.Input.Header("Authorization"); authHeader == "" {
		return errors.New("No 'Authorization' header in request")
	} else if strings.Index(authHeader, "Bearer ") != 0 {
		return errors.New("'Authorization' header is not for Bearer Token authorization: " + authHeader)
	} else {
		token := strings.Replace(authHeader, "Bearer ", "", 1)
		if token != a.token {
			return errors.New("Invalid credentials for Bearer Token authorization")
		}
	}

	return nil
}
