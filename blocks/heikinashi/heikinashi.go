package heikinashi

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/bar"
	"github.com/jamestunnell/marketanalysis/models"
)

const (
	Type  = "HeikinAshi"
	Descr = "Heikin-Ashi bar data values"
)

type HeikinAshi struct {
	prev  *models.OHLC
	open  *blocks.TypedOutput[float64]
	high  *blocks.TypedOutput[float64]
	low   *blocks.TypedOutput[float64]
	close *blocks.TypedOutput[float64]
	hl2   *blocks.TypedOutput[float64]
	hlc3  *blocks.TypedOutput[float64]
	occ3  *blocks.TypedOutput[float64]
	ohlc4 *blocks.TypedOutput[float64]
	hlcc4 *blocks.TypedOutput[float64]
}

func New() blocks.Block {
	return &HeikinAshi{
		prev:  nil,
		open:  blocks.NewTypedOutput[float64](),
		high:  blocks.NewTypedOutput[float64](),
		low:   blocks.NewTypedOutput[float64](),
		close: blocks.NewTypedOutput[float64](),
		hl2:   blocks.NewTypedOutput[float64](),
		hlc3:  blocks.NewTypedOutput[float64](),
		occ3:  blocks.NewTypedOutput[float64](),
		ohlc4: blocks.NewTypedOutput[float64](),
		hlcc4: blocks.NewTypedOutput[float64](),
	}
}

func (blk *HeikinAshi) GetType() string {
	return Type
}

func (blk *HeikinAshi) GetDescription() string {
	return Descr
}

func (blk *HeikinAshi) GetParams() models.Params {
	return models.Params{}
}

func (blk *HeikinAshi) GetInputs() blocks.Inputs {
	return blocks.Inputs{}
}

func (blk *HeikinAshi) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		bar.NameOpen:  blk.open,
		bar.NameHigh:  blk.high,
		bar.NameLow:   blk.low,
		bar.NameClose: blk.close,
		bar.NameHL2:   blk.hl2,
		bar.NameHLC3:  blk.hlc3,
		bar.NameOCC3:  blk.occ3,
		bar.NameOHLC4: blk.ohlc4,
		bar.NameHLCC4: blk.hlcc4,
	}
}

func (blk *HeikinAshi) GetWarmupPeriod() int {
	return 1
}

func (blk *HeikinAshi) IsWarm() bool {
	return blk.prev != nil
}

func (blk *HeikinAshi) Init() error {
	return nil
}

func (blk *HeikinAshi) Update(cur *models.Bar, isLast bool) {
	defer blk.updatePrev(cur)

	if blk.prev == nil {
		return
	}

	ha := cur.HeikinAshi(blk.prev)

	blk.open.SetValue(ha.Open)
	blk.high.SetValue(ha.High)
	blk.low.SetValue(ha.Low)
	blk.close.SetValue(ha.Close)

	blk.hl2.SetIfConnected(ha.HL2)
	blk.hlc3.SetIfConnected(ha.HLC3)
	blk.occ3.SetIfConnected(ha.OCC3)
	blk.ohlc4.SetIfConnected(ha.OHLC4)
	blk.hlcc4.SetIfConnected(ha.HLCC4)
}

func (blk *HeikinAshi) updatePrev(cur *models.Bar) {
	blk.prev = cur.OHLC
}
