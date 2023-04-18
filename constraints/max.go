package constraints

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"golang.org/x/exp/constraints"
)

const TypeMax = "Max"

type TypedMax[T constraints.Ordered] struct {
	Max T
}

func NewMax[T constraints.Ordered](max T) models.Constraint {
	return &TypedMax[T]{Max: max}
}

func (m *TypedMax[T]) Type() string {
	return TypeMax
}

func (m *TypedMax[T]) Check(val any) error {
	tVal, ok := val.(T)
	if !ok {
		actual := reflect.TypeOf(val).String()
		expected := reflect.TypeOf(m.Max).String()

		return commonerrs.NewErrWrongType(actual, expected)
	}

	if tVal > m.Max {
		return commonerrs.NewErrMoreThanMax("value", tVal, m.Max)
	}

	return nil
}

func (m *TypedMax[T]) Bounds() []any {
	return []any{m.Max}
}
