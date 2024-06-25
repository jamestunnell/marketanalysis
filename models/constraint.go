package models

import (
	"fmt"
	"slices"

	"golang.org/x/exp/constraints"
)

type Constraint[T constraints.Ordered] interface {
	GetLimits() []T
	CheckVal(T) error
}

type Unconstrained[T constraints.Ordered] struct{}
type Less[T constraints.Ordered] struct{ Max T }
type LessEq[T constraints.Ordered] struct{ Max T }
type Greater[T constraints.Ordered] struct{ Min T }
type GreaterEq[T constraints.Ordered] struct{ Min T }
type RangeIncl[T constraints.Ordered] struct{ Min, Max T }
type RangeExcl[T constraints.Ordered] struct{ Min, Max T }
type OneOf[T constraints.Ordered] struct{ Allowed []T }

func NewUnconstrained[T constraints.Ordered]() *Unconstrained[T] {
	return &Unconstrained[T]{}
}

func NewLess[T constraints.Ordered](max T) *Less[T] {
	return &Less[T]{Max: max}
}

func NewLessEq[T constraints.Ordered](max T) *LessEq[T] {
	return &LessEq[T]{Max: max}
}

func NewGreater[T constraints.Ordered](min T) *Greater[T] {
	return &Greater[T]{Min: min}
}

func NewGreaterEq[T constraints.Ordered](min T) *GreaterEq[T] {
	return &GreaterEq[T]{Min: min}
}

func NewRangeExcl[T constraints.Ordered](min, max T) *RangeExcl[T] {
	return &RangeExcl[T]{Min: min, Max: max}
}

func NewRangeIncl[T constraints.Ordered](min, max T) *RangeIncl[T] {
	return &RangeIncl[T]{Min: min, Max: max}
}

func NewOneOf[T constraints.Ordered](allowed []T) *OneOf[T] {
	return &OneOf[T]{Allowed: allowed}
}

func (c *Unconstrained[T]) GetLimits() []T {
	return []T{}
}

func (c *Less[T]) GetLimits() []T {
	return []T{c.Max}
}

func (c *LessEq[T]) GetLimits() []T {
	return []T{c.Max}
}

func (c *Greater[T]) GetLimits() []T {
	return []T{c.Min}
}

func (c *GreaterEq[T]) GetLimits() []T {
	return []T{c.Min}
}

func (c *RangeExcl[T]) GetLimits() []T {
	return []T{c.Min, c.Max}
}

func (c *RangeIncl[T]) GetLimits() []T {
	return []T{c.Min, c.Max}
}

func (c *OneOf[T]) GetLimits() []T {
	return slices.Clone(c.Allowed)
}

func (c *Unconstrained[T]) CheckVal(val T) error {
	return nil
}

func (c *Less[T]) CheckVal(val T) error {
	if val >= c.Max {
		return fmt.Errorf("%v is not < %v", val, c.Max)
	}

	return nil
}

func (c *LessEq[T]) CheckVal(val T) error {
	if val > c.Max {
		return fmt.Errorf("%v is not <= %v", val, c.Max)
	}

	return nil
}

func (c *Greater[T]) CheckVal(val T) error {
	if val <= c.Min {
		return fmt.Errorf("%v is not > %v", val, c.Min)
	}

	return nil
}

func (c *GreaterEq[T]) CheckVal(val T) error {
	if val < c.Min {
		return fmt.Errorf("%v is not >= %v", val, c.Min)
	}

	return nil
}

func (c *RangeExcl[T]) CheckVal(val T) error {
	if val < c.Min || val >= c.Max {
		return fmt.Errorf("%v is not in range [%v, %v)", val, c.Min, c.Max)
	}

	return nil
}

func (c *RangeIncl[T]) CheckVal(val T) error {
	if val < c.Min || val > c.Max {
		return fmt.Errorf("%v is not in range [%v, %v]", val, c.Min, c.Max)
	}

	return nil
}

func (c *OneOf[T]) CheckVal(val T) error {
	if !slices.Contains(c.Allowed, val) {
		return fmt.Errorf("%v is not one of %v", val, c.Allowed)
	}

	return nil
}
