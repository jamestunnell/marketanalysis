package optimization

import (
	"math/rand"
	"sort"

	"github.com/MaxHalford/eaopt"
	"golang.org/x/exp/slices"
)

type OrderedInts struct {
	Length   int
	Ints     eaopt.IntSlice
	Min, Max int
}

func RandomOrderedInts(len int, rng *rand.Rand, min, max int) PartialGenome {
	ints := make([]int, len)

	for i := 0; i < len; i++ {
		ints[i] = rng.Intn(max-min) + min
	}

	sort.Ints(ints)

	return &OrderedInts{
		Length: len,
		Ints:   eaopt.IntSlice(ints),
		Min:    min,
		Max:    max,
	}
}

func (p *OrderedInts) Mutate(rng *rand.Rand) {
	idx := rng.Intn(p.Length)

	p.Ints[idx] = rng.Intn(p.Max-p.Min) + p.Min

	sort.Ints(p.Ints)
}

func (p *OrderedInts) Crossover(genome PartialGenome, rng *rand.Rand) {
	other := genome.(*OrderedInts)

	eaopt.CrossOX(p.Ints, other.Ints, rng)

	sort.Ints(p.Ints)
	sort.Ints(other.Ints)
}

func (p *OrderedInts) Clone() PartialGenome {
	return &OrderedInts{
		Length: p.Length,
		Ints:   slices.Clone(p.Ints),
		Min:    p.Min,
		Max:    p.Max,
	}
}
