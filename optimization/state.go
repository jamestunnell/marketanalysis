package optimization

import (
	"math/rand"

	"github.com/jamestunnell/marketanalysis/models"
)

type State[T any] interface {
	GetMeasureVal() T
	Clone() State[T]
	Mutate()
}

type ParameterState struct {
	RNG    *rand.Rand
	Values Values
}

func NewParameterState(rng *rand.Rand, values Values) *ParameterState {
	for _, v := range values {
		v.Init(rng)
	}

	return &ParameterState{RNG: rng, Values: values}
}

func (s *ParameterState) GetMeasureVal() models.ParamVals {
	paramVals := models.ParamVals{}

	for name, v := range s.Values {
		paramVals[name] = v.GetValue()
	}

	return paramVals
}

func (s *ParameterState) Clone() State[models.ParamVals] {
	values := Values{}

	for name, v := range s.Values {
		values[name] = v.Clone()
	}

	return &ParameterState{RNG: s.RNG, Values: values}
}

func (s *ParameterState) Mutate() {
	for _, v := range s.Values {
		v.Mutate(s.RNG)
	}
}
