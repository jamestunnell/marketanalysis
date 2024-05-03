package blocks

import (
	"reflect"
	"slices"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type StringEnum struct {
	Value, Default string
	Enum           []string
}

func (p *StringEnum) GetDefault() any {
	return p.Default
}

func (p *StringEnum) GetSchema() map[string]any {
	return map[string]any{
		"type": "string",
		"enum": p.Enum,
	}
}

func (p *StringEnum) GetVal() any {
	return p.Value
}

func (p *StringEnum) SetVal(val any) error {
	strVal, ok := val.(string)
	if !ok {
		actual := reflect.TypeOf(val).String()

		return commonerrs.NewErrWrongType(actual, "string")
	}

	if !slices.Contains(p.Enum, strVal) {
		return commonerrs.NewErrNotOneOf[string]("value", strVal, p.Enum)
	}

	p.Value = strVal

	return nil
}
