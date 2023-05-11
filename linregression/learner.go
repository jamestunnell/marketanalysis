package linregression

type Learner interface {
	Learn(src DataSource, alpha float64, numIter int) (Predictor, error)
}
