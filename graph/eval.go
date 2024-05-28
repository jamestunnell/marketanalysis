package graph

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/recorders"
)

const LocationNewYork = "America/New_York"

type EvalSlopeConfig struct {
	Date      date.Date `json:"date"`
	Window    int       `json:"window"`
	Source    *Address  `json:"source"`
	Predictor *Address  `json:"predictor"`
}

func EvalSlope(
	graphConfig *Configuration,
	loader bars.Loader,
	config *EvalSlopeConfig,
) (*models.TimeSeries, error) {
	graphConfig.ClearAllRecording()

	log.Debug().
		Stringer("date", config.Date).
		Stringer("source", config.Source).
		Stringer("predictor", config.Predictor).
		Int("window", config.Window).
		Msg("evaluating graph using slope")

	if err := graphConfig.SetRecording(config.Source); err != nil {
		return nil, fmt.Errorf("failed to set recording for source output: %w", err)
	}

	if err := graphConfig.SetRecording(config.Predictor); err != nil {
		return nil, fmt.Errorf("failed to set recording for predictor output: %w", err)
	}

	sourceName := config.Source.String()
	predName := config.Predictor.String()
	recorder := recorders.NewTimeSeries(LocationNewYork)

	if err := RunDay(config.Date, graphConfig, loader, recorder); err != nil {
		return nil, fmt.Errorf("failed to run graph on %s: %w", config.Date, err)
	}

	sourceQ, found := recorder.FindQuantity(sourceName)
	if !found {
		return nil, errors.New("failed to find source quantity")
	}

	predQ, found := recorder.FindQuantity(predName)
	if !found {
		return nil, errors.New("failed to find predictor quantity")
	}

	slices.SortFunc(sourceQ.Records, func(a, b *models.QuantityRecord) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	slopeQ := &models.Quantity{
		Name:    "Source Future Slope",
		Records: []*models.QuantityRecord{},
	}

	// Find the slope of future values in the window
	lr := indicators.NewLinRegression(config.Window)
	maxSlopeMagn := -math.MaxFloat64
	for i, record := range sourceQ.Records {
		lr.Update(record.Value)

		if !lr.Warm() {
			continue
		}

		slopeQ.Records = append(slopeQ.Records, &models.QuantityRecord{
			Timestamp: sourceQ.Records[i-(config.Window-1)].Timestamp,
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

	ts := &models.TimeSeries{
		Quantities: []*models.Quantity{sourceQ, predQ, slopeQ, evalQ},
	}

	return ts, nil
}
