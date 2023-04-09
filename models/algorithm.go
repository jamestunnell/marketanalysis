package models

import "github.com/jamestunnell/marketanalysis/models/bar"

type Algorithm interface {
	Type() string
	Params() map[string]any
	DirectionSignal() chan int

	WarmupPeriod() int
	WarmUp(bars []*bar.Bar) error
	Update(bar *bar.Bar)
	Current() map[string]float64
}
