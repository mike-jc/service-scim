package auth

import (
	"github.com/astaxie/beego/context"
)

type Interface interface {
	Auth(ctx *context.Context) error
}

type Abstract struct {
	Interface
}

func (a *Abstract) Auth(ctx *context.Context) error {
	return nil
}
