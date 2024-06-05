package dmi

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type DMI struct {
	period              *blocks.TypedParam[int]
	dmi                 *indicators.DMI
	pdi, ndi, trend, dx *blocks.TypedOutput[float64]
}

const (
	Type  = "DMI"
	Descr = "Directional Movement Index"

	NamePDI = "pdi"
	NameNDI = "ndi"
	NameDX  = "dx"
)

func New() blocks.Block {
	return &DMI{
		period: blocks.NewTypedParam(10, blocks.NewInclusiveMin(1)),
		dmi:    nil,
		pdi:    blocks.NewTypedOutput[float64](),
		ndi:    blocks.NewTypedOutput[float64](),
		trend:  blocks.NewTypedOutput[float64](),
		dx:     blocks.NewTypedOutput[float64](),
	}
}

func (blk *DMI) GetType() string {
	return Type
}

func (blk *DMI) GetDescription() string {
	return Descr
}

func (blk *DMI) GetParams() blocks.Params {
	return blocks.Params{
		blocks.NamePeriod: blk.period,
	}
}

func (blk *DMI) GetInputs() blocks.Inputs {
	return blocks.Inputs{}
}

func (blk *DMI) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameDX:           blk.dx,
		NameNDI:          blk.ndi,
		NamePDI:          blk.pdi,
		blocks.NameTrend: blk.trend,
	}
}

func (blk *DMI) GetWarmupPeriod() int {
	return blk.dmi.WarmupPeriod()
}

func (blk *DMI) IsWarm() bool {
	return blk.dmi.Warm()
}

func (blk *DMI) Init() error {
	blk.dmi = indicators.NewDMI(blk.period.CurrentVal)

	return nil
}

func (blk *DMI) Update(cur *models.Bar) {
	blk.dmi.Update(cur)

	if !blk.dmi.Warm() {
		return
	}

	blk.dx.SetValue(blk.dmi.DX())
	blk.ndi.SetValue(blk.dmi.NDI())
	blk.pdi.SetValue(blk.dmi.PDI())
	blk.trend.SetValue(blk.dmi.Trend())
}
