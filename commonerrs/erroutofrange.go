package commonerrs

import "fmt"

type ErrOutOfRange struct {
	Name          string
	Val, Min, Max any
}

func NewErrOutOfRange(name string, val, min, max any) *ErrOutOfRange {
	return &ErrOutOfRange{
		Name: name,
		Val:  val,
		Min:  min,
		Max:  max,
	}
}

func (err *ErrOutOfRange) Error() string {
	const strFmt = "%s %v is not in range %v..%v"
	return fmt.Sprintf(strFmt, err.Name, err.Val, err.Min, err.Max)
}
