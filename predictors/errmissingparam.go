package predictors

import "fmt"

type ErrMissingParam struct {
	Name string
}

func (err *ErrMissingParam) Error() string {
	return fmt.Sprintf("missing param %s", err.Name)
}
