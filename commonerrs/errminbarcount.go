package commonerrs

import "fmt"

type ErrMinBarCount struct {
	Purpose     string
	Actual, Min int
}

func NewErrMinBarCount(
	purpose string, min, actual int) *ErrMinBarCount {
	return &ErrMinBarCount{
		Purpose: purpose,
		Min:     min,
		Actual:  actual,
	}
}

func (err *ErrMinBarCount) Error() string {
	const strFmt = "expected at least %d %s bars, got %d"
	return fmt.Sprintf(strFmt, err.Min, err.Purpose, err.Actual)
}
