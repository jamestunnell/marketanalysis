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

	*Collector
}

var (
	errExists = errors.New("collection already exists, use add command")
)

func (cmd *CollectNew) Init() error {
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

	cmd.Collector = NewCollector(c, loc)

	return nil
}

func (cmd *CollectNew) Run() error {
	log.Info().Msg("collecting all bars")

	return cmd.CollectBars(cmd.StartDate.In(cmd.loc), time.Now())
}
