package prediction

import (
	"errors"
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models/bar"
	"github.com/patrikeh/go-deep/training"
)

type BarPredictorCore struct {
	barDur    time.Duration
	nPrevBars int
	ATR       *indicators.ATR
	prev      []*BarMeasure
	warm      bool
}

type BarPredictor struct {
	*BarPredictorCore
	predictor Predictor
}

var errNotWarmedUp = errors.New("not warmed up")

func NewBarPredictor(
	barDur time.Duration,
	atrLen int,
	p Predictor) (*BarPredictor, error) {
	const nBarMeasureFeatures = 3

	nIn := p.InputCount()
	if (nIn % nBarMeasureFeatures) != 0 {
		return nil, fmt.Errorf("predictor input count %d is not a multiple of 3", nIn)
	}

	nOut := p.OutputCount()
	if nOut != nBarMeasureFeatures {
		return nil, fmt.Errorf("predictor output count %d is not 3", nOut)
	}

	nInBars := nIn / nBarMeasureFeatures
	nPrevBars := nInBars - 1
	if nPrevBars <= 0 {
		return nil, fmt.Errorf("prev bar count %d is not positive", nPrevBars)
	}

	core := &BarPredictorCore{
		barDur:    barDur,
		nPrevBars: nPrevBars,
		ATR:       indicators.NewATR(atrLen),
		prev:      []*BarMeasure{},
		warm:      false,
	}

	bp := &BarPredictor{
		BarPredictorCore: core,
		predictor:        p,
	}

	return bp, nil
}

func (bp *BarPredictorCore) WarmUp(bars []*bar.Bar) error {
	wp := bp.TotalWarmupPeriod()
	if len(bars) != wp {
		return commonerrs.NewErrExactBarCount("warmup", wp, len(bars))
	}

	atrWP := bp.ATR.WarmupPeriod()
	atrWarmupBars := bars[:atrWP]

	err := bp.ATR.WarmUp(atrWarmupBars)
	if err != nil {
		err = fmt.Errorf("failed to init ATR indicator: %w", err)

		return err
	}

	for _, bar := range bars[atrWP:] {
		m := NewBarMeasure(bar, bp.ATR.Current())

		bp.prev = append(bp.prev, m)

		_ = bp.ATR.Update(bar)
	}

	bp.warm = true

	return nil
}

func (bp *BarPredictorCore) TotalWarmupPeriod() int {
	return bp.ATR.WarmupPeriod() + bp.nPrevBars
}

func (bp *BarPredictor) Train(bars []*bar.Bar, nIter int) error {
	trainingCore := &BarPredictorCore{
		barDur:    bp.barDur,
		nPrevBars: bp.nPrevBars,
		ATR:       indicators.NewATR(bp.ATR.Length()),
		prev:      []*BarMeasure{},
		warm:      false,
	}

	wp := trainingCore.TotalWarmupPeriod()
	if len(bars) < (wp + 1) {
		return commonerrs.NewErrMinBarCount("warmup+training", wp+1, len(bars))
	}

	warmupBars := bars[:wp]
	trainingBars := bars[wp:]

	if err := trainingCore.WarmUp(warmupBars); err != nil {
		return fmt.Errorf("failed to ")
	}

	examples := training.Examples{}
	cur := NewBarMeasure(trainingBars[0], trainingCore.ATR.Current())

	for i := 1; i < len(trainingBars); i++ {
		ins := combine(trainingCore.prev, cur)
		atr := trainingCore.ATR.Update(trainingBars[i])
		pred := NewBarMeasure(trainingBars[i], atr)

		example := training.Example{
			Input:    ins,
			Response: pred.ToFloat64s(),
		}

		examples = append(examples, example)

		// shift prev
		for i := 0; i < (trainingCore.nPrevBars - 1); i++ {
			trainingCore.prev[i] = trainingCore.prev[i+1]
		}
		trainingCore.prev[trainingCore.nPrevBars-1] = cur
		cur = pred
	}

	return bp.predictor.Train(examples, nIter)
}

func (bp *BarPredictor) Predict(curBar *bar.Bar) (*bar.Bar, error) {
	if !bp.warm {
		return nil, errNotWarmedUp
	}

	atr := bp.ATR.Current()
	cur := NewBarMeasure(curBar, atr)
	ins := combine(bp.prev, cur)

	outs, err := bp.predictor.Predict(ins)
	if err != nil {
		return nil, fmt.Errorf("failed to predict: %w", err)
	}

	predM := &BarMeasure{
		Body:   outs[0],
		Top:    outs[1],
		Bottom: outs[2],
	}

	atr = bp.ATR.Update(curBar)

	tNext := curBar.Timestamp.Add(bp.barDur)
	predBar := bar.NewFromOHLC(tNext, predM.ToOHLC(atr, curBar.Close))

	// shift prev
	for i := 0; i < (bp.nPrevBars - 1); i++ {
		bp.prev[i] = bp.prev[i+1]
	}
	bp.prev[bp.nPrevBars-1] = cur

	return predBar, nil
}

func combine(prev []*BarMeasure, cur *BarMeasure) []float64 {
	all := []float64{}

	for _, m := range prev {
		all = append(all, m.ToFloat64s()...)
	}

	all = append(all, cur.ToFloat64s()...)

	return all
}
