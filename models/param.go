package models

import (
	"reflect"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	goconstraints "golang.org/x/exp/constraints"
)

type Param interface {
	GetValueType() string
	GetDefaultVal() any
	GetConstraintInfo() *ConstraintInfo
	GetCurrentVal() any
	SetCurrentVal(any) error
}

type ConstraintInfo struct {
	Type   string `json:"type"`
	Limits []any  `json:"limits"`
}

type typedParam[T goconstraints.Ordered] struct {
	CurrentVal, DefaultVal T
	Constraint             Constraint[T]
	ValueType              string
}

type IntParam struct {
	*typedParam[int]
}

type FloatParam struct {
	*typedParam[float64]
}

func newTypedParam[T goconstraints.Ordered](
	defaultVal T,
	constr Constraint[T],
) *typedParam[T] {
	var zeroVal T

	return &typedParam[T]{
		ValueType:  reflect.TypeOf(zeroVal).String(),
		CurrentVal: zeroVal,
		DefaultVal: defaultVal,
		Constraint: constr,
	}
}

func NewIntParam(defaultVal int, constr Constraint[int]) *IntParam {
	return &IntParam{
		typedParam: newTypedParam(defaultVal, constr),
	}
}

func NewFloatParam(defaultVal float64, constr Constraint[float64]) *FloatParam {
	return &FloatParam{
		typedParam: newTypedParam[float64](defaultVal, constr),
	}
}

func (p *typedParam[T]) GetDefaultVal() any {
	return p.DefaultVal
}

func (p *typedParam[T]) GetConstraintInfo() *ConstraintInfo {
	return &ConstraintInfo{
		Type:   reflect.TypeOf(p.Constraint).String(),
		Limits: sliceutils.Map(p.Constraint.GetLimits(), func(t T) any { return t }),
	}
}

func (p *typedParam[T]) GetValueType() string {
	return p.ValueType
}

func (p *typedParam[T]) GetCurrentVal() any {
	return p.CurrentVal
}

func (p *typedParam[T]) SetCurrentVal(val any) error {
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

func (p *IntParam) SetCurrentVal(val any) error {
	if fltVal, ok := val.(float64); ok {
		intVal := int(fltVal)

		if float64(intVal) != fltVal {
			return commonerrs.NewErrWrongType("float64", "int")
		}

		return p.typedParam.SetCurrentVal(intVal)
	}

	return p.typedParam.SetCurrentVal(val)
}
