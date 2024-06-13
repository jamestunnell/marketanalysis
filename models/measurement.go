package models

import (
	"slices"

	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
)

const (
	MeasureMin    = "min"
	MeasureMax    = "max"
	MeasureMean   = "mean"
	MeasureStddev = "stddev"
)

type MeasureFunc func([]float64) float64

func GetMeasureFunc(typ string) (MeasureFunc, bool) {
	switch typ {
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
