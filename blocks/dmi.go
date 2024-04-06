package blocks

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type DMI struct {
	period       *models.TypedParam[int]
	dmi          *indicators.DMI
	in           *models.TypedInput[*models.Bar]
	pdi, ndi, dx *models.TypedOutput[float64]
}

const (
	DescrDMI = "Directional Movement Index"
	NamePDI = "pdi"
	NameNDI = "ndi"
	NameDX = "dx"
	TypeDMI = "DMI"
)

func NewDMI() models.Block {
	periodRange := constraints.NewValRange(1, 200)

	return &DMI{
		period: models.NewParam[int](periodRange),
		dmi:    nil,
		in:     models.NewTypedInput[*models.Bar](),
		pdi: models.NewTypedOutput[float64](),
		ndi: models.NewTypedOutput[float64](),
		dx: models.NewTypedOutput[float64](),
	}
}

func (blk *DMI) GetType() string {
	return TypeDMI
}

func (blk *DMI) GetDescription() string {
	return DescrDMI
}

func (blk *DMI) GetParams() models.Params {
	return models.Params{
		NamePeriod: blk.period,
	}
}

func (blk *DMI) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: blk.in,
	}
}

func (blk *DMI) GetOutputs() models.Outputs {
	return models.Outputs{
		NameDX: blk.dx,
		NameNDI: blk.ndi,
		NamePDI: blk.pdi,
	}
}

func (blk *DMI) IsWarm() bool {
	return blk.dmi.Warm()
}

func (blk *DMI) Init() error {
	blk.dmi = indicators.NewDMI(blk.period.Value)

	return nil
}

func (blk *DMI) Update() {
	if !blk.in.IsSet() {
		return
	}

	blk.dmi.Update(blk.in.Get())

	blk.dx.Set(blk.dmi.DX())
	blk.ndi.Set(blk.dmi.NDI())
	blk.pdi.Set(blk.dmi.PDI())
}
