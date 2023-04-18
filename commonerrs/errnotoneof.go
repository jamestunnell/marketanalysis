package commonerrs

import "fmt"

type ErrNotOneOf[T any] struct {
	Name    string
	Actual  any
	Allowed []T
}

func NewErrNotOneOf[T any](name string, actual any, allowed []T) *ErrNotOneOf[T] {
	return &ErrNotOneOf[T]{
		Name:    name,
		Actual:  actual,
		Allowed: allowed,
	}
}

func (err *ErrNotOneOf[T]) Error() string {
	const strFmt = "%s %v is not one of %v"
	return fmt.Sprintf(strFmt, err.Name, err.Actual, err.Allowed)
}
