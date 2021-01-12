package sdks

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"log"
	"service-scim/models/sdks/jwt"
	"sync"
	"time"
)

type Jwt struct {
	jwtSecret        string
	onceGetJwtSecret sync.Once
}

type JwtPayload struct {
	Domain   string `json:"domain"`
	Id       string `json:"id"`
	Instance string `json:"instance"`
}

type JwtClaims struct {
	jwt.StandardClaims

	Recipient string     `json:"rec"`
	Payload   JwtPayload `json:"payload"`
}

func (*Jwt) CreateUserPayload(userId string, instance string) *JwtPayload {
	return &JwtPayload{
		Domain:   modelsSdkJwt.DomainUserId,
		Id:       userId,
		Instance: instance,
	}
}

func (o *Jwt) CreateSystemUserPayload(instance string) *JwtPayload {
	return o.CreateUserPayload("-1", instance)
}

func (o *Jwt) GenerateToken(payload *JwtPayload, ttl time.Duration, recipient string) (tokenStr string, err error) {
	// Create the Claims
	claims := &JwtClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
		},
		recipient,
		*payload,
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = tokenObj.SignedString([]byte(o.getJwtSecret()))
	return
}

func (o *Jwt) CreateFromClaims(claimsPayload *modelsSdkJwt.ClaimsUserPayload, recipient string, liveTime time.Duration) (token string, err error) {
	payload := o.CreateUserPayload(claimsPayload.Id, claimsPayload.Instance)
	token, err = o.GenerateToken(payload, liveTime, recipient)
	if err != nil {
		err = errors.New(fmt.Sprintf("Can not generate token for %s service (%s:%s): %s", recipient, claimsPayload.Id, claimsPayload.Instance, err.Error()))
		LogMain.Log(logger.CreateError(err.Error()).SetCode("sdk.jwt.generate.error"))
	}
	return
}

func (o *Jwt) getJwtSecret() string {
	o.onceGetJwtSecret.Do(func() {
		o.jwtSecret = beego.AppConfig.String("jwtSecret")
		if o.jwtSecret == "" {
			log.Panic("No jwtSecret in config")
		}
	})
	return o.jwtSecret
}
