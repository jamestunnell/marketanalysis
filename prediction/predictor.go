package prediction

//go:generate mockgen -destination=mock_prediction/mocks.go . Predictor

type Predictor interface {
	Train(elems []*TrainingElem)
}

type TrainingElem struct {
	Inputs  []float64
	Outputs []float64
}