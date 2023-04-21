package optimization

import "fmt"

type ErrUnsupportedType struct {
	Type string
}

func (err *ErrUnsupportedType) Error() string {
	return fmt.Sprintf("type '%s' is not supported", err.Type)
}
