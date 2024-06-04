package emv

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMV struct {
	emv     *indicators.EMV
	current *blocks.TypedOutput[float64]
	average *blocks.TypedOutput[float64]

	period *blocks.IntRange
	scale  *blocks.FltRange
}

const (
	Type  = "EMV"
	Descr = "Ease of Movement Value"

	NameScale = "scale"

	NameCurrent = "current"
	NameAverage = "average"
)

func New() blocks.Block {
	return &EMV{
		emv:     nil,
		current: blocks.NewTypedOutput[float64](),
		average: blocks.NewTypedOutput[float64](),
		period:  &blocks.IntRange{Default: 14, Min: 1, Max: 1000},
		scale:   &blocks.FltRange{Default: 10000.0, Min: 1000.0, Max: 1000000.0},
	}
}

func (blk *EMV) GetType() string {
	return Type
}

func (blk *EMV) GetDescription() string {
	return Descr
}

func (blk *EMV) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NamePeriod: blk.period,
		NameScale:         blk.scale,
	}
}

func (blk *EMV) GetInputs() blocks.Inputs {
	return blocks.Inputs{}
}

func (blk *EMV) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameCurrent: blk.current,
		NameAverage: blk.average,
	}
}

func (blk *EMV) GetWarmupPeriod() int {
	return blk.emv.WarmupPeriod()
}

func (blk *EMV) IsWarm() bool {
	return blk.emv.FullyWarm()
}

func (blk *EMV) Init() error {
	blk.emv = indicators.NewEMV(blk.period.Value, blk.scale.Value)

	return nil
}

func (blk *EMV) Update(cur *models.Bar) {
	blk.emv.Update(cur)

	if !blk.emv.PartlyWarm() {
		return
	}

	blk.current.SetValue(blk.emv.Current())

	if !blk.emv.FullyWarm() {
		return
	}

	blk.average.SetValue(blk.emv.Average())
}
