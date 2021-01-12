package middleware

import (
	"errors"
	"gitlab.com/24sessions/sdk-go-configurator/data"
)

type Factory struct {
}

func (f *Factory) Engine(scimConfig *sdksData.ScimContainer) (e Interface, err error) {
	switch scimConfig.Middleware() {
	case "", "default":
		return new(Default), nil
	case "rabobank":
		return new(Rabobank), nil
	default:
		return nil, errors.New("Unknown middleware: " + scimConfig.Middleware())
	}
}
