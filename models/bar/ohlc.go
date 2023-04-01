package bar

type OHLC struct {
	Open  float64 `json:"o"`
	High  float64 `json:"h"`
	Low   float64 `json:"l"`
	Close float64 `json:"c"`
}

func (ohlc *OHLC) Float64s() []float64 {
	return []float64{ohlc.Open, ohlc.High, ohlc.Low, ohlc.Close}
}
