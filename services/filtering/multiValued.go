package filtering

import (
	"service-scim/models/config"
	"service-scim/services/navigation"
)

type MultiValued []interface{}

func (c MultiValued) Get(index int) interface{} {
	return c[index]
}

func (c MultiValued) Set(index int, value interface{}) {
	c[index] = value
}

func (c MultiValued) Len() int {
	return len(c)
}

func (c MultiValued) Add(value ...interface{}) MultiValued {
	return MultiValued(append([]interface{}(c), value...))
}

func (c MultiValued) Remove(index int) MultiValued {
	// TODO
	return nil
}

func (c MultiValued) Filter(root navigation.FilterNode, attr *modelsConfig.Attribute) chan interface{} {
	output := make(chan interface{})
	go func() {
		for _, elem := range c {
			if m, ok := elem.(map[string]interface{}); ok {
				if Complex(m).Evaluate(root, attr) {
					output <- elem
				}
			}
		}
		close(output)
	}()
	return output
}
