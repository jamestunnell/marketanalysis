package prediction

//go:generate mockgen -destination=mock_prediction/mocks.go . Predictor

type Predictor interface {
	Train(elems []*TrainingElem)
	Predict(ins []float64) []float64
}

type TrainingElem struct {
	Inputs  []float64
	Outputs []float64
}
