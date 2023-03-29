package bar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Bar struct {
	Timestamp time.Time
	*BarCommon
}

type BarJSON struct {
	Timestamp string `json:"t"`
	*BarCommon
}

type BarCommon struct {
	Open       float64 `json:"o"`
	High       float64 `json:"h"`
	Low        float64 `json:"l"`
	Close      float64 `json:"c"`
	Volume     uint64  `json:"v"`
	TradeCount uint64  `json:"n"`
	VWAP       float64 `json:"vw"`
}

func New(t time.Time, o, h, l, c float64, vol, trades uint64, vwap float64) *Bar {
	return &Bar{
		Timestamp: t,
		BarCommon: &BarCommon{
			Open:       o,
			High:       h,
			Low:        l,
			Close:      c,
			Volume:     vol,
			TradeCount: trades,
			VWAP:       vwap,
		},
	}
}

func NewFromAlpacaBar(alpacaBar marketdata.Bar) *Bar {
	bc := &BarCommon{
		Open:       alpacaBar.Open,
		High:       alpacaBar.High,
		Low:        alpacaBar.Low,
		Close:      alpacaBar.Close,
		Volume:     alpacaBar.Volume,
		TradeCount: alpacaBar.TradeCount,
		VWAP:       alpacaBar.VWAP,
	}

	return &Bar{
		Timestamp: alpacaBar.Timestamp,
		BarCommon: bc,
	}
}

func (b *Bar) OpenCloseLowHigh() OpenCloseLowHigh {
	return OpenCloseLowHigh{b.Open, b.Close, b.Low, b.High}
}

func (b *Bar) MarshalJSON() ([]byte, error) {
	bj := &BarJSON{
		Timestamp: b.Timestamp.Format(time.RFC3339),
		BarCommon: b.BarCommon,
	}

	return json.Marshal(bj)
}

func (b *Bar) UnmarshalJSON(d []byte) error {
	var bj BarJSON

	if err := json.Unmarshal(d, &bj); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	ts, err := time.Parse(time.RFC3339, bj.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	b.BarCommon = bj.BarCommon
	b.Timestamp = ts

	return nil
}
