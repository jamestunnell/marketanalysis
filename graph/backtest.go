package graph

import (
	"errors"
	"fmt"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/models"
)

func Backtest(
	graphConfig *Configuration,
	symbol string,
	testDate date.Date,
	loc *time.Location,
	loadBars bars.LoadBarsFunc,
	predictor *Address,
	threshold float64,
) (*models.TimeSeries, error) {
	graphConfig.ClearAllRecording()

	log.Debug().
		Stringer("date", testDate).
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

	timeSeries, err := RunDay(graphConfig, symbol, testDate, loc, loadBars)
	if err != nil {
		return nil, fmt.Errorf("failed to run graph on %s: %w", testDate, err)
	}

	sourceQ, found := timeSeries.FindQuantity(sourceBlkCfg.Name + ".close")
	if !found {
		return nil, errors.New("failed to find source quantity")
	}

	predQ, found := timeSeries.FindQuantity(predictor.String())
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

	closePosition := func(t time.Time, value float64, reason string) {
		position.Close(t, value, reason)

		currentEquity += position.ClosedPL

		equityQ.Records = append(equityQ.Records, &models.QuantityRecord{
			Timestamp: t,
			Value:     currentEquity,
		})

		log.Debug().
			Stringer("entryTime", position.Entry.Time).
			Stringer("exitTime", position.Exit.Time).
			Float64("profitLoss", position.ClosedPL).
			Str("reason", position.ExitReason).
			Msg("closed position")

		position = nil
	}

	for i, r := range predQ.Records {
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
			if r.Value < -threshold {
				break
			}

			if r.Value > threshold {
				dir = models.DirUp

				break
			}

			dir = models.DirNone
		case models.DirUp:
			if r.Value > threshold {
				break
			}

			if r.Value < -threshold {
				dir = models.DirDown

				break
			}

			dir = models.DirNone
		case models.DirNone:
			if r.Value > threshold {
				dir = models.DirUp
			} else if r.Value < -threshold {
				dir = models.DirDown
			}
		}

		if i == (len(predQ.Records) - 1) {
			// log.Debug().Msg("last pred record")

			if position != nil {
				closePosition(t, rSource.Value, "end of run")
			}

			continue
		}

		if dir == prevDir {
			// if position != nil {
			// 	pl, _ := position.OpenProfitLoss(rSource.Value)
			// 	if pl < -0.5 {
			// 		closePosition(t, rSource.Value, "stop loss")
			// 	} else if pl > 1.0 {
			// 		closePosition(t, rSource.Value, "take profit")
			// 	}
			// }

			continue
		}

		if position != nil {
			reason := fmt.Sprintf("dir changed from %s", prevDir)

			closePosition(t, rSource.Value, reason)
		}

		switch dir {
		case models.DirUp:
			position = models.NewLongPosition(t, rSource.Value)
		case models.DirDown:
			position = models.NewShortPosition(t, rSource.Value)
		}

		if position != nil {
			// log.Debug().
			// 	Stringer("entryTime", t).
			// 	Str("type", position.Type).
			// 	Msg("opened position")
		}
	}

	timeSeries.AddQuantity(equityQ)

	return timeSeries, nil
}
