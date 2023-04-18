package constraints

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"golang.org/x/exp/constraints"
)

const TypeOneOf = "OneOf"

type TypedOneOf[T constraints.Ordered] struct {
	Allowed []T
}

func NewOneOf[T constraints.Ordered](first T, more ...T) models.Constraint {
	allowed := append([]T{first}, more...)

	return &TypedOneOf[T]{Allowed: allowed}
}

func (m *TypedOneOf[T]) Type() string {
	return TypeOneOf
}

func (m *TypedOneOf[T]) Check(val any) error {
	tVal, ok := val.(T)
	if !ok {
		actual := reflect.TypeOf(val).String()
		expected := reflect.TypeOf(m.Allowed[0]).String()

		return commonerrs.NewErrWrongType(actual, expected)
	}

	for _, allowedVal := range m.Allowed {
		if tVal == allowedVal {
			return nil
		}
	}

	return commonerrs.NewErrNotOneOf("value", tVal, m.Allowed)
}

func (m *TypedOneOf[T]) Bounds() []any {
	return sliceutils.Map(m.Allowed, func(val T) any {
		return val
	})
}
