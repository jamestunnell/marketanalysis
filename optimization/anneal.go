package optimization

import (
	"github.com/ccssmnn/hego"
)

type SAState[T any] struct {
	Objective Objective[T]
	Base      State[T]
}

func (s *SAState[T]) Neighbor() hego.AnnealingState {
	n := &SAState[T]{
		Objective: s.Objective,
		Base:      s.Base.Clone(),
	}

	n.Base.Mutate()

	return n
}

// Energy returns the energy of the current state. Lower is better
func (s *SAState[T]) Energy() float64 {
	return s.Objective.Measure(s.Base.GetMeasureVal())
}
