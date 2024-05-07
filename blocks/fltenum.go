package blocks

import (
	"reflect"
	"slices"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type FltEnum struct {
	Value, Default float64
	Enum           []float64
}

func (p *FltEnum) GetDefault() any {
	return p.Default
}

func (p *FltEnum) GetSchema() map[string]any {
	return map[string]any{
		"type": "number",
		"enum": p.Enum,
	}
}

func (p *FltEnum) GetType() string {
	return "float64"
}

func (p *FltEnum) GetVal() any {
	return p.Value
}

func (p *FltEnum) SetVal(val any) error {
	fltVal, ok := val.(float64)
	if !ok {
		actual := reflect.TypeOf(val).String()

		return commonerrs.NewErrWrongType(actual, "float64")
	}

	if !slices.Contains(p.Enum, fltVal) {
		return commonerrs.NewErrNotOneOf[float64]("value", fltVal, p.Enum)
	}

	p.Value = fltVal

	return nil
}