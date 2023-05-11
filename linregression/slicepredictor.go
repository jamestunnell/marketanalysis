package linregression

type SlicePredictor struct {
	NumInputs int
	Theta     []float64
	M         []float64
	S         []float64
}

func (p *SlicePredictor) PredictOne(ins []float64) float64 {
	// normalize inputs
	for i := 0; i < p.NumInputs; i++ {
		ins[i] = (ins[i] - p.M[i]) / p.S[i]
	}

	return ComputeHypothesis(ins, p.Theta)
}
