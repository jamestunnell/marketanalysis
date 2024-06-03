package pivots

import (
	"time"

	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type StateNew struct {
	parent *Pivots
}

func (state *StateNew) GetName() string {
	return "new"
}

func (state *StateNew) Enter() {}
func (state *StateNew) Exit() {
	state.parent.warm = true
}

func (state *StateNew) Run(t time.Time, val float64) statemachine.State[float64] {
	return NewStateNeutral(t, val, state.parent)
}
