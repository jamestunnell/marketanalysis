package mlregression

import (
	"fmt"
)

type SliceLearner struct {
}

func NewSliceLearner() Learner {
	return &SliceLearner{}
}

func (l *SliceLearner) Learn(src Data, alpha float64, numIter int) (Predictor, error) {
	inputs := src.Inputs()
	y := src.Output()
	theta := make([]float64, len(inputs[0])+1)

	// Normalize all the elements to keep an identical scale between different data
	XNorm, M, S, err := Normalize(inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize: %w", err)
	}

	// Perform gradient descent to calculate Theta
	theta, err = LinearGradient(XNorm, y, theta, alpha, numIter, false)
	if err != nil {
		return nil, fmt.Errorf("failed to run linear gradient: %w", err)
	}

	p := &SlicePredictor{
		NumInputs: len(inputs),
		Theta:     theta,
		M:         M,
		S:         S,
	}

	return p, nil
}
