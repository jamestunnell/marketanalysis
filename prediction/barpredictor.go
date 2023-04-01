package prediction

import (
	"errors"
	"fmt"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models/bar"
	"golang.org/x/exp/slices"
)

type BarPredictor struct {
	depth     int
	atr       *indicators.ATR
	predictor Predictor
}

type warmupResult struct {
	Prev      [][]float64
	Remaining []*bar.Bar
}

func NewBarPredictor(
	depth int,
	atr *indicators.ATR,
	p Predictor) *BarPredictor {
	return &BarPredictor{
		depth:     depth,
		atr:       atr,
		predictor: p,
	}
}

func (bp *BarPredictor) warmup(bars []*bar.Bar) (*warmupResult, error) {
	wp := bp.atr.WarmupPeriod()
	initLen := bp.depth + wp

	if len(bars) <= initLen {
		err := errors.New("too few bars")

		return nil, err
	}

	err := bp.atr.Initialize(bars[:wp])
	if err != nil {
		err = fmt.Errorf("failed to init ATR indicator: %w", err)

		return nil, err
	}

	atr := bp.atr.Current()
	prev := [][]float64{}

	for i := wp; i < initLen; i++ {
		bar := bars[i]
		body, top, bottom := MeasureBar(bar)

		body /= atr
		top /= atr
		bottom /= atr

		prev = append(prev, []float64{body, top, bottom})

		atr = bp.atr.Update(bar)
	}

	result := &warmupResult{
		Prev:      prev,
		Remaining: bars[initLen:],
	}

	return result, nil
}

func (bp *BarPredictor) Train(allBars []*bar.Bar) error {
	result, err := bp.warmup(allBars)
	if err != nil {
		return fmt.Errorf("failed to warm up: %w")
	}

	atr := bp.atr.Current()
	prev := result.Prev
	elems := []*TrainingElem{}

	for _, bar := range result.Remaining {
		body, top, bottom := MeasureBar(bar)

		body /= atr
		top /= atr
		bottom /= atr

		cur := []float64{body, top, bottom}
		elem := &TrainingElem{
			Inputs:  combine(prev),
			Outputs: cur,
		}

		elems = append(elems, elem)

		atr = bp.atr.Update(bar)

		// shift prev
		for i := 0; i < (bp.depth - 1); i++ {
			prev[i] = prev[i+1]
		}
		result.Prev[bp.depth-1] = cur
	}

	bp.predictor.Train(elems)

	return nil
}

func (bp *BarPredictor) Predict(allBars []*bar.Bar) ([]*bar.Bar, error) {
	result, err := bp.warmup(allBars)
	if err != nil {
		return []*bar.Bar{}, fmt.Errorf("failed to warm up: %w")
	}

	prev := result.Prev
	for _, bar := range result.Remaining {
		// TODO
	}

	return []*bar.Bar{}, nil
}

func combine(fSlices [][]float64) []float64 {
	all := []float64{}

	for _, fSlice := range fSlices {
		all = append(all, slices.Clone(fSlice)...)
	}

	return all
}
