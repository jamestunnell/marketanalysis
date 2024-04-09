package blocks

import "github.com/jamestunnell/marketanalysis/models"

const (
	DescrHeikinAshi = "Heikin-Ashi bar data values"
	TypeHeikinAshi  = "HeikinAshi"
)

type HeikinAshi struct {
	prev  *models.OHLC
	open  *models.TypedOutput[float64]
	high  *models.TypedOutput[float64]
	low   *models.TypedOutput[float64]
	close *models.TypedOutput[float64]
	hl2   *models.TypedOutput[float64]
	hlc3  *models.TypedOutput[float64]
	occ3  *models.TypedOutput[float64]
	ohlc4 *models.TypedOutput[float64]
	hlcc4 *models.TypedOutput[float64]
}

func NewHeikinAshi() models.Block {
	return &HeikinAshi{
		prev:  nil,
		open:  models.NewTypedOutput[float64](),
		high:  models.NewTypedOutput[float64](),
		low:   models.NewTypedOutput[float64](),
		close: models.NewTypedOutput[float64](),
		hl2:   models.NewTypedOutput[float64](),
		hlc3:  models.NewTypedOutput[float64](),
		occ3:  models.NewTypedOutput[float64](),
		ohlc4: models.NewTypedOutput[float64](),
		hlcc4: models.NewTypedOutput[float64](),
	}
}

func (blk *HeikinAshi) GetType() string {
	return TypeHeikinAshi
}

func (blk *HeikinAshi) GetDescription() string {
	return DescrHeikinAshi
}

func (blk *HeikinAshi) GetParams() models.Params {
	return models.Params{}
}

func (blk *HeikinAshi) GetInputs() models.Inputs {
	return models.Inputs{}
}

func (blk *HeikinAshi) GetOutputs() models.Outputs {
	return models.Outputs{
		NameOpen:  blk.open,
		NameHigh:  blk.high,
		NameLow:   blk.low,
		NameClose: blk.close,
		NameHL2:   blk.hl2,
		NameHLC3:  blk.hlc3,
		NameOCC3:  blk.occ3,
		NameOHLC4: blk.ohlc4,
		NameHLCC4: blk.hlcc4,
	}
}

func (blk *HeikinAshi) IsWarm() bool {
	return blk.prev != nil
}

func (blk *HeikinAshi) Init() error {
	return nil
}

func (blk *HeikinAshi) Update(cur *models.Bar) {
	defer blk.updatePrev(cur)

	if blk.prev == nil {
		return
	}

	ha := cur.HeikinAshi(blk.prev)

	blk.open.Set(ha.Open)
	blk.high.Set(ha.High)
	blk.low.Set(ha.Low)
	blk.close.Set(ha.Close)

	blk.hl2.SetIfConnected(ha.HL2)
	blk.hlc3.SetIfConnected(ha.HLC3)
	blk.occ3.SetIfConnected(ha.OCC3)
	blk.ohlc4.SetIfConnected(ha.OHLC4)
	blk.hlcc4.SetIfConnected(ha.HLCC4)
}

func (blk *HeikinAshi) updatePrev(cur *models.Bar) {
	blk.prev = cur.OHLC
}
