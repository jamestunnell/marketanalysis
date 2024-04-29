package blocks

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type FltRange struct {
	Value, Min, Max, Default float64
}

func (p *FltRange) GetDefault() any {
	return p.Default
}

func (p *FltRange) GetSchema() map[string]any {
	return map[string]any{
		"type":    "number",
		"minimum": p.Min,
		"maximum": p.Max,
	}
}

func (p *FltRange) GetVal() any {
	return p.Value
}

func (p *FltRange) SetVal(val any) error {
	fltVal, ok := val.(float64)
	if !ok {
		actual := reflect.TypeOf(val).String()

		return commonerrs.NewErrWrongType(actual, "float64")
	}

	if fltVal < p.Min {
		return commonerrs.NewErrLessThanMin("value", fltVal, p.Min)
	}

	if fltVal > p.Max {
		return commonerrs.NewErrMoreThanMax("value", fltVal, p.Max)
	}

	p.Value = fltVal

	return nil
}
