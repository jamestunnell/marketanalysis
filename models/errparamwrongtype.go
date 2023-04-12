package models

import (
	"fmt"
	"reflect"
)

type ErrParamWrongType struct {
	Name  string
	Value any
}

func (err *ErrParamWrongType) Error() string {
	return fmt.Sprintf("param '%s' has wrong type '%s'", err.Name, reflect.TypeOf(err.Value).String())
}
