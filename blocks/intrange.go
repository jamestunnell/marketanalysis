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

func (p *IntRange) GetSchema() map[string]any {
	return map[string]any{
		"type":    "integer",
		"minimum": p.Min,
		"maximum": p.Max,
	}
}

func (p *IntRange) GetVal() any {
	return p.Value
}

func (p *IntRange) SetVal(val any) error {
	intVal, ok := val.(int)
	if !ok {
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
