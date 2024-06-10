package supertrend

import (
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type WaitForFirst struct {
	*State
}

func NewWaitForFirst(parent *Supertrend) *WaitForFirst {
	return &WaitForFirst{
		State: NewState(parent),
	}
}

func (state *WaitForFirst) GetName() string {
	return "waitforfirst"
}

func (state *WaitForFirst) Enter() {
}

func (state *WaitForFirst) Run(
	t time.Time,
	cur *models.OHLC,
) statemachine.State[*models.OHLC] {
	state.parent.atr.Update(cur)

	if !state.parent.atr.Warm() || !state.parent.in.IsValueSet() {
		return nil
	}

	atr := state.parent.atr.Current() * state.parent.atrMul.CurrentVal
	in := state.parent.in.GetValue()
	upper := in + atr
	lower := in - atr

	state.parent.trend.SetValue(0.0)
	state.parent.upper.SetValue(upper)
	state.parent.lower.SetValue(lower)

	prevVals := &values{in: in, upper: upper, lower: lower, trend: 0.0}

	return NewNormal(state.parent, prevVals)
}

func (state *WaitForFirst) Exit() {

}
