package services

import (
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/sdk-go-configurator"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"strings"
	"sync"
)

const MainDomain = "24sessions.com"

var testConfig *sdksData.InstanceContainer
var configurationsLock sync.Mutex
var configurator *sdks.Configurator
var configuratorOnce sync.Once

func prodConfiguration(domain string) (config *sdksData.InstanceContainer, err error) {
	configuratorOnce.Do(func() {
		configurator = new(sdks.Configurator)
	})
	return configurator.GetInstance(domain)
}

func testConfiguration() *sdksData.InstanceContainer {

	if testConfig == nil {
		testConfig = &sdksData.InstanceContainer{
			CompanyLocale:      "en",
			CompanyName:        "24sessions Test",
			CompanyTimezone:    "UTC",
			CompanyStatus:      "trial",
			ScimEnabled:        []byte("true"),
			ScimAuthType:       "basic",
			ScimResponseFormat: "json",
		}
	}

	return testConfig
}

func NewConfig(domain string) (config *sdksData.ScimContainer, err error) {
	configurationsLock.Lock()
	defer configurationsLock.Unlock()

	if beego.BConfig.RunMode == "test" {
		config = sdksData.NewScim(testConfiguration())
	} else {
		if pConfig, cErr := prodConfiguration(domain); cErr != nil {
			return nil, cErr
		} else {
			config = sdksData.NewScim(pConfig)
		}
	}
	return
}

func InstanceDomainFromScimDomain(domain string) string {
	return strings.Replace(domain, "scim.", "", 1)
}
