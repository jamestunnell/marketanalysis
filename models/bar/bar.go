package bar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Bar struct {
	Timestamp time.Time
	*OHLC
}

type BarJSON struct {
	Timestamp string `json:"t"`
	*OHLC
}

// type BarCommon struct {
// 	Volume     uint64  `json:"v"`
// 	TradeCount uint64  `json:"n"`
// 	VWAP       float64 `json:"vw"`
// }

func New(t time.Time, o, h, l, c float64) *Bar {
	return &Bar{
		Timestamp: t,
		OHLC: &OHLC{
			Open:  o,
			High:  h,
			Low:   l,
			Close: c,
			// Volume:     vol,
			// TradeCount: trades,
			// VWAP:       vwap,
		},
	}
}

func NewFromOHLC(t time.Time, ohlc *OHLC) *Bar {
	return &Bar{
		Timestamp: t,
		OHLC:      ohlc,
	}
}

func NewFromAlpacaBar(alpacaBar marketdata.Bar) *Bar {
	ohlc := &OHLC{
		Open:  alpacaBar.Open,
		High:  alpacaBar.High,
		Low:   alpacaBar.Low,
		Close: alpacaBar.Close,
		// Volume:     alpacaBar.Volume,
		// TradeCount: alpacaBar.TradeCount,
		// VWAP:       alpacaBar.VWAP,
	}

	return &Bar{
		Timestamp: alpacaBar.Timestamp,
		OHLC:      ohlc,
	}
}

func (b *Bar) MarshalJSON() ([]byte, error) {
	bj := &BarJSON{
		Timestamp: b.Timestamp.Format(time.RFC3339),
		OHLC:      b.OHLC,
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

	b.OHLC = bj.OHLC
	b.Timestamp = ts

	return nil
}
