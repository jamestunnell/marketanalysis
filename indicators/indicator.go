package indicators

import "github.com/jamestunnell/marketanalysis/models"

type Indicator interface {
	WarmupPeriod() int
	Initialize(bars []*models.Bar) error
	Current() float64
	Update(bar *models.Bar) float64
}
