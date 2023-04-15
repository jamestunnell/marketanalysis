package indicators

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type EMA struct {
	len        int
	warm       bool
	current, k float64
}

type MovingAvgJSON struct {
	Length int `json:"length"`
}

func NewEMA(len int) *EMA {
	return &EMA{
		len:     len,
		warm:    false,
		current: 0.0,
		k:       2.0 / (float64(len) + 1),
	}
}

func (ema *EMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(&MovingAvgJSON{Length: ema.len})
}

func (ema *EMA) UnmarshalJSON(d []byte) error {
	var maj MovingAvgJSON

	if err := json.Unmarshal(d, &maj); err != nil {
		return fmt.Errorf("failed to unmarshal EMM JSON: %w", err)
	}

	ema.len = maj.Length
	ema.warm = false
	ema.current = 0.0
	ema.k = 2.0 / (float64(maj.Length) + 1)

	return nil
}

func (ema *EMA) Period() int {
	return ema.len
}

func (ema *EMA) WarmUp(vals []float64) error {
	if len(vals) != ema.len {
		return commonerrs.NewErrExactCount("warmup values", ema.len, len(vals))
	}

	sum := 0.0
	for _, close := range vals {
		sum += close
	}

	ema.current = sum / float64(ema.len)
	ema.warm = true

	return nil
}

func (ema *EMA) Update(val float64) {
	if !ema.warm {
		return
	}

	ema.current = (val * ema.k) + (ema.current * (1.0 - ema.k))
}

func (ema *EMA) Current() float64 {
	if !ema.warm {
		return 0.0
	}

	return ema.current
}
