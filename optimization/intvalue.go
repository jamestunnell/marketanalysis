package optimization

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

type IntConstMutator struct {
	Value int
}

func NewIntValue(m IntMutator) *IntValue {
	return &IntValue{
		Val:     0,
		Mutator: m,
	}
}

func NewIntConstMutator(val int) *IntConstMutator {
	return &IntConstMutator{Value: val}
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

func (v *IntValue) Init(rng *rand.Rand) {
	v.Val = v.Mutator.Start(rng)
}

func (v *IntValue) GetValue() any {
	return v.Val
}

func (p *IntValue) Mutate(rng *rand.Rand) {
	p.Val = p.Mutator.Mutate(p.Val, rng)
}

func (p *IntValue) Clone() Value {
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

func (m *IntConstMutator) Start(rng *rand.Rand) int {
	return m.Value
}

func (m *IntConstMutator) Mutate(currentVal int, rng *rand.Rand) int {
	return m.Value
}
