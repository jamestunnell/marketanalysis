package evaluation

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/mlregression"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/processing"
	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type SlopeEvaluator struct {
	chain  *processing.Chain
	config *SlopeEvalConfig

	inSlopes     []float64
	inIntercepts []float64
	outSlopes    []float64
}

type SlopeEvalConfig struct {
	InHorizon, OutHorizon int
	MLRAlpha              float64
	MLRIterations         int
}

type mlrData struct {
	Ins  [][]float64
	Outs []float64
}

func NewSlopeEvaluator(
	chain *processing.Chain, config *SlopeEvalConfig) *SlopeEvaluator {
	return &SlopeEvaluator{
		chain:        chain,
		config:       config,
		inSlopes:     []float64{},
		inIntercepts: []float64{},
		outSlopes:    []float64{},
	}
}

func (e *SlopeEvaluator) IncorporateBarSet(bars models.Bars) error {
	if err := e.chain.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize chain: %w", err)
	}

	inBuf := buffer.NewCircularBuffer[float64](e.config.InHorizon)
	outBuf := buffer.NewCircularBuffer[float64](e.config.OutHorizon)
	inXs := sliceutils.New(e.config.InHorizon, func(idx int) float64 {
		return float64(idx)
	})
	outXs := sliceutils.New(e.config.OutHorizon, func(idx int) float64 {
		return float64(idx)
	})

	for _, b := range bars {
		e.chain.Update(b)

		if !e.chain.SourceWarm() {
			continue
		}

		inBuf.Add(e.chain.SourceOutput())

		if !e.chain.ProcsWarm() {
			continue
		}

		outBuf.Add(e.chain.ProcsOutput())

		if !inBuf.Full() || !outBuf.Full() {
			continue
		}

		inLine, ok := indicators.LinearRegression(inXs, inBuf.Array())
		if !ok {
			log.Fatal().Msg("input linear regression failed")
		}

		outLine, ok := indicators.LinearRegression(outXs, outBuf.Array())
		if !ok {
			log.Fatal().Msg("output linear regression failed")
		}

		e.inSlopes = append(e.inSlopes, inLine.Slope)
		e.inIntercepts = append(e.inIntercepts, inLine.Intercept)
		e.outSlopes = append(e.outSlopes, outLine.Slope)
	}

	return nil
}

func (e *SlopeEvaluator) MakePredictor() (mlregression.Predictor, error) {
	l := mlregression.NewSliceLearner()
	d := &mlrData{Ins: [][]float64{e.inSlopes, e.inIntercepts}, Outs: e.outSlopes}

	pred, err := l.Learn(d, e.config.MLRAlpha, e.config.MLRIterations)
	if err != nil {
		return nil, fmt.Errorf("ML regression failed: %w", err)
	}

	return pred, nil
}

func (d *mlrData) Inputs() [][]float64 {
	return d.Ins
}

func (d *mlrData) Outputs() []float64 {
	return d.Outs
}
