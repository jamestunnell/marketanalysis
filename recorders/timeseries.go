package recorders

import (
	"sort"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rs/zerolog/log"
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

	log.Debug().Strs("names", valNames).Msg("recording values")

	rec.Quantities = sliceutils.Map(valNames, func(name string) *models.Quantity {
		return &models.Quantity{
			Name:    name,
			Records: []models.QuantityRecord{},
		}
	})

	return nil
}

func (rec *TimeSeries) Process(tvs map[string]models.TimeValue[float64]) {
	for _, q := range rec.Quantities {
		if tv, found := tvs[q.Name]; found {
			if rec.loc != nil {
				tv.Time = tv.Time.In(rec.loc)
			}

			q.Records = append(q.Records, tv)
		}
	}
}

func (rec *TimeSeries) Finalize() error {
	rec.SortByTime()

	return nil
}
