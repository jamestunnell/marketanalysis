package blocks

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	goconstraints "golang.org/x/exp/constraints"
)

type Param interface {
	GetValueType() string
	GetDefaultVal() any
	GetConstraints() []Constraint
	GetCurrentVal() any
	SetCurrentVal(any) error
}

type TypedParam[T goconstraints.Ordered] struct {
	CurrentVal, DefaultVal T
	Constraints            []*TypedConstraint[T]
	ValueType              string
}

func NewTypedParam[T goconstraints.Ordered](
	defaultVal T,
	constraints ...*TypedConstraint[T],
) *TypedParam[T] {
	var zeroVal T

	return &TypedParam[T]{
		ValueType:   reflect.TypeOf(zeroVal).String(),
		CurrentVal:  zeroVal,
		DefaultVal:  defaultVal,
		Constraints: constraints,
	}
}

func (p *TypedParam[T]) GetDefaultVal() any {
	return p.DefaultVal
}

func (p *TypedParam[T]) GetConstraints() []Constraint {
	return sliceutils.Map(p.Constraints, func(c *TypedConstraint[T]) Constraint { return c })
}

func (p *TypedParam[T]) GetValueType() string {
	return p.ValueType
}

func (p *TypedParam[T]) GetCurrentVal() any {
	return p.CurrentVal
}

func (p *TypedParam[T]) SetCurrentVal(val any) error {
	t, ok := val.(T)
	if !ok {
		actual := reflect.TypeOf(val).String()
		expected := reflect.TypeOf(t).String()

		return commonerrs.NewErrWrongType(actual, expected)
	}

	for _, c := range p.Constraints {
		if err := c.CheckVal(t); err != nil {
			return err
		}
	}

	p.CurrentVal = t

	return nil
}
