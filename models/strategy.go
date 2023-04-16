package models

type Strategy interface {
	Type() string
	Params() Params
	ClosedPositions() Positions

	WarmupPeriod() int

	Initialize(bars Bars) error
	Update(bar *Bar)
	Close(bar *Bar, reason string)
}
