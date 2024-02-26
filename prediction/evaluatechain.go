package prediction

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/processing"
	"github.com/jamestunnell/marketanalysis/regression/linregression"
	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type EvalChainStepResult struct {
	BarIndex                               int
	InSlope, OutSlope, OutVal, FutureSlope float64
}

type EvalChainStepFunc func(*EvalChainStepResult)

func EvaluateChain(
	chain *processing.Chain,
	bars models.Bars,
	inHorizon, outHorizon, futureHorizon int,
	step EvalChainStepFunc) error {
	if err := chain.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize chain: %w", err)
	}

	inBuf := buffer.NewCircularBuffer[float64](inHorizon)
	outBuf := buffer.NewCircularBuffer[float64](outHorizon)
	result := &EvalChainStepResult{}

	for i, b := range bars {
		chain.Update(b)

		if !chain.SourceWarm() {
			continue
		}

		inBuf.Add(chain.SourceOutput())

		if !chain.ProcsWarm() {
			continue
		}

		outBuf.Add(chain.ProcsOutput())

		if !inBuf.Full() || !outBuf.Full() {
			continue
		}

		futureBars := bars.NextN(i, futureHorizon)
		if len(futureBars) != futureHorizon {
			continue
		}

		futurePrices := sliceutils.Map(
			futureBars, func(b *models.Bar) float64 { return b.Close })

		inSlope, err := linregression.Slope(inBuf.Array())
		if err != nil {
			return fmt.Errorf("input linear regression failed: %w", err)
		}

		outSlope, err := linregression.Slope(outBuf.Array())
		if err != nil {
			return fmt.Errorf("output linear regression failed: %w", err)
		}

		futureSlope, err := linregression.Slope(futurePrices)
		if err != nil {
			return fmt.Errorf("future linear regression failed: %w", err)
		}

		result.BarIndex = i
		result.InSlope = inSlope
		result.OutSlope = outSlope
		result.OutVal = chain.ProcsOutput()
		result.FutureSlope = futureSlope

		step(result)
	}

	return nil
}
