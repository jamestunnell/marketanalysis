package predictors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type MAOrdering struct {
	direction  models.Direction
	params     models.Params
	maOrdering *indicators.MAOrdering
	threshold  float64
}

const (
	ParamMinPeriod  = "minPeriod"
	ParamMaxPeriod  = "maxPeriod"
	ParamNumPeriods = "numPeriods"
	ParamSignalLen  = "signalLen"
	ParamThreshold  = "threshold"

	TypeMAOrdering = "MAOrdering"
)

func NewMAOrdering(params models.Params) (models.Predictor, error) {
	periodMin, err := params.GetInt(ParamMinPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get period min param: %w", err)
	}

	periodMax, err := params.GetInt(ParamMaxPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get period max param: %w", err)
	}

	nPeriods, err := params.GetInt(ParamNumPeriods)
	if err != nil {
		return nil, fmt.Errorf("failed to get num periods param: %w", err)
	}

	signalLen, err := params.GetInt(ParamSignalLen)
	if err != nil {
		return nil, fmt.Errorf("failed to get signal len param: %w", err)
	}

	threshold, err := params.GetFloat(ParamThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get threshold param: %w", err)
	}

	maOrdering, err := indicators.NewMAOrdering(periodMin, periodMax, nPeriods, signalLen)
	if err != nil {
		return nil, fmt.Errorf("failed to make MA ordering indicator: %w", err)
	}

	if threshold < 0.0 || threshold > 1.0 {
		return nil, commonerrs.NewErrOutOfRange("threshold", threshold, 0.0, 1.0)
	}

	mao := &MAOrdering{
		maOrdering: maOrdering,
		threshold:  threshold,
	}

	return mao, nil
}

func (mao *MAOrdering) Type() string {
	return TypeMAOrdering
}

func (mao *MAOrdering) Params() models.Params {
	return mao.params
}

func (mao *MAOrdering) WarmupPeriod() int {
	return mao.maOrdering.WarmupPeriod()
}

func (mao *MAOrdering) WarmUp(bars models.Bars) error {
	if err := mao.maOrdering.WarmUp(bars.ClosePrices()); err != nil {
		return fmt.Errorf("failed to warm up MA ordering indicator: %w", err)
	}

	return nil
}

func (mao *MAOrdering) Update(bar *models.Bar) {
	mao.maOrdering.Update(bar.Close)

	corr := mao.maOrdering.Correlation()
	if corr > mao.threshold {
		mao.direction = models.DirUp
	} else if corr < -mao.threshold {
		mao.direction = models.DirDown
	} else {
		mao.direction = models.DirNone
	}
}

func (mao *MAOrdering) Direction() models.Direction {
	return mao.direction
}
