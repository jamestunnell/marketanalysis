package recorders

import (
	"fmt"
	"sort"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type TimeSeries struct {
	*models.TimeSeries

	loc         *time.Location
	localTZ     string
	recordCount int
}

func NewTimeSeries(localTZ string) *TimeSeries {
	return &TimeSeries{
		TimeSeries: &models.TimeSeries{
			Quantities: []*models.Quantity{},
		},
		loc:         nil,
		localTZ:     localTZ,
		recordCount: 0,
	}
}

func (rec *TimeSeries) Init(valNames []string) error {
	if rec.localTZ != "" {
		loc, err := time.LoadLocation(rec.localTZ)
		if err != nil {
			return fmt.Errorf("failed to load location from local time zone '%s': %w", rec.localTZ, err)
		}

		rec.loc = loc
	}

	sort.Strings(valNames)

	rec.Quantities = sliceutils.Map(valNames, func(name string) *models.Quantity {
		return &models.Quantity{
			Name:    name,
			Records: []*models.QuantityRecord{},
		}
	})
	rec.recordCount = 0

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

	rec.recordCount++
}

func (rec *TimeSeries) Finalize() error {
	return nil
}
