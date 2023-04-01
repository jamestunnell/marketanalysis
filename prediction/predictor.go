package prediction

//go:generate mockgen -destination=mock_prediction/mocks.go . Predictor

type Predictor interface {
	InputCount() int
	OutputCount() int
	Trained() bool

	Train(elems []*TrainingElem) error
	Predict(ins []float64) ([]float64, error)
}

type TrainingElem struct {
	Inputs  []float64
	Outputs []float64
}
