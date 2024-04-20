package collect

import (
	"errors"
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"
)

type CollectNew struct {
	StartDate date.Date `json:"startDate"`
	Dir       string    `json:"dir"`
	Symbol    string    `json:"symbol"`
	TimeZone  string    `json:"timeZone"`
}

var (
	errExists = errors.New("collection already exists, use add command")
)

func (cmd *CollectNew) Run() error {
	const fmtMakeStoreFailed = "failed to make store for collection dir '%s': %w"

	loc, err := time.LoadLocation(cmd.TimeZone)
	if err != nil {
		return fmt.Errorf("failed to load location from time zone '%s': %w", cmd.TimeZone, err)
	}

	store, err := collection.NewDirStore(cmd.Dir)
	if err != nil {
		return err
	}

	if collection.Exists(store) {
		return errExists
	}

	start := cmd.StartDate.In(loc)

	bars, err := GetAlpacaBars(start, cmd.Symbol, loc)
	if err != nil {
		return err
	}

	info := &models.CollectionInfo{
		Symbol:     cmd.Symbol,
		Resolution: models.Resolution1Min,
		TimeZone:   cmd.TimeZone,
		StartDate:  cmd.StartDate,
	}

	c, err := collection.New(info, store)
	if err != nil {
		return fmt.Errorf("failed to create new collection: %w", err)
	}

	log.Info().
		Str("dir", cmd.Dir).
		Interface("info", info).
		Msg("created new collection")

	if err = c.StoreBars(bars); err != nil {
		return fmt.Errorf("failed to store bars: %w", err)
	}

	log.Info().Msg("stored bars in collection")

	return nil
}
