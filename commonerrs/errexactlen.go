package commonerrs

import "fmt"

type ErrExactLen struct {
	Name             string
	Actual, Expected int
}

func NewErrExactLen(
	name string, actual, expected int) *ErrExactLen {
	return &ErrExactLen{
		Name:     name,
		Expected: expected,
		Actual:   actual,
	}
}

func (err *ErrExactLen) Error() string {
	const strFmt = "%s len %d is not %d"
	return fmt.Sprintf(strFmt, err.Name, err.Actual, err.Expected)
}
