package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type DownDebounce struct {
	*State

	count int
}

func NewDownDebounce(parent *TrendGate) *DownDebounce {
	return &DownDebounce{
		State: NewState(parent),
		count: 1,
	}
}

func (state *DownDebounce) GetName() string {
	return "downdebounce"
}

func (state *DownDebounce) Enter() {
	state.parent.out.SetValue(0.0)
}

func (state *DownDebounce) Run(
	t time.Time,
	cur float64,
) statemachine.State[float64] {
	state.parent.out.SetValue(0.0)

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

	if state.count >= debouncePer {
		return NewDown(state.parent)
	}

	state.count++

	return nil
}

func (state *DownDebounce) Exit() {

}
