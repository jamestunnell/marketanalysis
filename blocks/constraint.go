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
	ExclusiveMax ConstraintType = iota
	ExclusiveMin
	InclusiveMax
	InclusiveMin
	OneOf

	StrExclusiveMax = "ExclusiveMax"
	StrExclusiveMin = "ExclusiveMin"
	StrInclusiveMax = "InclusiveMax"
	StrInclusiveMin = "InclusiveMin"
	StrOneOf        = "OneOf"
)

func (t ConstraintType) String() string {
	var s string

	switch t {
	case ExclusiveMax:
		s = StrExclusiveMax
	case ExclusiveMin:
		s = StrExclusiveMin
	case InclusiveMax:
		s = StrInclusiveMax
	case InclusiveMin:
		s = StrInclusiveMin
	case OneOf:
		s = StrOneOf
	}

	return s
}

func NewExclusiveMax[T goconstraints.Ordered](max T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: ExclusiveMax, Limits: []T{max}}
}

func NewInclusiveMin[T goconstraints.Ordered](min T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: InclusiveMin, Limits: []T{min}}
}

func NewExclusiveMin[T goconstraints.Ordered](min T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: ExclusiveMin, Limits: []T{min}}
}

func NewInclusiveMax[T goconstraints.Ordered](max T) *TypedConstraint[T] {
	return &TypedConstraint[T]{Type: InclusiveMax, Limits: []T{max}}
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
	case ExclusiveMax:
		if val >= c.Limits[0] {
			err = fmt.Errorf("%v is not less than %v", val, c.Limits[0])
		}
	case ExclusiveMin:
		if val <= c.Limits[0] {
			err = fmt.Errorf("%v is not greater than %v", val, c.Limits[0])
		}
	case InclusiveMax:
		if val > c.Limits[0] {
			err = fmt.Errorf("%v is greater than %v", val, c.Limits[0])
		}
	case InclusiveMin:
		if val < c.Limits[0] {
			err = fmt.Errorf("%v is less than %v", val, c.Limits[0])
		}
	case OneOf:
		if !slices.Contains(c.Limits, val) {
			err = fmt.Errorf("%v is not one of %v", val, c.Limits)
		}
	}

	return err
}
