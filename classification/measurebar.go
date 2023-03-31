package classification

import (
	"github.com/jamestunnell/marketanalysis/models/bar"
)

type BarMeasurement struct {
	Bullish                   bool
	Body, TopWick, BottomWick float64
}

func MeasureBar(b *bar.Bar) *BarMeasurement {
	body := b.Close - b.Open
	var m *BarMeasurement
	if body > 0 {
		m = &BarMeasurement{
			Bullish:    true,
			Body:       body,
			TopWick:    b.High - b.Close,
			BottomWick: b.Open - b.Low,
		}
	} else {
		m = &BarMeasurement{
			Bullish:    false,
			Body:       -body,
			TopWick:    b.High - b.Open,
			BottomWick: b.Close - b.Low,
		}
	}

	return m
}
