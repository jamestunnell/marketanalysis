package constraints

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"golang.org/x/exp/constraints"
)

const TypeRange = "Range"

type TypedRange[T constraints.Ordered] struct {
	Min, Max T
}

func NewRange[T constraints.Ordered](min, max T) models.Constraint {
	return &TypedRange[T]{Min: min, Max: max}
}

func (m *TypedRange[T]) Type() string {
	return TypeRange
}

func (m *TypedRange[T]) Check(val any) error {
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

func (m *TypedRange[T]) Bounds() []any {
	return []any{m.Min, m.Max}
}
