package supertrend

import (
	"math"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type Normal struct {
	*State

	prev *values
}

type values struct {
	in    float64
	upper float64
	lower float64
	trend float64
}

func NewNormal(parent *Supertrend, prev *values) *Normal {
	return &Normal{
		State: NewState(parent),
		prev:  prev,
	}
}

func (state *Normal) GetName() string {
	return "normal"
}

func (state *Normal) Enter() {
}

func (state *Normal) Run(
	t time.Time,
	cur *models.OHLC,
) statemachine.State[*models.OHLC] {
	state.parent.atr.Update(cur)

	atr := state.parent.atr.Current() * state.parent.atrMul.CurrentVal
	in := state.parent.in.GetValue()
	upper := in + atr
	lower := in - atr

	var newTrend float64
	var newUpper float64
	var newLower float64

	if state.prev.trend <= 0.0 && in > state.prev.upper {
		newTrend = 1.0
	} else if state.prev.trend >= 0.0 && in < state.prev.lower {
		newTrend = -1.0
	} else {
		newTrend = state.prev.trend
	}

	if state.prev.in < state.prev.upper {
		newUpper = math.Min(upper, state.prev.upper)
	} else {
		newUpper = upper
	}

	if state.prev.in > state.prev.lower {
		newLower = math.Max(lower, state.prev.lower)
	} else {
		newLower = lower
	}

	state.parent.trend.SetValue(newTrend)
	state.parent.upper.SetValue(newUpper)
	state.parent.lower.SetValue(newLower)

	state.prev.trend = newTrend
	state.prev.upper = newUpper
	state.prev.lower = newLower
	state.prev.in = in

	return nil
}

func (state *Normal) Exit() {

}
