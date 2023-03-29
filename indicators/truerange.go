package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models/bar"
)

func TrueRange(cur, prev *bar.Bar) float64 {
	a := cur.High - cur.Low
	b := math.Abs(cur.High - prev.Close)
	c := math.Abs(cur.Low - prev.Close)

	return math.Max(math.Max(a, b), c)
}
