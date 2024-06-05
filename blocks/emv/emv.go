package emv

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMV struct {
	emv     *indicators.EMV
	current *blocks.TypedOutput[float64]
	average *blocks.TypedOutput[float64]

	period *blocks.TypedParam[int]
	scale  *blocks.TypedParam[float64]
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
		period:  blocks.NewTypedParam(10, blocks.NewGreaterEqual(1)),
		scale:   blocks.NewTypedParam(100000.0, blocks.NewGreater(0.0)),
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
	emv, err := indicators.NewEMV(blk.period.CurrentVal, blk.scale.CurrentVal)
	if err != nil {
		return fmt.Errorf("failed to make EMV indicator: %w", err)
	}

	blk.emv = emv

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
