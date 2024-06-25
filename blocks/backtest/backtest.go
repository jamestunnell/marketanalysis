package backtest

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Backtest struct {
	source    *blocks.TypedInput[float64]
	predictor *blocks.TypedInput[float64]
	equity    *blocks.TypedOutput[float64]
	threshold *models.FloatParam

	firstUpdate   bool
	currentEquity float64
	dir           models.Direction
	position      *models.Position
}

const (
	Type  = "Backtest"
	Descr = "Backtest a source using a predictor."

	NameEquity    = "equity"
	NameSource    = "source"
	NamePredictor = "predictor"
	NameThreshold = "threshold"
)

func New() blocks.Block {
	return &Backtest{
		source:    blocks.NewTypedInput[float64](),
		predictor: blocks.NewTypedInput[float64](),
		equity:    blocks.NewTypedOutput[float64](),
		threshold: models.NewFloatParam(0.5, models.NewRangeExcl(0.0, 1.0)),
	}
}

func (blk *Backtest) GetType() string {
	return Type
}

func (blk *Backtest) GetDescription() string {
	return Descr
}

func (blk *Backtest) GetParams() models.Params {
	return models.Params{
		NameThreshold: blk.threshold,
	}
}

func (blk *Backtest) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		NameSource:    blk.source,
		NamePredictor: blk.predictor,
	}
}

func (blk *Backtest) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameEquity: blk.equity,
	}
}

func (blk *Backtest) GetWarmupPeriod() int {
	return 0
}

func (blk *Backtest) IsWarm() bool {
	return true
}

func (blk *Backtest) Init() error {
	blk.currentEquity = 0.0
	blk.dir = models.DirDown
	blk.position = nil
	blk.firstUpdate = true

	return nil
}

func (blk *Backtest) Update(cur *models.Bar, isLast bool) {
	if !blk.source.IsValueSet() || !blk.predictor.IsValueSet() {
		return
	}

	defer func() {
		blk.equity.SetValue(blk.currentEquity)
	}()

	t := cur.Timestamp
	source := blk.source.GetValue()
	pred := blk.predictor.GetValue()
	thresh := blk.threshold.CurrentVal

	if isLast && (blk.position != nil) {
		blk.closePosition(t, source, "end of run")

		return
	}

	prevDir := blk.dir

	switch blk.dir {
	case models.DirDown:
		if pred < -thresh {
			break
		}

		if pred > thresh {
			blk.dir = models.DirUp

			break
		}

		blk.dir = models.DirNone
	case models.DirUp:
		if pred > thresh {
			break
		}

		if pred < -thresh {
			blk.dir = models.DirDown

			break
		}

		blk.dir = models.DirNone
	case models.DirNone:
		if pred > thresh {
			blk.dir = models.DirUp
		} else if pred < -thresh {
			blk.dir = models.DirDown
		}
	}

	if blk.dir == prevDir {
		// if position != nil {
		// 	pl, _ := position.OpenProfitLoss(rSource.Value)
		// 	if pl < -0.5 {
		// 		closePosition(t, rSource.Value, "stop loss")
		// 	} else if pl > 1.0 {
		// 		closePosition(t, rSource.Value, "take profit")
		// 	}
		// }

		return
	}

	if blk.position != nil {
		reason := fmt.Sprintf("dir changed from %s", blk.dir)

		blk.closePosition(t, source, reason)
	}

	switch blk.dir {
	case models.DirUp:
		blk.position = models.NewLongPosition(t, source)
	case models.DirDown:
		blk.position = models.NewShortPosition(t, source)
	}
}

func (blk *Backtest) closePosition(t time.Time, value float64, reason string) {
	blk.position.Close(t, value, reason)

	blk.currentEquity += blk.position.ClosedPL

	log.Trace().
		Stringer("entryTime", blk.position.Entry.Time).
		Stringer("exitTime", blk.position.Exit.Time).
		Float64("profitLoss", blk.position.ClosedPL).
		Str("reason", blk.position.ExitReason).
		Msg("closed position")

	blk.position = nil
}
