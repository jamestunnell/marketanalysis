package trendgate

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type OpenNeg struct {
	*State
}

func NewOpenNeg(parent *TrendGate) *OpenNeg {
	return &OpenNeg{State: NewState(parent)}
}

func (state *OpenNeg) GetName() string {
	return "open(-)"
}

func (state *OpenNeg) Enter() {
	state.setOutput(-1.0)
}

func (state *OpenNeg) Run(t time.Time, val float64) statemachine.State[float64] {
	if val > state.parent.openThresh.CurrentVal {
		return NewOpenPos(state.parent)
	}

	if val > -state.parent.closeThresh.CurrentVal {
		return NewClosed(state.parent)
	}

	state.setOutput(-1.0)

	return nil
}

func (state *OpenNeg) Exit() {

}
