package filtering

import (
	"service-scim/errors"
	"service-scim/models/config"
	"service-scim/services/navigation"
)

type Complex map[string]interface{}

func (c Complex) Get(p navigation.Path, attr *modelsConfig.Attribute) chan interface{} {
	output := make(chan interface{})
	go func() {
		c.get(p, attr, output)
		close(output)
	}()
	return output
}

func (c Complex) get(p navigation.Path, attr *modelsConfig.Attribute, output chan interface{}) {
	subAttr := attr.GetAttribute(p, false)
	if subAttr == nil {
		return
	}

	if v, ok := c[subAttr.Name]; ok && v != nil {
		if p.FilterRoot() != nil {
			if mv, ok := v.([]interface{}); ok && mv != nil {
				matches := MultiValued(mv).Filter(p.FilterRoot(), subAttr)
				for match := range matches {
					if p.Next() != nil {
						if v0, ok := match.(map[string]interface{}); ok && v0 != nil {
							Complex(v0).get(p.Next(), subAttr, output)
						}
					} else {
						output <- match
					}
				}
			}
		} else {
			if p.Next() != nil {
				if v0, ok := v.(map[string]interface{}); ok && v0 != nil {
					Complex(v0).get(p.Next(), subAttr, output)
				}
			} else {
				output <- v
			}
		}
	}
}

// Evaluate given predicate
func (c Complex) Evaluate(filter navigation.FilterNode, attr *modelsConfig.Attribute) bool {
	return newPredicate(filter, attr).evaluate(c)
}

// Set the value at the specified Path.
// Path is a dot separated string that may contain filter
func (c Complex) Set(p navigation.Path, value interface{}, attr *modelsConfig.Attribute) error {
	subAttr := attr.GetAttribute(p, true)
	if subAttr == nil {
		return scimErrors.InvalidPathError(p.CollectValue(), "no attribute found")
	}

	// TODO validate

	base, last := p.SeparateAtLast()
	itemsToSet := make(chan interface{})
	go func() {
		if base != nil {
			c.get(base, attr, itemsToSet)
		} else {
			itemsToSet <- c
		}
		close(itemsToSet)
	}()

	for item := range itemsToSet {
		if m, ok := item.(Complex); ok && m != nil {
			m.set(last, value, subAttr)
		} else if m, ok := item.(map[string]interface{}); ok && m != nil {
			Complex(m).set(last, value, subAttr)
		}
	}
	return nil
}

func (c Complex) set(p navigation.Path, value interface{}, attr *modelsConfig.Attribute) {
	if attr.MultiValued && p.FilterRoot() != nil {
		if mv, ok := c[attr.Name].([]interface{}); ok {
			for i, v := range mv {
				if c0, ok := v.(map[string]interface{}); ok {
					if Complex(c0).Evaluate(p.FilterRoot(), attr) {
						MultiValued(mv).Set(i, value)
					}
				}
			}
		}
	} else {
		c[attr.Name] = value
	}
}
