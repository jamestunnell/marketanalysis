package commonerrs

import "fmt"

type ErrExactCount struct {
	ItemsName        string
	Actual, Expected int
}

func NewErrExactCount(
	purpose string, expected, actual int) *ErrExactCount {
	return &ErrExactCount{
		ItemsName: purpose,
		Expected:  expected,
		Actual:    actual,
	}
}

func (err *ErrExactCount) Error() string {
	const strFmt = "expected exactly %d %s, got %d"
	return fmt.Sprintf(strFmt, err.Expected, err.ItemsName, err.Actual)
}
