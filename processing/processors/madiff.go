package processors

import (
	"sort"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type MADiff struct {
	output  float64
	fastEMA *indicators.EMA
	slowEMA *indicators.EMA
	period1 *models.TypedParam[int]
	period2 *models.TypedParam[int]
}

const (
	Period1Name = "period1"
	Period2Name = "period2"

	TypeMADiff = "MADiff"
)

func NewMADiff() *MADiff {
	periodRange := constraints.NewValRange(1, 200)

	return &MADiff{
		output:  0.0,
		fastEMA: nil,
		slowEMA: nil,
		period1: models.NewParam[int](periodRange),
		period2: models.NewParam[int](periodRange),
	}
}

func (mad *MADiff) Type() string {
	return TypeMADiff
}

func (mad *MADiff) Params() models.Params {
	return models.Params{
		Period1Name: mad.period1,
		Period2Name: mad.period2,
	}
}

func (mad *MADiff) Initialize() error {
	periods := []int{mad.period1.Value, mad.period2.Value}

	sort.Ints(periods)

	mad.output = 0.0
	mad.fastEMA = indicators.NewEMA(periods[0])
	mad.slowEMA = indicators.NewEMA(periods[1])

	return nil
}

func (mad *MADiff) WarmupPeriod() int {
	return mad.slowEMA.Period()
}

func (mad *MADiff) Output() float64 {
	return mad.output
}

func (mad *MADiff) Warm() bool {
	return mad.slowEMA.Warm()
}

func (mad *MADiff) Update(val float64) {
	mad.slowEMA.Update(val)
	mad.fastEMA.Update(val)

	if mad.slowEMA.Warm() {
		mad.output = mad.fastEMA.Current() - mad.slowEMA.Current()
	}
}
