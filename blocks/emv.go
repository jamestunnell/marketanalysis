package blocks

import (
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMV struct {
	emv *indicators.EMV
	out *models.TypedOutput[float64]
}

const (
	DescrEMV = "Ease of Movement Value"
	TypeEMV  = "EMV"
)

func NewEMV() models.Block {
	return &EMV{
		emv: nil,
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
	return models.Inputs{}
}

func (blk *EMV) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOut: blk.out,
	}
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

	blk.out.Set(blk.emv.EMV())
}
