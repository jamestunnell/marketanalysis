package optimization

import (
	"fmt"
	"math/rand"

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

	// Add a callback to stop when the problem is solved
	ga.EarlyStop = func(ga *eaopt.GA) bool {
		return ga.HallOfFame[0].Fitness == 0.0
	}

	// Run the GA
	ga.Minimize(newGenome)

	return ga.HallOfFame, nil
}
