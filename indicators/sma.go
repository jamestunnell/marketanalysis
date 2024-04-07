package indicators

import (
	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type SMA struct {
	len     int
	buf     *buffer.CircularBuffer[float64]
	current float64
}

func NewSMA(len int) *SMA {
	return &SMA{
		len:     len,
		buf:     buffer.NewCircularBuffer[float64](len),
		current: 0.0,
	}
}

func (sma *SMA) Period() int {
	return sma.len
}

func (sma *SMA) Warm() bool {
	return sma.buf.Full()
}

func (sma *SMA) Update(val float64) {
	sma.buf.Add(val)

	if sma.buf.Full() {
		sum := 0.0
		sma.buf.Each(func(val float64) {
			sum += val
		})

		sma.current = sum / float64(sma.len)
	}
}

func (sma *SMA) Current() float64 {
	return sma.current
}
