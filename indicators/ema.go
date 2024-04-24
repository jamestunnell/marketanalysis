package indicators

import (
	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type EMA struct {
	current, k float64
	len        int
	startBuf   *buffer.CircularBuffer[float64]
	warm       bool
}

func NewEMA(len int) *EMA {
	return &EMA{
		current:  0.0,
		k:        EMAWeightingMultiplier(len),
		len:      len,
		startBuf: buffer.NewCircularBuffer[float64](len),
		warm:     false,
	}
}

func (ema *EMA) Period() int {
	return ema.len
}

func (ema *EMA) Warm() bool {
	return ema.warm
}

func (ema *EMA) Update(val float64) {
	if ema.warm {
		ema.current = (val * ema.k) + (ema.current * (1.0 - ema.k))

		return
	}

	ema.startBuf.Add(val)

	if ema.startBuf.Full() {
		sum := 0.0
		ema.startBuf.Each(func(val float64) {
			sum += val
		})

		ema.current = sum / float64(ema.len)
		ema.warm = true
	}
}

func (ema *EMA) Current() float64 {
	return ema.current
}

func EMAWeightingMultiplier(len int) float64 {
	return 2.0 / (float64(len) + 1)
}
