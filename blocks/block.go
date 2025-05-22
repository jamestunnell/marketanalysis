package blocks

import (
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketdata"
)

type Block interface {
	GetType() string
	GetDescription() string
	GetParams() models.Params
	GetInputs() Inputs
	GetOutputs() Outputs
	GetWarmupPeriod() int

	IsWarm() bool

	Init() error
	Update(current *marketdata.Bar, isLast bool)
}

func ClearOutputs(blk Block) {
	for _, out := range blk.GetOutputs() {
		out.ClearValue()
	}
}
