package linregression

type Predictor interface {
	PredictOne(ins []float64) float64
}
