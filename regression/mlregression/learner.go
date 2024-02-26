package mlregression

type Learner interface {
	Learn(src Data, alpha float64, numIter int) (Predictor, error)
}
