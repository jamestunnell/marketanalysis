package blocks

import (
	"fmt"
	"slices"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	goconstraints "golang.org/x/exp/constraints"
)

type Constraint interface {
	GetType() ConstraintType
	GetLimits() []any
}

type ConstraintType int

type TypedConstraint[T goconstraints.Ordered] struct {
	Type   ConstraintType
	Limits []T
}

const (
	None ConstraintType = iota
	Less
	LessEqual
	Greater
	GreaterEqual
	RangeIncl
	RangeExcl
	OneOf

	StrNone         = "none"
	StrLess         = "less"
	StrLessEqual    = "lessEqal"
	StrGreater      = "greater"
	StrGreaterEqual = "greaterEqual"
	StrRangeIncl    = "rangeIncl"
	StrRangeExcl    = "rangeExcl"
	StrOneOf        = "oneOf"
)

func (t ConstraintType) String() string {
	var s string

	switch t {
	case None:
		s = StrNone
	case Less:
		s = StrLess
	case LessEqual:
		s = StrLessEqual
	case Greater:
		s = StrGreater
	case GreaterEqual:
		s = StrGreaterEqual
	case RangeExcl:
		s = StrRangeExcl
	case RangeIncl:
		s = StrRangeIncl
	case OneOf:
		s = StrOneOf
	}

	return s
}

func NewNone[T goconstraints.Ordered]() *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: None, Limits: []T{}}
}

func NewLess[T goconstraints.Ordered](val T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: Less, Limits: []T{val}}
}

func NewLessEqual[T goconstraints.Ordered](val T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: LessEqual, Limits: []T{val}}
}

func NewGreater[T goconstraints.Ordered](val T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: Greater, Limits: []T{val}}
}

func NewGreaterEqual[T goconstraints.Ordered](val T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: GreaterEqual, Limits: []T{val}}
}

func NewRangeExcl[T goconstraints.Ordered](start, end T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: RangeExcl, Limits: []T{start, end}}
}

func NewRangeIncl[T goconstraints.Ordered](start, end T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: RangeIncl, Limits: []T{start, end}}
}

func NewOneOf[T goconstraints.Ordered](vals []T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: OneOf, Limits: vals}
}

func (c *TypedConstraint[T]) GetType() ConstraintType {
	return c.Type
}

func (c *TypedConstraint[T]) GetLimits() []any {
	return sliceutils.Map(c.Limits, func(val T) any { return val })
}

func (c *TypedConstraint[T]) CheckVal(val T) error {
	var err error

	switch c.Type {
	case None:
	case Less:
		if val >= c.Limits[0] {
			err = fmt.Errorf("%v is not < %v", val, c.Limits[0])
		}
	case Greater:
		if val <= c.Limits[0] {
			err = fmt.Errorf("%v is not > %v", val, c.Limits[0])
		}
	case LessEqual:
		if val > c.Limits[0] {
			err = fmt.Errorf("%v is not <= %v", val, c.Limits[0])
		}
	case GreaterEqual:
		if val < c.Limits[0] {
			err = fmt.Errorf("%v is not >= %v", val, c.Limits[0])
		}
	case RangeExcl:
		if val < c.Limits[0] || val >= c.Limits[1] {
			err = fmt.Errorf("%v is not in range [%v, %v)", val, c.Limits[0], c.Limits[1])
		}
	case RangeIncl:
		if val < c.Limits[0] || val > c.Limits[1] {
			err = fmt.Errorf("%v is not in range [%v, %v]", val, c.Limits[0], c.Limits[1])
		}
	case OneOf:
		if !slices.Contains(c.Limits, val) {
			err = fmt.Errorf("%v is not one of %v", val, c.Limits)
		}
	}

	return err
}
