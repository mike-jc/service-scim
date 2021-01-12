package auth_test

import (
	"encoding/base64"
	"service-scim/services/auth"
	"service-scim/tests"
	"testing"
)

var b *auth.HttpBasic
var user = "test-user"
var password = "Sec!r3t"

func init() {
	tests.Init()

	b = new(auth.HttpBasic)
	b.SetCredentials(user, password)
}

func TestBasicAuthOK(t *testing.T) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
	ctx := tests.NewContext()
	ctx.Input.Context.Request.Header.Set("Authorization", authHeader)

	if err := b.Auth(ctx); err != nil {
		t.Errorf("HTTP Basic Authentication failed, should be successful: %s", err.Error())
	}
}

func TestBasicAuthFailed(t *testing.T) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":some-stuff-"+password))
	ctx := tests.NewContext()
	ctx.Input.Context.Request.Header.Set("Authorization", authHeader)

	if err := b.Auth(ctx); err == nil {
		t.Errorf("HTTP Basic Authentication is successful, should failed: %s", err.Error())
	} else if err.Error() != "Invalid credentials for HTTP Basic authorization" {
		t.Errorf("HTTP Basic Authentication failed with wrong error: %s", err.Error())
	}
}
