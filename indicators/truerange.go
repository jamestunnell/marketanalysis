package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models"
)

func TrueRange(cur, prev *models.OHLC) float64 {
	a := cur.High - cur.Low
	b := cur.High - prev.Close
	c := prev.Close - cur.Low

	return math.Max(math.Max(a, b), c)
}
