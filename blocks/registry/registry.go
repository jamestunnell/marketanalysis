package registry

import (
	"golang.org/x/exp/maps"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/blocks/aroon"
	"github.com/jamestunnell/marketanalysis/blocks/atr"
	"github.com/jamestunnell/marketanalysis/blocks/bar"
	"github.com/jamestunnell/marketanalysis/blocks/dmi"
	"github.com/jamestunnell/marketanalysis/blocks/ema"
	"github.com/jamestunnell/marketanalysis/blocks/emv"
	"github.com/jamestunnell/marketanalysis/blocks/heikinashi"
	"github.com/jamestunnell/marketanalysis/blocks/maorder"
	"github.com/jamestunnell/marketanalysis/blocks/math"
	"github.com/jamestunnell/marketanalysis/blocks/multitrend"
	"github.com/jamestunnell/marketanalysis/blocks/pivots"
	"github.com/jamestunnell/marketanalysis/blocks/sma"
	"github.com/jamestunnell/marketanalysis/blocks/supertrend"
)

type NewBlockFunc func() blocks.Block

type Registry interface {
	Types() []string
	Add(typ string, new NewBlockFunc)
	Get(typ string) (NewBlockFunc, bool)
}

var reg = map[string]NewBlockFunc{}

func init() {
	Add(aroon.Type, aroon.New)
	Add(atr.Type, atr.New)
	Add(bar.Type, bar.New)
	Add(dmi.Type, dmi.New)
	Add(ema.Type, ema.New)
	Add(emv.Type, emv.New)
	Add(heikinashi.Type, heikinashi.New)
	Add(maorder.Type, maorder.New)
	Add(multitrend.Type, multitrend.New)
	Add(pivots.Type, pivots.New)
	Add(sma.Type, sma.New)
	Add(supertrend.Type, supertrend.New)

	Add(math.TypeAdd2, math.NewAdd2)
	Add(math.TypeAddK, math.NewAddK)
	Add(math.TypeDiv2, math.NewDiv2)
	Add(math.TypeSub2, math.NewSub2)
	Add(math.TypeMul2, math.NewMul2)
	Add(math.TypeMulK, math.NewMulK)
}

func Types() []string {
	return maps.Keys(reg)
}

func Add(typ string, new NewBlockFunc) {
	reg[typ] = new
}

func Get(typ string) (NewBlockFunc, bool) {
	entry, found := reg[typ]
	if !found {
		return nil, false
	}

	return entry, true
}
