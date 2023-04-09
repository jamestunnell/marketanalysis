package models

import "fmt"

type ErrParamWrongType struct {
	Name string
}

func (err *ErrParamWrongType) Error() string {
	return fmt.Sprintf("param '%s' has wrong type", err.Name)
}
