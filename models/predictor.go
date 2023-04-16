package models

type Predictor interface {
	Type() string
	Params() Params
	WarmupPeriod() int
	WarmUp(bars Bars) error
	Update(bar *Bar)
	Direction() Direction
}
