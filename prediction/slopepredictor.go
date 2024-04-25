package prediction

// import (
// 	"fmt"

// 	"github.com/rs/zerolog/log"
// 	"github.com/sajari/regression"

// 	"github.com/jamestunnell/marketanalysis/models"
// 	"github.com/jamestunnell/marketanalysis/processing"
// 	"github.com/jamestunnell/marketanalysis/provision"
// 	"github.com/jamestunnell/marketanalysis/regression/linregression"
// 	"github.com/jamestunnell/marketanalysis/util/buffer"
// )

// type SlopePredictor struct {
// 	inHorizon, outHorizon, futureHorizon int
// 	regr                                 *regression.Regression
// }

// type evalChainStepResult struct {
// 	Bar                                    *models.Bar
// 	InSlope, OutSlope, OutVal, FutureSlope float64
// }

// type evalChainStepFunc func(*evalChainStepResult)

// func NewSlopePredictor(inHorizon, outHorizon, futureHorizon int) *SlopePredictor {
// 	return &SlopePredictor{
// 		inHorizon:     inHorizon,
// 		outHorizon:    outHorizon,
// 		futureHorizon: futureHorizon,
// 		regr:          nil,
// 	}
// }

// func (pred *SlopePredictor) evaluateChain(
// 	chain *processing.Chain,
// 	bars provision.BarSequence,
// 	step evalChainStepFunc,
// ) error {
// 	inBuf := buffer.NewCircularBuffer[float64](pred.inHorizon)
// 	outBuf := buffer.NewCircularBuffer[float64](pred.outHorizon)
// 	futureBuf := buffer.NewCircularBuffer[float64](pred.futureHorizon)
// 	result := &evalChainStepResult{}

// 	err := processing.Evaluate(chain, bars, func(bar *models.Bar, sourceOut, procsOut float64) error {
// 		futureBuf.Add(bar.Close)

// 		// This makes sure that the future data leads the chain input and output
// 		if !futureBuf.Full() {
// 			return nil
// 		}

// 		inBuf.Add(sourceOut)
// 		outBuf.Add(procsOut)

// 		if !inBuf.Full() || !outBuf.Full() {
// 			return nil
// 		}

// 		inSlope, err := linregression.Slope(inBuf.Array())
// 		if err != nil {
// 			return fmt.Errorf("input linear regression failed: %w", err)
// 		}

// 		outSlope, err := linregression.Slope(outBuf.Array())
// 		if err != nil {
// 			return fmt.Errorf("output linear regression failed: %w", err)
// 		}

// 		futureSlope, err := linregression.Slope(futureBuf.Array())
// 		if err != nil {
// 			return fmt.Errorf("future linear regression failed: %w", err)
// 		}

// 		result.Bar = bar
// 		result.InSlope = inSlope
// 		result.OutSlope = outSlope
// 		result.OutVal = chain.ProcsOutput()
// 		result.FutureSlope = futureSlope

// 		step(result)

// 		return nil
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to evaulate chain: %w", err)
// 	}

// 	return nil
// }

// func (pred *SlopePredictor) Train(
// 	chain *processing.Chain,
// 	seqs provision.BarSequences) error {
// 	dataPoints := regression.DataPoints{}
// 	eachSeq := func(seq provision.BarSequence) error {
// 		step := func(result *evalChainStepResult) {
// 			vars := []float64{result.InSlope, result.OutSlope, result.OutVal}
// 			dp := regression.DataPoint(result.FutureSlope, vars)

// 			dataPoints = append(dataPoints, dp)
// 		}

// 		err := pred.evaluateChain(chain, seq, step)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}

// 	if err := seqs.EachSequence(eachSeq); err != nil {
// 		return fmt.Errorf("each sequence failed: %w", err)
// 	}

// 	r := new(regression.Regression)

// 	r.SetObserved("future slope")
// 	r.SetVar(0, "input slope")
// 	r.SetVar(1, "output slope")
// 	r.SetVar(2, "output val")
// 	r.Train(dataPoints...)
// 	r.Run()

// 	pred.regr = r

// 	return nil
// }

// func (pred *SlopePredictor) Test(
// 	chain *processing.Chain,
// 	seqs provision.BarSequences) (models.Positions, error) {
// 	positions := models.Positions{}

// 	var pos *models.Position

// 	eachSeq := func(seq provision.BarSequence) error {
// 		var bar *models.Bar

// 		step := func(result *evalChainStepResult) {
// 			predFutureSlope, err := pred.PredictFutureSlope(result.InSlope, result.OutSlope, result.OutVal)
// 			if err != nil {
// 				log.Fatal().Err(err).Msg("failed to predict slope")
// 			}

// 			bar = result.Bar

// 			if pos != nil {
// 				if (pos.Type == models.PositionTypeLong && predFutureSlope < 0.0) ||
// 					(pos.Type == models.PositionTypeShort && predFutureSlope > 0.0) {
// 					pos.Close(bar.Timestamp, bar.Close, "direction change")

// 					positions = append(positions, pos)

// 					pos = nil
// 				}
// 			}

// 			if pos == nil {
// 				if predFutureSlope > 0.0 {
// 					pos = models.NewLongPosition(bar.Timestamp, bar.Close)
// 				} else if predFutureSlope < 0.0 {
// 					pos = models.NewShortPosition(bar.Timestamp, bar.Close)
// 				}
// 			}
// 		}

// 		err := pred.evaluateChain(chain, seq, step)
// 		if err != nil {
// 			return err
// 		}

// 		if pos != nil {
// 			pos.Close(bar.Timestamp, bar.Close, "end-of-day")

// 			positions = append(positions, pos)

// 			pos = nil
// 		}

// 		return nil
// 	}

// 	if err := seqs.EachSequence(eachSeq); err != nil {
// 		return models.Positions{}, fmt.Errorf("each sequence failed: %w", err)
// 	}

// 	return positions, nil
// }

// func (pred *SlopePredictor) PredictFutureSlope(inSlope, outSlope, outVal float64) (float64, error) {
// 	predIns := []float64{inSlope, outSlope, outVal}

// 	futureSlope, err := pred.regr.Predict(predIns)
// 	if err != nil {
// 		return 0.0, fmt.Errorf("failed to predict: %w", err)
// 	}

// 	return futureSlope, nil
// }
