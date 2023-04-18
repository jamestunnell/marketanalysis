package constraints

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"golang.org/x/exp/constraints"
)

const TypeMin = "Min"

type TypedMin[T constraints.Ordered] struct {
	Min T
}

func NewMin[T constraints.Ordered](min T) models.Constraint {
	return &TypedMin[T]{Min: min}
}

func (m *TypedMin[T]) Type() string {
	return TypeMin
}

func (m *TypedMin[T]) Check(val any) error {
	tVal, ok := val.(T)
	if !ok {
		actual := reflect.TypeOf(val).String()
		expected := reflect.TypeOf(m.Min).String()

		return commonerrs.NewErrWrongType(actual, expected)
	}

	if tVal < m.Min {
		return commonerrs.NewErrLessThanMin("value", tVal, m.Min)
	}

	return nil
}

func (m *TypedMin[T]) Bounds() []any {
	return []any{m.Min}
}
