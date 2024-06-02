package models

import (
	"slices"
	"time"
)

type MarketDataset struct {
	Symbol string `json:"symbol" bson:"_id"`
	Bars   `json:"bars"`
}

func NewMarketDataset(symbol string) *MarketDataset {
	return &MarketDataset{
		Symbol: symbol,
		Bars:   Bars{},
	}
}

func (ds *MarketDataset) Validate() []error {
	return []error{}
}

func (ds *MarketDataset) GetKey() string {
	return ds.Symbol
}

func (ds *MarketDataset) AddBars(bars Bars) {
	bars.Sort()

	for _, bar := range bars {
		idx, found := bars.BinarySearch(bar.Timestamp)
		if found {
			bars[idx] = bar
		} else {
			slices.Insert(ds.Bars, idx, bar)
		}
	}
}

func (ds *MarketDataset) GetNBarsBefore(n int, t time.Time) Bars {
	idx, found := ds.Bars.BinarySearch(t)

	if !found {
		return Bars{}
	}

	if idx <= n {
		return ds.Bars[:idx]
	}

	return ds.Bars[idx-n : idx]
}
