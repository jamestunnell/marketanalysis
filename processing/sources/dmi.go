package sources

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type DMI struct {
	period       *models.TypedParam[int]
	dmiValueType *models.TypedParam[string]
	output       float64
	dmi          *indicators.DMI
}

const (
	DMIValueName = "DMIValue"
	DMIValueDX   = "DX"
	DMIValueMDI  = "MDI"
	DMIValuePDI  = "PDI"

	PeriodName = "period"

	TypeDMI = "DMI"
)

func DMIValueTypes() []string {
	return []string{DMIValueDX, DMIValueMDI, DMIValuePDI}
}

func NewDMI() *DMI {
	periodRange := constraints.NewValRange(1, 200)
	dmiValueTypesEnum := constraints.NewValOneOf(DMIValueTypes())

	return &DMI{
		output:       0.0,
		period:       models.NewParam[int](periodRange),
		dmiValueType: models.NewParam[string](dmiValueTypesEnum),
		dmi:          nil,
	}
}

func (dmi *DMI) Type() string {
	return TypeDMI
}

func (dmi *DMI) Params() models.Params {
	return models.Params{
		PeriodName:   dmi.period,
		DMIValueName: dmi.dmiValueType,
	}
}

func (dmi *DMI) Initialize() error {
	dmi.dmi = indicators.NewDMI(dmi.period.Value)
	dmi.output = 0.0

	return nil
}

func (dmi *DMI) WarmupPeriod() int {
	return dmi.WarmupPeriod()
}

func (dmi *DMI) Output() float64 {
	return dmi.output
}

func (dmi *DMI) Warm() bool {
	return dmi.dmi.Warm()
}

func (dmi *DMI) Update(bar *models.Bar) {
	dmi.dmi.Update(bar)

	if dmi.Warm() {
		switch dmi.dmiValueType.Value {
		case DMIValueDX:
			dmi.output = dmi.dmi.DX()
		case DMIValueMDI:
			dmi.output = dmi.dmi.MDI()
		case DMIValuePDI:
			dmi.output = dmi.dmi.PDI()
		}
	}
}
