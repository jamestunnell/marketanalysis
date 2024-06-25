package bar

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
)

const (
	Type  = "Bar"
	Descr = "Bar data values"

	NameOpen  = "open"
	NameHigh  = "high"
	NameLow   = "low"
	NameClose = "close"

	NameHL2   = "hl2"
	NameHLC3  = "hlc3"
	NameOCC3  = "occ3"
	NameOHLC4 = "ohlc4"
	NameHLCC4 = "hlcc4"

	NameVWAP = "vwap"
)

type Bar struct {
	open  *blocks.TypedOutput[float64]
	high  *blocks.TypedOutput[float64]
	low   *blocks.TypedOutput[float64]
	close *blocks.TypedOutput[float64]
	hl2   *blocks.TypedOutput[float64]
	hlc3  *blocks.TypedOutput[float64]
	occ3  *blocks.TypedOutput[float64]
	ohlc4 *blocks.TypedOutput[float64]
	hlcc4 *blocks.TypedOutput[float64]
	vwap  *blocks.TypedOutput[float64]
}

func New() blocks.Block {
	return &Bar{
		open:  blocks.NewTypedOutput[float64](),
		high:  blocks.NewTypedOutput[float64](),
		low:   blocks.NewTypedOutput[float64](),
		close: blocks.NewTypedOutput[float64](),
		hl2:   blocks.NewTypedOutput[float64](),
		hlc3:  blocks.NewTypedOutput[float64](),
		occ3:  blocks.NewTypedOutput[float64](),
		ohlc4: blocks.NewTypedOutput[float64](),
		hlcc4: blocks.NewTypedOutput[float64](),
		vwap:  blocks.NewTypedOutput[float64](),
	}
}

func (blk *Bar) GetType() string {
	return Type
}

func (blk *Bar) GetDescription() string {
	return Descr
}

func (blk *Bar) GetParams() models.Params {
	return models.Params{}
}

func (blk *Bar) GetInputs() blocks.Inputs {
	return blocks.Inputs{}
}

func (blk *Bar) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameOpen:  blk.open,
		NameHigh:  blk.high,
		NameLow:   blk.low,
		NameClose: blk.close,
		NameHL2:   blk.hl2,
		NameHLC3:  blk.hlc3,
		NameOCC3:  blk.occ3,
		NameOHLC4: blk.ohlc4,
		NameHLCC4: blk.hlcc4,
		NameVWAP:  blk.vwap,
	}
}

func (blk *Bar) GetWarmupPeriod() int {
	return 0
}

func (blk *Bar) IsWarm() bool {
	return true
}

func (blk *Bar) Init() error {
	return nil
}

func (blk *Bar) Update(cur *models.Bar, isLast bool) {
	blk.close.SetValue(cur.Close)

	blk.open.SetIfConnected(cur.GetOpen)
	blk.high.SetIfConnected(cur.GetHigh)
	blk.low.SetIfConnected(cur.GetLow)

	blk.hl2.SetIfConnected(cur.HL2)
	blk.hlc3.SetIfConnected(cur.HLC3)
	blk.occ3.SetIfConnected(cur.OCC3)
	blk.ohlc4.SetIfConnected(cur.OHLC4)
	blk.hlcc4.SetIfConnected(cur.HLCC4)

	blk.vwap.SetIfConnected(cur.GetVWAP)
}
