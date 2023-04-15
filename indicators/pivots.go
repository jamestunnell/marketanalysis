package indicators

import (
	"time"

	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/util/buffer"
)

type Pivots struct {
	Length             int
	NPivots            int
	Pivots             *buffer.CircularBuffer[*Pivot]
	primaryCandidate   *Pivot
	secondaryCandidate *Pivot
	maxAge, tDelta     time.Duration
}

type Pivot struct {
	Type  PivotType
	Time  time.Time
	Value float64
}

type PivotType int

const (
	PivotNeutral PivotType = iota
	PivotLow
	PivotHigh
)

func NewPivots(length, nPivots int) (*Pivots, error) {
	if length < 2 {
		return nil, commonerrs.NewErrLessThanMin("length", length, 2)
	}

	if nPivots < 2 {
		return nil, commonerrs.NewErrLessThanMin("nPivots", nPivots, 2)
	}

	pivs := &Pivots{
		Length:             length,
		NPivots:            nPivots,
		Pivots:             nil,
		primaryCandidate:   nil,
		secondaryCandidate: nil,
		maxAge:             0,
		tDelta:             0,
	}

	return pivs, nil
}

func NewPivotNeutral(t time.Time, val float64) *Pivot {
	return &Pivot{
		Type:  PivotNeutral,
		Time:  t,
		Value: val,
	}
}

func NewPivotHigh(t time.Time, val float64) *Pivot {
	return &Pivot{
		Type:  PivotHigh,
		Time:  t,
		Value: val,
	}
}

func NewPivotLow(t time.Time, val float64) *Pivot {
	return &Pivot{
		Type:  PivotLow,
		Time:  t,
		Value: val,
	}
}

func (zz *Pivots) WarmupPeriod() int {
	return zz.Length
}

func (zz *Pivots) WarmUp(times []time.Time, vals []float64) error {
	if len(vals) != zz.Length {
		return commonerrs.NewErrExactCount("warmup vals", zz.Length, len(vals))
	}

	if len(times) != zz.Length {
		return commonerrs.NewErrExactCount("warmup times", zz.Length, len(times))
	}

	zz.Pivots = buffer.NewCircularBuffer[*Pivot](zz.NPivots)
	zz.primaryCandidate = nil
	zz.secondaryCandidate = nil
	zz.tDelta = times[1].Sub(times[0])
	zz.maxAge = zz.tDelta * time.Duration(zz.Length)

	zz.Pivots.Add(NewPivotNeutral(times[0], vals[0]))

	min, max, minIdx, maxIdx := minMax(vals)

	// flat signal - no cand-date for next pivot yet
	if minIdx == maxIdx {
		return nil
	}

	var startUpdates int

	if minIdx == 0 || (maxIdx < minIdx) {
		zz.primaryCandidate = NewPivotHigh(times[maxIdx], max)

		startUpdates = maxIdx + 1
	} else {
		zz.primaryCandidate = NewPivotLow(times[minIdx], min)

		startUpdates = minIdx + 1
	}

	// this should take care of remaining warmup values
	for i := startUpdates; i < zz.Length; i++ {
		_ = zz.Update(times[i], vals[i])
	}

	return nil
}

func (zz *Pivots) Direction() Direction {
	if zz.primaryCandidate == nil {
		return DirNone
	}

	switch zz.primaryCandidate.Type {
	case PivotHigh:
		return DirUp
	case PivotLow:
		return DirDown
	}

	return DirNone
}

func (zz *Pivots) Update(t time.Time, val float64) bool {
	// make sure it's warmed up
	if zz.Pivots == nil {
		return false
	}

	dir := zz.Direction()
	prev, _ := zz.Pivots.Newest()

	switch dir {
	case DirNone:
		if val > prev.Value {
			zz.Pivots.Add(NewPivotNeutral(t.Add(-zz.tDelta), prev.Value))
			zz.primaryCandidate = NewPivotHigh(t, val)
		} else if val < prev.Value {
			zz.Pivots.Add(NewPivotNeutral(t.Add(-zz.tDelta), prev.Value))
			zz.primaryCandidate = NewPivotLow(t, val)
		}
	case DirUp:
		if val >= zz.primaryCandidate.Value {
			zz.primaryCandidate.Value = val
			zz.primaryCandidate.Time = t
			zz.secondaryCandidate = nil
		} else if (val < prev.Value) || (t.Sub(zz.primaryCandidate.Time) >= zz.maxAge) {
			zz.Pivots.Add(zz.primaryCandidate)

			zz.primaryCandidate = zz.secondaryCandidate
		}
	case DirDown:
		if val <= zz.primaryCandidate.Value {
			zz.primaryCandidate.Value = val
			zz.primaryCandidate.Time = t
			zz.secondaryCandidate = nil
		} else if (val > prev.Value) || (t.Sub(zz.primaryCandidate.Time) >= zz.maxAge) {
			zz.Pivots.Add(zz.primaryCandidate)
			zz.primaryCandidate = zz.secondaryCandidate
		}
	}

	return dir != zz.Direction()
}

func (zz *Pivots) PivotsAfter(t time.Time) []*Pivot {
	pivs := []*Pivot{}
	zz.Pivots.Each(func(piv *Pivot) {
		if piv.Time.After(t) {
			pivs = append(pivs, piv)
		}
	})

	return pivs
}

func (zz *Pivots) NewestPivot() (*Pivot, bool) {
	return zz.Pivots.Newest()
}

func minMax(vals []float64) (min, max float64, minIdx, maxIdx int) {
	min, max = vals[0], vals[0]
	minIdx, maxIdx = 0, 0

	for i := 1; i < len(vals); i++ {
		x := vals[i]

		if x < min {
			min = x
			minIdx = i
		}

		if x > max {
			max = x
			maxIdx = i
		}
	}

	return
}
