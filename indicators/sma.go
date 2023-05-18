package indicators

import (
	"encoding/json"
	"fmt"

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

func (sma *SMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(&MovingAvgJSON{Length: sma.len})
}

func (sma *SMA) UnmarshalJSON(d []byte) error {
	var maj MovingAvgJSON

	if err := json.Unmarshal(d, &maj); err != nil {
		return fmt.Errorf("failed to unmarshal EMM JSON: %w", err)
	}

	sma.len = maj.Length
	sma.buf = buffer.NewCircularBuffer[float64](maj.Length)
	sma.current = 0.0

	return nil
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
