package pivots

import (
	"time"
)

type Pivot struct {
	Type      PivotType
	Timestamp time.Time
	Value     float64
}

type PivotType int

const (
	PivotNeutral PivotType = iota
	PivotLow
	PivotHigh
)

func (typ PivotType) String() string {
	switch typ {
	case PivotHigh:
		return "high"
	case PivotNeutral:
		return "neutral"
	case PivotLow:
		return "low"
	}

	return ""
}

func NewPivotNeutral(t time.Time, val float64) *Pivot {
	return &Pivot{
		Type:      PivotNeutral,
		Timestamp: t,
		Value:     val,
	}
}

func NewPivotHigh(t time.Time, val float64) *Pivot {
	return &Pivot{
		Type:      PivotHigh,
		Timestamp: t,
		Value:     val,
	}
}

func NewPivotLow(t time.Time, val float64) *Pivot {
	return &Pivot{
		Type:      PivotLow,
		Timestamp: t,
		Value:     val,
	}
}
