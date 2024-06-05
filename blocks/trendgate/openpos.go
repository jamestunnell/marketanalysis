package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type OpenPos struct {
	*State
}

func NewOpenPos(parent *TrendGate) *OpenPos {
	return &OpenPos{State: NewState(parent)}
}

func (state *OpenPos) GetName() string {
	return "open(+)"
}

func (state *OpenPos) Enter() {
	state.setOutput(1.0)
}

func (state *OpenPos) Run(t time.Time, val float64) statemachine.State[float64] {
	if val < -state.parent.openThresh.CurrentVal {
		return NewOpenNeg(state.parent)
	}

	if val < state.parent.closeThresh.CurrentVal {
		return NewClosed(state.parent)
	}

	state.setOutput(1.0)

	return nil
}

func (state *OpenPos) Exit() {

}
