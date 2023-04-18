package collection

import "fmt"

type ErrNotDir struct {
	Path string
}

func (err *ErrNotDir) Error() string {
	return fmt.Sprintf("'%s' is not a dir", err.Path)
}
