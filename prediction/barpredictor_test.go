package prediction_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models/bar"
	"github.com/jamestunnell/marketanalysis/prediction"
	"github.com/jamestunnell/marketanalysis/prediction/mock_prediction"
	"github.com/stretchr/testify/assert"
)

func TestBarPredictor(t *testing.T) {
	ctrl := gomock.NewController(t)
	p := mock_prediction.NewMockPredictor(ctrl)
	atr := indicators.NewATR(5)

	bp := prediction.NewBarPredictor(2, atr, p)

	// no bars - can't train
	assert.Error(t, bp.Train([]*bar.Bar{}))
}
