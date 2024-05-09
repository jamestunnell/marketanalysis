package recorders

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type CSV struct {
	writer     *csv.Writer
	loc        *time.Location
	localTZ    string
	valNames   []string
	record     []string
	notFlushed int
}

const timeCol = "time"

func NewCSV(w io.Writer, localTZ string) *CSV {
	return &CSV{
		writer:     csv.NewWriter(w),
		loc:        nil,
		localTZ:    localTZ,
		valNames:   []string{},
		record:     []string{},
		notFlushed: 0,
	}
}

func (rec *CSV) Init(valNames []string) error {
	if rec.localTZ != "" {
		loc, err := time.LoadLocation(rec.localTZ)
		if err != nil {
			return fmt.Errorf("failed to load location from local time zone '%s': %w", rec.localTZ, err)
		}

		rec.loc = loc
	}

	rec.valNames = valNames
	rec.record = make([]string, 1+len(valNames))

	header := append([]string{timeCol}, valNames...)

	if err := rec.writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	rec.writer.Flush()
	rec.notFlushed = 0

	return nil
}

func (rec *CSV) Record(t time.Time, vals map[string]float64) {
	if rec.loc != nil {
		t = t.In(rec.loc)
	}

	rec.record[0] = t.String()

	missing := []string{}
	for i, valName := range rec.valNames {
		var valStr string

		if val, found := vals[valName]; found {
			valStr = strconv.FormatFloat(val, 'g', -1, 64)
		} else {
			valStr = ""

			missing = append(missing, valName)
		}

		rec.record[1+i] = valStr
	}

	if len(missing) > 0 {
		log.Debug().Strs("names", missing).Msg("CSV: values missing")
	}

	if err := rec.writer.Write(rec.record); err != nil {
		log.Warn().Err(err).Msg("CSV: failed to write record")

		return
	}

	rec.notFlushed++
}

func (rec *CSV) Flush() {
	log.Debug().Int("count", rec.notFlushed).Msg("flushed CSV records")

	rec.writer.Flush()

	rec.notFlushed = 0
}
