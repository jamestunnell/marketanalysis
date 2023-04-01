package prediction

import (
	"github.com/jamestunnell/marketanalysis/models/bar"
)

func MeasureBar(b *bar.Bar) (body, top, bottom float64) {
	body = (b.Close - b.Open)

	if body > 0 {
		top = b.High - b.Close
		bottom = b.Open - b.Low
	} else {
		top = b.High - b.Open
		bottom = b.Close - b.Low
	}

	return
}
