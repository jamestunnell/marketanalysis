package optimization

import (
	"math/rand"

	"github.com/MaxHalford/eaopt"
)

type FitFunc func(eaopt.Genome) (float64, error)

type CompositeGenome struct {
	Fit   FitFunc
	Parts []PartialGenome
}

type PartialGenome interface {
	Mutate(rng *rand.Rand)
	Crossover(genome PartialGenome, rng *rand.Rand)
	Clone() PartialGenome
	Value() any
}

func NewCompositeGenome(fit FitFunc, parts ...PartialGenome) eaopt.Genome {
	return &CompositeGenome{
		Fit:   fit,
		Parts: parts,
	}
}

func (g *CompositeGenome) Evaluate() (float64, error) {
	return g.Fit(g)
}

func (g *CompositeGenome) Mutate(rng *rand.Rand) {
	for _, part := range g.Parts {
		part.Mutate(rng)
	}
}

func (g *CompositeGenome) Crossover(genome eaopt.Genome, rng *rand.Rand) {
	other := genome.(*CompositeGenome)

	for i, part := range g.Parts {
		part.Crossover(other.Parts[i], rng)
	}
}

func (g *CompositeGenome) Clone() eaopt.Genome {
	parts := make([]PartialGenome, len(g.Parts))
	for i, part := range g.Parts {
		parts[i] = part.Clone()
	}

	return NewCompositeGenome(g.Fit, parts...)
}
