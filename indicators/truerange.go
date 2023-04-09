package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models"
)

func TrueRange(cur, prev *models.Bar) float64 {
	a := cur.High - cur.Low
	b := math.Abs(cur.High - prev.Close)
	c := math.Abs(cur.Low - prev.Close)

	return math.Max(math.Max(a, b), c)
}
