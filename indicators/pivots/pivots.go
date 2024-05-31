package pivots

import (
	"time"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type Pivots struct {
	Length       int
	stateMachine *statemachine.StateMachine[float64]
	warm         bool
	prev         *Pivot
	current      *Pivot
	all          []*Pivot
}

func New(length int) (*Pivots, error) {
	if length < 2 {
		return nil, commonerrs.NewErrLessThanMin("length", length, 2)
	}

	pivs := &Pivots{
		Length:       length,
		warm:         false,
		stateMachine: nil,
		prev:         nil,
		current:      nil,
		all:          []*Pivot{},
	}

	pivs.stateMachine = statemachine.New("pivots", &StateNew{parent: pivs})

	return pivs, nil
}

func (pivs *Pivots) WarmupPeriod() int {
	return 1
}

func (pivs *Pivots) IsWarm() bool {
	return pivs.warm
}

func (pivs *Pivots) Update(t time.Time, val float64) bool {
	prev := pivs.prev

	pivs.stateMachine.Run(t, val)

	if pivs.prev != prev {
		pivs.all = append(pivs.all, pivs.prev)

		return true
	}

	return false
}

func (pivs *Pivots) GetAll() []*Pivot {
	return pivs.all
}

func (pivs *Pivots) GetLatest() *Pivot {
	if len(pivs.all) == 0 {
		return nil
	}

	return sliceutils.Last(pivs.all)
}
