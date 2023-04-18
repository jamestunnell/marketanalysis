package collection

import "fmt"

type ErrNotFile struct {
	Path string
}

func (err *ErrNotFile) Error() string {
	return fmt.Sprintf("'%s' is not a file", err.Path)
}
