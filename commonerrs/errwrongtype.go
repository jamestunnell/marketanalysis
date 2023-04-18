package commonerrs

import "fmt"

type ErrWrongType struct {
	Actual, Expected string
}

func NewErrWrongType(actual, expected string) *ErrWrongType {
	return &ErrWrongType{
		Actual:   actual,
		Expected: expected,
	}
}

func (err *ErrWrongType) Error() string {
	return fmt.Sprintf("type %s is not %s", err.Actual, err.Expected)
}
