package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type Down struct {
	*State
}

func NewDown(parent *TrendGate) *Down {
	return &Down{
		State: NewState(parent),
	}
}

func (state *Down) GetName() string {
	return "down"
}

func (state *Down) Enter() {
	state.parent.out.SetValue(-1.0)
}

func (state *Down) Run(
	t time.Time,
	cur float64,
) statemachine.State[float64] {
	state.parent.out.SetValue(-1.0)

	debouncePer := state.DebouncePeriod()

	if cur > state.OpenThresh() {
		if debouncePer > 0 {
			return NewUpDebounce(state.parent)
		}

		return NewUp(state.parent)
	}

	if cur > -state.CloseThresh() {
		return NewNone(state.parent)
	}

	return nil
}

func (state *Down) Exit() {

}
