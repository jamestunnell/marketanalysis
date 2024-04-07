package blocks

import (
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMV struct {
	EMV *indicators.EMV
	in  *models.TypedInput[*models.Bar]
	out *models.TypedOutput[float64]
}

const (
	DescrEMV = "Ease of Movement Value"
	TypeEMV  = "EMV"
)

func NewEMV() models.Block {
	return &EMV{
		EMV: nil,
		in:  models.NewTypedInput[*models.Bar](),
		out: models.NewTypedOutput[float64](),
	}
}

func (blk *EMV) GetType() string {
	return TypeEMV
}

func (blk *EMV) GetDescription() string {
	return DescrEMV
}

func (blk *EMV) GetParams() models.Params {
	return models.Params{}
}

func (blk *EMV) GetInputs() models.Inputs {
	return models.Inputs{
		NameIn: blk.in,
	}
}

func (blk *EMV) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: blk.out,
	}
}

func (blk *EMV) IsWarm() bool {
	return blk.EMV.Warm()
}

func (blk *EMV) Init() error {
	blk.EMV = indicators.NewEMV()

	return nil
}

func (blk *EMV) Update() {
	if !blk.in.IsSet() {
		return
	}

	blk.EMV.Update(blk.in.Get())

	blk.out.Set(blk.EMV.EMV())
}
