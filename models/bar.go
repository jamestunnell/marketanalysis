package models

import (
	"math"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Bar struct {
	Timestamp  time.Time `json:"t"`
	Volume     uint64    `json:"v"`
	TradeCount uint64    `json:"n"`
	VWAP       float64   `json:"vw"`

	*OHLC
}

type OHLC struct {
	Open  float64 `json:"o"`
	High  float64 `json:"h"`
	Low   float64 `json:"l"`
	Close float64 `json:"c"`
}

func NewBarFromAlpaca(alpacaBar marketdata.Bar) *Bar {
	ohlc := &OHLC{
		Open:  alpacaBar.Open,
		High:  alpacaBar.High,
		Low:   alpacaBar.Low,
		Close: alpacaBar.Close,
	}

	return &Bar{
		Timestamp:  alpacaBar.Timestamp,
		Volume:     alpacaBar.Volume,
		TradeCount: alpacaBar.TradeCount,
		VWAP:       alpacaBar.VWAP,
		OHLC:       ohlc,
	}
}

func (b *Bar) Localize() {
	b.Timestamp = b.Timestamp.Local()
}

func (b *Bar) HeikinAshi(prev *Bar) *OHLC {
	return &OHLC{
		Open:  0.5 * (prev.Open + prev.Close),
		Close: 0.25 * (b.Open + b.High + b.Low + b.Close),
		High:  math.Max(math.Max(b.High, b.Open), b.Close),
		Low:   math.Max(math.Max(b.Low, b.Open), b.Close),
	}
}

func (ohlc *OHLC) Float64s() []float64 {
	return []float64{ohlc.Open, ohlc.High, ohlc.Low, ohlc.Close}
}
