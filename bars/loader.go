package bars

import (
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
)

type Loader interface {
	Init() error
	GetBars(timespan.TimeSpan) (models.Bars, error)
	GetDayBars(d date.Date) (models.Bars, error)
	GetRunBars(d date.Date, warmupPeriod int) (models.Bars, error)
}
