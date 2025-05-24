package optimization

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/ccssmnn/hego"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Results[T any] struct {
	Last       *Result[T]    `json:"last"`
	Best       *Result[T]    `json:"best"`
	Runtime    time.Duration `json:"runtime"`
	Iterations int           `json:"iterations"`
	History    []*Result[T]  `json:"resultHistory"`
}

type Result[T any] struct {
	Value T       `json:"value"`
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
	resultHook func(*Result[models.ParamVals]),
) (*Results[models.ParamVals], error) {
	if settings.Algorithm != AlgorithmSA {
		err := errors.New("unsupported optimization algorithm " + settings.Algorithm)

		return nil, err
	}

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	return OptimizeSA(settings, objective, resultHook, NewParameterState(rng, values))
}

func OptimizeSA[T any](
	settings *Settings,
	objective Objective[T],
	resultHook func(*Result[T]),
	base State[T],
) (*Results[T], error) {
	saSettings := hego.SASettings{
		Temperature:     10.0,
		AnnealingFactor: 0.999,
		Settings: hego.Settings{
			MaxIterations: settings.MaxIterations,
			KeepHistory:   settings.KeepHistory,
		},
	}
	history := []*Result[T]{}
	initialState := &SAState[T]{
		Base:      base,
		Objective: objective,
		MeasureHook: func(val T, score float64) {
			result := &Result[T]{
				Value: val,
				Score: score,
			}

			resultHook(result)

			history = append(history, result)
		},
	}

	log.Debug().Interface("initial state", initialState).Msg("optimizeSA: starting")

	r, err := hego.SA(initialState, saSettings)
	if err != nil {
		return nil, fmt.Errorf("simulated annealing failed: %w", err)
	}

	log.Debug().Msg("optimizeSA: complete")

	results := &Results[T]{
		Runtime:    r.Runtime,
		Iterations: r.Iterations,
		Last:       history[len(history)-1],
		Best: slices.MinFunc(history, func(a, b *Result[T]) int {
			return cmp.Compare(a.Score, b.Score)
		}),
		History: history,
	}

	return results, nil
}
