package processors

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

type Diff struct {
	output, prev float64
}

const TypeDiff = "Diff"

func NewDiff() *Diff {
	return &Diff{output: 0.0, prev: 0.0}
}

func (d *Diff) Type() string {
	return TypeDiff
}

func (d *Diff) Params() models.Params {
	return models.Params{}
}

func (d *Diff) Initialize() error {
	d.prev = 0.0
	d.output = 0.0

	return nil
}

func (d *Diff) WarmupPeriod() int {
	return 2
}

func (d *Diff) Output() float64 {
	return d.output
}

func (d *Diff) WarmUp(vals []float64) error {
	if len(vals) != 2 {
		return commonerrs.NewErrExactLen("warmup vals", len(vals), 2)
	}

	d.prev = vals[0]

	d.Update(vals[1])
}

func (d *Diff) Update(val float64) {
	d.output = val - d.prev
	d.prev = val
}
