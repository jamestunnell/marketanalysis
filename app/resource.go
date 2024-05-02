package app

import (
	"fmt"
	"reflect"
)

type Resource interface {
	Validate() []error
	GetKey() string
}

type ResourceInfo struct {
	KeyName    string
	Name       string
	NamePlural string
}

func NewResource[T Resource]() T {
	var val T

	return reflect.New(reflect.TypeOf(val).Elem()).Interface().(T)
}

func (rdef *ResourceInfo) NewErrNotFound(key string) Error {
	what := fmt.Sprintf("%s with %s '%s'", rdef.Name, rdef.KeyName, key)

	return NewErrNotFound(what)
}
