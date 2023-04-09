package models

type Indicator interface {
	WarmupPeriod() int
	WarmUp(bars Bars) error
	Update(bar *Bar) float64
	Current() float64
}
