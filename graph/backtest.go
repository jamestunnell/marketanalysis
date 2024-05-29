package graph

import (
	"errors"
	"fmt"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/recorders"
)

func Backtest(
	graphConfig *Configuration,
	loader bars.Loader,
	date date.Date,
	predictor *Address,
	threshold float64,
) (*models.TimeSeries, error) {
	graphConfig.ClearAllRecording()

	log.Debug().
		Stringer("date", date).
		Stringer("predictor", predictor).
		Msg("backtesting graph")

	sourceBlkCfg := &BlockConfig{
		Name:      "source-" + nanoid.Must(),
		Type:      "Bar",
		ParamVals: map[string]any{},
		Recording: []string{"close"},
	}

	graphConfig.Blocks = append(graphConfig.Blocks, sourceBlkCfg)

	if err := graphConfig.SetRecording(predictor); err != nil {
		return nil, fmt.Errorf("failed to set recording for predictor output: %w", err)
	}

	sourceAddr := sourceBlkCfg.Name + ".close"
	predAddr := predictor.String()
	recorder := recorders.NewTimeSeries(LocationNewYork)

	if err := RunDay(date, graphConfig, loader, recorder); err != nil {
		return nil, fmt.Errorf("failed to run graph on %s: %w", date, err)
	}

	sourceQ, found := recorder.FindQuantity(sourceAddr)
	if !found {
		return nil, errors.New("failed to find source quantity")
	}

	predQ, found := recorder.FindQuantity(predAddr)
	if !found {
		return nil, errors.New("failed to find predictor quantity")
	}

	currentEquity := 0.0
	equityQ := &models.Quantity{
		Name: "Equity",
		Records: []*models.QuantityRecord{
			{
				Timestamp: sourceQ.Records[0].Timestamp,
				Value:     0.0,
			},
		},
	}
	dir := models.DirNone

	var position *models.Position

	for _, r := range predQ.Records {
		t := r.Timestamp

		rSource, found := sourceQ.FindRecord(t)
		if !found {
			log.Warn().
				Time("timestamp", t).
				Msg("failed to find matching source record")

			continue
		}

		prevDir := dir

		switch dir {
		case models.DirDown:
			if r.Value > -threshold {
				if threshold == 0.0 {
					dir = models.DirUp
				} else {
					dir = models.DirNone
				}
			}
		case models.DirUp:
			if r.Value < threshold {
				if threshold == 0.0 {
					dir = models.DirDown
				} else {
					dir = models.DirNone
				}
			}
		case models.DirNone:
			if r.Value > threshold {
				dir = models.DirUp
			} else if r.Value < -threshold {
				dir = models.DirDown
			}
		}

		if dir == prevDir {
			continue
		}

		if position != nil {
			position.Close(t, rSource.Value, fmt.Sprintf("dir changed from %s", prevDir))

			currentEquity += position.ClosedPL

			equityQ.Records = append(equityQ.Records, &models.QuantityRecord{
				Timestamp: t,
				Value:     currentEquity,
			})

			position = nil
		}

		switch dir {
		case models.DirUp:
			position = models.NewLongPosition(t, rSource.Value)
		case models.DirDown:
			position = models.NewShortPosition(t, rSource.Value)
		}
	}

	ts := &models.TimeSeries{
		Quantities: []*models.Quantity{sourceQ, predQ, equityQ},
	}

	return ts, nil
}
