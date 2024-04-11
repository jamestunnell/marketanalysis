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
	writer   *csv.Writer
	valNames []string
	record   []string
}

const timeCol = "time"

func NewCSV(w io.Writer) *CSV {
	return &CSV{writer: csv.NewWriter(w)}
}

func (rec *CSV) Init(valNames []string) error {
	rec.valNames = valNames
	rec.record = make([]string, 1+len(valNames))

	header := append([]string{timeCol}, valNames...)

	if err := rec.writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	return nil
}

func (rec *CSV) Record(t time.Time, vals map[string]float64) {
	rec.record[0] = t.Local().String()

	for i, valName := range rec.valNames {
		var valStr string

		if val, found := vals[valName]; found {
			valStr = strconv.FormatFloat(val, 'g', -1, 64)
		} else {
			log.Warn().Str("name", valName).Msg("value not recorded")
		}

		rec.record[1+i] = valStr
	}

	if err := rec.writer.Write(rec.record); err != nil {
		log.Warn().Err(err).Msg("failed to write CSV record")
	}
}

func (rec *CSV) Flush() {
	rec.writer.Flush()
}
