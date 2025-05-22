package collect

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jamestunnell/marketanalysis/loading"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketdata"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

type Collect struct {
	StartDate date.Date `json:"startDate"`
	Dir       string    `json:"dir"`
	Symbol    string    `json:"symbol"`
	TimeZone  string    `json:"timeZone"`
	Format    string    `json:"format"`

	loc *time.Location
}

var (
	errExists = errors.New("collection already exists, use add command")
)

func (cmd *Collect) Init() error {
	info, err := os.Stat(cmd.Dir)
	if err != nil {
		return fmt.Errorf("failed to stat dir '%s': %w", cmd.Dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("'%s' is not a dir")
	}

	loc, err := time.LoadLocation(cmd.TimeZone)
	if err != nil {
		return fmt.Errorf("failed to load location '%s': %w", cmd.TimeZone, err)
	}

	switch cmd.Format {
	case "csv", "CSV", "ndjson", "NDJSON":
	default:
		return fmt.Errorf("unsupported format %s", cmd.Format)
	}

	cmd.loc = loc

	return nil
}

func (cmd *Collect) Run() error {
	log.Info().Msg("collecting bars")

	ts := timespan.NewTimeSpan(cmd.StartDate.In(cmd.loc), time.Now())

	var err error
	var bars marketdata.Bars

	if bars, err = loading.GetBarsOneMin(cmd.Symbol, ts, cmd.loc); err != nil {
		return fmt.Errorf("failed to load bars: %w", err)
	}

	var fpath string

	switch cmd.Format {
	case "csv", "CSV":
		fpath = path.Join(cmd.Dir, cmd.Symbol+".csv")

		err = models.StoreInFile(fpath, bars.StoreToCSV)
	case "ndjson", "NDJSON":
		fpath = path.Join(cmd.Dir, cmd.Symbol+".ndjson")

		err = models.StoreInFile(fpath, bars.StoreToNDJSON)
	}

	if err != nil {
		return fmt.Errorf("failed to store bars in '%s': %w", fpath, err)
	}

	log.Info().Str("fpath", fpath).Msg("bars stored")

	return nil
}
