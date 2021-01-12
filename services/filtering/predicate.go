package filtering

import (
	"reflect"
	"service-scim/models/config"
	"service-scim/services/navigation"
	"strconv"
	"strings"
)

func newPredicate(filter navigation.FilterNode, attr *modelsConfig.Attribute) PredicateInterface {
	return &Predicate{filter, attr}
}

type PredicateInterface interface {
	evaluate(Complex) bool
}

type Predicate struct {
	filter navigation.FilterNode
	attr   *modelsConfig.Attribute
}

func (p *Predicate) evaluate(c Complex) bool {
	return p.getFunc(p.filter)(c)
}

type predicateFunc func(Complex) bool

func (p *Predicate) getFunc(filter navigation.FilterNode) predicateFunc {
	switch filter.Data() {
	case navigation.And:
		return p.andFunc(filter)
	case navigation.Or:
		return p.orFunc(filter)
	case navigation.Not:
		return p.notFunc(filter)
	case navigation.Eq:
		return p.eqFunc(filter)
	case navigation.Ne:
		return p.neFunc(filter)
	case navigation.Sw:
		return p.swFunc(filter)
	case navigation.Ew:
		return p.ewFunc(filter)
	case navigation.Co:
		return p.coFunc(filter)
	case navigation.Pr:
		return p.prFunc(filter)
	case navigation.Gt:
		return p.gtFunc(filter)
	case navigation.Ge:
		return p.geFunc(filter)
	case navigation.Lt:
		return p.ltFunc(filter)
	case navigation.Le:
		return p.leFunc(filter)
	}
	return nil
}

func (p *Predicate) andFunc(filter navigation.FilterNode) predicateFunc {
	lhs := p.getFunc(filter.Left())
	rhs := p.getFunc(filter.Right())
	return func(c Complex) bool {
		return lhs(c) && rhs(c)
	}
}

func (p *Predicate) orFunc(filter navigation.FilterNode) predicateFunc {
	lhs := p.getFunc(filter.Left())
	rhs := p.getFunc(filter.Right())
	return func(c Complex) bool {
		return lhs(c) || rhs(c)
	}
}

func (p *Predicate) notFunc(filter navigation.FilterNode) predicateFunc {
	lhs := p.getFunc(filter.Left())
	return func(c Complex) bool {
		return !lhs(c)
	}
}

func (p *Predicate) eqFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.compare(filter.Left(), filter.Right(), c) == equal
	}
}

func (p *Predicate) neFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.compare(filter.Left(), filter.Right(), c) != equal
	}
}

func (p *Predicate) gtFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.compare(filter.Left(), filter.Right(), c) == greater
	}
}

func (p *Predicate) geFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		r := p.compare(filter.Left(), filter.Right(), c)
		return r == greater || r == equal
	}
}

func (p *Predicate) ltFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.compare(filter.Left(), filter.Right(), c) == less
	}
}

func (p *Predicate) leFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		r := p.compare(filter.Left(), filter.Right(), c)
		return r == less || r == equal
	}
}

func (p *Predicate) prFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		if filter.Left().Type() != navigation.PathOperand {
			return false
		}

		key := filter.Left().Data().(navigation.Path)
		lVal := reflect.ValueOf(<-c.Get(key, p.attr))
		if !lVal.IsValid() {
			return false
		}

		if lVal.Kind() == reflect.Interface {
			lVal = lVal.Elem()
		}

		switch lVal.Kind() {
		case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
			return lVal.Len() > 0
		default:
			return true
		}
	}
}

func (p *Predicate) swFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.stringOp(filter.Left(), filter.Right(), c, func(a, b string) bool {
			return strings.HasPrefix(a, b)
		})
	}
}

func (p *Predicate) ewFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.stringOp(filter.Left(), filter.Right(), c, func(a, b string) bool {
			return strings.HasSuffix(a, b)
		})
	}
}

func (p *Predicate) coFunc(filter navigation.FilterNode) predicateFunc {
	return func(c Complex) bool {
		return p.stringOp(filter.Left(), filter.Right(), c, func(a, b string) bool {
			return strings.Contains(a, b)
		})
	}
}

func (p *Predicate) stringOp(lhs, rhs navigation.FilterNode, c Complex, op func(a, b string) bool) bool {
	if lhs.Type() != navigation.PathOperand || rhs.Type() != navigation.ConstantOperand {
		return false
	}
	if p.attr == nil {
		return false
	}

	key := lhs.Data().(navigation.Path)
	attr := p.attr.GetAttribute(key, false)
	if attr == nil || attr.MultiValued || attr.Type == "complex" {
		return false
	}

	lVal := reflect.ValueOf(<-c.Get(key, p.attr))
	if !lVal.IsValid() {
		return false
	} else if lVal.Kind() == reflect.Interface {
		lVal = lVal.Elem()
	}

	rVal := reflect.ValueOf(rhs.Data())
	if !rVal.IsValid() {
		return false
	} else if rVal.Kind() == reflect.Interface {
		lVal = lVal.Elem()
	}

	if !p.kindOf(lVal, reflect.String) || !p.kindOf(rVal, reflect.String) {
		return false
	} else {
		if *attr.CaseExact {
			return op(lVal.String(), rVal.String())
		} else {
			return op(strings.ToLower(lVal.String()), strings.ToLower(rVal.String()))
		}
	}
}

func (p *Predicate) compare(lhs, rhs navigation.FilterNode, c Complex) comparison {
	if lhs.Type() != navigation.PathOperand || rhs.Type() != navigation.ConstantOperand {
		return invalid
	}

	key := lhs.Data().(navigation.Path)
	attr := p.attr.GetAttribute(key, true)
	if attr == nil || attr.MultiValued || attr.Type == "complex" {
		return invalid
	}

	lVal := reflect.ValueOf(<-c.Get(key, p.attr))
	if !lVal.IsValid() {
		return invalid
	} else if lVal.Kind() == reflect.Interface {
		lVal = lVal.Elem()
	}

	rVal := reflect.ValueOf(rhs.Data())
	if !rVal.IsValid() {
		return invalid
	} else if rVal.Kind() == reflect.Interface {
		lVal = lVal.Elem()
	}

	switch attr.Type {
	case "integer":
		if !p.kindOf(lVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) ||
			!p.kindOf(rVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
			return invalid
		} else {
			switch {
			case lVal.Int() == rVal.Int():
				return equal
			case lVal.Int() < rVal.Int():
				return less
			case lVal.Int() > rVal.Int():
				return greater
			}
		}

	case "decimal":
		if !p.kindOf(lVal, reflect.Float32, reflect.Float64) ||
			!p.kindOf(rVal, reflect.Float32, reflect.Float64) {
			return invalid
		} else {
			switch {
			case lVal.Float() == rVal.Float():
				return equal
			case lVal.Float() < rVal.Float():
				return less
			case lVal.Float() > rVal.Float():
				return greater
			}
		}

	case "boolean":
		if !p.kindOf(lVal, reflect.Bool) ||
			!p.kindOf(rVal, reflect.Bool) {
			return invalid
		} else {
			if lVal.Bool() == rVal.Bool() {
				return equal
			} else {
				return invalid
			}
		}

	case "string", "binary", "datetime", "reference":
		// int can be convert to string without data lost
		if p.kindOf(lVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
			lVal = reflect.ValueOf(strconv.FormatInt(lVal.Int(), 10))
		}
		if p.kindOf(rVal, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64) {
			rVal = reflect.ValueOf(strconv.FormatInt(rVal.Int(), 10))
		}

		if !p.kindOf(lVal, reflect.String) ||
			!p.kindOf(rVal, reflect.String) {
			return invalid
		} else {
			var a, b string
			if *attr.CaseExact {
				a, b = lVal.String(), rVal.String()
			} else {
				a, b = strings.ToLower(lVal.String()), strings.ToLower(rVal.String())
			}
			switch {
			case a == b:
				return equal
			case a < b:
				return less
			case a > b:
				return greater
			}
		}
	}

	return invalid
}

func (p *Predicate) kindOf(v reflect.Value, kinds ...reflect.Kind) bool {
	for _, kind := range kinds {
		if kind == v.Kind() {
			return true
		}
	}
	return false
}

type comparison int

const (
	invalid = comparison(-2)
	less    = comparison(-1)
	equal   = comparison(0)
	greater = comparison(1)
)
