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

func (rdef *ResourceDef[T]) NewErrNotFound(key string) Error {
	what := fmt.Sprintf("%s with %s '%s'", rdef.Name, rdef.KeyName, key)

	return NewErrNotFound(what)
}
