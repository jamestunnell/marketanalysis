package commonerrs

import "fmt"

type ErrNotFound struct {
	Type, Name string
}

func NewErrNotFound(typ, name string) *ErrNotFound {
	return &ErrNotFound{
		Name: name,
		Type: typ,
	}
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("%s %s not found", err.Type, err.Name)
}
