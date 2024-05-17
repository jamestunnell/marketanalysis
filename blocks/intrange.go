package blocks

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type IntRange struct {
	Value, Min, Max, Default int
}

func (p *IntRange) GetDefault() any {
	return p.Default
}

func (p *IntRange) GetLimits() []any {
	return []any{p.Min, p.Max}
}

func (p *IntRange) GetType() string {
	return "IntRange"
}

func (p *IntRange) GetVal() any {
	return p.Value
}

func (p *IntRange) SetVal(val any) error {
	var intVal int

	switch vv := val.(type) {
	case float64:
		intVal = int(vv)

		if float64(intVal) != vv {
			return commonerrs.NewErrWrongType("float64", "int")
		}
	case int:
		intVal = vv
	default:
		actual := reflect.TypeOf(val).String()

		return commonerrs.NewErrWrongType(actual, "int")
	}

	if intVal < p.Min {
		return commonerrs.NewErrLessThanMin("value", intVal, p.Min)
	}

	if intVal > p.Max {
		return commonerrs.NewErrMoreThanMax("value", intVal, p.Max)
	}

	p.Value = intVal

	return nil
}
