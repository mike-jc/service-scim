package serviceSp

import (
	"fmt"
	"github.com/astaxie/beego"
	logger "gitlab.com/24sessions/lib-go-logger/logger/services"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"service-scim/models/config"
	"service-scim/services"
	"service-scim/system"
)

const SpConfigCreated = "2019-04-04T17:43:00Z"
const SpConfigLastModified = "2019-04-04T17:43:00Z"

type ServiceConfig struct {
	scimConfig *sdksData.ScimContainer
	baseUrl    string
}

func (c *ServiceConfig) SetScimConfig(scimConfig *sdksData.ScimContainer) {
	c.scimConfig = scimConfig
}

func (c *ServiceConfig) SetBaseUrl(baseUrl string) {
	c.baseUrl = baseUrl
}

func (c *ServiceConfig) Config() (config modelsConfig.SpConfig, err error) {
	var authSchemes []*modelsConfig.AuthenticationScheme
	if authSchemes, err = c.AuthenticationSchemesConfig(); err != nil {
		return
	}

	config = modelsConfig.SpConfig{
		Schemas:               []string{"urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"},
		Id:                    "ServiceProviderConfig",
		Patch:                 c.PatchOperationConfig(),
		Bulk:                  c.BulkOperationConfig(),
		Filter:                c.FilterOperationConfig(),
		ChangePassword:        c.ChangePasswordOperationConfig(),
		Sort:                  c.SortOperationConfig(),
		ETag:                  c.ETagConfig(),
		AuthenticationSchemes: authSchemes,
		Meta:                  c.MetaConfig(),
	}
	return
}

func (c *ServiceConfig) PatchOperationConfig() *modelsConfig.PatchOperation {
	enabled, _ := beego.AppConfig.Bool("sp.patch.enabled")

	return &modelsConfig.PatchOperation{
		Supported: enabled,
	}
}

func (c *ServiceConfig) BulkOperationConfig() *modelsConfig.BulkOperation {
	var maxOperations int
	var maxPayloadSize int64

	enabled, _ := beego.AppConfig.Bool("sp.bulk.enabled")
	if enabled {
		maxOperations, _ = beego.AppConfig.Int("sp.bulk.maxOperations")
		if maxOperations < 0 {
			maxOperations = 0
		}
		maxPayloadSize, _ = beego.AppConfig.Int64("sp.bulk.maxPayloadSize")
		if maxPayloadSize < 0 {
			maxPayloadSize = 0
		}
		if maxOperations == 0 || maxPayloadSize == 0 {
			enabled = false
		}
	}

	return &modelsConfig.BulkOperation{
		Supported:      enabled,
		MaxOperations:  maxOperations,
		MaxPayloadSize: maxPayloadSize,
	}
}

func (c *ServiceConfig) FilterOperationConfig() *modelsConfig.FilterOperation {
	var maxResults int

	enabled, _ := beego.AppConfig.Bool("sp.filter.enabled")
	if enabled {
		maxResults, _ = beego.AppConfig.Int("sp.filter.maxResults")
		if maxResults < 0 {
			maxResults = 0
		}
		if maxResults == 0 {
			enabled = false
		}
	}

	return &modelsConfig.FilterOperation{
		Supported:  enabled,
		MaxResults: maxResults,
	}
}

func (c *ServiceConfig) ChangePasswordOperationConfig() *modelsConfig.ChangePasswordOperation {
	enabled, _ := beego.AppConfig.Bool("sp.changePassword.enabled")

	return &modelsConfig.ChangePasswordOperation{
		Supported: enabled,
	}
}

func (c *ServiceConfig) SortOperationConfig() *modelsConfig.SortOperation {
	enabled, _ := beego.AppConfig.Bool("sp.sort.enabled")

	return &modelsConfig.SortOperation{
		Supported: enabled,
	}
}

func (c *ServiceConfig) ETagConfig() *modelsConfig.ETag {
	enabled, _ := beego.AppConfig.Bool("sp.etag.enabled")

	return &modelsConfig.ETag{
		Supported: enabled,
	}
}

func (c *ServiceConfig) AuthenticationSchemesConfig() ([]*modelsConfig.AuthenticationScheme, error) {
	schemes := make([]*modelsConfig.AuthenticationScheme, 0)
	authType := c.scimConfig.AuthType()

	if authType == "none" {
		return nil, nil
	} else if authType == "basic" {
		schemes = append(schemes, &modelsConfig.AuthenticationScheme{
			Name:        "HTTP Basic",
			Description: "Authentication scheme using the HTTP Basic Standard",
			SpecUri:     "http://www.rfc-editor.org/info/rfc2617",
			Type:        "httpbasic",
			Primary:     true,
		})
	} else if authType == "token" {
		schemes = append(schemes, &modelsConfig.AuthenticationScheme{
			Name:        "App Bearer Token",
			Description: "Authentication scheme using the long-term Bearer Token",
			Type:        "bearertoken",
			Primary:     true,
		})
	} else if authType == "certificate" {
		schemes = append(schemes, &modelsConfig.AuthenticationScheme{
			Name:        "Certificate-based Authorization",
			Description: "Authentication scheme using the SSL certificates",
			Type:        "certificate",
			Primary:     true,
		})
	} else {
		err := fmt.Errorf("Invalid authentication type")
		services.LogMain.Log(logger.CreateError(err.Error()).SetCode("app.invalid_auth_type"))
		return nil, err
	}

	return schemes, nil
}

func (c *ServiceConfig) MetaConfig() *modelsConfig.SpConfigMeta {
	return &modelsConfig.SpConfigMeta{
		ResourceType: "ServiceProviderConfig",
		Location:     c.baseUrl + "/" + system.AppVersion() + "/ServiceProviderConfig",
		Created:      SpConfigCreated,
		LastModified: SpConfigLastModified,
		Version:      system.AppVersion(),
	}
}
