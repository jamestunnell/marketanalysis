package models

import (
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rs/zerolog/log"
)

type Quantity struct {
	Name         string             `json:"name"`
	Records      []QuantityRecord   `json:"records"`
	Measurements map[string]float64 `json:"measurements"`
	Attributes   map[string]any     `json:"attributes"`
}

type QuantityRecord = TimeValue[float64]

const AttrCluster = "cluster"

func NewQuantity(name string, records ...QuantityRecord) *Quantity {
	return &Quantity{
		Name:         name,
		Records:      records,
		Measurements: map[string]float64{},
		Attributes:   map[string]any{},
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

func (q *Quantity) RecordValues() []float64 {
	return sliceutils.Map(q.Records, func(r QuantityRecord) float64 {
		return r.Value
	})
}

func (q *Quantity) MeasureAll(mTypes []string) {
	for _, mType := range mTypes {
		val, found := q.Measure(mType)
		if !found {
			log.Warn().Str("type", mType).Msg("measurement failed")

			continue
		}

		log.Trace().
			Str("quantity", q.Name).
			Str("type", mType).
			Float64("value", val).
			Msg("measured quantity")

	}
}

func (q *Quantity) Measure(typ string) (float64, bool) {
	measure, found := GetMeasureFunc(typ)
	if !found {
		return 0.0, false
	}

	val := measure(q.RecordValues())

	q.Measurements[typ] = val

	return val, true
}
