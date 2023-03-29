package indicators

import "github.com/jamestunnell/marketanalysis/models/bar"

type Indicator interface {
	WarmupPeriod() int
	Initialize(bars []*bar.Bar) error
	Current() float64
	Update(bar *bar.Bar) float64
}
