package models

type Constraint interface {
	Type() string
	Check(any) error
	ValueBounds() []any
}
