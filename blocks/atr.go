package blocks

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type ATR struct {
	period *models.TypedParam[int]
	atr    *indicators.ATR
	out    *models.TypedOutput[float64]
}

const (
	DescrATR = "Average True Range"
	TypeATR  = "ATR"
)

func NewATR() models.Block {
	periodRange := constraints.NewValRange(1, 200)

	return &ATR{
		period: models.NewParam[int](1, periodRange),
		atr:    nil,
		out:    models.NewTypedOutput[float64](),
	}
}

func (blk *ATR) GetType() string {
	return TypeATR
}

func (blk *ATR) GetDescription() string {
	return DescrATR
}

func (blk *ATR) GetParams() models.Params {
	return models.Params{
		NamePeriod: blk.period,
	}
}

func (blk *ATR) GetInputs() models.Inputs {
	return models.Inputs{}
}

func (blk *ATR) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: blk.out,
	}
}

func (blk *ATR) GetWarmupPeriod() int {
	return blk.atr.Period()
}

func (blk *ATR) IsWarm() bool {
	return blk.atr.Warm()
}

func (blk *ATR) Init() error {
	blk.atr = indicators.NewATR(blk.period.Value)

	return nil
}

func (blk *ATR) Update(cur *models.Bar) {
	blk.atr.Update(cur.OHLC)

	if !blk.atr.Warm() {
		return
	}

	blk.out.SetValue(blk.atr.Current())
}
