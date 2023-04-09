package prediction

import (
	"github.com/jamestunnell/marketanalysis/models"
)

type BarMeasure struct {
	Body, Top, Bottom float64
}

func NewBarMeasure(b *models.Bar, atr float64) *BarMeasure {
	body := (b.Close - b.Open)

	var top float64
	var bottom float64

	if body > 0 {
		top = b.High - b.Close
		bottom = b.Open - b.Low
	} else {
		top = b.High - b.Open
		bottom = b.Close - b.Low
	}

	body /= atr
	top /= atr
	bottom /= atr

	return &BarMeasure{
		Body:   body,
		Top:    top,
		Bottom: bottom,
	}
}

func (m *BarMeasure) ToFloat64s() []float64 {
	return []float64{m.Body, m.Top, m.Bottom}
}

func (m *BarMeasure) ToOHLC(atr, lastClose float64) *models.OHLC {
	body := m.Body * atr
	top := m.Top * atr
	bottom := m.Bottom * atr

	o := lastClose
	c := o + body

	var l float64
	var h float64

	if body > 0 {
		h = c + top
		l = o - bottom
	} else {
		h = o + top
		l = c - bottom
	}

	ohlc := &models.OHLC{
		Open:  o,
		High:  h,
		Low:   l,
		Close: c,
	}

	return ohlc
}
