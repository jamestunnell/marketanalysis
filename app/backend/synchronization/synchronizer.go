package synchronization

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/alpaca"
	"github.com/jamestunnell/marketanalysis/app/backend"
	"github.com/jamestunnell/marketanalysis/app/backend/stores"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketanalysis/util/sliceutils"
)

type synchronizer struct {
	triggers chan *backend.SyncTrigger
	db       *mongo.Database
	stop     chan struct{}
}

func NewSynchronizer(db *mongo.Database) backend.Synchronizer {
	return &synchronizer{
		triggers: make(chan *backend.SyncTrigger, 10),
		db:       db,
		stop:     make(chan struct{}),
	}
}

func CompareDates(a, b date.Date) int {
	return cmp.Compare(a.DaysSinceEpoch(), b.DaysSinceEpoch())
}

func (sync *synchronizer) Trigger(trig *backend.SyncTrigger) {
	sync.triggers <- trig
}

func (sync *synchronizer) Start() {
	go sync.runUntilStopped()
}

func (sync *synchronizer) Stop() {
	sync.stop <- struct{}{}
}

func (sync *synchronizer) scanAllExisting() {
	log.Debug().Msg("sync: scanning all existing")

	store := stores.NewSecurities(sync.db)

	securities, appErr := store.GetAll(context.Background())
	if appErr != nil {
		log.Warn().
			Err(appErr).
			Msg("sync: failed to get securities")

		return
	}

	for _, sec := range securities {
		log.Debug().
			Str("symbol", sec.Symbol).
			Msg("sync: scan existing")

		if err := sync.scan(sec); err != nil {
			log.Warn().
				Err(err).
				Str("symbol", sec.Symbol).
				Msg("sync: scan existing failed")
		}
	}
}

const scanInterval = 30 * time.Second

func (sync *synchronizer) runUntilStopped() {
	log.Debug().Msg("sync: running until stopped")

	sync.scanAllExisting()

	keepGoing := true

	for keepGoing {
		select {
		case trig := <-sync.triggers:
			go sync.Run(trig.Op, trig.Security)
		case <-sync.stop:
			keepGoing = false
		case <-time.After(scanInterval):
			sync.scanAllExisting()
		}
	}

	log.Debug().Msg("sync: stopped")
}

func (sync *synchronizer) Run(op backend.SyncOp, security *models.Security) {
	id := nanoid.Must()

	log.Debug().
		Str("id", id).
		Stringer("op", op).
		Str("symbol", security.Symbol).
		Msg("sync: run triggered")

	var err error

	switch op {
	case backend.SyncAdd:
		sync.add(security)
	case backend.SyncRemove:
		sync.remove(security.Symbol)
	case backend.SyncScan:
		sync.scan(security)
	default:
		err = fmt.Errorf("sync: unknown op type %d", op)
	}

	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("sync: op failed")
	}

	log.Debug().
		Str("id", id).
		Msg("sync: run complete")
}

func (sync *synchronizer) add(sec *models.Security) error {
	loc, err := time.LoadLocation(sec.TimeZone)
	if err != nil {
		return fmt.Errorf("failed to load time zone location '%s': %w", sec.TimeZone, err)
	}

	store := stores.NewBarSets(sec.Symbol, sync.db)

	return sync.addNStartingToday(store, sec, loc, sec.Days)
}

func (sync *synchronizer) remove(symbol string) error {
	err := sync.db.Collection(symbol).Drop(context.Background())
	if err != nil {
		return fmt.Errorf("failed to drop collection for symbol '%s': %w", symbol, err)
	}

	return nil
}

func (sync *synchronizer) scan(sec *models.Security) error {
	loc, err := time.LoadLocation(sec.TimeZone)
	if err != nil {
		return fmt.Errorf("failed to load time zone location '%s': %w", sec.TimeZone, err)
	}

	store := stores.NewBarSets(sec.Symbol, sync.db)

	keys, appErr := store.GetAllKeys(context.Background())
	if appErr != nil {
		return appErr
	}

	if len(keys) == 0 {
		return sync.addNStartingToday(store, sec, loc, sec.Days)
	}

	dates, err := sliceutils.MapErr(keys, func(key string) (date.Date, error) {
		return date.Parse(date.RFC3339, key)
	})
	if err != nil {
		return fmt.Errorf("existing key not formatted as date: %w", err)
	}

	slices.SortFunc(dates, CompareDates)

	oldest := dates[0]
	newest := dates[len(dates)-1]
	today := date.Today()

	err = sync.updateDay(store, sec, loc, newest)
	if err != nil {
		return fmt.Errorf("failed to update newest: %w", err)
	}

	// fill in days after the newest existing
	for d := newest.Add(1); !d.After(today); d = d.Add(1) {
		added, err := sync.addDay(store, sec, loc, d)
		if err != nil {
			return err
		}

		if added {
			dates = append(dates, d)
		}
	}

	// fill in days before the oldest if needed
	for d := oldest.Add(-1); len(dates) < sec.Days; d = d.Add(-1) {
		added, err := sync.addDay(store, sec, loc, d)
		if err != nil {
			return err
		}

		if added {
			dates = append(dates, d)
		}
	}

	if len(dates) == sec.Days {
		return nil
	}

	slices.SortFunc(dates, CompareDates)

	// trim extra
	for _, d := range dates[:len(dates)-sec.Days] {
		if appErr := store.Delete(context.Background(), d.String()); appErr != nil {
			log.Warn().
				Err(appErr).
				Stringer("date", d).
				Msg("sync: failed to trim extra date")

			continue
		}

		log.Debug().
			Stringer("date", d).
			Str("symbol", sec.Symbol).
			Msg("sync: removed extra day")
	}

	return nil
}

func (sync *synchronizer) addNStartingToday(
	store backend.Store[*models.BarSet],
	sec *models.Security,
	loc *time.Location,
	n int,
) error {

	dates := []date.Date{}

	for d := date.Today(); len(dates) < n; d = d.Add(-1) {
		added, err := sync.addDay(store, sec, loc, d)
		if err != nil {
			return err
		}

		if added {
			dates = append(dates, d)
		}
	}

	return nil
}

func (sync *synchronizer) addDay(
	store backend.Store[*models.BarSet],
	sec *models.Security,
	loc *time.Location,
	d date.Date,
) (bool, error) {
	bars, err := sync.loadDay(sec, loc, d)
	if err != nil {
		return false, err
	}

	if len(bars) == 0 {
		return false, nil
	}

	barSet := &models.BarSet{Bars: bars, Date: d.String()}

	if appErr := store.Create(context.Background(), barSet); appErr != nil {
		return false, fmt.Errorf("failed to create bar set for '%s': %w", d, appErr)
	}

	log.Debug().
		Stringer("date", d).
		Int("bars", len(barSet.Bars)).
		Str("symbol", sec.Symbol).
		Msg("sync: added day")

	return true, nil
}

func (sync *synchronizer) updateDay(
	store backend.Store[*models.BarSet],
	sec *models.Security,
	loc *time.Location,
	d date.Date,
) error {
	bars, err := sync.loadDay(sec, loc, d)
	if err != nil {
		return err
	}

	if len(bars) == 0 {
		return nil
	}

	barSet := &models.BarSet{Bars: bars, Date: d.String()}

	if appErr := store.Update(context.Background(), barSet); appErr != nil {
		return fmt.Errorf("failed to update bar set for '%s': %w", d, appErr)
	}

	log.Debug().
		Stringer("date", d).
		Int("bars", len(barSet.Bars)).
		Str("symbol", sec.Symbol).
		Msg("sync: updated day")

	return nil
}

func (sync *synchronizer) loadDay(
	sec *models.Security,
	loc *time.Location,
	d date.Date,
) (models.Bars, error) {
	weekDay := d.Weekday()
	if weekDay == time.Saturday || weekDay == time.Sunday {
		return models.Bars{}, nil
	}

	ts := timespan.NewTimeSpan(d.In(loc), d.Add(1).In(loc))

	bars, err := alpaca.GetBarsOneMin(sec.Symbol, ts, loc)
	if err != nil {
		return models.Bars{}, fmt.Errorf("failed to load bars for '%s': %w", d, err)
	}

	return bars, nil
}
