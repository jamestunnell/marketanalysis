package prediction

// import (
// 	"github.com/jamestunnell/marketanalysis/models"
// 	"github.com/jamestunnell/marketanalysis/processing"
// 	"github.com/patrikeh/go-deep"
// 	"github.com/patrikeh/go-deep/training"
// )

// type DeepPredictor struct {
// 	inCount, outCount int
// 	examples          training.Examples
// 	nIter             int
// 	deepNeuralNet     *deep.Neural
// }

// func NewDeepPredictor(cfg *deep.Config) *DeepPredictor {
// 	return &DeepPredictor{
// 		deepNeuralNet: deep.NewNeural(cfg),
// 		inCount:       cfg.Inputs,
// 		outCount:      cfg.Layout[len(cfg.Layout)-1],
// 	}
// }

// func (dp *DeepPredictor) InputCount() int {
// 	return dp.inCount
// }

// func (dp *DeepPredictor) OutputCount() int {
// 	return dp.outCount
// }

// func (dp *DeepPredictor) Train(chain *processing.Chain,
// 	provider provision.BarProvider) error {
// 	// params: learning rate, momentum, alpha decay, nesterov
// 	optimizer := training.NewSGD(0.05, 0.1, 1e-6, true)
// 	// params: optimizer, verbosity (print stats at every 50th iteration)
// 	trainer := training.NewTrainer(optimizer, 50)

// 	training, heldout := examples.Split(0.5)
// 	trainer.Train(dp.deepNeuralNet, training, heldout, nIter) // training, validation, iterations

// 	return nil
// }

// func (dp *DeepPredictor) Test(chain *processing.Chain,
// 	provider provision.BarProvider) (models.Positions, error) {
// }

// func (dp *DeepPredictor) Predict(ins []float64) ([]float64, error) {
// 	return dp.deepNeuralNet.Predict(ins), nil
// }
