package optimization

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ccssmnn/hego"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Results struct {
	Result        *Result       `json:"result"`
	Runtime       time.Duration `json:"runtime"`
	Iterations    int           `json:"iterations"`
	ResultHistory []*Result     `json:"resultHistory"`
}

type Result struct {
	Value any     `json:"value"`
	Score float64 `json:"score"`
}

type Objective[T any] interface {
	Measure(T) float64
}

const (
	AlgorithmSA = "Anneal"
)

func OptimizeParameters(
	settings *Settings,
	values Values,
	objective Objective[models.ParamVals],
) (*Results, error) {
	if settings.Algorithm != AlgorithmSA {
		err := errors.New("unsupported optimization algorithm " + settings.Algorithm)

		return nil, err
	}

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	initialState := &SAState[models.ParamVals]{
		Objective: objective,
		Base:      NewParameterState(rng, values),
	}

	log.Debug().Interface("initial state", initialState).Msg("starting SA optimization")

	return OptimizeSA(settings, initialState)
}

func OptimizeSA[T any](settings *Settings, initialState *SAState[T]) (*Results, error) {
	saSettings := hego.SASettings{
		Temperature:     10.0,
		AnnealingFactor: 0.999,
		Settings: hego.Settings{
			MaxIterations: settings.MaxIterations,
			KeepHistory:   settings.KeepHistory,
		},
	}

	r, err := hego.SA(initialState, saSettings)
	if err != nil {
		return nil, fmt.Errorf("simulated annealing failed: %w", err)
	}

	makeOptResult := func(state hego.AnnealingState, energy float64) *Result {
		return &Result{
			Value: state.(*SAState[T]).Base.GetMeasureVal(),
			Score: energy,
		}
	}
	history := make([]*Result, len(r.States))

	for i := 0; i < len(r.States); i++ {
		history[i] = makeOptResult(r.States[i], r.Energies[i])
	}

	results := &Results{
		Runtime:       r.Runtime,
		Iterations:    r.Iterations,
		Result:        makeOptResult(r.State, r.Energy),
		ResultHistory: history,
	}

	return results, nil
}
