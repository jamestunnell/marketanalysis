package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type Up struct {
	*State
}

func NewUp(parent *TrendGate) *Up {
	return &Up{
		State: NewState(parent),
	}
}

func (state *Up) GetName() string {
	return "up"
}

func (state *Up) Enter() {
	state.parent.out.SetValue(1.0)
}

func (state *Up) Run(
	t time.Time,
	cur float64,
) statemachine.State[float64] {
	state.parent.out.SetValue(1.0)

	debouncePer := state.DebouncePeriod()

	if cur < -state.OpenThresh() {
		if debouncePer > 0 {
			return NewDownDebounce(state.parent)
		}

		return NewDown(state.parent)
	}

	if cur < state.CloseThresh() {
		return NewNone(state.parent)
	}

	return nil
}

func (state *Up) Exit() {

}
