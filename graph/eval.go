package graph

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/indicators/pivots"
	"github.com/jamestunnell/marketanalysis/models"
)

func EvalSlope(
	graphConfig *Configuration,
	symbol string,
	evalDate date.Date,
	loc *time.Location,
	loadBars bars.LoadBarsFunc,
	source, predictor *Address,
	horizon int,
) (*models.TimeSeries, error) {
	graphConfig.ClearAllRecording()

	log.Debug().
		Stringer("date", evalDate).
		Stringer("source", source).
		Stringer("predictor", predictor).
		Int("horizon", horizon).
		Msg("evaluating graph using slope")

	if err := graphConfig.SetRecording(source); err != nil {
		return nil, fmt.Errorf("failed to set recording for source output: %w", err)
	}

	if err := graphConfig.SetRecording(predictor); err != nil {
		return nil, fmt.Errorf("failed to set recording for predictor output: %w", err)
	}

	timeSeries, err := RunDay(graphConfig, symbol, evalDate, loc, loadBars)
	if err != nil {
		return nil, fmt.Errorf("failed to run graph on %s: %w", evalDate, err)
	}

	sourceQ, found := timeSeries.FindQuantity(source.String())
	if !found {
		return nil, errors.New("failed to find source quantity")
	}

	predQ, found := timeSeries.FindQuantity(predictor.String())
	if !found {
		return nil, errors.New("failed to find predictor quantity")
	}

	slopeQ := &models.Quantity{
		Name:    "Source Future Slope",
		Records: []*models.QuantityRecord{},
	}

	pivotsQ := &models.Quantity{
		Name:    "Source Pivots",
		Records: []*models.QuantityRecord{},
	}

	pivots, err := pivots.New(horizon * 2)
	if err != nil {
		return nil, fmt.Errorf("failed to make pivots indicator: %w", err)
	}

	log.Debug().Msg("eval: finding source pivot points")

	for _, record := range sourceQ.Records {
		added := pivots.Update(record.Timestamp, record.Value)
		if added {
			pivot := pivots.GetLatest()

			log.Debug().
				Stringer("type", pivot.Type).
				Time("timestamp", pivot.Timestamp).
				Float64("value", pivot.Value).
				Msg("found pivot")

			pivotsQ.Records = append(pivotsQ.Records, &models.QuantityRecord{
				Timestamp: pivot.Timestamp,
				Value:     pivot.Value,
			})
		}
	}

	log.Debug().Msg("eval: calculating future source slopes")

	// Find the slope of future values in the window
	lr := indicators.NewLinRegression(horizon)
	maxSlopeMagn := -math.MaxFloat64

	for i, record := range sourceQ.Records {
		lr.Update(record.Value)

		if !lr.Warm() {
			continue
		}

		slopeQ.Records = append(slopeQ.Records, &models.QuantityRecord{
			Timestamp: sourceQ.Records[i-(horizon-1)].Timestamp,
			Value:     lr.Slope(),
		})

		magn := math.Abs(lr.Slope())
		if magn > maxSlopeMagn {
			maxSlopeMagn = magn
		}
	}

	// Normalize slope to [-1,1]
	for _, record := range slopeQ.Records {
		record.Value /= maxSlopeMagn
	}

	evalQ := &models.Quantity{
		Name:    "Predictor Slope Agreement",
		Records: []*models.QuantityRecord{},
	}

	log.Debug().Int("pred records", len(predQ.Records)).Msg("eval: evaluating predictor slope agreement")

	// Evaluate predictor when it crosses threshold
	for _, record := range predQ.Records {
		slope, found := slopeQ.FindRecord(record.Timestamp)
		if !found {
			continue
		}

		evalQ.Records = append(evalQ.Records, &models.QuantityRecord{
			Value:     slope.Value * record.Value,
			Timestamp: record.Timestamp,
		})
	}

	timeSeries.AddQuantity(slopeQ)
	timeSeries.AddQuantity(evalQ)
	timeSeries.AddQuantity(pivotsQ)

	return timeSeries, nil
}
