package recorders

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type TimeSeriesJSON struct {
	*TimeSeries

	writer *bufio.Writer
}

func NewTimeSeriesJSON(w io.Writer, localTZ string) *TimeSeriesJSON {
	return &TimeSeriesJSON{
		TimeSeries: NewTimeSeries(localTZ),
		writer:     bufio.NewWriter(w),
	}
}

func (rec *TimeSeriesJSON) Finalize() error {
	d, _ := json.Marshal(rec.TimeSeries)

	if _, err := rec.writer.Write(d); err != nil {
		return fmt.Errorf("failed to write recording: %w", err)
	}

	log.Debug().Int("records", rec.recordCount).Msg("Wrote time series recording JSON")

	return nil
}
