package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
)

type LinRegression struct {
	length           int
	prevVals         *buffer.CircularBuffer[float64]
	slope, intercept float64
	coords           []stats.Coordinate
}

func NewLinRegression(length int) *LinRegression {
	coords := make([]stats.Coordinate, length)
	for i := 0; i < length; i++ {
		coords[i] = stats.Coordinate{X: float64(i), Y: 0.0}
	}

	return &LinRegression{
		length:   length,
		prevVals: buffer.NewCircularBuffer[float64](length),
		coords:   coords,
		slope:    0.0,
	}
}

func (trs *LinRegression) WarmupPeriod() int {
	return trs.length
}

func (trs *LinRegression) WarmUp(vals []float64) error {
	wp := trs.WarmupPeriod()
	if len(vals) < wp {
		return commonerrs.NewErrMinCount("warmup values", len(vals), wp)
	}

	vals = sliceutils.LastN(vals, wp) // no benefit from having more

	for i := 0; i < wp; i++ {
		trs.prevVals.Add(vals[i])
	}

	trs.updateSlope()

	return nil
}

func (trs *LinRegression) Slope() float64 {
	return trs.slope
}

func (trs *LinRegression) Intercept() float64 {
	return trs.intercept
}

func (trs *LinRegression) updateSlope() {
	trs.prevVals.EachWithIndex(func(i int, f float64) {
		trs.coords[i].Y = f
	})

	reg, err := stats.LinearRegression(trs.coords)
	if err != nil {
		log.Warn().Err(err).Msg("linear regression failed")

		return
	}

	dY := (reg[trs.length-1].Y - reg[0].Y)
	dX := float64(trs.length)

	trs.slope = dY / dX
	trs.intercept = reg[0].Y - (reg[0].X * trs.slope)
}

func (trs *LinRegression) Update(val float64) {
	trs.prevVals.Add(val)

	trs.updateSlope()
}
