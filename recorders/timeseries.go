package recorders

import (
	"sort"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type TimeSeries struct {
	*models.TimeSeries

	loc *time.Location
}

func NewTimeSeries(loc *time.Location) *TimeSeries {
	return &TimeSeries{
		TimeSeries: &models.TimeSeries{
			Quantities: []*models.Quantity{},
		},
		loc: nil,
	}
}

func (rec *TimeSeries) Init(valNames []string) error {
	sort.Strings(valNames)

	rec.Quantities = sliceutils.Map(valNames, func(name string) *models.Quantity {
		return &models.Quantity{
			Name:    name,
			Records: []*models.QuantityRecord{},
		}
	})

	return nil
}

func (rec *TimeSeries) Process(t time.Time, vals map[string]float64) {
	if rec.loc != nil {
		t = t.In(rec.loc)
	}

	for _, q := range rec.Quantities {
		if val, found := vals[q.Name]; found {
			record := &models.QuantityRecord{Timestamp: t, Value: val}

			q.Records = append(q.Records, record)
		}
	}
}

func (rec *TimeSeries) Finalize() error {
	rec.SortByTime()

	return nil
}
