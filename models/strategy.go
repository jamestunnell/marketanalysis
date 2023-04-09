package models

type Strategy interface {
	Type() string
	Params() Params

	ClosedPositions() []Position
	Close(bar *Bar)

	WarmupPeriod() int
	WarmUp(bars []*Bar) error

	Update(bar *Bar)
}
