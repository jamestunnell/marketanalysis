package processors

import (
	"github.com/jamestunnell/marketanalysis/models"
)

type Diff struct {
	hasPrev      bool
	output, prev float64
	warm         bool
}

const TypeDiff = "Diff"

func NewDiff() *Diff {
	return &Diff{
		hasPrev: false,
		output:  0.0,
		prev:    0.0,
		warm:    false,
	}
}

func (d *Diff) Type() string {
	return TypeDiff
}

func (d *Diff) Params() models.Params {
	return models.Params{}
}

func (d *Diff) Initialize() error {
	d.hasPrev = false
	d.output = 0.0
	d.prev = 0.0
	d.warm = false

	return nil
}

func (d *Diff) WarmupPeriod() int {
	return 2
}

func (d *Diff) Output() float64 {
	return d.output
}

func (d *Diff) Warm() bool {
	return d.warm
}

func (d *Diff) Update(val float64) {
	if !d.hasPrev {
		d.hasPrev = true
		d.prev = val

		return
	}

	d.output = val - d.prev
	d.prev = val
	d.warm = true
}
