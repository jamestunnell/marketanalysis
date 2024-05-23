package model

import (
	"fmt"
	"slices"
	"time"

	"github.com/jamestunnell/marketanalysis/blocks"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
)

type ModelRun struct {
	Collection models.Collection
	ModelFile  string
	Recorder   models.Recorder
	Date       date.Date

	model models.Model
	bars  models.Bars
}

func (cmd *ModelRun) Init() error {
	dr := timespan.NewDateRange(
		cmd.Collection.GetFirstDate(),
		cmd.Collection.GetLastDate().Add(1))
	if !dr.Contains(cmd.Date) {
		return fmt.Errorf("%s is not in collection date range %s", cmd.Date, dr)
	}

	m, err := blocks.LoadGraphModel(cmd.ModelFile)
	if err != nil {
		return fmt.Errorf("failed to load model from file '%s': %w", cmd.ModelFile, err)
	}

	err = m.Init(cmd.Recorder)
	if err != nil {
		return fmt.Errorf("failed to initialize model: %w", err)
	}

	wuPeriod := m.GetWarmupPeriod()
	startTime := cmd.Collection.GetInfo().CoreHours.Open.On(cmd.Date, cmd.Collection.GetLocation())

	bars, err := cmd.Collection.LoadBars(cmd.Date, cmd.Date)
	if err != nil {
		return fmt.Errorf("failed to load bars: %w", err)
	}

	startIdx, found := slices.BinarySearchFunc(bars, startTime, func(b *models.Bar, tgt time.Time) int {
		return b.Timestamp.Compare(tgt)
	})
	if !found {
		return fmt.Errorf("start time %s not in loaded bars", startTime)
	}

	prevDate := cmd.Date.Add(-1)

	// load more warmup bars if needed
	for startIdx < wuPeriod {
		if !dr.Contains(prevDate) {
			log.Warn().
				Int("missing", wuPeriod-startIdx).
				Msg("not enough bars could be loaded for warmup")

			break
		}

		moreBars, err := cmd.Collection.LoadBars(prevDate, prevDate)
		if err != nil {
			return fmt.Errorf("failed to load more bars: %w", err)
		}

		startIdx += len(moreBars)

		bars = append(moreBars, bars...)
	}

	// truncate the slice if we have enough bars for warmup
	if startIdx >= wuPeriod {
		bars = bars[startIdx-wuPeriod:]
	}

	cmd.bars = bars
	cmd.model = m

	return nil
}

func (cmd *ModelRun) Run() error {
	log.Debug().
		Stringer("date", cmd.Date).
		Msgf("running model with %d bars", len(cmd.bars))

	for _, bar := range cmd.bars {
		cmd.model.Update(bar)
	}

	if err := cmd.Recorder.Finalize(); err != nil {
		return fmt.Errorf("failed to finalize recording: %w", err)
	}

	return nil
}
