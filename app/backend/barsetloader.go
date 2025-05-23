package backend

import (
	"context"
	"fmt"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/jamestunnell/marketanalysis/loading"
	"github.com/jamestunnell/marketanalysis/models"
	"github.com/jamestunnell/marketdata"
	"github.com/jamestunnell/marketdata/alpaca"
)

type BarSetLoader struct {
	Symbol string
	Store  Store[*models.BarSet]
	todays marketdata.Bars
}

var locNY = loading.GetLocationNY()

func NewBarSetLoader(
	db *mongo.Database,
	symbol string,
) (*BarSetLoader, error) {
	info := &ResourceInfo{
		KeyName:    "date",
		Name:       "barset",
		NamePlural: "barsets",
	}
	col := db.Collection(symbol)
	store := NewMongoStore[*models.BarSet](info, col)

	today := date.Today().In(locNY)
	ts := timespan.NewTimeSpan(today, today.Add(1))
	bc := alpaca.NewFreeBarCollector(locNY)

	bars, err := bc.Collect(symbol, ts)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's bars: %w", err)
	}

	l := &BarSetLoader{
		Symbol: symbol,
		Store:  store,
		todays: bars,
	}

	return l, nil
}

func (l *BarSetLoader) Load(ctx context.Context, d date.Date) (marketdata.Bars, error) {
	log.Trace().Stringer("date", d).Msg("loading bars")

	if d.Equal(date.TodayIn(locNY)) {
		log.Trace().Msg("skipped loading today's bars")

		return l.todays, nil
	}

	dayBars, appErr := l.Store.Get(ctx, d.String())
	if appErr == nil {
		log.Trace().
			Int("count", len(dayBars.Bars)).
			Stringer("date", d).
			Msg("found bars in store")

		return dayBars.Bars, nil
	}

	ts := timespan.NewTimeSpan(d.In(locNY), d.Add(1).In(locNY))
	bc := alpaca.NewFreeBarCollector(locNY)

	bars, err := bc.Collect(l.Symbol, ts)
	if err != nil {
		return nil, fmt.Errorf("failed to get aplaca bars: %w", err)
	}

	for _, bar := range bars {
		bar.Timestamp = bar.Timestamp.In(locNY)
	}

	log.Trace().
		Int("count", len(bars)).
		Stringer("date", d).
		Msg("loaded bars from alpaca")

	dayBars = &models.BarSet{
		Bars: bars,
		Date: d.String(),
	}

	appErr = l.Store.Create(ctx, dayBars)
	if appErr != nil {
		log.Warn().
			Err(appErr).
			Stringer("date", d).
			Msg("failed to store day bars")
	} else {
		log.Debug().
			Int("count", len(bars)).
			Stringer("date", d).
			Msg("stored bars")
	}

	return dayBars.Bars, nil
}
