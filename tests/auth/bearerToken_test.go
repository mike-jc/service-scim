package auth_test

import (
	"service-scim/services/auth"
	"service-scim/tests"
	"testing"
)

var tk *auth.BearerToken
var token = "512a2d09d6c500"

func init() {
	tests.Init()

	tk = new(auth.BearerToken)
	tk.SetToken(token)
}

func TestTokenAuthOK(t *testing.T) {
	authHeader := "Bearer " + token
	ctx := tests.NewContext()
	ctx.Input.Context.Request.Header.Set("Authorization", authHeader)

	if err := tk.Auth(ctx); err != nil {
		t.Errorf("Bearer Token Authentication failed, should be successful: %s", err.Error())
	}
}

func TestTokenAuthFailed(t *testing.T) {
	authHeader := "Bearer some-stuff-" + token
	ctx := tests.NewContext()
	ctx.Input.Context.Request.Header.Set("Authorization", authHeader)

	if err := tk.Auth(ctx); err == nil {
		t.Errorf("Bearer Token Authentication is successful, should failed: %s", err.Error())
	} else if err.Error() != "Invalid credentials for Bearer Token authorization" {
		t.Errorf("Bearer Token Authentication failed with wrong error: %s", err.Error())
	}
}
