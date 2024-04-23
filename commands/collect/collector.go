package collect

import (
	"fmt"
	"time"

	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rs/zerolog/log"
)

type Collector struct {
	loc  *time.Location
	coll models.Collection
	sym  string
}

func NewCollector(coll models.Collection) *Collector {
	return &Collector{
		loc:  coll.GetLocation(),
		coll: coll,
		sym:  coll.GetInfo().Symbol,
	}
}

func (c *Collector) CollectBars(start, end time.Time) error {
	bars, err := GetAlpacaBars(start, end, c.sym, c.loc)
	if err != nil {
		return fmt.Errorf("failed to get bars: %w", err)
	}

	if len(bars) == 0 {
		return nil
	}

	if err = c.coll.StoreBars(bars); err != nil {
		return fmt.Errorf("failed to store bars: %w", err)
	}

	log.Info().Int("count", len(bars)).Msg("stored bars in collection")

	return nil
}
