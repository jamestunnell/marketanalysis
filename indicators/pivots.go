package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
	"github.com/jamestunnell/marketanalysis/models"
)

type Pivots struct {
	Length             int
	NPivots            int
	prev               *Pivot
	primaryCandidate   *Pivot
	secondaryCandidate *Pivot
	age                int
	warm               bool
}

type Pivot struct {
	Type   PivotType
	Length int
	Value  float64
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
		prev:               nil,
		primaryCandidate:   nil,
		secondaryCandidate: nil,
		age:                0,
	}

	return pivs, nil
}

func NewPivotNeutral(length int, val float64) *Pivot {
	return &Pivot{
		Type:   PivotNeutral,
		Length: length,
		Value:  val,
	}
}

func NewPivotHigh(length int, val float64) *Pivot {
	return &Pivot{
		Type:   PivotHigh,
		Length: length,
		Value:  val,
	}
}

func NewPivotLow(length int, val float64) *Pivot {
	return &Pivot{
		Type:   PivotLow,
		Length: length,
		Value:  val,
	}
}

func (zz *Pivots) WarmupPeriod() int {
	return zz.Length
}

func (zz *Pivots) WarmUp(vals []float64) error {
	nVals := len(vals)
	if nVals < zz.Length {
		return commonerrs.NewErrMinCount("warmup vals", nVals, zz.Length)
	}

	zz.primaryCandidate = nil
	zz.secondaryCandidate = nil
	zz.prev = NewPivotNeutral(0, vals[0])
	zz.age = 0
	zz.warm = true

	min, max, minIdx, maxIdx := minMax(vals)

	var startUpdates int

	// flat signal - no candidate for next pivot yet
	if minIdx == maxIdx {
		zz.primaryCandidate = NewPivotNeutral(nVals, vals[0])

		startUpdates = nVals
	} else if minIdx == 0 || (maxIdx < minIdx) {
		zz.primaryCandidate = NewPivotHigh(maxIdx, max)

		startUpdates = maxIdx + 1
	} else {
		zz.primaryCandidate = NewPivotLow(minIdx, min)

		startUpdates = minIdx + 1
	}

	// this should take care of all remaining warmup values
	for i := startUpdates; i < zz.Length; i++ {
		_, _, _ = zz.Update(vals[i])
	}

	return nil
}

func (zz *Pivots) Direction() models.Direction {
	switch zz.primaryCandidate.Type {
	case PivotHigh:
		return models.DirUp
	case PivotLow:
		return models.DirDown
	}

	return models.DirNone
}

func (zz *Pivots) Update(val float64) (piv *Pivot, age int, pivotDetected bool) {
	if !zz.warm {
		return
	}

	dir := zz.Direction()

	switch dir {
	case models.DirNone:
		if val > zz.prev.Value {
			zz.prev = zz.primaryCandidate
			zz.primaryCandidate = NewPivotHigh(0, val)
		} else if val < zz.prev.Value {
			zz.prev = zz.primaryCandidate
			zz.primaryCandidate = NewPivotLow(0, val)
		} else {
			zz.primaryCandidate.Length++
		}
	case models.DirUp:
		if val >= zz.primaryCandidate.Value {
			zz.primaryCandidate.Value = val
			zz.primaryCandidate.Length += zz.age
			zz.secondaryCandidate = nil
			zz.age = 0
		} else if (val < zz.prev.Value) || (zz.age >= zz.Length) {
			piv = zz.primaryCandidate
			age = zz.age

			zz.prev = zz.primaryCandidate
			zz.primaryCandidate = zz.secondaryCandidate
			zz.age = 0
		} else {
			if zz.secondaryCandidate == nil {
				zz.secondaryCandidate = NewPivotLow(0, val)
			} else {
				zz.secondaryCandidate.Length++
				zz.secondaryCandidate.Value = val
			}

			zz.age++
		}
	case models.DirDown:
		if val <= zz.primaryCandidate.Value {
			zz.primaryCandidate.Value = val
			zz.primaryCandidate.Length += zz.age
			zz.secondaryCandidate = nil
			zz.age = 0
		} else if (val > zz.prev.Value) || (zz.age >= zz.Length) {
			piv = zz.primaryCandidate
			age = zz.age

			zz.prev = zz.primaryCandidate
			zz.primaryCandidate = zz.secondaryCandidate
			zz.age = 0
		} else {
			if zz.secondaryCandidate == nil {
				zz.secondaryCandidate = NewPivotHigh(0, val)
			} else {
				zz.secondaryCandidate.Length++
				zz.secondaryCandidate.Value = val
			}

			zz.age++
		}
	}

	return piv, age, piv != nil
}

func (zz *Pivots) LastPivot() *Pivot {
	return zz.prev
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
