package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/buffer"
	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
)

type TrueRangeSlope struct {
	length  int
	prevTRs *buffer.CircularBuffer[float64]
	prevBar *models.Bar
	slope   float64
	coords  []stats.Coordinate
}

func NewTrueRangeSlope(length int) *TrueRangeSlope {
	coords := make([]stats.Coordinate, length)
	for i := 0; i < length; i++ {
		coords[i] = stats.Coordinate{X: float64(i), Y: 0.0}
	}

	return &TrueRangeSlope{
		length:  length,
		prevTRs: buffer.NewCircularBuffer[float64](length),
		prevBar: nil,
		coords:  coords,
		slope:   0.0,
	}
}

func (trs *TrueRangeSlope) Length() int {
	return trs.length
}

func (trs *TrueRangeSlope) WarmupPeriod() int {
	return trs.length + 1
}

func (trs *TrueRangeSlope) WarmUp(bars models.Bars) error {
	wp := trs.WarmupPeriod()
	if len(bars) < wp {
		return commonerrs.NewErrMinCount("warmup bars", len(bars), wp)
	}

	bars = bars.LastN(wp) // no benefit from having more
	trs.prevBar = bars[0]

	for i := 1; i < wp; i++ {
		bar := bars[i]

		tr := TrueRange(bar, trs.prevBar)

		trs.prevBar = bar

		trs.prevTRs.Add(tr)
	}

	trs.updateSlope()

	return nil
}

func (trs *TrueRangeSlope) Slope() float64 {
	return trs.slope
}

func (trs *TrueRangeSlope) updateSlope() {
	trs.prevTRs.EachWithIndex(func(i int, f float64) {
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
}

func (trs *TrueRangeSlope) Update(bar *models.Bar) {
	tr := TrueRange(bar, trs.prevBar)

	trs.prevBar = bar

	trs.prevTRs.Add(tr)

	trs.updateSlope()
}
