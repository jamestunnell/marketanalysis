package processors

import (
	"fmt"

	"github.com/jamestunnell/marketanalysis/constraints"
	"github.com/jamestunnell/marketanalysis/indicators"
	"github.com/jamestunnell/marketanalysis/models"
)

type Slope struct {
	length  *models.TypedParam[int]
	linregr *indicators.LinRegression
	output  float64
}

const (
	LengthName = "length"

	TypeSlope = "Slope"
)

func NewSlope() *Slope {
	lengthRange := constraints.NewValRange(2, 100)

	return &Slope{
		length:  models.NewParam[int](lengthRange),
		linregr: nil,
		output:  0.0,
	}
}

func (s *Slope) Type() string {
	return TypeSlope
}

func (s *Slope) Params() models.Params {
	return models.Params{
		LengthName: s.length,
	}
}

func (s *Slope) Initialize() error {
	s.linregr = indicators.NewLinRegression(s.length.Value)
	s.output = 0.0

	return nil
}

func (s *Slope) WarmupPeriod() int {
	return s.linregr.WarmupPeriod()
}

func (s *Slope) Output() float64 {
	return s.linregr.Slope()
}

func (s *Slope) WarmUp(vals []float64) error {
	if err := s.linregr.WarmUp(vals); err != nil {
		return fmt.Errorf("failed to warm up linear regression indicator: %w", err)
	}

	return nil
}

func (s *Slope) Update(val float64) {
	s.linregr.Update(val)
}
