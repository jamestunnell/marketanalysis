package constraints

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"golang.org/x/exp/constraints"
)

const TypeValRange = "ValRange"

type TypedValRange[T constraints.Ordered] struct {
	Min, Max T
}

func NewValRange[T constraints.Ordered](min, max T) models.Constraint {
	return &TypedValRange[T]{Min: min, Max: max}
}

func (m *TypedValRange[T]) Type() string {
	return TypeValRange
}

func (m *TypedValRange[T]) Check(val any) error {
	tVal, ok := val.(T)
	if !ok {
		actual := reflect.TypeOf(val).String()
		expected := reflect.TypeOf(m.Min).String()

		return commonerrs.NewErrWrongType(actual, expected)
	}

	if tVal < m.Min {
		return commonerrs.NewErrLessThanMin("value", tVal, m.Min)
	}

	if tVal > m.Max {
		return commonerrs.NewErrMoreThanMax("value", tVal, m.Max)
	}

	return nil
}

func (m *TypedValRange[T]) ValueBounds() []any {
	return []any{m.Min, m.Max}
}
