package optimization

import (
	"math/rand"
)

type FloatValue struct {
	Value      float64
	Min, Range float64
}

func RandomFloatValue(rng *rand.Rand, min, max float64) *FloatValue {
	fv := &FloatValue{
		Value: 0.0,
		Min:   min,
		Range: max - min,
	}

	fv.Mutate(rng)

	return fv
}

func (p *FloatValue) Mutate(rng *rand.Rand) {
	p.Value = p.Min + p.Range*rng.Float64()
}

func (p *FloatValue) Crossover(genome PartialGenome, rng *rand.Rand) {
	other := genome.(*FloatValue)

	if rng.Float64() < 0.5 {
		p.Value, other.Value = other.Value, p.Value
	}
}

func (p *FloatValue) Clone() PartialGenome {
	return &FloatValue{
		Value: p.Value,
		Min:   p.Min,
		Range: p.Range,
	}
}
