package bar

func (b *Bar) HeikinAshi(prev *Bar) *Bar {
	open := 0.5 * (prev.Open + prev.Close)
	close := 0.25 * (b.Open + b.High + b.Low + b.Close)
	high := max3(b.High, b.Open, b.Close)
	low := min3(b.Low, b.Open, b.Close)

	bar := New(
		b.Timestamp, open, high, low, close, b.Volume, b.TradeCount, b.VWAP)

	return bar
}

func max3(a, b, c float32) float32 {
	if a > b {
		if a > c {
			return a
		} else {
			return c
		}
	}

	if b > c {
		return b
	}

	return c
}

func min3(a, b, c float32) float32 {
	if a < b {
		if a < c {
			return a
		} else {
			return c
		}
	}

	if b < c {
		return b
	}

	return c
}
