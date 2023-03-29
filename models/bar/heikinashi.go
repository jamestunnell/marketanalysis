package bar

import "math"

func (b *Bar) HeikinAshi(prev *Bar) *Bar {
	open := 0.5 * (prev.Open + prev.Close)
	close := 0.25 * (b.Open + b.High + b.Low + b.Close)
	high := math.Max(math.Max(b.High, b.Open), b.Close)
	low := math.Max(math.Max(b.Low, b.Open), b.Close)

	bar := New(
		b.Timestamp, open, high, low, close, b.Volume, b.TradeCount, b.VWAP)

	return bar
}
