package graph

import (
	"context"
	"fmt"
	"slices"

	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/optimization"
)

type SourceQuantity struct {
	Address     *Address `json:"address"`
	Measurement string   `json:"measurement"`
}

type TargetParam struct {
	Address    *Address               `json:"address"`
	Constraint *models.ConstraintInfo `json:"constraint"`
}

type Objective struct {
	eval       func(models.ParamVals) []float64
	reduce     func([]float64) float64
	resultHook func(*optimization.Result)
}

type reduceFunc func([]float64) float64

const (
	ObjectiveMaximizeMean = "MaximizeMean"
	ObjectiveMaximizeSum  = "MaximizeSum"
	ObjectiveMinimizeMean = "MinimizeMean"
	ObjectiveMinimizeSum  = "MinimizeSum"
)

var reduceFuncs = map[string]reduceFunc{
	ObjectiveMaximizeMean: MaximizeMean,
	ObjectiveMaximizeSum:  MaximizeSum,
	ObjectiveMinimizeMean: MinimizeMean,
	ObjectiveMinimizeSum:  MinimizeSum,
}

func OptimizeParameters(
	ctx context.Context,
	cfg *Config,
	days int,
	source *SourceQuantity,
	targets []*TargetParam,
	objectiveType string,
	settings *optimization.Settings,
	load models.LoadBarsFunc,
	resultHook func(*optimization.Result),
) (*optimization.Results, error) {
	blks, errs := cfg.MakeBlocks()
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to make blocks: %w", errs[0])
	}

	values := map[string]optimization.Value{}

	for _, tgt := range targets {
		param, found := blks.FindParam(tgt.Address)
		if !found {
			return nil, fmt.Errorf("failed to find param '%s'", tgt.Address)
		}

		value, err := optimization.MakeValue(param.GetConstraintInfo())
		if err != nil {
			err = fmt.Errorf(
				"failed to make parameter optimization value for target param '%s': %w",
				tgt.Address, err)

			return nil, err
		}

		values[tgt.Address.String()] = value
	}

	reduce, found := reduceFuncs[objectiveType]
	if !found {
		return nil, fmt.Errorf("failed to find reduce function for %s objective", objectiveType)
	}

	eval := func(paramVals models.ParamVals) []float64 {
		mVals, err := EvaluateParameters(ctx, cfg, days, source, paramVals, load)
		if err != nil {
			log.Warn().Err(err).Msg("failed to evaluate param vals")

			return []float64{}
		}

		return mVals
	}
	objective := &Objective{eval: eval, reduce: reduce, resultHook: resultHook}

	return optimization.OptimizeParameters(settings, values, objective)
}

func EvaluateParameters(
	ctx context.Context,
	cfg *Config,
	days int,
	source *SourceQuantity,
	paramVals models.ParamVals,
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

func MaximizeSum(vals []float64) float64 {
	return -sum(vals)
}

func MinimizeSum(vals []float64) float64 {
	return sum(vals)
}

func MaximizeMean(vals []float64) float64 {
	return -sum(vals) / float64(len(vals))
}

func MinimizeMean(vals []float64) float64 {
	return sum(vals) / float64(len(vals))
}

func (obj *Objective) Measure(values optimization.Values) float64 {
	paramVals := models.ParamVals{}
	for name, optVal := range values {
		paramVals[name] = optVal.GetValue()
	}

	score := obj.reduce(obj.eval(paramVals))

	obj.resultHook(&optimization.Result{Score: score, Value: paramVals})

	return score
}

func sum(vals []float64) (sum float64) {
	for _, n := range vals {
		sum += n
	}

	return
}
