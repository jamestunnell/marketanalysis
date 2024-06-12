package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type None struct {
	*State
}

func NewNone(parent *TrendGate) *None {
	return &None{
		State: NewState(parent),
	}
}

func (state *None) GetName() string {
	return "none"
}

func (state *None) Enter() {
	state.parent.out.SetValue(0.0)
}

func (state *None) Run(
	t time.Time,
	cur float64,
) statemachine.State[float64] {
	state.parent.out.SetValue(0.0)

	thresh := state.OpenThresh()
	debouncePer := state.DebouncePeriod()

	if cur > thresh {
		if debouncePer > 0 {
			return NewUpDebounce(state.parent)
		}

		return NewUp(state.parent)
	}

	if cur < -thresh {
		if debouncePer > 0 {
			return NewDownDebounce(state.parent)
		}

		return NewDown(state.parent)
	}

	return nil
}

func (state *None) Exit() {

}
