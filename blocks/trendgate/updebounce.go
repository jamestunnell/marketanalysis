package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
	"github.com/rs/zerolog/log"
)

type UpDebounce struct {
	*State

	count int
}

func NewUpDebounce(parent *TrendGate) *UpDebounce {
	return &UpDebounce{
		State: NewState(parent),
		count: 1,
	}
}

func (state *UpDebounce) GetName() string {
	return "updebounce"
}

func (state *UpDebounce) Enter() {
	state.parent.out.SetValue(0.0)
}

func (state *UpDebounce) Run(
	t time.Time,
	cur float64,
) statemachine.State[float64] {
	state.parent.out.SetValue(0.0)

	debouncePer := state.DebouncePeriod()

	log.Debug().
		Float64("cur", cur).
		Time("timestamp", t).
		Float64("-openThresh", -state.OpenThresh()).
		Float64("closeThresh", state.CloseThresh()).
		Int("debouncePer", debouncePer).
		Msg("run updebounce")

	if cur < -state.OpenThresh() {

		if debouncePer > 0 {
			return NewDownDebounce(state.parent)
		}

		return NewDown(state.parent)
	}

	if cur < state.CloseThresh() {
		return NewNone(state.parent)
	}

	if state.count >= debouncePer {
		return NewUp(state.parent)
	}

	state.count++

	return nil
}

func (state *UpDebounce) Exit() {

}
