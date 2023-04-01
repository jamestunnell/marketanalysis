package commonerrs

import "fmt"

type ErrExactBarCount struct {
	Purpose          string
	Actual, Expected int
}

func NewErrExactBarCount(
	purpose string, expected, actual int) *ErrExactBarCount {
	return &ErrExactBarCount{
		Purpose:  purpose,
		Expected: expected,
		Actual:   actual,
	}
}

func (err *ErrExactBarCount) Error() string {
	const strFmt = "expected exactly %d %s bars, got %d"
	return fmt.Sprintf(strFmt, err.Expected, err.Purpose, err.Actual)
}
