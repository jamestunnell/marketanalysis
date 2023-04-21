package backtesting

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

type EvalFunc func(dir models.Direction, bars models.Bars)

func EvaluatePredictor(
	pred models.Predictor,
	bars models.Bars,
	eval EvalFunc) error {
	if err := pred.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize predictor: %w", err)
	}

	wp := pred.WarmupPeriod()
	if len(bars) <= wp {
		return commonerrs.NewErrMinCount("bars", len(bars), wp+1)
	}

	wuBars := bars[:wp]
	remBars := bars[wp:]

	if err := pred.WarmUp(wuBars); err != nil {
		return fmt.Errorf("failed to warm up predictor")
	}

	prevDir := pred.Direction()
	predBars := models.Bars{wuBars.Last()}

	for _, bar := range remBars {
		pred.Update(bar)

		predBars = append(predBars, bar)

		dir := pred.Direction()
		if dir != prevDir {
			eval(prevDir, predBars)

			predBars = models.Bars{bar}
			prevDir = dir
		}
	}

	// The last prediction segment ends when we run out of bars,
	// unless there's only one pred bar, which means the direction
	// changed just after the last bar.
	if len(predBars) > 1 {
		eval(prevDir, predBars)
	}

	return nil
}
