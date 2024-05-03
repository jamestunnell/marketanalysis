package commonerrs

import "fmt"

type ErrMoreThanMax struct {
	Name        string
	Max, Actual any
}

func NewErrMoreThanMax(name string, actual, max any) *ErrMoreThanMax {
	return &ErrMoreThanMax{
		Name:   name,
		Max:    max,
		Actual: actual,
	}
}

func (err *ErrMoreThanMax) Error() string {
	const strFmt = "%s %v is more than max %v"
	return fmt.Sprintf(strFmt, err.Name, err.Actual, err.Max)
}
