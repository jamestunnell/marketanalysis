package dateutils

import "github.com/rickb777/date"

type Controller interface {
	Current() date.Date
	Advance()
	AnyLeft() bool
}
