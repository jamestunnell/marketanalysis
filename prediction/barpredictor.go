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

func (bc *BarPredictor) Train(td []*bar.Bar) error {
	wp := bc.atr.WarmupPeriod()
	initLen := bc.depth + wp

	if len(td) <= initLen {
		return errors.New("not enough bars for init and training")
	}

	err := bc.atr.Initialize(td[:wp])
	if err != nil {
		return fmt.Errorf("failed to init ATR indicator: %w", err)
	}

	atr := bc.atr.Current()
	prev := [][]float64{}

	for i := wp; i < initLen; i++ {
		bar := td[i]
		body, top, bottom := MeasureBar(bar)

		body /= atr
		top /= atr
		bottom /= atr

		prev = append(prev, []float64{body, top, bottom})

		atr = bc.atr.Update(bar)
	}

	elems := []*TrainingElem{}
	for i := initLen; i < len(td); i++ {
		bar := td[i]
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

		atr = bc.atr.Update(bar)

		for i := 0; i < (bc.depth - 1); i++ {
			prev[i] = prev[i+1]
		}
		prev[bc.depth-1] = cur
	}

	bc.predictor.Train(elems)

	return nil
}

func combine(fSlices [][]float64) []float64 {
	all := []float64{}

	for _, fSlice := range fSlices {
		all = append(all, slices.Clone(fSlice)...)
	}

	return all
}
