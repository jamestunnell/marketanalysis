package collect

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/collection"
	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"
)

type CollectUpdate struct {
	Dir string `json:"dir"`
}

func (cmd *CollectUpdate) Run() error {
	store, err := collection.NewDirStore(cmd.Dir)
	if err != nil {
		return err
	}

	c, err := collection.Load(store)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	i := c.GetInfo()

	log.Info().
		Str("dir", cmd.Dir).
		Interface("info", i).
		Msg("loaded existing collection")

	loc, err := time.LoadLocation(i.TimeZone)
	if err != nil {
		return fmt.Errorf("failed to load location from time zone '%s': %w", i.TimeZone, err)
	}

	var startDate date.Date

	if c.IsEmpty() {
		startDate = i.StartDate
	} else {
		startDate = c.GetLastDate().Add(1)
	}

	start := startDate.In(loc)

	bars, err := GetAlpacaBars(start, i.Symbol, loc)
	if err != nil {
		return err
	}

	if err = c.StoreBars(bars); err != nil {
		return fmt.Errorf("failed to store bars: %w", err)
	}

	log.Info().Msg("stored bars in collection")

	return nil
}
