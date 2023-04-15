package commonerrs

import "fmt"

type ErrLessThanMin struct {
	Name        string
	Min, Actual any
}

func NewErrLessThanMin(name string, actual, min any) *ErrLessThanMin {
	return &ErrLessThanMin{
		Name:   name,
		Min:    min,
		Actual: actual,
	}
}

func (err *ErrLessThanMin) Error() string {
	const strFmt = "%s %v is less than min %v"
	return fmt.Sprintf(strFmt, err.Name, err.Actual, err.Min)
}
