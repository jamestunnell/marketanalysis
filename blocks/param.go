package blocks

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	goconstraints "golang.org/x/exp/constraints"
)

type Param interface {
	GetValueType() string
	GetDefaultVal() any
	GetConstraint() Constraint
	GetCurrentVal() any
	SetCurrentVal(any) error
}

type TypedParam[T goconstraints.Ordered] struct {
	CurrentVal, DefaultVal T
	Constraint             *TypedConstraint[T]
	ValueType              string
}

func NewTypedParam[T goconstraints.Ordered](
	defaultVal T,
	constr *TypedConstraint[T],
) *TypedParam[T] {
	var zeroVal T

	return &TypedParam[T]{
		ValueType:  reflect.TypeOf(zeroVal).String(),
		CurrentVal: zeroVal,
		DefaultVal: defaultVal,
		Constraint: constr,
	}
}

func (p *TypedParam[T]) GetDefaultVal() any {
	return p.DefaultVal
}

func (p *TypedParam[T]) GetConstraint() Constraint {
	return p.Constraint
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

	if err := p.Constraint.CheckVal(t); err != nil {
		return err
	}

	p.CurrentVal = t

	return nil
}
