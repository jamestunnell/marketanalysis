package blocks

import (
	"reflect"
	"slices"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type IntEnum struct {
	Value, Default int
	Enum           []int
}

func (p *IntEnum) GetDefault() any {
	return p.Default
}

func (p *IntEnum) GetLimits() []any {
	return sliceutils.Map(
		p.Enum, func(s int) any { return s })
}

func (p *IntEnum) GetType() string {
	return "IntEnum"
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
