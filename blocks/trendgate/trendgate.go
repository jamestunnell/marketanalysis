package trendgate

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/statemachine"
)

type TrendGate struct {
	openThresh     *models.FloatParam
	closeThresh    *models.FloatParam
	debouncePeriod *models.IntParam
	in             *blocks.TypedInput[float64]
	out            *blocks.TypedOutput[float64]
	stateMachine   *statemachine.StateMachine[float64]
}

type StateInput struct{}

const (
	Type  = "TrendGate"
	Descr = "Silence signal below threshold"

	NameDebouncePeriod = "debouncePeriod"
	NameOpenThresh     = "openThreshold"
	NameCloseThresh    = "closeThreshold"
)

func New() blocks.Block {
	return &TrendGate{
		openThresh:     models.NewFloatParam(0.5, models.NewRangeExcl(0.0, 1.0)),
		closeThresh:    models.NewFloatParam(0.25, models.NewRangeExcl(0.0, 1.0)),
		debouncePeriod: models.NewIntParam(0, models.NewGreaterEq(0)),
		in:             blocks.NewTypedInput[float64](),
		out:            blocks.NewTypedOutput[float64](),
		stateMachine:   nil,
	}
}

func (blk *TrendGate) GetType() string {
	return Type
}

func (blk *TrendGate) GetDescription() string {
	return Descr
}

func (blk *TrendGate) GetParams() models.Params {
	return models.Params{
		NameOpenThresh:     blk.openThresh,
		NameCloseThresh:    blk.closeThresh,
		NameDebouncePeriod: blk.debouncePeriod,
	}
}

func (blk *TrendGate) GetInputs() blocks.Inputs {
	return blocks.Inputs{
		blocks.NameIn: blk.in,
	}
}

func (blk *TrendGate) GetOutputs() blocks.Outputs {
	return blocks.Outputs{
		blocks.NameOut: blk.out,
	}
}

func (blk *TrendGate) GetWarmupPeriod() int {
	return 0
}

func (blk *TrendGate) IsWarm() bool {
	return true
}

func (blk *TrendGate) Init() error {
	if blk.closeThresh.CurrentVal > blk.openThresh.CurrentVal {
		return fmt.Errorf("close thresh %f is more than open thresh %f", blk.closeThresh.CurrentVal, blk.openThresh.CurrentVal)
	}

	log.Trace().
		Float64("openThresh", blk.openThresh.CurrentVal).
		Float64("closeThresh", blk.closeThresh.CurrentVal).
		Msg("trend gate initialized")

	blk.stateMachine = statemachine.New("gate", NewNone(blk))

	return nil
}

func (blk *TrendGate) Update(cur *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.stateMachine.Run(cur.Timestamp, blk.in.GetValue())
}
