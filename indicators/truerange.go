package indicators

import (
	"math"

	"github.com/jamestunnell/marketdata"
)

func TrueRange(cur, prev *marketdata.OHLC) float64 {
	return math.Max(cur.High, prev.Close) - math.Min(cur.Low, prev.Close)
}
