package models

import (
	"time"
)

type TimeSeries struct {
	Quantities []*Quantity `json:"quantities"`
}

type Quantity struct {
	Name    string            `json:"name"`
	Records []*QuantityRecord `json:"records"`
}

type QuantityRecord struct {
	Timestamp time.Time `json:"t"`
	Value     float64   `json:"v"`
}

func (ts *TimeSeries) FindQuantity(name string) (*Quantity, bool) {
	for _, q := range ts.Quantities {
		if q.Name == name {
			return q, true
		}
	}

	return nil, false
}

func (q *Quantity) FindRecord(t time.Time) (*QuantityRecord, bool) {
	for _, record := range q.Records {
		if record.Timestamp == t {
			return record, true
		}
	}

	return nil, false
}
