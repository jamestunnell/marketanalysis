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

	numPeriods  *blocks.IntRange
	periodStart *blocks.IntRange
	periodSpan  *blocks.IntRange

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
		numPeriods:  &blocks.IntRange{Default: 2, Min: 2, Max: 25},
		periodStart: &blocks.IntRange{Default: 2, Min: 2, Max: 100},
		periodSpan:  &blocks.IntRange{Default: 2, Min: 2, Max: 900},
		maOrdering:  nil,
	}
}

func (blk *MAOrder) GetType() string {
	return Type
}

func (blk *MAOrder) GetDescription() string {
	return Descr
}

func (blk *MAOrder) GetParams() blocks.Params {
	return blocks.Params{
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
	start := blk.periodStart.Value
	span := blk.periodSpan.Value
	periods := util.LinSpaceInts(start, start+span, blk.numPeriods.Value)

	blk.maOrdering = indicators.NewMAOrdering(periods)

	return nil
}

func (blk *MAOrder) Update(_ *models.Bar) {
	if !blk.in.IsValueSet() {
		return
	}

	blk.maOrdering.Update(blk.in.GetValue())

	if !blk.maOrdering.Warm() {
		return
	}

	blk.out.SetValue(blk.maOrdering.Correlation())
}
