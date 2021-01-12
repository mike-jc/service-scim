package resources

import (
	"github.com/astaxie/beego/context"
	sdksData "gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/services/middleware"
	"service-scim/services/validation"
	"strings"
	"sync"
)

type Abstract struct {
	includedAttributes []string
	excludedAttributes []string

	format string
	config *sdksData.ScimContainer

	typeValidator       *validation.AttributeType
	caseValidator       *validation.CorrectCase
	requireValidator    *validation.RequiredAttribute
	readOnlyValidator   *validation.ReadOnlyAttribute
	mutabilityValidator *validation.AttributeMutability
	uniquenessValidator *validation.Uniqueness

	middlewareFactory     *middleware.Factory
	middlewareFactoryOnce sync.Once
}

func (a *Abstract) Init() {
	a.typeValidator = new(validation.AttributeType)
	a.caseValidator = new(validation.CorrectCase)
	a.requireValidator = new(validation.RequiredAttribute)
	a.readOnlyValidator = new(validation.ReadOnlyAttribute)
	a.mutabilityValidator = new(validation.AttributeMutability)
	a.uniquenessValidator = new(validation.Uniqueness)
}

func (a *Abstract) AttributesFromRequest(ctx *context.Context) {
	a.includedAttributes = strings.Split(ctx.Input.Query("attributes"), ",")
	a.excludedAttributes = strings.Split(ctx.Input.Query("excludedAttributes"), ",")
}

func (a *Abstract) SetFormat(format string) {
	a.format = format
}

func (a *Abstract) SetConfig(config *sdksData.ScimContainer) {
	a.config = config
}

func (a *Abstract) MiddlewareEngine() (middleware.Interface, error) {
	a.middlewareFactoryOnce.Do(func() {
		a.middlewareFactory = new(middleware.Factory)
	})

	if engine, err := a.middlewareFactory.Engine(a.config); err != nil {
		return nil, err
	} else {
		engine.SetFormat(a.format)
		engine.SetConfig(a.config)
		return engine, nil
	}
}
