package models

type Strategy interface {
	Type() string
	Params() map[string]any
	Algos() []Algorithm
	PositionSignal() chan PositionType

	WarmupPeriod() int
	WarmUp(bars []*Bar) error

	Update(bar *Bar)
}
