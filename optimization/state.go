package optimization

import (
	"math/rand"
)

type State[T any] interface {
	GetMeasureVal() T
	Clone() State[T]
	Mutate()
}

type ParameterState struct {
	RNG       *rand.Rand
	Objective Objective[Values]
	Values    Values
}

func NewParameterState(
	rng *rand.Rand,
	objective Objective[Values],
) *ParameterState {
	return &ParameterState{
		Objective: objective,
		RNG:       rng,
		Values:    map[string]Value{},
	}
}

func (s *ParameterState) GetMeasureVal() Values {
	return s.Values
}

func (s *ParameterState) Clone() State[Values] {
	s2 := NewParameterState(s.RNG, s.Objective)

	for name, v := range s.Values {
		s2.Values[name] = v.Clone()
	}

	return s2
}

func (s *ParameterState) Mutate() {
	for _, v := range s.Values {
		v.Mutate(s.RNG)
	}
}
