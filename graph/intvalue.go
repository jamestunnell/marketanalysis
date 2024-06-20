package graph

import (
	"math"
	"math/rand"

	"golang.org/x/exp/slices"
)

type IntValue struct {
	Val     int
	Mutator IntMutator
}

type IntMutator interface {
	Start(rng *rand.Rand) int
	Mutate(current int, rng *rand.Rand) int
}

type IntRangeMutator struct {
	Min, Max  int
	valRange  int
	valStdDev float64
}

type IntEnumMutator struct {
	Values    []int
	idxStdDev float64
}

func NewIntValue(m IntMutator, rng *rand.Rand) *IntValue {
	return &IntValue{
		Val:     m.Start(rng),
		Mutator: m,
	}
}

func NewIntEnumMutator(vals []int) *IntEnumMutator {
	return &IntEnumMutator{
		Values:    vals,
		idxStdDev: float64(len(vals)-1) / 10.0,
	}
}

func NewIntRangeMutator(min, max int) *IntRangeMutator {
	valRange := max - min

	return &IntRangeMutator{
		Min:       min,
		Max:       max,
		valRange:  valRange,
		valStdDev: float64(valRange) / 10.0,
	}
}

func (v *IntValue) Value() any {
	return v.Val
}

func (p *IntValue) Mutate(rng *rand.Rand) {
	p.Val = p.Mutator.Mutate(p.Val, rng)
}

func (p *IntValue) Clone() *IntValue {
	return &IntValue{
		Val:     p.Val,
		Mutator: p.Mutator,
	}
}

func (m *IntRangeMutator) Start(rng *rand.Rand) int {
	return m.Min + rng.Intn(m.valRange)
}

func (m *IntRangeMutator) Mutate(current int, rng *rand.Rand) int {
	new := int(float64(current) + rng.NormFloat64()*m.valStdDev)

	if new < m.Min {
		new += m.valRange
	} else if new > m.Max {
		new -= m.valRange
	}

	return new
}

func (m *IntEnumMutator) Start(rng *rand.Rand) int {
	return m.Values[rng.Intn(len(m.Values))]
}

func (m *IntEnumMutator) Mutate(currentVal int, rng *rand.Rand) int {
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
