package models

import (
	"math"
	"slices"

	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
)

const (
	MeasureFirst  = "first"
	MeasureLast   = "last"
	MeasureMin    = "min"
	MeasureMax    = "max"
	MeasureMean   = "mean"
	MeasureStddev = "stddev"
)

type MeasureFunc func([]float64) float64

func GetMeasureFunc(typ string) (MeasureFunc, bool) {
	switch typ {
	case MeasureFirst:
		return First, true
	case MeasureLast:
		return Last, true
	case MeasureMin:
		return slices.Min[[]float64, float64], true
	case MeasureMax:
		return slices.Max[[]float64, float64], true
	case MeasureMean:
		return Mean, true
	case MeasureStddev:
		return Stddev, true
	}

	return nil, false
}

func First(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	return values[0]
}

func Last(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	return values[len(values)-1]
}

func Mean(values []float64) float64 {
	mean, err := stats.Mean(values)
	if err != nil {
		log.Warn().Err(err).Msg("failed to compute mean")
	}

	return mean
}

func Stddev(values []float64) float64 {
	sd, err := stats.StandardDeviation(values)
	if err != nil {
		log.Warn().Err(err).Msg("failed to compute stddev")
	}

	return sd
}
