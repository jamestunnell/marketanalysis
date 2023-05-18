package processors

import (
	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type AroonOsc struct {
	period     *models.TypedParam[int]
	aroonValue *models.TypedParam[string]
	osc        *indicators.AroonOsc
	output     float64
}

const (
	AroonValueUp   = "up"
	AroonValueDown = "down"
	AroonValueDiff = "diff"
	AroonValueName = "aroonValue"

	TypeAroonOsc = "AroonOsc"
)

func AroonValueTypes() []string {
	return []string{AroonValueUp, AroonValueDown, AroonValueDiff}
}

func NewAroonOsc() *AroonOsc {
	periodRange := constraints.NewValRange(1, 200)
	valueTypeEnum := constraints.NewValOneOf(AroonValueTypes())

	return &AroonOsc{
		period:     models.NewParam[int](periodRange),
		aroonValue: models.NewParam[string](valueTypeEnum),
		osc:        nil,
		output:     0.0,
	}
}

func (osc *AroonOsc) Type() string {
	return TypeAroonOsc
}

func (osc *AroonOsc) Params() models.Params {
	return models.Params{
		PeriodName:     osc.period,
		AroonValueName: osc.aroonValue,
	}
}

func (osc *AroonOsc) Initialize() error {
	osc.output = 0.0
	osc.osc = indicators.NewAroonOsc(osc.period.Value)

	return nil
}

func (osc *AroonOsc) WarmupPeriod() int {
	return osc.osc.WarmupPeriod()
}

func (osc *AroonOsc) Warm() bool {
	return osc.osc.Warm()
}

func (osc *AroonOsc) Update(val float64) {
	osc.osc.Update(val)

	if osc.osc.Warm() {
		osc.updateOutput()
	}
}

func (osc *AroonOsc) Output() float64 {
	return osc.output
}

func (osc *AroonOsc) updateOutput() {
	var val float64

	switch osc.aroonValue.Value {
	case AroonValueDown:
		val = osc.osc.Down()
	case AroonValueUp:
		val = osc.osc.Up()
	case AroonValueDiff:
		val = osc.osc.Diff()
	}

	osc.output = val
}
