package graph

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/ccssmnn/hego"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

type SourceQuantity struct {
	Address     *Address `json:"address"`
	Measurement string   `json:"measurement"`
}

type TargetParam struct {
	Address *Address `json:"address"`
	Min     any      `json:"min"`
	Max     any      `json:"max"`
}

type OptimizeSettings struct {
	Algorithm     string `json:"algorithm"`
	Objective     string `json:"objective"`
	RandomSeed    int64  `json:"randomSeed"`
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
	ParamVals map[string]any `json:"paramVals"`
	Score     float64        `json:"score"`
}

const (
	MaximizeSum = "MaximizeSum"
	MinimizeSum = "MinimizeSum"
)

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
		mVals, err := EvaluateParamVals(ctx, cfg, days, source, paramVals, load)
		if err != nil {
			log.Warn().Err(err).Msg("failed to evaluate param vals")

			return 0.0
		}

		// log.Info().Floats64("mVals", mVals).Msg("eval complete")

		switch settings.Objective {
		case MaximizeSum:
			return -sum(mVals)
		case MinimizeSum:
			return sum(mVals)
		}

		log.Warn().Str("objective", settings.Objective).Msg("unknown optimize objective")

		return 0.0
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

		if constraint.GetType() == blocks.OneOf {
			return nil, fmt.Errorf("unsupported constaint type '%s'", constraint.GetType())
		}

		switch param.GetValueType() {
		case "int":
			mut := NewIntRangeMutator(int(tgt.Min.(float64)), int(tgt.Max.(float64)))

			initialState.Ints[tgt.Address.String()] = NewIntValue(mut, rng)
		case "float64":
			mut := NewFloatRangeMutator(tgt.Min.(float64), tgt.Max.(float64))

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
			Verbose:       100,
			KeepHistory:   settings.KeepHistory,
		},
	}

	r, err := hego.SA(initialState, saSettings)
	if err != nil {
		return nil, fmt.Errorf("simulated annealing failed: %w", err)
	}

	makeOptResult := func(state hego.AnnealingState, energy float64) *OptResult {
		return &OptResult{
			ParamVals: state.(*OptState).ParamVals(),
			Score:     energy,
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
) ([]float64, error) {
	for addrStr, val := range paramVals {
		addr, err := ParseAddress(addrStr)

		if err != nil {
			return []float64{}, fmt.Errorf("failed to parse param address '%s': %w", addrStr, err)
		}

		blkConfig, found := cfg.FindBlock(addr.A)
		if !found {
			return []float64{}, fmt.Errorf("failed to find param block '%s'", addr.A)
		}

		blkConfig.SetParamVal(addr.B, val)
	}

	blkConfig, found := cfg.FindBlock(source.Address.A)
	if !found {
		return []float64{}, fmt.Errorf("failed to find source block '%s'", source.Address.A)
	}

	sourceOut, found := blkConfig.FindOutput(source.Address.B)
	if !found {
		return []float64{}, fmt.Errorf("failed to find source output '%s'", source.Address.B)
	}

	if !slices.Contains(sourceOut.Measurements, source.Measurement) {
		sourceOut.Measurements = append(sourceOut.Measurements, source.Measurement)
	}

	// Add two days for every week, since weekend days will be no-op
	days += 2 * (days / 7)

	startDate := date.Today().Add(date.PeriodOfDays(-days))

	summaryTS, err := RunMultiDaySummary(ctx, cfg, startDate, load)
	if err != nil {
		return []float64{}, fmt.Errorf("failed to run multi-day summary: %w", err)
	}

	sourceMeasurementName := source.Address.String() + ":" + source.Measurement

	q, found := summaryTS.FindQuantity(sourceMeasurementName)
	if !found {
		return []float64{}, fmt.Errorf("failed to find source measurement quantity '%s' in results", sourceMeasurementName)
	}

	return q.RecordValues(), nil
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

func sum(vals []float64) (sum float64) {
	for _, n := range vals {
		sum += n
	}

	return
}
