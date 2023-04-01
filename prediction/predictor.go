package prediction

import "github.com/patrikeh/go-deep/training"

//go:generate mockgen -destination=mock_prediction/mocks.go . Predictor

type Predictor interface {
	InputCount() int
	OutputCount() int

	Train(examples training.Examples, nIter int) error
	Predict(ins []float64) ([]float64, error)
}

type TrainingElem struct {
	Inputs  []float64
	Outputs []float64
}
