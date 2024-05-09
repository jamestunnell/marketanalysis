package graph

import (
	"fmt"

	"github.com/rickb777/date"
	"github.com/rs/zerolog/log"

	"github.com/jamestunnell/marketanalysis/bars"
	"github.com/jamestunnell/marketanalysis/blocks"
)

func RunDay(
	d date.Date,
	c *Configuration,
	l bars.Loader,
	r blocks.Recorder,
) error {
	g := New(c)

	if err := g.Init(r); err != nil {
		return fmt.Errorf("failed to init graph: %w", err)
	}

	if err := l.Init(); err != nil {
		return fmt.Errorf("failed to init bars loader: %w", err)
	}

	bars, err := l.GetRunBars(d, g.GetWarmupPeriod())
	if err != nil {
		return fmt.Errorf("failed to get run bars: %w", err)
	}

	log.Debug().
		Stringer("date", d).
		Msgf("running model with %d bars", len(bars))

	for _, bar := range bars {
		g.Update(bar)
	}

	r.Flush()

	return nil
}
