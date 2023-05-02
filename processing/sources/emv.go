package sources

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMV struct {
	emv *indicators.EMV
}

const (
	TypeEMV = "EMV"
)

func NewEMV() *EMV {
	return &EMV{
		emv: nil,
	}
}

func (emv *EMV) Type() string {
	return TypeEMV
}

func (emv *EMV) Params() models.Params {
	return models.Params{}
}

func (emv *EMV) Initialize() error {
	emv.emv = indicators.NewEMV()

	return nil
}

func (emv *EMV) WarmupPeriod() int {
	return emv.WarmupPeriod()
}

func (emv *EMV) Output() float64 {
	return emv.emv.EMV()
}

func (emv *EMV) WarmUp(bars models.Bars) error {
	if err := emv.emv.WarmUp(bars); err != nil {
		return fmt.Errorf("failed to warm up EMV indicator: %w", err)
	}

	return nil
}

func (emv *EMV) Update(bar *models.Bar) {
	emv.emv.Update(bar)
}
