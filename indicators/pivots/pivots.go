package pivots

import (
	"time"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
	"github.com/rs/zerolog/log"
)

type Pivots struct {
	Length       int
	stateMachine *statemachine.StateMachine[float64]
	warm         bool
	current      *Pivot
	completed    []*Pivot
}

func New(length int) (*Pivots, error) {
	if length < 2 {
		return nil, commonerrs.NewErrLessThanMin("length", length, 2)
	}

	pivs := &Pivots{
		Length:       length,
		warm:         false,
		stateMachine: nil,
		completed:    []*Pivot{},
		current:      nil,
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
	count := len(pivs.completed)

	pivs.stateMachine.Run(t, val)

	return len(pivs.completed) > count
}

func (pivs *Pivots) GetInProgress() *Pivot {
	return pivs.current
}

func (pivs *Pivots) GetCompleted() []*Pivot {
	return pivs.completed
}

func (pivs *Pivots) GetLastCompleted() *Pivot {
	n := len(pivs.completed)
	if n == 0 {
		log.Warn().Msg("no pivots completed yet")

		return nil
	}

	return pivs.completed[n-1]
}

func (pivs *Pivots) addCompleted(piv *Pivot) {
	log.Debug().Interface("pivot", piv).Msg("completed pivot")

	pivs.completed = append(pivs.completed, piv)
}
