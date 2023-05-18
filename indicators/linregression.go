package indicators

import (
	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
)

type LinRegression struct {
	coords    []stats.Coordinate
	intercept float64
	length    int
	prevVals  *buffer.CircularBuffer[float64]
	slope     float64
}

func NewLinRegression(length int) *LinRegression {
	coords := make([]stats.Coordinate, length)
	for i := 0; i < length; i++ {
		coords[i] = stats.Coordinate{X: float64(i), Y: 0.0}
	}

	return &LinRegression{
		coords:    coords,
		intercept: 0.0,
		length:    length,
		prevVals:  buffer.NewCircularBuffer[float64](length),
		slope:     0.0,
	}
}

func (lr *LinRegression) WarmupPeriod() int {
	return lr.length
}

func (lr *LinRegression) Warm() bool {
	return lr.prevVals.Full()
}

func (lr *LinRegression) Slope() float64 {
	return lr.slope
}

func (lr *LinRegression) Intercept() float64 {
	return lr.intercept
}

func (lr *LinRegression) updateSlopeIntercept() {
	lr.prevVals.EachWithIndex(func(i int, f float64) {
		lr.coords[i].Y = f
	})

	reg, err := stats.LinearRegression(lr.coords)
	if err != nil {
		log.Warn().Err(err).Msg("linear regression failed")

		return
	}

	dY := (reg[lr.length-1].Y - reg[0].Y)
	dX := float64(lr.length)

	lr.slope = dY / dX
	lr.intercept = reg[0].Y - (reg[0].X * lr.slope)
}

func (lr *LinRegression) Update(val float64) {
	lr.prevVals.Add(val)
	if !lr.prevVals.Full() {
		return
	}

	lr.updateSlopeIntercept()
}
