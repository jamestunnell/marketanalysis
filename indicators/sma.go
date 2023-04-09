package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type SMA struct {
	len     int
	buf     *buffer.FullCircularBuffer
	current float64
}

func NewSMA(len int) models.Indicator {
	return &SMA{
		len:     len,
		buf:     nil,
		current: 0.0,
	}
}

func (sma *SMA) WarmupPeriod() int {
	return sma.len
}

func (sma *SMA) WarmUp(bars models.Bars) error {
	if len(bars) != sma.len {
		return commonerrs.NewErrExactBarCount("warm up", sma.len, len(bars))
	}

	sma.buf = buffer.NewFullCircularBuffer(bars.ClosePrices())
	sma.current = sma.buf.Sum()

	return nil
}

func (sma *SMA) Update(bar *models.Bar) float64 {
	if sma.buf == nil {
		return 0.0
	}

	sma.buf.Add(bar.Close)

	sma.current = sma.buf.Sum()

	return sma.current
}

func (sma *SMA) Current() float64 {
	if sma.buf == nil {
		return 0.0
	}

	return sma.current
}
