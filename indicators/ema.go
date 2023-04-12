package indicators

import (
	"encoding/json"
	"fmt"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

type EMA struct {
	len        int
	warm       bool
	current, k float64
}

type MovingAvgJSON struct {
	Length int `json:"length"`
}

func NewEMA(len int) models.Indicator {
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

func (ema *EMA) WarmupPeriod() int {
	return ema.len
}

func (ema *EMA) WarmUp(bars models.Bars) error {
	if len(bars) != ema.len {
		return commonerrs.NewErrExactBarCount("warm up", ema.len, len(bars))
	}

	sum := 0.0
	for _, close := range bars.ClosePrices() {
		sum += close
	}

	ema.current = sum / float64(ema.len)
	ema.warm = true

	return nil
}

func (ema *EMA) Update(bar *models.Bar) float64 {
	if !ema.warm {
		return 0.0
	}

	ema.current = (bar.Close * ema.k) + (ema.current * (1.0 - ema.k))

	return ema.current
}

func (ema *EMA) Current() float64 {
	if !ema.warm {
		return 0.0
	}

	return ema.current
}
