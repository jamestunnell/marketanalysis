package recorders

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
)

type NDJSON struct {
	writer     *bufio.Writer
	loc        *time.Location
	localTZ    string
	valNames   []string
	notFlushed int
}

type Record struct {
	Time   time.Time          `json:"timestamp"`
	Values map[string]float64 `json:"values"`
}

func NewNDJSON(w io.Writer, localTZ string) *NDJSON {
	return &NDJSON{
		writer:     bufio.NewWriter(w),
		loc:        nil,
		localTZ:    localTZ,
		valNames:   []string{},
		notFlushed: 0,
	}
}

func (rec *NDJSON) Init(valNames []string) error {
	if rec.localTZ != "" {
		loc, err := time.LoadLocation(rec.localTZ)
		if err != nil {
			return fmt.Errorf("failed to load location from local time zone '%s': %w", rec.localTZ, err)
		}

		rec.loc = loc
	}

	slices.Sort(valNames)

	rec.valNames = valNames
	rec.notFlushed = 0

	return nil
}

func (rec *NDJSON) Record(t time.Time, vals map[string]float64) {
	if rec.loc != nil {
		t = t.In(rec.loc)
	}

	record := &Record{
		Time:   t,
		Values: vals,
	}
	d, _ := json.Marshal(record)

	if _, err := rec.writer.Write(d); err != nil {
		log.Warn().Err(err).Msg("NDJSON: failed to write record, aborting")

		return
	}

	if _, err := rec.writer.WriteRune('\n'); err != nil {
		log.Warn().Err(err).Msg("NDJSON: failed to write newline, aborting")

		return
	}

	rec.notFlushed++
}

func (rec *NDJSON) Flush() {
	log.Debug().Int("count", rec.notFlushed).Msg("flushed NDJSON records")

	rec.writer.Flush()

	rec.notFlushed = 0
}
