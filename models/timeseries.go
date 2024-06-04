package models

import (
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type TimeSeries struct {
	Quantities []*Quantity `json:"quantities"`
}

type Quantity struct {
	Name    string           `json:"name"`
	Records []QuantityRecord `json:"records"`
}

type QuantityRecord = TimeValue[float64]

func NewTimeSeries() *TimeSeries {
	return &TimeSeries{
		Quantities: []*Quantity{},
	}
}

func (ts *TimeSeries) IsEmpty() bool {
	for _, q := range ts.Quantities {
		if !q.IsEmpty() {
			return false
		}
	}

	return true
}

func (ts *TimeSeries) SortByTime() bool {
	for _, q := range ts.Quantities {
		q.SortByTime()
	}

	return true
}

func (ts *TimeSeries) AddQuantity(q *Quantity) {
	ts.Quantities = append(ts.Quantities, q)
}

func (ts *TimeSeries) FindQuantity(name string) (*Quantity, bool) {
	for _, q := range ts.Quantities {
		if q.Name == name {
			return q, true
		}
	}

	return nil, false
}

func (ts *TimeSeries) DropRecordsBefore(t time.Time) {
	for _, q := range ts.Quantities {
		q.DropRecordsBefore(t)
	}
}

func (q *Quantity) AddRecord(r QuantityRecord) {
	q.Records = append(q.Records, r)
}

func (q *Quantity) IsEmpty() bool {
	return len(q.Records) == 0
}

func (q *Quantity) SortByTime() {
	slices.SortStableFunc(q.Records, func(a, b QuantityRecord) int {
		return a.Time.Compare(b.Time)
	})
}

func (q *Quantity) FindRecord(t time.Time) (QuantityRecord, bool) {
	for _, record := range q.Records {
		if record.Time == t {
			return record, true
		}
	}

	return QuantityRecord{}, false
}

func (q *Quantity) FindRecordsAfter(t time.Time) []QuantityRecord {
	return sliceutils.Where(q.Records, func(r QuantityRecord) bool {
		return r.Time.After(t)
	})
}

func (q *Quantity) DropRecordsBefore(t time.Time) {
	q.Records = sliceutils.Where(q.Records, func(r QuantityRecord) bool {
		return !r.Time.Before(t)
	})
}
