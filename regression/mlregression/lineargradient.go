package mlregression

func LinearGradient(data [][]float64, y []float64, theta []float64, alpha float64, num_iters int, printCostFunction bool) ([]float64, error) {
	for i := 0; i < num_iters; i++ {
		//Number of training examples
		m := len(y)
		//Slice helper to calculate our new versions of theta
		thetaTemp := make([]float64, len(theta))

		//Sum (hi-yx)xi
		for rowI := 0; rowI < m; rowI++ {
			hi := ComputeHypothesis(data[rowI], theta)
			sumRowI := computeSumRowI(data[rowI], hi, y[rowI])
			for t := 0; t < len(theta); t++ {
				thetaTemp[t] += sumRowI[t]
			}
		}
		//Update theta
		for t := 0; t < len(theta); t++ {
			theta[t] = theta[t] - (alpha/float64(m))*thetaTemp[t]
		}
	}

	return theta, nil
}

// Multiply by xi only if not theta0
func computeSumRowI(x []float64, hi float64, yi float64) []float64 {
	theta := make([]float64, len(x)+1)
	theta[0] = hi - yi
	for i := 1; i < len(theta); i++ {
		theta[i] = (hi - yi) * x[i-1]
	}
	return theta
}
