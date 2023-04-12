package optimization

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MaxHalford/eaopt"
)

type NewGenomeFunc func(rng *rand.Rand) eaopt.Genome

// EAOpt uses evolutionary algorithms to optimize
func EAOpt(newGenome NewGenomeFunc, config eaopt.GAConfig) (eaopt.Individuals, error) {
	var ga, err = config.NewGA()
	if err != nil {
		err = fmt.Errorf("failed to make GA: %w", err)

		return eaopt.Individuals{}, err
	}

	// just used for rand ID in cloning indivuals
	rng := rand.New(rand.NewSource(time.Now().Unix()))

	best := eaopt.Individuals{}

	// Add a callback to stop when the problem is solved
	ga.Callback = func(ga *eaopt.GA) {
		newBest := make(eaopt.Individuals, len(ga.HallOfFame))

		for i, indiv := range ga.HallOfFame {
			newBest[i] = indiv.Clone(rng)
		}

		best = newBest
	}

	// Run the GA
	ga.Minimize(newGenome)

	return best, nil
}
