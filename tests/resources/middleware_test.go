package resources_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"math"
	"service-scim/services/middleware"
	"service-scim/tests"
	"testing"
)

var factory *middleware.Factory

func init() {
	tests.Init()

	factory = new(middleware.Factory)
}

func TestDefaultMappingForEntity(t *testing.T) {
	testMiddlewareForUser(t, nil, nil, "default",
		"defaultMapping.json", "defaultExpected.json",
	)
}

func TestDefaultConcatenatingForEntity(t *testing.T) {
	testMiddlewareForUser(t, nil, nil, "default",
		"defaultConcatenating.json", "defaultExpected.json",
	)
}

func TestRabobankDefaultsForUser(t *testing.T) {
	testMiddlewareForUser(t, nil, nil, "rabobank",
		"rabobankUserDefaults.json", "rabobankUserExpected.json",
	)
}

func TestRabobank(t *testing.T) {
	testMiddlewareForUser(t, nil, nil, "rabobank",
		"rabobankUserDefaults.json", "rabobankUserExpected.json",
	)
}

func TestRabobankDefaultsForGroup(t *testing.T) {
	groupId := "123"
	testMiddlewareForGroup(t, tests.NewConfig(), &groupId, "rabobank",
		"rabobankGroupDefaults.json", "rabobankGroupExpected.json",
	)
}

func readProcessingMap(fileName string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if _, err := tests.ReadFromFixture(&data, "fixtures/middleware/"+fileName); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func middlewareEngine(name string) (middleware.Interface, error) {
	config := new(sdksData.InstanceContainer)
	config.ScimMiddleware = name

	if engine, err := factory.Engine(sdksData.NewScim(config)); err != nil {
		return nil, err
	} else {
		engine.SetFormat("json")
		return engine, nil
	}
}

func testMiddlewareForEntity(t *testing.T, config *sdksData.ScimContainer, engineType, inputFixture, expectedFixture string, process func(engine middleware.Interface, input map[string]interface{}) (map[string]interface{}, error)) {
	if input, err := readProcessingMap(inputFixture); err != nil {
		t.Errorf(err.Error())
	} else if expected, err := readProcessingMap(expectedFixture); err != nil {
		t.Errorf(err.Error())
	} else if engine, err := middlewareEngine(engineType); err != nil {
		t.Errorf("Cannot get middleware for its `%s` type: %s", engineType, err.Error())
	} else {
		engine.SetConfig(config)
		if input, err := process(engine, input); err != nil {
			t.Errorf(err.Error())
		} else {
			assert.Equal(t, expected, dataIntersection(input, expected).(map[string]interface{}))
		}
	}
}

func testMiddlewareForUser(t *testing.T, config *sdksData.ScimContainer, id *string, engineType, inputFixture, expectedFixture string) {
	testMiddlewareForEntity(t, config, engineType, inputFixture, expectedFixture, func(engine middleware.Interface, input map[string]interface{}) (output map[string]interface{}, err error) {
		if output, err = engine.ProcessUser(input, id); err != nil {
			err = fmt.Errorf("Cannot apply %s middleware for user: %s", engineType, err.Error())
		}
		return
	})
}

func testMiddlewareForGroup(t *testing.T, config *sdksData.ScimContainer, id *string, engineType, inputFixture, expectedFixture string) {
	testMiddlewareForEntity(t, config, engineType, inputFixture, expectedFixture, func(engine middleware.Interface, input map[string]interface{}) (output map[string]interface{}, err error) {
		if output, err = engine.ProcessGroup(input, id); err != nil {
			err = fmt.Errorf("Cannot apply %s middleware for group: %s", engineType, err.Error())
		}
		return
	})
}

// Intersection of two maps by their keys with values of the first map
func dataIntersection(d1, d2 interface{}) interface{} {
	switch d1.(type) {
	case []interface{}:
		switch d2.(type) {
		case []interface{}:
			newD1 := d1.([]interface{})
			newD2 := d2.([]interface{})
			l := int(math.Min(float64(len(newD1)), float64(len(newD2))))

			s := make([]interface{}, 0)
			for i := 0; i < l; i++ {
				s = append(s, dataIntersection(newD1[i], newD2[i]))
			}
			return s
		default:
			return nil
		}
	case []map[string]interface{}:
		switch d2.(type) {
		case []interface{}:
			newD1 := d1.([]map[string]interface{})
			newD2 := d2.([]interface{})
			l := int(math.Min(float64(len(newD1)), float64(len(newD2))))

			s := make([]interface{}, 0)
			for i := 0; i < l; i++ {
				s = append(s, dataIntersection(newD1[i], newD2[i]))
			}
			return s
		default:
			return nil
		}
	case map[string]interface{}:
		switch d2.(type) {
		case map[string]interface{}:
			newD1 := d1.(map[string]interface{})
			newD2 := d2.(map[string]interface{})

			m := make(map[string]interface{})
			for k1, v1 := range newD1 {
				if v2, ok := newD2[k1]; ok {
					m[k1] = dataIntersection(v1, v2)
				}
			}
			return m
		default:
			return nil
		}
	default:
		switch d2.(type) {
		case []interface{}:
			return nil
		case map[string]interface{}:
			return nil
		default:
			return d1
		}
	}
}
