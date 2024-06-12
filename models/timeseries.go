package models

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/rs/zerolog/log"
	"github.com/soniakeys/cluster"
)

type TimeSeries struct {
	Quantities []*Quantity `json:"quantities"`
}

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

type MakePointFunc func(q *Quantity) (cluster.Point, error)

func (ts *TimeSeries) Cluster(k int, makePoint MakePointFunc) error {
	if k < 1 {
		return fmt.Errorf("k %d is not positive", k)
	}

	if k == 1 {
		for _, q := range ts.Quantities {
			q.Attributes[AttrCluster] = 0
		}

		return nil
	}

	points := []cluster.Point{}
	qs := []*Quantity{}
	for _, q := range ts.Quantities {
		p, err := makePoint(q)
		if err != nil {
			log.Warn().Err(err).Str("name", q.Name).Msg("cannot make clustering points for quantity")

			continue
		}

		qs = append(qs, q)
		points = append(points, p)
	}

	if len(qs) == 0 {
		return errors.New("failed to make clustering points for all quantities")
	}

	centers, cNums, _, _ := cluster.KMPP(points, k)

	// clone and reverse sort center points
	revSortedCenters := slices.Clone(centers)

	slices.SortFunc(revSortedCenters, func(a, b cluster.Point) int {
		return slices.Compare(a, b)
	})
	slices.Reverse(revSortedCenters)

	for i, cNum := range cNums {
		tgt := centers[cNum]

		// find the sorted index
		sortedIdx := slices.IndexFunc(revSortedCenters, func(c cluster.Point) bool {
			return slices.Compare(c, tgt) == 0
		})

		qs[i].Attributes[AttrCluster] = sortedIdx
	}

	return nil
}

func QuantityMeanStddev(q *Quantity) (cluster.Point, error) {
	vals := q.RecordValues()

	if len(vals) == 0 {
		return cluster.Point{}, fmt.Errorf("quantity %s has no values", q.Name)
	}

	mean, err := stats.Mean(vals)
	if err != nil {
		return cluster.Point{}, fmt.Errorf("failed to calc mean for quantity %s: %w", q.Name, err)
	}

	sd, err := stats.StandardDeviation(vals)
	if err != nil {
		return cluster.Point{}, fmt.Errorf("failed to calc stddev for quantity %s: %w", q.Name, err)
	}

	return cluster.Point{mean, sd}, nil
}
