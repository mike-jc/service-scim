package modelsSdkJwt

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type ClaimsWithLeeway struct {
	jwt.StandardClaims

	leeway int64
}

// Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c ClaimsWithLeeway) Valid() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		log.Println("token is expired by " + delta.String())
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if c.VerifyIssuedAt(now, false) == false {
		vErr.Inner = fmt.Errorf("Token used before issued")
		log.Println("Token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		log.Println("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (c *ClaimsWithLeeway) VerifyExpiresAt(cmp int64, req bool) bool {
	return c.StandardClaims.VerifyExpiresAt(cmp-c.getLeeway(), req)
}

// Compares "iat" claim against "cmp".
// If required is false, this method will return true if the value matches or is unset
func (c *ClaimsWithLeeway) VerifyIssuedAt(cmp int64, req bool) bool {
	return c.StandardClaims.VerifyIssuedAt(cmp+c.getLeeway(), req)
}

func (c *ClaimsWithLeeway) getLeeway() int64 {
	if c.leeway <= 0 {
		leeway, _ := beego.AppConfig.Int("jwtLeeway")
		if leeway <= 0 {
			leeway = 60
		}
		c.leeway = int64(leeway)
	}
	return c.leeway
}
