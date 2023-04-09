package models

import "github.com/jamestunnell/marketanalysis/models/bar"

type Strategy interface {
	Type() string
	Params() map[string]any
	Algos() []Algorithm
	PositionSignal() chan PositionType

	WarmupPeriod() int
	WarmUp(bars []*bar.Bar) error

	Update(bar *bar.Bar)
}
