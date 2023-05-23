package linregression

import "errors"

type Line struct {
	Slope, Intercept float64
}

var (
	errLenMismatch = errors.New("input len mismatch")
	errNoInput     = errors.New("no inputs")
)

func LinearRegression(xs, ys []float64) (*Line, error) {
	n := len(xs)
	if n != len(ys) {
		return nil, errLenMismatch
	}

	if n == 0 {
		return nil, errNoInput
	}

	// Placeholder for the math to be done
	var sum [5]float64

	// Loop over data keeping index in place
	i := 0
	for ; i < n; i++ {
		sum[0] += xs[i]
		sum[1] += ys[i]
		sum[2] += xs[i] * xs[i]
		sum[3] += xs[i] * ys[i]
		sum[4] += ys[i] * ys[i]
	}

	// Find gradient and intercept
	f := float64(i)
	slope := (f*sum[3] - sum[0]*sum[1]) / (f*sum[2] - sum[0]*sum[0])
	line := &Line{
		Slope:     slope,
		Intercept: (sum[1] / f) - (slope * sum[0] / f),
	}

	return line, nil
}
