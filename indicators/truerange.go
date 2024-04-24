package indicators

import (
	"math"

	"github.com/jamestunnell/marketanalysis/models"
)

func TrueRange(cur, prev *models.OHLC) float64 {
	return math.Max(cur.High, prev.Close) - math.Min(cur.Low, prev.Close)
}
