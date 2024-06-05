package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type Closed struct {
	*State
}

func NewClosed(parent *TrendGate) *Closed {
	return &Closed{
		State: NewState(parent),
	}
}

func (state *Closed) GetName() string {
	return "closed"
}

func (state *Closed) Enter() {
	state.setOutput(0.0)
}

func (state *Closed) Run(t time.Time, val float64) statemachine.State[float64] {
	if val > state.parent.openThresh.CurrentVal {
		return NewOpenPos(state.parent)
	}

	if val < -state.parent.openThresh.CurrentVal {
		return NewOpenNeg(state.parent)
	}

	state.setOutput(0.0)

	return nil
}

func (state *Closed) Exit() {

}
