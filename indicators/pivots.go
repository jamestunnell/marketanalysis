package indicators

import (
	"github.com/jamestunnell/marketanalysis/commonerrs"
)

type Pivots struct {
	Length       int
	NPivots      int
	PrevPivots   []float64
	Candidate    float64
	CandidateAge int
	Direction    Direction
}

type Direction int

const (
	Up   = 1
	Flat = 0
	Down = -1
)

func NewPivots(length, nPivots int) *Pivots {
	return &Pivots{
		Length:       length,
		NPivots:      nPivots,
		PrevPivots:   []float64{},
		Candidate:    0.0,
		CandidateAge: 0,
		Direction:    Flat,
	}
}

func (zz *Pivots) WarmupPeriod() int {
	return zz.Length
}

func (zz *Pivots) addPivot(pivot float64) {
	zz.PrevPivots = append(zz.PrevPivots, pivot)
}

func (zz *Pivots) WarmUp(vals []float64) error {
	if len(vals) != zz.Length {
		return commonerrs.NewErrExactCount("warmup vals", zz.Length, len(vals))
	}

	lastIdx := len(vals) - 1

	// find most recent pivot(s)
	min, max, minIdx, maxIdx := minMax(vals)
	if minIdx < maxIdx {
		zz.addPivot(min)

		if maxIdx < lastIdx {
			zz.addPivot(max)

			zz.Candidate = vals[lastIdx]
			zz.Direction = Down
		} else {
			zz.Candidate = max
			zz.Direction = Up
		}
	} else {
		zz.addPivot(max)

		if minIdx < lastIdx {
			zz.addPivot(min)

			zz.Candidate = vals[lastIdx]
			zz.Direction = Up
		} else {
			zz.Candidate = min
			zz.Direction = Down
		}
	}

	zz.CandidateAge = 0

	return nil
}

func (zz *Pivots) Update(val float64) {
	switch zz.Direction {
	case Up:
		if val >= zz.Candidate {
			zz.Candidate = val
			zz.CandidateAge = 0
		} else {
			zz.CandidateAge++

			if zz.CandidateAge >= zz.Length {
				zz.addPivot(zz.Candidate)

				zz.Candidate = val
				zz.CandidateAge = 0
				zz.Direction = Down
			}
		}
	case Down:
		if val <= zz.Candidate {
			zz.Candidate = val
			zz.CandidateAge = 0
		} else {
			zz.CandidateAge++

			if zz.CandidateAge >= zz.Length {
				zz.addPivot(zz.Candidate)

				zz.Candidate = val
				zz.CandidateAge = 0
				zz.Direction = Up
			}
		}
	}
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
