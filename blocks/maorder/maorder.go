package maorder

import (
	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util"
)

type MAOrder struct {
	in  *blocks.TypedInput[float64]
	out *blocks.TypedOutput[float64]

	numPeriods  *models.IntParam
	periodStart *models.IntParam
	periodSpan  *models.IntParam

	maOrdering *indicators.MAOrdering
}

const (
	Type  = "MAOrder"
	Descr = `Generates multiple MAs over a range of periods.
Compares sorted order of MA outputs to the sorted order by MA period to determine trend direction.`

	NameNumPeriods  = "numPeriods"
	NamePeriodSpan  = "periodSpan"
	NamePeriodStart = "periodStart"
)

func New() blocks.Block {
	return &MAOrder{
		in:          blocks.NewTypedInput[float64](),
		out:         blocks.NewTypedOutput[float64](),
		numPeriods:  models.NewIntParam(5, models.NewGreaterEq(2)),
		periodStart: models.NewIntParam(10, models.NewGreaterEq(1)),
		periodSpan:  models.NewIntParam(20, models.NewGreaterEq(2)),
		maOrdering:  nil,
	}
}

func (blk *MAOrder) GetType() string {
	return Type
}

func (blk *MAOrder) GetDescription() string {
	return Descr
}

func (blk *MAOrder) GetParams() models.Params {
	return models.Params{
		NameNumPeriods:  blk.numPeriods,
		NamePeriodSpan:  blk.periodSpan,
		NamePeriodStart: blk.periodStart,
	}
}

func (blk *MAOrder) GetInputs() blocks.Inputs {
	return blocks.Inputs{blocks.NameIn: blk.in}
}

func (blk *MAOrder) GetOutputs() blocks.Outputs {
	return blocks.Outputs{blocks.NameOut: blk.out}
}

func (blk *MAOrder) GetWarmupPeriod() int {
	return blk.maOrdering.WarmupPeriod()
}

func (blk *MAOrder) IsWarm() bool {
	return blk.maOrdering.Warm()
}

func (blk *MAOrder) Init() error {
	start := blk.periodStart.CurrentVal
	span := blk.periodSpan.CurrentVal

	periods := util.LinSpaceInts(start, start+span, blk.numPeriods.CurrentVal)

	blk.maOrdering = indicators.NewMAOrdering(periods)

	return nil
}

func (blk *MAOrder) Update(_ *models.Bar, isLast bool) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.maOrdering.Update(blk.in.GetValue())

	if !blk.maOrdering.Warm() {
		return
	}

	blk.out.SetValue(blk.maOrdering.Correlation())
}
