package constraints

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"golang.org/x/exp/constraints"
)

const TypeValOneOf = "ValOneOf"

type TypedValOneOf[T constraints.Ordered] struct {
	Allowed []T
}

func NewValOneOf[T constraints.Ordered](first T, more ...T) models.Constraint {
	allowed := append([]T{first}, more...)

	return &TypedValOneOf[T]{Allowed: allowed}
}

func (m *TypedValOneOf[T]) Type() string {
	return TypeValOneOf
}

func (m *TypedValOneOf[T]) Check(val any) error {
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

func (m *TypedValOneOf[T]) ValueBounds() []any {
	return sliceutils.Map(m.Allowed, func(val T) any {
		return val
	})
}
