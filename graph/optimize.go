package graph

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ccssmnn/hego"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type SourceQuantity struct {
	Address     *Address `json:"address"`
	Measurement string   `json:"measurment"`
}

type TargetParam struct {
	Address *Address `json:"address"`
	Min     any      `json:"min"`
	Max     any      `json:"max"`
}

type OptimizeSettings struct {
	RandomSeed    int64  `json:"randomSeed"`
	Algorithm     string `json:"algorithm"`
	MaxIterations int    `json:"maxIterations"`
	KeepHistory   bool   `json:"keepHistory"`
}

type OptimizeResults struct {
	Result        *OptResult    `json:"result"`
	Runtime       time.Duration `json:"runtime"`
	Iterations    int           `json:"iterations"`
	ResultHistory []*OptResult  `json:"resultHistory"`
}

type OptResult struct {
	ParamVals         map[string]any `json:"paramVals"`
	SourceMeasurement float64        `json:"sourceMeasurement"`
}

func Optimize(
	ctx context.Context,
	cfg *Config,
	days int,
	source *SourceQuantity,
	params []*TargetParam,
	settings *OptimizeSettings,
	load models.LoadBarsFunc,
) (*OptimizeResults, error) {
	if settings.Algorithm != "SimulatedAnnealing" {
		err := errors.New("unsupported optimization algorithm " + settings.Algorithm)

		return nil, err
	}

	return OptimizeSA(ctx, cfg, days, source, params, settings, load)
}

func OptimizeSA(
	ctx context.Context,
	cfg *Config,
	days int,
	source *SourceQuantity,
	targetParams []*TargetParam,
	settings *OptimizeSettings,
	load models.LoadBarsFunc,
) (*OptimizeResults, error) {
	rng := rand.New(rand.NewSource(settings.RandomSeed))
	eval := func(paramVals map[string]any) float64 {
		result, err := EvaluateParamVals(ctx, cfg, days, source, paramVals, load)
		if err != nil {
			log.Warn().Err(err).Msg("failed to evaluate param vals")

			return 0.0
		}

		return result
	}

	blks, errs := cfg.MakeBlocks()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to make blocks: %w", errs[0])
	}

	initialState := NewOptState(rng, eval)
	for _, tgt := range targetParams {
		param, found := blks.FindParam(tgt.Address)
		if !found {
			return nil, fmt.Errorf("failed to find param '%s'", tgt.Address)
		}

		constraint := param.GetConstraint()

		switch param.GetValueType() {
		case "int":
			limits := sliceutils.Map(constraint.GetLimits(), func(v any) int { return v.(int) })
			var mut IntMutator
			if constraint.GetType() == blocks.OneOf {
				mut = NewIntEnumMutator(limits)
			} else {
				mut = NewIntRangeMutator(limits[0], limits[1])
			}

			initialState.Ints[tgt.Address.String()] = NewIntValue(mut, rng)
		case "float64":
			limits := sliceutils.Map(constraint.GetLimits(), func(v any) float64 { return v.(float64) })
			var mut FloatMutator
			if constraint.GetType() == blocks.OneOf {
				mut = NewFloatEnumMutator(limits)
			} else {
				mut = NewFloatRangeMutator(limits[0], limits[1])
			}

			initialState.Floats[tgt.Address.String()] = NewFloatValue(mut, rng)
		default:
			return nil, fmt.Errorf("param %s has unsupported type %s", tgt.Address, param.GetValueType())
		}
	}

	saSettings := hego.SASettings{
		Temperature:     10.0,
		AnnealingFactor: 0.999,
		Settings: hego.Settings{
			MaxIterations: settings.MaxIterations,
			Verbose:       1000,
			KeepHistory:   settings.KeepHistory,
		},
	}

	r, err := hego.SA(initialState, saSettings)
	if err != nil {
		return nil, fmt.Errorf("simulated annealing failed: %w", err)
	}

	makeOptResult := func(state hego.AnnealingState, energy float64) *OptResult {
		return &OptResult{
			ParamVals:         state.(*OptState).ParamVals(),
			SourceMeasurement: energy,
		}
	}
	history := make([]*OptResult, len(r.States))

	for i := 0; i < len(r.States); i++ {
		history[i] = makeOptResult(r.States[i], r.Energies[i])
	}

	results := &OptimizeResults{
		Runtime:       r.Runtime,
		Iterations:    r.Iterations,
		Result:        makeOptResult(r.State, r.Energy),
		ResultHistory: history,
	}

	return results, nil
}

func EvaluateParamVals(
	ctx context.Context,
	cfg *Config,
	days int,
	source *SourceQuantity,
	paramVals map[string]any,
	load models.LoadBarsFunc,
) (float64, error) {
	for addrStr, val := range paramVals {
		addr, err := ParseAddress(addrStr)

		if err != nil {
			return 0.0, fmt.Errorf("failed to parase param address '%s': %w", addrStr, err)
		}

		blkConfig, found := cfg.FindBlock(addr.A)
		if !found {
			return 0.0, fmt.Errorf("failed to find block %s", addr.A)
		}

		blkConfig.SetParamVal(addr.B, val)
	}

	// Add two days for every week, since weekend days will be no-op
	days += 2 * (days / 7)

	startDate := date.Today().Add(date.PeriodOfDays(-days))

	ts, err := RunMultiDaySummary(ctx, cfg, startDate, load)
	if err != nil {
		return 0.0, fmt.Errorf("failed to run multi-day summary: %w", err)
	}

	q, found := ts.FindQuantity(source.Address.String())
	if !found {
		return 0.0, fmt.Errorf("failed to find source quantity '%s' in results", source.Address)
	}

	mVal, found := q.Measurements[source.Measurement]
	if !found {
		return 0.0, fmt.Errorf("failed to find source quantity measurement '%s'", source.Measurement)
	}

	return mVal, nil
}

type OptState struct {
	rng    *rand.Rand
	eval   func(map[string]any) float64
	Ints   map[string]*IntValue
	Floats map[string]*FloatValue
}

func NewOptState(rng *rand.Rand, eval func(map[string]any) float64) *OptState {
	return &OptState{
		eval:   eval,
		rng:    rng,
		Ints:   map[string]*IntValue{},
		Floats: map[string]*FloatValue{},
	}
}

func (s *OptState) AddIntValue(name string, val *IntValue) {
	s.Ints[name] = val
}

func (s *OptState) AddFloatValue(name string, val *FloatValue) {
	s.Floats[name] = val
}

func (s *OptState) ParamVals() map[string]any {
	paramVals := map[string]any{}

	for name, v := range s.Ints {
		paramVals[name] = v.Value()
	}

	for name, v := range s.Floats {
		paramVals[name] = v.Value()
	}

	return paramVals
}

func (s *OptState) Clone() *OptState {
	s2 := NewOptState(s.rng, s.eval)

	for name, v := range s.Ints {
		s2.AddIntValue(name, v.Clone())
	}

	for name, v := range s.Floats {
		s2.AddFloatValue(name, v.Clone())
	}

	return s2
}

func (s *OptState) Mutate() {
	for _, v := range s.Ints {
		v.Mutate(s.rng)
	}

	for _, v := range s.Floats {
		v.Mutate(s.rng)
	}
}

func (s *OptState) Neighbor() hego.AnnealingState {
	n := s.Clone()

	n.Mutate()

	return n
}

// Energy returns the energy of the current state. Lower is better
func (s *OptState) Energy() float64 {
	return s.eval(s.ParamVals())
}
