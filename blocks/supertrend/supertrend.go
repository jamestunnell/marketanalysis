package supertrend

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type Supertrend struct {
	in                  *blocks.TypedInput[float64]
	trend, lower, upper *blocks.TypedOutput[float64]

	stateMachine *statemachine.StateMachine[*models.OHLC]

	atrPeriod *blocks.IntParam
	atrMul    *blocks.FloatParam
	atr       *indicators.ATR
}

const (
	Type  = "Supertrend"
	Descr = "Creates dynamic support and resistance levels from ATR. Uses these levels to determine trend."

	NameATRPeriod = "atrPeriod"
	NameATRMul    = "atrMul"
	NameLower     = "lower"
	NameUpper     = "upper"
	NameTrend     = "trend"
)

func New() blocks.Block {
	return &Supertrend{
		in:           blocks.NewTypedInput[float64](),
		trend:        blocks.NewTypedOutput[float64](),
		lower:        blocks.NewTypedOutput[float64](),
		upper:        blocks.NewTypedOutput[float64](),
		stateMachine: nil,
		atrPeriod:    blocks.NewIntParam(20, blocks.NewGreaterEqual(1)),
		atrMul:       blocks.NewFloatParam(5.0, blocks.NewGreater(0.0)),
		atr:          nil,
	}
}

func (blk *Supertrend) GetType() string {
	return Type
}

func (blk *Supertrend) GetDescription() string {
	return Descr
}

func (blk *Supertrend) GetParams() blocks.Params {
	return blocks.Params{
		NameATRPeriod: blk.atrPeriod,
		NameATRMul:    blk.atrMul,
	}
}

func (blk *Supertrend) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *Supertrend) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		NameLower: blk.lower,
		NameUpper: blk.upper,
		NameTrend: blk.trend,
	}
}

func (blk *Supertrend) GetWarmupPeriod() int {
	return blk.atr.Period()
}

func (blk *Supertrend) IsWarm() bool {
	return blk.atr.Warm()
}

func (blk *Supertrend) Init() error {
	blk.atr = indicators.NewATR(blk.atrPeriod.CurrentVal)
	blk.stateMachine = statemachine.New("supertrend", NewWaitForFirst(blk))

	return nil
}

func (blk *Supertrend) Update(cur *models.Bar, isLast bool) {
	blk.stateMachine.Run(cur.Timestamp, cur.OHLC)
}
