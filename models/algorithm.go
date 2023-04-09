package models

type Algorithm interface {
	Type() string
	Params() map[string]any
	DirectionSignal() chan int

	WarmupPeriod() int
	WarmUp(bars []*Bar) error
	Update(bar *Bar)
	Current() map[string]float64
}
