package app

import (
	"fmt"
)

type ResourceDef[T any] struct {
	KeyName    string
	Name       string
	NamePlural string
	Validate   func(*T) []error
	GetKey     func(*T) string
}

func (rdef *ResourceDef[T]) NewNotFoundError(key string) *Error {
	descr := fmt.Sprintf("%s with %s '%s'", rdef.Name, rdef.KeyName, key)

	return NewNotFoundError(descr)
}
