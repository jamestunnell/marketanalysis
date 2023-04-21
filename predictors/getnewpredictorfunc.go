package predictors

import "github.com/jamestunnell/marketanalysis/models"

type NewPredictorFunc func() models.Predictor

func GetNewPredictorFunc(
	typ string) (NewPredictorFunc, bool) {
	var newPredictor NewPredictorFunc
	switch typ {
	case TypeMACross:
		newPredictor = NewMACross
	case TypeMAOrdering:
		newPredictor = NewMAOrdering
	// case TypeMAPivot:
	// 	newPredictor = NewMAPivot
	case TypePivot:
		newPredictor = NewPivot
	}

	if newPredictor == nil {
		return nil, false
	}

	return newPredictor, true
}
