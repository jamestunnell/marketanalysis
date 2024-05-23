package recorders

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/jamestunnell/marketanalysis/util/sliceutils"
	"github.com/rs/zerolog/log"
)

type JSON struct {
	writer      *bufio.Writer
	loc         *time.Location
	localTZ     string
	quants      []*TimeSeriesQuantity
	recordCount int
}

type TimeSeriesRecording struct {
	Quantities []*TimeSeriesQuantity `json:"quantities"`
}

type TimeSeriesQuantity struct {
	Name    string            `json:"name"`
	Records []*QuantityRecord `json:"records"`
}

type QuantityRecord struct {
	Timestamp time.Time `json:"t"`
	Value     float64   `json:"v"`
}

func NewJSON(w io.Writer, localTZ string) *JSON {
	return &JSON{
		writer:      bufio.NewWriter(w),
		loc:         nil,
		localTZ:     localTZ,
		quants:      []*TimeSeriesQuantity{},
		recordCount: 0,
	}
}

func (rec *JSON) Init(valNames []string) error {
	if rec.localTZ != "" {
		loc, err := time.LoadLocation(rec.localTZ)
		if err != nil {
			return fmt.Errorf("failed to load location from local time zone '%s': %w", rec.localTZ, err)
		}

		rec.loc = loc
	}

	sort.Strings(valNames)

	rec.quants = sliceutils.Map(valNames, func(name string) *TimeSeriesQuantity {
		return &TimeSeriesQuantity{
			Name:    name,
			Records: []*QuantityRecord{},
		}
	})
	rec.recordCount = 0

	return nil
}

func (rec *JSON) Process(t time.Time, vals map[string]float64) {
	if rec.loc != nil {
		t = t.In(rec.loc)
	}

	for _, q := range rec.quants {
		if val, found := vals[q.Name]; found {
			record := &QuantityRecord{Timestamp: t, Value: val}

			q.Records = append(q.Records, record)
		}
	}

	rec.recordCount++
}

func (rec *JSON) Finalize() error {
	recording := &TimeSeriesRecording{
		Quantities: rec.quants,
	}

	d, _ := json.Marshal(recording)

	if _, err := rec.writer.Write(d); err != nil {
		return fmt.Errorf("failed to write recording: %w", err)
	}

	log.Debug().Int("records", rec.recordCount).Msg("Wrote time series recording JSON")

	return nil
}
