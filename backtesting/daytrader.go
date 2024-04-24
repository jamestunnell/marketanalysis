package backtesting

import (
	"fmt"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"

	"github.com/jamestunnell/marketanalysis/models"
)

type DayTrader struct {
	Collection     models.Collection
	Predictor      models.Predictor
	DateController DateController
	Eval           EvalFunc
}

type DateController interface {
	Current() date.Date
	Advance()
	AnyLeft() bool
}

const (
	DayTradeMarketOpenLocalMin  = 390
	DayTradeMarketCloseLocalMin = 780
)

func NewDayTrader(
	c models.Collection,
	p models.Predictor,
	dc DateController,
	eval EvalFunc) *DayTrader {
	return &DayTrader{
		Collection:     c,
		Predictor:      p,
		DateController: dc,
		Eval:           eval,
	}
}

func (t *DayTrader) RunTest() error {
	open := t.DateController.Current().Local().Add(DayTradeMarketOpenLocalMin * time.Minute)
	close := t.DateController.Current().Local().Add(DayTradeMarketCloseLocalMin * time.Minute)
	ts := timespan.NewTimeSpan(open, close)

	bars, err := t.Collection.LoadBars(ts)
	if err != nil {
		return fmt.Errorf("failed to load bars: %w", err)
	}

	if len(bars) == 0 {
		return nil
	}

	err = EvaluatePredictor(t.Predictor, bars, t.Eval)
	if err != nil {
		return fmt.Errorf("failed to evaluate predictor: %w", err)
	}

	return nil
}

func (t *DayTrader) Advance() {
	t.DateController.Advance()
}

func (t *DayTrader) AnyLeft() bool {
	return t.DateController.AnyLeft()
}
