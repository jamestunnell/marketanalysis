package blocks

import (
	"reflect"
	"slices"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type IntEnum struct {
	Value, Default int
	Enum           []int
}

func (p *IntEnum) GetDefault() any {
	return p.Default
}

func (p *IntEnum) GetSchema() map[string]any {
	return map[string]any{
		"type": "number",
		"enum": p.Enum,
	}
}

func (p *IntEnum) GetType() string {
	return "int"
}

func (p *IntEnum) GetVal() any {
	return p.Value
}

func (p *IntEnum) SetVal(val any) error {
	intVal, ok := val.(int)
	if !ok {
		actual := reflect.TypeOf(val).String()

		return commonerrs.NewErrWrongType(actual, "int")
	}

	if !slices.Contains(p.Enum, intVal) {
		return commonerrs.NewErrNotOneOf[int]("value", intVal, p.Enum)
	}

	p.Value = intVal

	return nil
}
