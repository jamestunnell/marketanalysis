package optimization

import (
	"github.com/ccssmnn/hego"
)

type SAState[T any] struct {
	Objective   Objective[T]
	Base        State[T]
	MeasureHook func(val T, score float64)
}

func (s *SAState[T]) Neighbor() hego.AnnealingState {
	n := &SAState[T]{
		Objective:   s.Objective,
		Base:        s.Base.Clone(),
		MeasureHook: s.MeasureHook,
	}

	n.Base.Mutate()

	return n
}

// Energy returns the energy of the current state. Lower is better
func (s *SAState[T]) Energy() float64 {
	val := s.Base.GetMeasureVal()
	energy := s.Objective.Measure(val)

	s.MeasureHook(val, energy)

	return energy
}
