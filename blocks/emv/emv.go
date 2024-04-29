package emv

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMV struct {
	emv *indicators.EMV
	out *blocks.TypedOutput[float64]
}

const (
	Type  = "EMV"
	Descr = "Ease of Movement Value"
)

func New() blocks.Block {
	return &EMV{
		emv: nil,
		out: blocks.NewTypedOutput[float64](),
	}
}

func (blk *EMV) GetType() string {
	return Type
}

func (blk *EMV) GetDescription() string {
	return Descr
}

func (blk *EMV) GetParams() blocks.Params {
	return blocks.Params{}
}

func (blk *EMV) GetInputs() blocks.Inputs {
	return blocks.Inputs{}
}

func (blk *EMV) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *EMV) GetWarmupPeriod() int {
	return blk.emv.WarmupPeriod()
}

func (blk *EMV) IsWarm() bool {
	return blk.emv.Warm()
}

func (blk *EMV) Init() error {
	blk.emv = indicators.NewEMV()

	return nil
}

func (blk *EMV) Update(cur *models.Bar) {
	blk.emv.Update(cur)

	if !blk.emv.Warm() {
		return
	}

	blk.out.SetValue(blk.emv.EMV())
}
