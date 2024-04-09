package blocks

import "github.com/jamestunnell/marketanalysis/models"

const (
	DescrBar = "Bar data values"
	TypeBar  = "Bar"
)

type Bar struct {
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

func NewBar() models.Block {
	return &Bar{
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

func (blk *Bar) GetType() string {
	return TypeBar
}

func (blk *Bar) GetDescription() string {
	return DescrBar
}

func (blk *Bar) GetParams() models.Params {
	return models.Params{}
}

func (blk *Bar) GetInputs() models.Inputs {
	return models.Inputs{}
}

func (blk *Bar) GetOutputs() models.Outputs {
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

func (blk *Bar) IsWarm() bool {
	return true
}

func (blk *Bar) Init() error {
	return nil
}

func (blk *Bar) Update(cur *models.Bar) {
	blk.open.Set(cur.Open)
	blk.high.Set(cur.High)
	blk.low.Set(cur.Low)
	blk.close.Set(cur.Close)

	blk.hl2.SetIfConnected(cur.HL2)
	blk.hlc3.SetIfConnected(cur.HLC3)
	blk.occ3.SetIfConnected(cur.OCC3)
	blk.ohlc4.SetIfConnected(cur.OHLC4)
	blk.hlcc4.SetIfConnected(cur.HLCC4)
}
