package indicators

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type SMA struct {
	len     int
	buf     *buffer.FullCircularBuffer
	current float64
}

func NewSMA(len int) *SMA {
	return &SMA{
		len:     len,
		buf:     nil,
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
	sma.buf = nil
	sma.current = 0.0

	return nil
}

func (sma *SMA) Period() int {
	return sma.len
}

func (sma *SMA) WarmUp(vals []float64) error {
	if len(vals) != sma.len {
		return commonerrs.NewErrExactCount("warmup vals", sma.len, len(vals))
	}

	sma.buf = buffer.NewFullCircularBuffer(vals)
	sma.current = sma.buf.Sum()

	return nil
}

func (sma *SMA) Update(val float64) {
	if sma.buf == nil {
		return
	}

	sma.buf.Add(val)

	sma.current = sma.buf.Sum() / float64(sma.len)
}

func (sma *SMA) Current() float64 {
	if sma.buf == nil {
		return 0.0
	}

	return sma.current
}
