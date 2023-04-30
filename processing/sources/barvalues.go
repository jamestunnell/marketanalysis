package sources

import "github.com/jamestunnell/marketanalysis/models"

const (
	BarValueName = "BarValue"

	BarValueClose = "close"
	BarValueOpen  = "open"
	BarValueHigh  = "high"
	BarValueLow   = "low"

	BarValueHL2   = "hl2"
	BarValueHLC3  = "hlc3"
	BarValueOCC3  = "occ3"
	BarValueOHLC4 = "ohlc4"
	BarValueHLCC4 = "hlcc4"

	OneHalf   = 1.0 / 2.0
	OneThird  = 1.0 / 3.0
	OneFourth = 1.0 / 4.0
)

func BarValueTypes() []string {
	return []string{
		BarValueClose,
		BarValueOpen,
		BarValueHigh,
		BarValueLow,
		BarValueHL2,
		BarValueHLC3,
		BarValueOCC3,
		BarValueOHLC4,
		BarValueHLCC4,
	}
}

func BarValue(typ string, ohlc *models.OHLC) float64 {
	var val float64

	switch typ {
	case BarValueClose:
		val = ohlc.Close
	case BarValueOpen:
		val = ohlc.Open
	case BarValueHigh:
		val = ohlc.High
	case BarValueLow:
		val = ohlc.Low
	case BarValueHL2:
		val = OneHalf * (ohlc.High + ohlc.Low)
	case BarValueHLC3:
		val = OneThird * (ohlc.High + ohlc.Low + ohlc.Close)
	case BarValueOCC3:
		val = OneThird * (ohlc.Open + ohlc.Close + ohlc.Close)
	case BarValueOHLC4:
		val = OneFourth * (ohlc.Open + ohlc.High + ohlc.Low + ohlc.Close)
	case BarValueHLCC4:
		val = OneFourth * (ohlc.High + ohlc.Low + ohlc.Close + ohlc.Close)
	}

	return val
}
