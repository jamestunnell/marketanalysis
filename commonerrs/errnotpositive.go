package commonerrs

import "fmt"

type ErrNotPositive struct {
	Name  string
	Value any
}

func NewErrNotPositive(name string, val any) *ErrNotPositive {
	return &ErrNotPositive{Name: name, Value: val}
}

func (err *ErrNotPositive) Error() string {
	return fmt.Sprintf("%s %v is not positive", err.Name, err.Value)
}
