package models

import "fmt"

type ErrParamNotFound struct {
	Name string
}

func (err *ErrParamNotFound) Error() string {
	return fmt.Sprintf("param '%s' not found", err.Name)
}
