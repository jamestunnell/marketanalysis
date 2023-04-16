package commonerrs

import "fmt"

type ErrMinCount struct {
	Name            string
	Count, MinCount int
}

func NewErrMinCount(name string, count, minCount int) *ErrMinCount {
	return &ErrMinCount{
		Name:     name,
		Count:    count,
		MinCount: minCount,
	}
}

func (err *ErrMinCount) Error() string {
	const strFmt = "%s count %d is less than min%d"
	return fmt.Sprintf(strFmt, err.Name, err.Count, err.MinCount)
}
