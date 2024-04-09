package blocks

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util"
)

type MAOrder struct {
	in  *models.TypedInput[float64]
	out *models.TypedOutput[float64]

	numPeriods  *models.TypedParam[int]
	periodStart *models.TypedParam[int]
	periodSpan  *models.TypedParam[int]

	maOrdering *indicators.MAOrdering
}

const (
	DescrMAOrder = `Generates multiple MAs over a range of periods.
Compares sorted order of MA outputs to the sorted order by MA period to determine trend direction.`
	NameNumPeriods  = "numPeriods"
	NamePeriodSpan  = "periodSpan"
	NamePeriodStart = "periodStart"
	TypeMAOrder     = "MAOrder"
)

func NewMAOrder() models.Block {
	numPeriodsRange := constraints.NewValRange(3, 25)
	startRange := constraints.NewValRange(2, 100)
	spanRange := constraints.NewValRange(5, 100)

	return &MAOrder{
		numPeriods:  models.NewParam[int](3, numPeriodsRange),
		periodStart: models.NewParam[int](2, startRange),
		periodSpan:  models.NewParam[int](5, spanRange),
		maOrdering:  nil,
	}
}

func (blk *MAOrder) GetType() string {
	return TypeMAOrder
}

func (blk *MAOrder) GetDescription() string {
	return DescrMAOrder
}

func (blk *MAOrder) GetParams() models.Params {
	return models.Params{
		NameNumPeriods:  blk.numPeriods,
		NamePeriodSpan:  blk.periodSpan,
		NamePeriodStart: blk.periodStart,
	}
}

func (blk *MAOrder) GetInputs() models.Inputs {
	return models.Inputs{NameIn: blk.in}
}

func (blk *MAOrder) GetOutputs() models.Outputs {
	return models.Outputs{NameOut: blk.out}
}

func (blk *MAOrder) Init() error {
	start := blk.periodStart.Value
	span := blk.periodSpan.Value
	periods := util.LinSpaceInts(start, start+span, blk.numPeriods.Value)

	maOrdering := indicators.NewMAOrdering(periods)

	blk.maOrdering = maOrdering

	return nil
}

func (blk *MAOrder) IsWarm() bool {
	return blk.maOrdering.Warm()
}

func (blk *MAOrder) Update(bar *models.Bar) {
	blk.maOrdering.Update(blk.in.Get())

	if !blk.maOrdering.Warm() {
		return
	}

	blk.out.Set(blk.maOrdering.Correlation())
}
