package optimization

import (
	"math"
	"math/rand"

	"golang.org/x/exp/slices"
)

type FloatValue struct {
	Val     float64
	Mutator FloatMutator
}

type FloatMutator interface {
	Start(rng *rand.Rand) float64
	Mutate(current float64, rng *rand.Rand) float64
}

type FloatRangeMutator struct {
	Min, Max  float64
	valRange  float64
	valStdDev float64
}

type FloatEnumMutator struct {
	Values    []float64
	idxStdDev float64
}

func NewFloatValue(m FloatMutator) *FloatValue {
	return &FloatValue{
		Val:     0.0,
		Mutator: m,
	}
}

func NewFloatEnumMutator(vals []float64) *FloatEnumMutator {
	return &FloatEnumMutator{
		Values:    vals,
		idxStdDev: float64(len(vals)-1) / 10.0,
	}
}

func NewFloatRangeMutator(min, max float64) *FloatRangeMutator {
	valRange := max - min

	return &FloatRangeMutator{
		Min:       min,
		Max:       max,
		valRange:  valRange,
		valStdDev: valRange / 10.0,
	}
}

func (v *FloatValue) Init(rng *rand.Rand) {
	v.Val = v.Mutator.Start(rng)
}

func (v *FloatValue) GetValue() any {
	return v.Val
}

func (p *FloatValue) Mutate(rng *rand.Rand) {
	p.Val = p.Mutator.Mutate(p.Val, rng)
}

func (p *FloatValue) Clone() Value {
	return &FloatValue{
		Val:     p.Val,
		Mutator: p.Mutator,
	}
}

func (m *FloatRangeMutator) Start(rng *rand.Rand) float64 {
	return m.Min + (m.valRange * rng.Float64())
}

func (m *FloatRangeMutator) Mutate(current float64, rng *rand.Rand) float64 {
	new := current + rng.NormFloat64()*m.valStdDev

	if new < m.Min {
		new += m.valRange
	} else if new > m.Max {
		new -= m.valRange
	}

	return new
}

func (m *FloatEnumMutator) Start(rng *rand.Rand) float64 {
	return m.Values[rng.Intn(len(m.Values))]
}

func (m *FloatEnumMutator) Mutate(currentVal float64, rng *rand.Rand) float64 {
	currentIdx := slices.Index(m.Values, currentVal)
	idx := int(math.Round(float64(currentIdx) + rng.NormFloat64()*m.idxStdDev))

	n := len(m.Values)
	if idx >= n {
		idx -= n
	} else if idx < 0 {
		idx += n
	}

	return m.Values[idx]
}
