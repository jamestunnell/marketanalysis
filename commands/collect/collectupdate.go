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

	*Collector

	startDate date.Date
}

func (cmd *CollectUpdate) Init() error {
	store, err := collection.NewDirStore(cmd.Dir)
	if err != nil {
		return err
	}

	c, err := collection.Load(store)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	i := c.GetInfo()

	cmd.Collector = NewCollector(c)
	cmd.startDate = i.StartDate

	return nil
}

func (cmd *CollectUpdate) Run() error {
	if cmd.coll.IsEmpty() {
		log.Info().Msg("collecting all bars")

		return cmd.CollectBars(cmd.startDate.In(cmd.loc), time.Now())
	}

	if firstDate := cmd.coll.GetFirstDate(); firstDate.After(cmd.startDate) {
		log.Info().Msg("collecting missing bars from the beginning")

		start := cmd.startDate.In(cmd.loc)
		end := firstDate.In(cmd.loc)

		if err := cmd.CollectBars(start, end); err != nil {
			return err
		}
	}

	log.Info().Msg("collecting bars missing from the end")

	start := cmd.coll.GetLastDate().In(cmd.loc)

	return cmd.CollectBars(start, time.Now())
}
