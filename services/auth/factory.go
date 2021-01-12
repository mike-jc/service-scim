package auth

import (
	"errors"
	"gitlab.com/24sessions/sdk-go-configurator/data"
)

type Factory struct {
}

func (f *Factory) Engine(scimConfig *sdksData.ScimContainer) (e Interface, err error) {
	switch scimConfig.AuthType() {
	case "none":
		return nil, nil
	case "basic":
		engine := new(HttpBasic)
		engine.SetCredentials(scimConfig.AuthBasicUser(), scimConfig.AuthBasicPassword())
		return engine, nil
	case "token":
		engine := new(BearerToken)
		engine.SetToken(scimConfig.AuthToken())
		return engine, nil
	default:
		return nil, errors.New("Unknown auth type: " + scimConfig.AuthType())
	}
}
